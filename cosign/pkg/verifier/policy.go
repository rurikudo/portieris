//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Implementation of verify against containers/image policy interface

package verifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	kube "github.com/IBM/portieris/cosign/pkg/kubernetes"
	"github.com/golang/glog"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/cosign/fulcio"
	"github.com/sigstore/sigstore/pkg/signature/payload"

	dconfig "github.com/docker/cli/cli/config"
)

func VerifyByPolicy(kWrapper kube.WrapperInterface, imageToVerify string, namespace string, key string, keyNamespace string) (*bytes.Buffer, error, error) {
	//  return digest, deny, err

	if key == "" {
		return nil, nil, fmt.Errorf("KeySecret missing in signedBy requirement")
	}

	secretNamespace := namespace
	// Override the default namespace behavior if a namespace was provided in this policy
	if keyNamespace != "" {
		secretNamespace = keyNamespace
	}
	secretBytes, err := kWrapper.GetSecretKey(secretNamespace, key)
	if err != nil {
		return nil, nil, err
	}

	//

	keyData, err := decodeArmoredKey(secretBytes)
	if err != nil {
		return nil, nil, err
	}
	glog.Infof("PublicKey... %v", keyData)
	glog.Infof("DOCKER_CONFIG... %v", os.Getenv("DOCKER_CONFIG"))

	// cosign verify
	co := cosign.CheckOpts{
		Claims: true,
		Tlog:   cosign.Experimental(),
		Roots:  fulcio.Roots,
		PubKey: keyData,
	}
	glog.Infof("CheckOpts... %v", co)

	glog.Infof("imageToVerify... %v", imageToVerify)
	ref, err := name.ParseReference(imageToVerify)
	if err != nil {
		return nil, nil, err
	}
	glog.Infof("ParseReference... %v", ref)

	// test
	cf, err := dconfig.Load(os.Getenv("DOCKER_CONFIG"))
	if err != nil {
		glog.Infof("err config.Load... %v", err)
	}
	glog.Infof("config... %v", cf)
	glog.Infof("start remote.....")
	op := remote.WithAuthFromKeychain(authn.DefaultKeychain)
	glog.Infof("WithAuthFromKeychain... %v", op)
	img, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		panic(err)
	}
	fmt.Println(img.Digest)
	// test

	verified, err := cosign.Verify(context.Background(), ref, co)
	if err != nil {
		glog.Infof("cosign.Verify err... %v", err)
		return nil, nil, err
	}

	// cosign return payload like this
	// [{"critical":{"identity":{"docker-reference":"us.icr.io/mutation-advisor/sigstore-image-sign-riko"},"image":{"docker-manifest-digest":"sha256:5403064f94b617f7975a19ba4d1a1299fd584397f6ee4393d0e16744ed11aab1"},"type":"cosign container image signature"},"optional":{"CommonName":"ruri.a28@gmail.com"}}]
	for _, vp := range verified {
		ss := payload.Simple{}
		err := json.Unmarshal(vp.Payload, &ss)
		if err != nil {
			fmt.Println("error decoding the payload:", err.Error())
			return nil, nil, err
		}
		digest := ss.Critical.Image.DockerManifestDigest
		glog.Infof("digest %v", digest)
		return bytes.NewBufferString(strings.TrimPrefix(digest, "sha256:")), nil, nil
	}
	return nil, nil, nil
}

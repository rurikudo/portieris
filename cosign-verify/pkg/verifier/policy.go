//
// Copyright 2021 IBM Corporation
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

package verifier

import (
	"context"
	"encoding/json"
	"fmt"

	kube "github.com/IBM/portieris/cosign-verify/pkg/kubernetes"
	"github.com/golang/glog"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/cosign/fulcio"
	"github.com/sigstore/sigstore/pkg/signature/payload"
)

func VerifyByPolicy(kWrapper kube.WrapperInterface, imageToVerify ImageToVerify) (string, string, error, error) {
	//  return common-name, digest, deny, err
	key := imageToVerify.Key
	keyNamespace := imageToVerify.KeyNamespace
	namespace := imageToVerify.Namespace
	imgName := imageToVerify.Image
	commonName := imageToVerify.CommonName
	cosign_Experimental := imageToVerify.TransparencyLog
	// cosign option
	co := cosign.CheckOpts{
		Claims: true,
		Tlog:   cosign_Experimental,
		Roots:  fulcio.Roots,
	}
	glog.Infof("cosign_Experimental... %v", cosign_Experimental)

	if key != "" {
		secretNamespace := namespace
		// Override the default namespace behavior if a namespace was provided in this policy
		if keyNamespace != "" {
			secretNamespace = keyNamespace
		}
		secretBytes, err := kWrapper.GetSecretKey(secretNamespace, key)
		if err != nil {
			return "", "", nil, err
		}

		keyData, err := decodeArmoredKey(secretBytes)
		if err != nil {
			return "", "", nil, err
		}
		co.PubKey = keyData
	}

	glog.Infof("image to verify: %v", imgName)
	ref, err := name.ParseReference(imgName)
	if err != nil {
		return "", "", nil, err
	}

	verified, err := cosign.Verify(context.Background(), ref, co)
	if err != nil {
		glog.Infof("cosign verify err: %v", err)
		return "", "", nil, err
	}

	for _, vp := range verified {
		ss := payload.Simple{}
		err := json.Unmarshal(vp.Payload, &ss)
		if err != nil {
			fmt.Println("error decoding the payload:", err.Error())
			return "", "", nil, err
		}
		digest := ss.Critical.Image.DockerManifestDigest
		glog.Infof("digest: %v", digest)

		cn := vp.Cert.Subject.CommonName
		// check signer
		if commonName != cn {
			glog.Infof("Not match with CommonName in CosignRequirement %v: %v", commonName, cn)
			return cn, digest, fmt.Errorf("Not match with CommonName in CosignRequirement %v: %v", commonName, cn), nil
		}
		return cn, digest, nil, nil
	}
	return "", "", fmt.Errorf("SignedPayload is empty"), nil
}
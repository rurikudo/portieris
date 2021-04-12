// Copyright 2020, 2021 Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Implementation of verify against containers/image policy interface

package cosign

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	policyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/IBM/portieris/pkg/kubernetes"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/cosign/fulcio"
	"github.com/sigstore/sigstore/pkg/signature/payload"
)

func (v verifier) VerifyByPolicy(kWrapper kubernetes.WrapperInterface, imageToVerify string, namespace string, inPolicies []policyv1.CosignRequirement) (*bytes.Buffer, error, error) {
	//  return digest, deny, err
	for _, inPolicy := range inPolicies {
		if inPolicy.KeySecret == "" {
			return nil, nil, fmt.Errorf("KeySecret missing in signedBy requirement")
		}

		secretNamespace := namespace
		// Override the default namespace behavior if a namespace was provided in this policy
		if inPolicy.KeySecretNamespace != "" {
			secretNamespace = inPolicy.KeySecretNamespace
		}
		secretBytes, err := kWrapper.GetSecretKey(secretNamespace, inPolicy.KeySecret)
		if err != nil {
			return nil, nil, err
		}

		keyData, err := decodeArmoredKey(secretBytes)
		if err != nil {
			return nil, nil, err
		}

		// cosign verify
		ref, err := name.ParseReference(imageToVerify)
		if err != nil {
			return nil, nil, err
		}
		verified, err := cosign.Verify(context.Background(), ref, cosign.CheckOpts{
			Roots:  fulcio.Roots,
			PubKey: keyData,
			Claims: true,
		})
		if err != nil {
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
			return bytes.NewBufferString(strings.TrimPrefix(digest, "sha256:")), nil, nil
		}

	}
	// to be fixed
	return nil, nil, nil
}

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

func VerifyByCosign(kWrapper kube.WrapperInterface, imgName, key, keyNamespace string, cosign_Experimental bool) ([]string, string, error, error) {
	//  return commonName, digest, deny, err
	// cosign option
	co := cosign.CheckOpts{
		Claims: true,
		Tlog:   cosign_Experimental,
		Roots:  fulcio.Roots,
	}
	glog.Infof("cosign_Experimental... %v", cosign_Experimental)

	if key != "" {
		secretNamespace := keyNamespace
		secretBytes, err := kWrapper.GetSecretKey(secretNamespace, key)
		if err != nil {
			return nil, "", nil, err
		}

		keyData, err := decodeArmoredKey(secretBytes)
		if err != nil {
			return nil, "", nil, err
		}
		co.PubKey = keyData
	}

	glog.Infof("image to verify: %v", imgName)
	ref, err := name.ParseReference(imgName)
	if err != nil {
		return nil, "", nil, err
	}

	verified, err := cosign.Verify(context.Background(), ref, co)
	if err != nil {
		glog.Infof("cosign verify err: %v", err)
		return nil, "", nil, err
	}

	if len(verified) == 0 {
		glog.Infof("[]cosign.SignedPayload is empty: no valid signature")
		return nil, "", fmt.Errorf("no valid signature"), nil
	}

	var commonNames []string
	var digest string
	for _, vp := range verified {
		ss := payload.Simple{}
		err := json.Unmarshal(vp.Payload, &ss)
		if err != nil {
			fmt.Println("error decoding the payload:", err.Error())
			return nil, "", nil, err
		}
		digest = ss.Critical.Image.DockerManifestDigest
		glog.Infof("digest: %v", digest)

		cn := vp.Cert.Subject.CommonName
		commonNames = append(commonNames, cn)
	}
	return commonNames, digest, nil, nil
}

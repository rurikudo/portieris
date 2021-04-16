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

package multi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	policyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/golang/glog"
)

const cosignVerifierUrl = "http://localhost:8081"

type cosignVerifyInput struct {
	Image              string `json:"image"`
	Namespace          string `json:"namespace"`
	Key                string `json:"key"`
	KeySecretNamespace string `json:"keyNamespace"`
	TransparencyLog    bool   `json:"transparencyLog"`
}

type cosignVerifyResult struct {
	Deny       error  `json:"deny"`
	Err        error  `json:"err"`
	Digest     string `json:"digest"`
	CommonName string `json:"commonName"`
}

func cosignVerify(img string, namespace string, policy policyv1.CosignRequirement, transparencyLog bool) (string, *bytes.Buffer, error, error) {
	cvInput := cosignVerifyInput{
		Image:              img,
		Namespace:          namespace,
		Key:                policy.KeySecret,
		KeySecretNamespace: policy.KeySecretNamespace,
		TransparencyLog:    transparencyLog,
	}
	cosignVerifyInputJson, _ := json.Marshal(cvInput)

	cvUrl := cosignVerifierUrl
	glog.Infof("http.Cosign %v", cosignVerifierUrl)
	cvResponse, err := http.Post(cvUrl, "application/json", bytes.NewBuffer([]byte(cosignVerifyInputJson)))
	if err != nil {
		return "", nil, nil, err
	}
	if cvResponse.StatusCode != 200 {
		errMsg := "Error reported from CosignVerifier"
		return "", nil, nil, errors.New(errMsg)
	}

	var cvres cosignVerifyResult
	body, err := ioutil.ReadAll(cvResponse.Body)
	if err != nil {
		fmt.Println("Request error:", err)
		return "", nil, nil, err
	}
	str_json := string(body)
	err = json.Unmarshal([]byte(str_json), &cvres)
	if err != nil {
		fmt.Println(err)
		return "", nil, nil, err
	}

	// glog.Infof("cosignVerifyResult: %v", cvres)
	dres := bytes.NewBufferString(strings.TrimPrefix(cvres.Digest, "sha256:"))
	return cvres.CommonName, dres, cvres.Deny, cvres.Err
}

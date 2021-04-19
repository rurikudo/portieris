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
	Key                string `json:"key"`
	KeySecretNamespace string `json:"keyNamespace"`
	TransparencyLog    bool   `json:"transparencyLog"`
}

type cosignVerifyResult struct {
	Deny       string   `json:"deny"`
	Err        string   `json:"err"`
	Digest     string   `json:"digest"`
	CommonName []string `json:"commonName"`
}

func CosignVerify(img string, namespace string, policy policyv1.CosignRequirement, transparencyLog bool) (string, *bytes.Buffer, error, error) {
	cvInput := cosignVerifyInput{
		Image:              img,
		Key:                policy.KeySecret,
		KeySecretNamespace: policy.KeySecretNamespace,
		TransparencyLog:    transparencyLog,
	}
	if policy.KeySecretNamespace == "" {
		cvInput.KeySecretNamespace = namespace
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
	digest_res := bytes.NewBufferString(strings.TrimPrefix(cvres.Digest, "sha256:"))
	var deny_res error
	var err_res error
	if cvres.Deny != "" {
		deny_res = fmt.Errorf(cvres.Deny)
	}
	if cvres.Err != "" {
		err_res = fmt.Errorf(cvres.Err)
	}
	cn, deny := checkCommonName(cvres.CommonName, policy.CommonName)
	if deny != nil {
		return "", nil, deny, nil
	}
	return cn, digest_res, deny_res, err_res
}

func checkCommonName(results []string, expected string) (string, error) {
	for _, cn := range results {
		if cn == expected {
			return cn, nil
		}
	}
	cn := strings.Join(results, ",")
	return "", fmt.Errorf("Not match with CommonName in CosignRequirement %v: %v", cn, expected)
}

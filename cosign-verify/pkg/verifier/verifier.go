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

package verifier

import (
	"encoding/json"

	kube "github.com/IBM/portieris/cosign-verify/pkg/kubernetes"
	"github.com/golang/glog"
)

func Verifier(imageToVerify ImageToVerify) *VerifyResult {
	//  Verifier(imageToVerify ImageToVerify,kubeconfig *string)
	// kubeconfig := flag.String("kubeconfig", "", "location of kubeconfig file to use for an out-of-cluster kube client configuration")
	kubeClientConfig := kube.GetKubeClientConfig()
	kubeClientset := kube.GetKubeClient(kubeClientConfig)
	kubeWrapper := kube.NewKubeClientsetWrapper(kubeClientset)
	commonName, digest, deny, err := VerifyByPolicy(kubeWrapper, imageToVerify)
	glog.Infof("digest... %v", digest)
	var err_s string
	var deny_s string
	if err != nil {
		err_s = err.Error()
	}
	if deny != nil {
		deny_s = deny.Error()
	}
	vres := &VerifyResult{
		Deny:       deny_s,
		Digest:     digest,
		CommonName: commonName,
		Err:        err_s,
	}
	e, err := json.Marshal(vres)
	glog.Infof("VerifyResult... %v", string(e))
	// glog.Infof("VerifyResult... %v", vres)
	return vres
}

type ImageToVerify struct {
	Image           string `json:"image"`
	Namespace       string `json:"namespace"`
	Key             string `json:"key"`
	KeyNamespace    string `json:"keyNamespace"`
	CommonName      string `json:"commonName"`
	TransparencyLog bool   `json:"transparencyLog"`
}

type VerifyResult struct {
	Deny       string `json:"deny"`
	Err        string `json:"err"`
	Digest     string `json:"digest"`
	CommonName string `json:"commonName"`
}

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

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/IBM/portieris/cosign-verify/pkg/verifier"
	"github.com/golang/glog"
)

var kubeconfig *string = flag.String("kubeconfig", "", "location of kubeconfig file to use for an out-of-cluster kube client configuration")

type ImageToVerify struct {
	Image           string `json:"image"`
	Key             string `json:"key"`
	KeyNamespace    string `json:"keyNamespace"`
	TransparencyLog bool   `json:"transparencyLog"`
}

type VerifyResult struct {
	Deny       string   `json:"deny"`
	Err        string   `json:"err"`
	Digest     string   `json:"digest"`
	CommonName []string `json:"commonName"`
}

func CosignVerify(w http.ResponseWriter, r *http.Request) {
	glog.Infof("cosign-verifier is called....")
	var imageToVerify ImageToVerify
	json.NewDecoder(r.Body).Decode(&imageToVerify)
	// input
	image := imageToVerify.Image
	key := imageToVerify.Key
	keyNamespace := imageToVerify.KeyNamespace
	cosign_Experimental := imageToVerify.TransparencyLog
	commonName, digest, deny, err := verifier.Verifier(image, key, keyNamespace, cosign_Experimental, kubeconfig)
	// output
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
	res, _ := json.Marshal(vres)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func handleRequests() {
	http.HandleFunc("/", CosignVerify)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	// flag.Parse() // glog flags
	handleRequests()
}

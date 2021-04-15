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
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/portieris/cosign/pkg/verifier"
	"github.com/golang/glog"
)

type ImageToVerify struct {
	Image        string                `json:"image"`
	Namespace    string                `json:"namespace"`
	Key          string                `json:"key"`
	KeyNamespace string                `json:"keyNamespace"`
	Credential   []verifier.Credential `json:"credential"`
}

type VerifyResult struct {
	Deny   error        `json:"deny"`
	Err    error        `json:"err"`
	Digest bytes.Buffer `json:"digest"`
}

func CosignVerify(w http.ResponseWriter, r *http.Request) {
	glog.Infof("cosign-verifier is called....")
	var imageToVerify ImageToVerify
	json.NewDecoder(r.Body).Decode(&imageToVerify)
	digest, deny, err := verifier.Verifier(imageToVerify.Image, imageToVerify.Namespace, imageToVerify.Key, imageToVerify.KeyNamespace)
	vres := VerifyResult{
		Deny:   deny,
		Err:    err,
		Digest: *digest,
	}
	res, err := json.Marshal(vres)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func handleRequests() {
	http.HandleFunc("/", CosignVerify)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}

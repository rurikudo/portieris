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
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/portieris/cosign-verify/pkg/verifier"
	"github.com/golang/glog"
)

// var kubeconfig *string = flag.String("kubeconfig", "", "location of kubeconfig file to use for an out-of-cluster kube client configuration")

func CosignVerify(w http.ResponseWriter, r *http.Request) {
	glog.Infof("cosign-verifier is called....")
	var imageToVerify verifier.ImageToVerify
	json.NewDecoder(r.Body).Decode(&imageToVerify)
	// vres := verifier.Verifier(imageToVerify, kubeconfig)
	vres := verifier.Verifier(imageToVerify)
	res, _ := json.Marshal(vres)
	// glog.Infof("cosign-verifier res.... %v", res)
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

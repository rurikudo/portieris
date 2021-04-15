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
	"bytes"
	"flag"

	kube "github.com/IBM/portieris/cosign/pkg/kubernetes"
)

func Verifier(imageToVerify string, namespace string, keySecret string, keySecretNamespace string) (*bytes.Buffer, error, error) {
	kubeconfig := flag.String("kubeconfig", "", "location of kubeconfig file to use for an out-of-cluster kube client configuration")
	kubeClientConfig := kube.GetKubeClientConfig(kubeconfig)
	kubeClientset := kube.GetKubeClient(kubeClientConfig)
	kubeWrapper := kube.NewKubeClientsetWrapper(kubeClientset)
	digest, deny, err := VerifyByPolicy(kubeWrapper, imageToVerify, namespace, keySecret, keySecretNamespace)
	return digest, deny, err
}

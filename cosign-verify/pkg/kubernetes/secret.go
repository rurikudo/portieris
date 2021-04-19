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

package kube

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Auths struct contains an embedded RegistriesStruct of name auths
type Auths struct {
	Registries RegistriesStruct `json:"auths"`
}

// RegistriesStruct is a map of registries to their credentials
type RegistriesStruct map[string]RegistryCredentials

// RegistryCredentials defines the fields stored per registry in an docker config secret
type RegistryCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     string `json:"auth"`
}

// GetSecretKey obtains the "key" data from the named secret
func (w *Wrapper) GetSecretKey(namespace, secretName string) ([]byte, error) {
	// Retrieve secret
	secret, err := w.CoreV1().Secrets(namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		glog.Error("Error: ", err)
		return nil, err
	}
	glog.Infof("Found secret %s", secretName)
	// Obtain the key data
	if key, ok := secret.Data["key"]; ok {
		return key, nil
	}
	return nil, fmt.Errorf("Secret %q in %q does not contain a \"key\" attribute", secretName, namespace)
}

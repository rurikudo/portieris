// Copyright 2018, 2021 Portieris Authors.
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

package cosign

import (
	"bytes"

	policyv1 "github.com/rurikudo/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	"github.com/rurikudo/portieris/pkg/kubernetes"
)

// Verifier is for verifying simple signing
type Verifier interface {
	VerifyByPolicy(kWrapper kubernetes.WrapperInterface, imageToVerify string, namespace string, inPolicies []policyv1.CosignRequirement) (*bytes.Buffer, error, error)
}

type verifier struct{}

// NewVerifier creates a new Verfier
func NewVerifier() Verifier {
	return &verifier{}
}

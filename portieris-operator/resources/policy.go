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

package resources

import (
	policyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	apiv1alpha1 "github.com/IBM/portieris/portieris-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BuildDefaultImagePolicyForPortieris(cr *apiv1alpha1.Portieris) *policyv1.ImagePolicy {
	policy := &policyv1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: cr.Namespace,
		},
		Spec: policyv1.ImagePolicySpec{
			Repositories: []policyv1.Repository{
				{
					Name: cr.Spec.Image.Host + "/" + cr.Spec.Image.Image,
				},
			},
		},
	}
	return policy
}

func BuildKubeSystemImagePolicyForPortieris(cr *apiv1alpha1.Portieris) *policyv1.ImagePolicy {
	var repo []policyv1.Repository
	if cr.Spec.IBMContainerService {
		repo = []policyv1.Repository{
			{
				Name: "*",
			},
			{
				Name: "registry*.bluemix.net/armada/*",
			},
			{
				Name: "registry*.bluemix.net/armada-worker/*",
			},
			{
				Name: "registry*.bluemix.net/armada-master/*",
			},
			{
				Name: "*.icr.io/armada/*",
			},
			{
				Name: "*.icr.io/armada-worker/*",
			},
			{
				Name: "*.icr.io/armada-master/*",
			},
			{
				Name: "icr.io/armada/*",
			},
			{
				Name: "icr.io/armada-worker/*",
			},
			{
				Name: "icr.io/armada-master/*",
			},
		}
	} else {
		repo = []policyv1.Repository{
			{
				Name: "*",
			},
		}
	}

	policy := &policyv1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "kube-system",
		},
		Spec: policyv1.ImagePolicySpec{
			Repositories: repo,
		},
	}
	return policy
}

func BuildIBMSystemImagePolicyForPortieris(cr *apiv1alpha1.Portieris) *policyv1.ImagePolicy {
	policy := &policyv1.ImagePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "ibm-system",
		},
		Spec: policyv1.ImagePolicySpec{
			Repositories: []policyv1.Repository{
				{
					Name: "registry*.bluemix.net/armada/*",
				},
				{
					Name: "registry*.bluemix.net/armada-worker/*",
				},
				{
					Name: "registry*.bluemix.net/armada-master/*",
				},
				{
					Name: "*.icr.io/armada/*",
				},
				{
					Name: "*.icr.io/armada-worker/*",
				},
				{
					Name: "*.icr.io/armada-master/*",
				},
				{
					Name: "icr.io/armada/*",
				},
				{
					Name: "icr.io/armada-worker/*",
				},
				{
					Name: "icr.io/armada-master/*",
				},
				{
					Name: "icr.io/ext/istio/*",
				},
			},
		},
	}
	return policy
}

func BuildClusterImagePolicyForPortieris(cr *apiv1alpha1.Portieris) *policyv1.ClusterImagePolicy {
	var repositories []policyv1.Repository
	for _, r := range cr.Spec.AllowedRepositories {
		repository := policyv1.Repository{
			Name: r,
		}
		repositories = append(repositories, repository)
	}

	policy := &policyv1.ClusterImagePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: cr.Namespace,
		},
		Spec: policyv1.ImagePolicySpec{
			Repositories: repositories,
		},
	}
	return policy
}

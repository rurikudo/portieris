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
	apiv1alpha1 "github.com/IBM/portieris/portieris-operator/api/v1alpha1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildCRD(name, namespace string, crdNames extv1.CustomResourceDefinitionNames) *extv1.CustomResourceDefinition {
	xPreserve := true
	newCRD := &extv1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: extv1.CustomResourceDefinitionSpec{
			Group: "portieris.cloud.ibm.com",
			//Version: "v1",
			Names: crdNames,
			Scope: "Namespaced",
			Validation: &extv1.CustomResourceValidation{
				OpenAPIV3Schema: &extv1.JSONSchemaProps{
					Type:                   "object",
					XPreserveUnknownFields: &xPreserve,
				},
			},
			Versions: []extv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
		},
	}
	return newCRD
}

func buildClusterScopeCRD(name, namespace string, crdNames extv1.CustomResourceDefinitionNames) *extv1.CustomResourceDefinition {
	xPreserve := true
	newCRD := &extv1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: extv1.CustomResourceDefinitionSpec{
			Group: "portieris.cloud.ibm.com",
			//Version: "v1",
			Names: crdNames,
			Scope: "Cluster",
			Validation: &extv1.CustomResourceValidation{
				OpenAPIV3Schema: &extv1.JSONSchemaProps{
					Type:                   "object",
					XPreserveUnknownFields: &xPreserve,
				},
			},
			Versions: []extv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
		},
	}
	return newCRD
}

//cluster image policy crd
func BuildClusterImagePolicyCRD(cr *apiv1alpha1.Portieris) *extv1.CustomResourceDefinition {
	crdNames := extv1.CustomResourceDefinitionNames{
		Kind:     "ClusterImagePolicy",
		Plural:   "clusterimagepolicies",
		ListKind: "ClusterImagePolicyList",
		Singular: "clusterimagepolicy",
	}
	return buildClusterScopeCRD("clusterimagepolicies.portieris.cloud.ibm.com", cr.Namespace, crdNames)
}

//image policy crd
func BuildImagePolicyCRD(cr *apiv1alpha1.Portieris) *extv1.CustomResourceDefinition {
	crdNames := extv1.CustomResourceDefinitionNames{
		Kind:     "ImagePolicy",
		Plural:   "imagepolicies",
		ListKind: "ImagePolicyList",
		Singular: "imagepolicy",
	}
	return buildCRD("imagepolicies.portieris.cloud.ibm.com", cr.Namespace, crdNames)
}

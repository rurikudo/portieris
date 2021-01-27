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
	apiv1alpha1 "github.com/rurikudo/portieris/portieris-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// //cert manager certificate
// func BuildCertificateForPortieris(cr *apiv1alpha1.Portieris) *corev1.Secret {
// 	apiVersion: cert-manager.io/v1alpha2
// 	kind: Certificate
// 	metadata:
// 	name: portieris-certs
// 	namespace: {{ .Release.Namespace }}
// 	spec:
// 	dnsNames:
// 		- portieris.{{ .Release.Namespace }}.svc
// 	secretName: portieris-certs
// 	issuerRef:
// 		name: portieris
// }

// //cert manager issuer
// func BuildIssuerForPortieris(cr *apiv1alpha1.Portieris) *corev1.Secret {
// 	apiVersion: cert-manager.io/v1alpha2
// 	kind: Issuer
// 	metadata:
// 	name: portieris
// 	namespace: {{ .Release.Namespace }}
// 	spec:
// 	selfSigned: {}
// }

// ishield-server-tls
func BuildTlsSecretForPortieris(cr *apiv1alpha1.Portieris) *corev1.Secret {
	var empty []byte
	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetWebhookServerTlsSecretName(),
			Namespace: cr.Namespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       empty, // "tls.crt"
			corev1.TLSPrivateKeyKey: empty,
			"ca.crt":                empty,
		},
		Type: corev1.SecretTypeTLS,
	}
	return sec
}

//func scc {}

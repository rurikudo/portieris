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
	admv1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//service
func BuildServiceForPortieris(cr *apiv1alpha1.Portieris) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "portieris",
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       cr.Spec.Service.Port,
					Name:       "https",
					Protocol:   v1.ProtocolTCP,
					TargetPort: cr.Spec.Service.TargetPort,
				},
				{
					Port:       cr.Spec.Service.MetricsPort,
					Name:       "metrics",
					Protocol:   v1.ProtocolTCP,
					TargetPort: cr.Spec.Service.MetricsTargetPort,
				},
			},
			Selector: cr.Spec.SelectorLabels,
		},
	}
	return svc
}

//webhook configuration
func BuildMutatingWebhookConfigurationForPortieris(cr *apiv1alpha1.Portieris) *admv1.MutatingWebhookConfiguration {

	var path *string
	mutate := "/admit"
	path = &mutate

	var empty []byte
	sideEffect := admv1.SideEffectClassNone

	var namespaceSelector *metav1.LabelSelector
	if cr.Spec.AllowAdmissionSkip {
		matchExpressions := []metav1.LabelSelectorRequirement{
			{
				Key:      "securityenforcement.admission.cloud.ibm.com/namespace",
				Operator: "NotIn",
				Values:   []string{"skip"},
			},
		}
		namespaceSelector.MatchExpressions = matchExpressions
	}

	wc := &admv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "image-admission-config",
			Namespace: cr.Namespace,
		},
		Webhooks: []admv1.MutatingWebhook{
			{
				Name: "trust.hooks.securityenforcement.admission.cloud.ibm.com",
				ClientConfig: admv1.WebhookClientConfig{
					Service: &admv1.ServiceReference{
						Name:      cr.GetWebhookServiceName(),
						Namespace: cr.Namespace,
						Path:      path, //"/admit"
					},
					CABundle: empty,
				},
				Rules: []admv1.RuleWithOperations{
					{
						Operations: []admv1.OperationType{
							admv1.Create, admv1.Update,
						},
						Rule: admv1.Rule{
							APIGroups:   []string{"*"},
							APIVersions: []string{"*"},
							Resources:   []string{"pods", "deployments", "replicationcontrollers", "replicasets", "daemonsets", "statefulsets", "jobs", "cronjobs"},
						},
					},
				},
				SideEffects:       &sideEffect,
				FailurePolicy:     cr.Spec.WebHooks.FailurePolicy,
				NamespaceSelector: namespaceSelector,
			},
		},
	}
	return wc
}

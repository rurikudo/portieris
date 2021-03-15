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

package v1alpha1

import (
	admv1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	CleanupFinalizerName = "cleanup.finalizers.portieris.io"
)

// PortierisSpec defines the desired state of Portieris
type PortierisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name                       string                  `json:"name,omitempty"`
	MetaLabels                 map[string]string       `json:"labels,omitempty"`
	SelectorLabels             map[string]string       `json:"selector,omitempty"`
	ReplicaCount               *int32                  `json:"replicaCount,omitempty"`
	Image                      Image                   `json:"image,omitempty"`
	Service                    Service                 `json:"service,omitempty"`
	SecurityContext            *v1.PodSecurityContext  `json:"securityContext,omitempty"`
	WebHooks                   WebHooks                `json:"webHooks,omitempty"`
	IBMContainerService        bool                    `json:"IBMContainerService,omitempty"`
	SkipSecretCreation         bool                    `json:"skipSecretCreation,omitempty"`
	SecretName                 string                  `json:"secretName,omitempty"`
	UseCertManager             bool                    `json:"useCertManager,omitempty"`
	Resources                  v1.ResourceRequirements `json:"resources,omitempty"`
	NodeSelector               map[string]string       `json:"nodeSelector,omitempty"`
	Tolerations                []v1.Toleration         `json:"tolerations,omitempty"`
	Affinity                   *v1.Affinity            `json:"affinity,omitempty"`
	AllowAdmissionSkip         bool                    `json:"allowAdmissionSkip,omitempty"`
	AllowedRepositories        []string                `json:"allowedRepositories,omitempty"`
	SecurityContextConstraints bool                    `json:"securityContextConstraints,omitempty"`
}

type Image struct {
	Host        string                    `json:"host,omitempty"`
	PullSecrets []v1.LocalObjectReference `json:"pullSecret,omitempty"`
	Image       string                    `json:"image,omitempty"`
	Tag         string                    `json:"tag,omitempty"`
	PullPolicy  v1.PullPolicy             `json:"pullPolicy,omitempty"`
}

type Service struct {
	Type              string             `json:"type,omitempty"`
	Port              int32              `json:"port,omitempty"`
	TargetPort        intstr.IntOrString `json:"targetPort,omitempty"`
	MetricsTargetPort intstr.IntOrString `json:"metricsTargetPort,omitempty"`
	MetricsPort       int32              `json:"metricsPort,omitempty"`
}

type WebHooks struct {
	FailurePolicy *admv1.FailurePolicyType `json:"failurePolicy,omitempty"`
}

type LabelSelector struct {
	MatchExpressions []MatchExpression `json:"matchExpressions,omitempty"`
}

type MatchExpression struct {
	Key      string   `json:"key,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Values   []string `json:"values,omitempty"`
}

// PortierisStatus defines the observed state of Portieris
type PortierisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Portieris is the Schema for the portieris API
type Portieris struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PortierisSpec   `json:"spec,omitempty"`
	Status PortierisStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PortierisList contains a list of Portieris
type PortierisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Portieris `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Portieris{}, &PortierisList{})
}

func (self *Portieris) GetServiceAccountName() string {
	return "portieris"
}

func (self *Portieris) GetClusterRoleName() string {
	return "portieris"
}

func (self *Portieris) GetClusterRoleBindingName() string {
	return "admission-portieris-webhook"
}

func (self *Portieris) GetWebhookServiceName() string {
	return "portieris"
}

func (self *Portieris) GetWebhookServerTlsSecretName() string {
	// "portieris-certs"
	return self.Spec.SecretName
}

func (self *Portieris) GetPodSecurityPolicyName() string {
	return "anyuid-portieris"
}

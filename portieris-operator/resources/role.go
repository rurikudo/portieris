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
	"strings"

	scc "github.com/openshift/api/security/v1"
	apiv1alpha1 "github.com/rurikudo/portieris/portieris-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//sa
func BuildServiceAccountForPortieris(cr *apiv1alpha1.Portieris) *corev1.ServiceAccount {
	labels := map[string]string{
		"app":                          cr.Name,
		"app.kubernetes.io/name":       cr.Name,
		"app.kubernetes.io/managed-by": "operator",
		"role":                         "security",
	}
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetServiceAccountName(),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
	}
	return sa
}

//cluster role
func BuildClusterRoleForPortieris(cr *apiv1alpha1.Portieris) *rbacv1.ClusterRole {
	labels := map[string]string{
		"app":                          cr.Name,
		"app.kubernetes.io/name":       cr.Name,
		"app.kubernetes.io/managed-by": "operator",
		"role":                         "security",
	}
	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetClusterRoleName(),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"portieris.cloud.ibm.com",
				},
				Resources: []string{
					"imagepolicies", "clusterimagepolicies",
				},
				Verbs: []string{
					"get", "list", "watch", "patch", "create",
				},
			},
			{
				APIGroups: []string{
					"apiextensions.k8s.io",
				},
				Resources: []string{
					"customresourcedefinitions",
				},
				Verbs: []string{
					"create", "delete", "get",
				},
			},
			{
				APIGroups: []string{
					"admissionregistration.k8s.io",
				},
				Resources: []string{
					"validatingwebhookconfigurations", "mutatingwebhookconfigurations",
				},
				Verbs: []string{
					"get", "create", "delete",
				},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"secrets", "serviceaccountss",
				},
				Verbs: []string{
					"get",
				},
			},
		},
	}
	return role
}

//cluster role-binding
func BuildClusterRoleBindingForPortieris(cr *apiv1alpha1.Portieris) *rbacv1.ClusterRoleBinding {
	labels := map[string]string{
		"app":                          cr.Name,
		"app.kubernetes.io/name":       cr.Name,
		"app.kubernetes.io/managed-by": "operator",
		"role":                         "security",
	}
	rolebinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetClusterRoleBindingName(),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      cr.GetServiceAccountName(),
				Namespace: cr.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     cr.GetClusterRoleName(),
		},
	}
	return rolebinding
}

//pod security policy
//scc
func BuildSecurityContextConstraints(cr *apiv1alpha1.Portieris) *scc.SecurityContextConstraints {
	user := strings.Join([]string{"system:serviceaccount", cr.Namespace, cr.GetServiceAccountName()}, ":")
	var priority int32 = 10
	metaLabels := map[string]string{
		"app":                          cr.Name,
		"app.kubernetes.io/name":       cr.Name,
		"app.kubernetes.io/managed-by": "operator",
		"role":                         "security",
	}

	return &scc.SecurityContextConstraints{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SecurityContextConstraints",
			APIVersion: scc.SchemeGroupVersion.Group + "/" + scc.SchemeGroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "anyuid-portieris",
			Namespace: cr.Namespace,
			Labels:    metaLabels,
		},
		AllowHostDirVolumePlugin: false,
		AllowHostIPC:             false,
		AllowHostNetwork:         false,
		AllowHostPID:             false,
		AllowHostPorts:           false,
		AllowPrivilegedContainer: false,
		AllowedCapabilities:      []corev1.Capability{},
		DefaultAddCapabilities:   []corev1.Capability{},
		FSGroup: scc.FSGroupStrategyOptions{
			Type: scc.FSGroupStrategyRunAsAny,
		},
		Groups: []string{
			"system:cluster-admins",
		},
		Priority:               &priority,
		ReadOnlyRootFilesystem: false,
		RequiredDropCapabilities: []corev1.Capability{
			"MKNOD",
		},
		RunAsUser: scc.RunAsUserStrategyOptions{
			Type: scc.RunAsUserStrategyMustRunAs,
			UID:  cr.Spec.SecurityContext.RunAsUser,
		},
		SELinuxContext:     scc.SELinuxContextStrategyOptions{Type: scc.SELinuxStrategyMustRunAs, SELinuxOptions: &corev1.SELinuxOptions{}},
		SupplementalGroups: scc.SupplementalGroupsStrategyOptions{Type: scc.SupplementalGroupsStrategyRunAsAny},
		Users:              []string{user},
		Volumes:            []scc.FSType{scc.FSTypeEmptyDir, scc.FSTypeSecret, scc.FSProjected, scc.FSTypeDownwardAPI, scc.FSTypePersistentVolumeClaim, scc.FSTypeConfigMap},
	}
}

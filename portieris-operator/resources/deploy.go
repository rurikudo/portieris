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
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	apiv1alpha1 "github.com/rurikudo/portieris/portieris-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

//deployment
func BuildDeploymentForPortieris(cr *apiv1alpha1.Portieris) *appsv1.Deployment {
	image := cr.Spec.Image.Host + "/" + cr.Spec.Image.Image + ":" + cr.Spec.Image.Tag
	labels := cr.Spec.MetaLabels
	servervolumemounts := []v1.VolumeMount{
		{
			MountPath: "/etc/certs",
			Name:      "portieris-certs",
			ReadOnly:  true,
		},
	}

	serverContainer := v1.Container{
		Name:            cr.Spec.Name,
		Image:           image,
		ImagePullPolicy: cr.Spec.Image.PullPolicy,
		ReadinessProbe: &v1.Probe{
			InitialDelaySeconds: 10,
			PeriodSeconds:       10,
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path:   "/health/readiness",
					Port:   intstr.IntOrString{IntVal: 8000},
					Scheme: v1.URISchemeHTTPS,
				},
			},
		},
		LivenessProbe: &v1.Probe{
			InitialDelaySeconds: 10,
			PeriodSeconds:       10,
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path:   "/health/liveness",
					Port:   intstr.IntOrString{IntVal: 8000},
					Scheme: v1.URISchemeHTTPS,
				},
			},
		},
		Ports: []v1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 80,
				Protocol:      v1.ProtocolTCP,
			},
			{
				Name:          "metrics-port",
				ContainerPort: 8080,
				Protocol:      v1.ProtocolTCP,
			},
		},
		VolumeMounts: servervolumemounts,
		Resources:    cr.Spec.Resources,
	}

	containers := []v1.Container{
		serverContainer,
	}

	volumes := []v1.Volume{
		{
			Name: "portieris-certs",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: "portieris-certs",
				},
			},
		},
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: cr.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: cr.Spec.SelectorLabels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: cr.Spec.SelectorLabels,
				},
				Spec: v1.PodSpec{
					ImagePullSecrets:   cr.Spec.Image.PullSecrets,
					ServiceAccountName: cr.GetServiceAccountName(),
					SecurityContext:    cr.Spec.SecurityContext,
					Containers:         containers,
					NodeSelector:       cr.Spec.NodeSelector,
					Affinity:           cr.Spec.Affinity,
					Tolerations:        cr.Spec.Tolerations,

					Volumes: volumes,
				},
			},
		},
	}
}

// EqualDeployments returns a Boolean
func EqualDeployments(expected *appsv1.Deployment, found *appsv1.Deployment) bool {
	if !EqualLabels(found.ObjectMeta.Labels, expected.ObjectMeta.Labels) {
		return false
	}
	if !EqualPods(expected.Spec.Template, found.Spec.Template) {
		return false
	}
	return true
}

// EqualPods returns a Boolean
func EqualPods(expected v1.PodTemplateSpec, found v1.PodTemplateSpec) bool {
	if !EqualLabels(found.ObjectMeta.Labels, expected.ObjectMeta.Labels) {
		return false
	}
	if !EqualAnnotations(found.ObjectMeta.Annotations, expected.ObjectMeta.Annotations) {
		return false
	}
	if !reflect.DeepEqual(found.Spec.ServiceAccountName, expected.Spec.ServiceAccountName) {
		return false
	}
	if len(found.Spec.Containers) != len(expected.Spec.Containers) {
		return false
	}
	if !EqualContainers(expected.Spec.Containers[0], found.Spec.Containers[0]) {
		return false
	}
	return true
}

// EqualContainers returns a Boolean
func EqualContainers(expected v1.Container, found v1.Container) bool {
	if !reflect.DeepEqual(found.Name, expected.Name) {
		return false
	}
	if !reflect.DeepEqual(found.Image, expected.Image) {
		return false
	}
	if !reflect.DeepEqual(found.ImagePullPolicy, expected.ImagePullPolicy) {
		return false
	}
	if !reflect.DeepEqual(found.VolumeMounts, expected.VolumeMounts) {
		return false
	}
	if !reflect.DeepEqual(found.SecurityContext, expected.SecurityContext) {
		return false
	}
	if !reflect.DeepEqual(found.Ports, expected.Ports) {
		return false
	}
	if !reflect.DeepEqual(found.Args, expected.Args) {
		return false
	}
	if !reflect.DeepEqual(found.Env, expected.Env) {
		return false
	}
	return true
}

func EqualLabels(found map[string]string, expected map[string]string) bool {
	return reflect.DeepEqual(found, expected)
}

func EqualAnnotations(found map[string]string, expected map[string]string) bool {
	return reflect.DeepEqual(found, expected)
}

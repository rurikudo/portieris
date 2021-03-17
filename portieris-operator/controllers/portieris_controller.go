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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	apisv1alpha1 "github.com/rurikudo/portieris/portieris-operator/api/v1alpha1"
	admv1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controller_portieris")

// PortierisReconciler reconciles a Portieris object
type PortierisReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apis.portieris.io,resources=portieris;portieris/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apis.portieris.io,resources=portieris/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods;services;serviceaccounts;services/finalizers;endpoints;persistentvolumeclaims;events;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,resourceNames=portieris-operator,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get
// +kubebuilder:rbac:groups=apps,resources=deployments;replicasets,verbs=get
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=*
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=*
// +kubebuilder:rbac:groups=policy,resources=podsecuritypolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=*
// +kubebuilder:rbac:groups=portieris.cloud.ibm.com,resources=clusterimagepolicies;imagepolicies,verbs=*
// +kubebuilder:rbac:groups=security.openshift.io,resources=securitycontextconstraints,verbs=*
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates;issuers;certificaterequests,verbs=*

func (r *PortierisReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("Request.Namespace", req.Namespace)

	// your logic here
	// Fetch the portieris instance
	instance := &apisv1alpha1.Portieris{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, " r.Get(ctx, req.NamespacedName, instance) ", instance.Name)
		return ctrl.Result{}, err
	}

	var recResult ctrl.Result
	var recErr error

	// Portieris is under deletion - finalizer step
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if containsString(instance.ObjectMeta.Finalizers, apisv1alpha1.CleanupFinalizerName) {
			if err := r.deleteClusterScopedChildrenResources(instance); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				reqLogger.Error(err, "Error occured during finalizer process. retrying soon.")
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, apisv1alpha1.CleanupFinalizerName)
			if err := r.Update(context.Background(), instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// otherwise, normal reconcile

	// Custom Resource Definition (CRD)
	recResult, recErr = r.createOrUpdateImagePolicyCRD(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}
	recResult, recErr = r.createOrUpdateClusterImagePolicyCRD(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	//Custom Resources (CR)
	recResult, recErr = r.createOrUpdateDefaultImagePolicyCR(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}
	recResult, recErr = r.createOrUpdateKubeSystemImagePolicyCR(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}
	if instance.Spec.IBMContainerService {
		recResult, recErr = r.createOrUpdateIBMSystemImagePolicyCR(instance)
		if recErr != nil || recResult.Requeue {
			return recResult, recErr
		}
		return ctrl.Result{}, nil
	}

	recResult, recErr = r.createOrUpdateClusterImagePolicyCR(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Secret
	if !instance.Spec.SkipSecretCreation {
		if instance.Spec.UseCertManager {
			// Issuer
			recResult, recErr = r.createOrUpdateIssuer(instance)
			if recErr != nil || recResult.Requeue {
				return recResult, recErr
			}
			// Certificate
			recResult, recErr = r.createOrUpdateCertificate(instance)
			if recErr != nil || recResult.Requeue {
				return recResult, recErr
			}
		} else {
			// Secret
			recResult, recErr = r.createOrUpdateTlsSecret(instance)
			if recErr != nil || recResult.Requeue {
				return recResult, recErr
			}
		}
	}

	//Service Account
	recResult, recErr = r.createOrUpdateServiceAccount(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	//Cluster Role
	recResult, recErr = r.createOrUpdateClusterRole(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	//Cluster Role Binding
	recResult, recErr = r.createOrUpdateClusterRoleBinding(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// SCC
	if instance.Spec.SecurityContextConstraints {
		recResult, recErr = r.createOrUpdateSecurityContextConstraints(instance)
		if recErr != nil || recResult.Requeue {
			return recResult, recErr
		}
	}

	//Deployment
	recResult, recErr = r.createOrUpdateWebhookDeployment(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	//Service
	recResult, recErr = r.createOrUpdateWebhookService(instance)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	//Webhook Configuration
	// wait until deployment is available
	if r.isDeploymentAvailable(instance) {
		recResult, recErr = r.createOrUpdateWebhook(instance)
		if recErr != nil || recResult.Requeue {
			return recResult, recErr
		}
	} else {
		recResult, recErr = r.deleteWebhook(instance)
		if recErr != nil || recResult.Requeue {
			return recResult, recErr
		}
	}

	reqLogger.Info("Reconciliation successful!", "Name", instance.Name)
	// since we updated the status in the CR, sleep 5 seconds to allow the CR to be refreshed.
	time.Sleep(5 * time.Second)

	return ctrl.Result{}, nil
}

func (r *PortierisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apisv1alpha1.Portieris{}).
		Owns(&apisv1alpha1.Portieris{}).
		Owns(&appsv1.Deployment{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&admv1.MutatingWebhookConfiguration{}).
		Complete(r)
}

// Owns(&imgpolicyv1.ClusterImagePolicy{}).
// Owns(&scc.SecurityContextConstraints{}).

func (r *PortierisReconciler) deleteClusterScopedChildrenResources(instance *apisv1alpha1.Portieris) error {
	// delete any cluster scope resources owned by the instance
	// (In Kubernetes 1.20 and later, a garbage collector ignore cluster scope children even if their owner is deleted)
	var err error
	_, err = r.deleteWebhook(instance)
	if err != nil {
		return err
	}

	// SCC
	if instance.Spec.SecurityContextConstraints {
		_, err = r.deleteSecurityContextConstraints(instance)
		if err != nil {
			return err
		}
	}

	_, err = r.deleteClusterRoleBinding(instance)
	if err != nil {
		return err
	}
	_, err = r.deleteClusterRole(instance)
	if err != nil {
		return err
	}

	_, err = r.deleteClusterImagePolicyCRD(instance)
	if err != nil {
		return err
	}

	_, err = r.deleteImagePolicyCRD(instance)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

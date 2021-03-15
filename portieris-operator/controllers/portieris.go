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
	"fmt"
	"reflect"
	"time"

	imgpolicyv1 "github.com/IBM/portieris/pkg/apis/portieris.cloud.ibm.com/v1"
	scc "github.com/openshift/api/security/v1"
	apiv1alpha1 "github.com/rurikudo/portieris/portieris-operator/api/v1alpha1"
	res "github.com/rurikudo/portieris/portieris-operator/resources"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	cert "github.com/rurikudo/portieris/portieris-operator/cert"
	admv1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/**********************************************

				CRD

***********************************************/

func (r *PortierisReconciler) createOrUpdateCRD(instance *apiv1alpha1.Portieris, expected *extv1.CustomResourceDefinition) (ctrl.Result, error) {
	ctx := context.Background()

	found := &extv1.CustomResourceDefinition{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"CRD.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If CRD does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: ""}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(expected.Spec, found.Spec) {
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

func (r *PortierisReconciler) deleteCRD(instance *apiv1alpha1.Portieris, expected *extv1.CustomResourceDefinition) (ctrl.Result, error) {
	ctx := context.Background()
	found := &extv1.CustomResourceDefinition{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"CustomResourceDefinition.Name", expected.Name)

	err := r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err == nil {
		reqLogger.Info(fmt.Sprintf("Deleting the IShield CustomResourceDefinition %s", expected.Name))
		err = r.Delete(ctx, found)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to delete the IShield CustomResourceDefinition %s", expected.Name))
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if errors.IsNotFound(err) {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else {
		return ctrl.Result{}, err
	}

}

func (r *PortierisReconciler) createOrUpdateImagePolicyCRD(
	instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildImagePolicyCRD(instance)
	return r.createOrUpdateCRD(instance, expected)
}

func (r *PortierisReconciler) createOrUpdateClusterImagePolicyCRD(
	instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildClusterImagePolicyCRD(instance)
	return r.createOrUpdateCRD(instance, expected)
}

func (r *PortierisReconciler) deleteImagePolicyCRD(
	instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildImagePolicyCRD(instance)
	return r.deleteCRD(instance, expected)
}

func (r *PortierisReconciler) deleteClusterImagePolicyCRD(
	instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildClusterImagePolicyCRD(instance)
	return r.deleteCRD(instance, expected)
}

/**********************************************

				CR

***********************************************/

func (r *PortierisReconciler) createOrUpdateDefaultImagePolicyCR(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()

	expected := res.BuildDefaultImagePolicyForPortieris(instance)
	found := &imgpolicyv1.ImagePolicy{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"DefaultImagePolicy.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: instance.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(expected.Spec, found.Spec) {
			// If spec is incorrect, update it and requeue
			reqLogger.Info("Try to update the resource...")
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil
}

func (r *PortierisReconciler) createOrUpdateKubeSystemImagePolicyCR(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()

	expected := res.BuildKubeSystemImagePolicyForPortieris(instance)
	found := &imgpolicyv1.ImagePolicy{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"Kube-SystemImagePolicy.Name", expected.Name)

	// If PodSecurityPolicy does not exist, create it and requeue
	err := r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(expected.Spec, found.Spec) {
			// If spec is incorrect, update it and requeue
			reqLogger.Info("Try to update the resource...")
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil
}

func (r *PortierisReconciler) createOrUpdateIBMSystemImagePolicyCR(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()

	expected := res.BuildIBMSystemImagePolicyForPortieris(instance)
	found := &imgpolicyv1.ImagePolicy{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"IBMSystemImagePolicy.Name", expected.Name)

	// If PodSecurityPolicy does not exist, create it and requeue
	err := r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: instance.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(expected.Spec, found.Spec) {
			// If spec is incorrect, update it and requeue
			reqLogger.Info("Try to update the resource...")
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil
}

func (r *PortierisReconciler) createOrUpdateClusterImagePolicyCR(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()

	expected := res.BuildClusterImagePolicyForPortieris(instance)
	found := &imgpolicyv1.ClusterImagePolicy{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"ClusterImagePolicy.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(expected.Spec, found.Spec) {
			// If spec is incorrect, update it and requeue
			reqLogger.Info("Try to update the resource...")
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil
}

/**********************************************

				Role

***********************************************/

func (r *PortierisReconciler) createOrUpdateServiceAccount(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	expected := res.BuildServiceAccountForPortieris(instance)
	found := &corev1.ServiceAccount{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"ServiceAccount.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: instance.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

func (r *PortierisReconciler) createOrUpdateClusterRole(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	expected := res.BuildClusterRoleForPortieris(instance)
	found := &rbacv1.ClusterRole{}

	reqLogger := r.Log.WithValues(
		"ClusterRole.Namespace", instance.Namespace,
		"Instance.Name", instance.Name,
		"ClusterRole.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !res.EqualRules(expected.Rules, found.Rules) {
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

func (r *PortierisReconciler) deleteClusterRole(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	found := &rbacv1.ClusterRole{}
	expected := res.BuildClusterRoleForPortieris(instance)
	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"ClusterRole.Name", expected.Name)

	err := r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err == nil {
		reqLogger.Info(fmt.Sprintf("Deleting the ClusterRole %s", expected.Name))
		err = r.Delete(ctx, found)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to delete the ClusterRole %s", expected.Name))
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if errors.IsNotFound(err) {

		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else {
		return ctrl.Result{}, err
	}

}

func (r *PortierisReconciler) createOrUpdateClusterRoleBinding(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	expected := res.BuildClusterRoleBindingForPortieris(instance)
	found := &rbacv1.ClusterRoleBinding{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"ClusterRoleBinding.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !res.EqualClusterRoleBindings(expected, found) {
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

func (r *PortierisReconciler) deleteClusterRoleBinding(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	found := &rbacv1.ClusterRoleBinding{}
	expected := res.BuildClusterRoleBindingForPortieris(instance)

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"ClusterRoleBinding.Name", expected.Name)

	err := r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err == nil {
		reqLogger.Info(fmt.Sprintf("Deleting the ClusterRoleBinding %s", expected.Name))
		err = r.Delete(ctx, found)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to delete the ClusterRoleBinding %s", expected.Name))
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if errors.IsNotFound(err) {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else {
		return ctrl.Result{}, err
	}

}

func (r *PortierisReconciler) createOrUpdateSecurityContextConstraints(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	expected := res.BuildSecurityContextConstraints(instance)
	found := &scc.SecurityContextConstraints{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"SecurityContextConstraints.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(expected, found) {
			expected.ObjectMeta = found.ObjectMeta
			err = r.Update(ctx, expected)
			if err != nil {
				reqLogger.Error(err, "Failed to update the resource")
				return ctrl.Result{}, err
			}
		}
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

func (r *PortierisReconciler) deleteSecurityContextConstraints(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	found := &scc.SecurityContextConstraints{}
	expected := res.BuildSecurityContextConstraints(instance)

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"SecurityContextConstraints.Name", expected.Name)

	err := r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err == nil {
		reqLogger.Info(fmt.Sprintf("Deleting the SCC %s", expected.Name))
		err = r.Delete(ctx, found)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to delete the SCC %s", expected.Name))
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if errors.IsNotFound(err) {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else {
		return ctrl.Result{}, err
	}

}

/**********************************************

				Secret

***********************************************/

func (r *PortierisReconciler) createOrUpdateCertSecret(instance *apiv1alpha1.Portieris, expected *corev1.Secret) (ctrl.Result, error) {
	ctx := context.Background()
	found := &corev1.Secret{}

	reqLogger := r.Log.WithValues(
		"Secret.Namespace", instance.Namespace,
		"Instance.Name", instance.Name,
		"Secret.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If CRD does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: instance.Namespace}, found)

	expected = addCertValues(instance, expected)

	if err != nil && errors.IsNotFound(err) {

		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

// func (r *PortierisReconciler) createOrUpdateCertificate(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
// 	ctx := context.Background()
// 	expected := res.BuildCertificateForPortieris(instance)
// 	found := &certmanager.Certificate{}

// 	reqLogger := r.Log.WithValues(
// 		"Instance.Name", instance.Name,
// 		"Certificate.Name", expected.Name)

// 	// Set CR instance as the owner and controller
// 	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
// 	if err != nil {
// 		reqLogger.Error(err, "Failed to define expected resource")
// 		return ctrl.Result{}, err
// 	}

// 	// If PodSecurityPolicy does not exist, create it and requeue
// 	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

// 	if err != nil && errors.IsNotFound(err) {
// 		reqLogger.Info("Creating a new resource")
// 		err = r.Create(ctx, expected)
// 		if err != nil && errors.IsAlreadyExists(err) {
// 			// Already exists from previous reconcile, requeue.
// 			reqLogger.Info("Skip reconcile: resource already exists")
// 			return ctrl.Result{Requeue: true}, nil
// 		} else if err != nil {
// 			reqLogger.Error(err, "Failed to create new resource")
// 			return ctrl.Result{}, err
// 		}
// 		// Created successfully - return and requeue
// 		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
// 	} else if err != nil {
// 		return ctrl.Result{}, err
// 	} else {
// 		if !reflect.DeepEqual(expected.Spec, found.Spec) {
// 			expected.ObjectMeta = found.ObjectMeta
// 			err = r.Update(ctx, expected)
// 			if err != nil {
// 				reqLogger.Error(err, "Failed to update the resource")
// 				return ctrl.Result{}, err
// 			}
// 		}
// 	}
// 	// No extra validation

// 	// No reconcile was necessary
// 	return ctrl.Result{}, nil

// }

// func (r *PortierisReconciler) createOrUpdateIssuer(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
// 	ctx := context.Background()
// 	expected := res.BuildIssuerForPortieris(instance)
// 	found := &certmanager.Issuer{}

// 	reqLogger := r.Log.WithValues(
// 		"Instance.Name", instance.Name,
// 		"Issuer.Name", expected.Name)

// 	// Set CR instance as the owner and controller
// 	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
// 	if err != nil {
// 		reqLogger.Error(err, "Failed to define expected resource")
// 		return ctrl.Result{}, err
// 	}

// 	// If PodSecurityPolicy does not exist, create it and requeue
// 	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

// 	if err != nil && errors.IsNotFound(err) {
// 		reqLogger.Info("Creating a new resource")
// 		err = r.Create(ctx, expected)
// 		if err != nil && errors.IsAlreadyExists(err) {
// 			// Already exists from previous reconcile, requeue.
// 			reqLogger.Info("Skip reconcile: resource already exists")
// 			return ctrl.Result{Requeue: true}, nil
// 		} else if err != nil {
// 			reqLogger.Error(err, "Failed to create new resource")
// 			return ctrl.Result{}, err
// 		}
// 		// Created successfully - return and requeue
// 		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
// 	} else if err != nil {
// 		return ctrl.Result{}, err
// 	} else {
// 		if !reflect.DeepEqual(expected.Spec, found.Spec) {
// 			expected.ObjectMeta = found.ObjectMeta
// 			err = r.Update(ctx, expected)
// 			if err != nil {
// 				reqLogger.Error(err, "Failed to update the resource")
// 				return ctrl.Result{}, err
// 			}
// 		}
// 	}
// 	// No extra validation

// 	// No reconcile was necessary
// 	return ctrl.Result{}, nil

// }

func addCertValues(instance *apiv1alpha1.Portieris, expected *corev1.Secret) *corev1.Secret {
	reqLogger := log.WithValues(
		"Secret.Namespace", instance.Namespace,
		"Instance.Name", instance.Name,
		"Secret.Name", expected.Name)

	// generate and put certsß
	ca, tlsKey, tlsCert, err := cert.GenerateCert(instance.GetWebhookServiceName(), instance.Namespace)
	if err != nil {
		reqLogger.Error(err, "Failed to generate certs")
	}

	_, ok_tc := expected.Data["tls.crt"]
	_, ok_tk := expected.Data["tls.key"]
	_, ok_ca := expected.Data["ca.crt"]
	if ok_ca && ok_tc && ok_tk {
		expected.Data["tls.crt"] = tlsCert
		expected.Data["tls.key"] = tlsKey
		expected.Data["ca.crt"] = ca
	}
	return expected
}

func (r *PortierisReconciler) createOrUpdateTlsSecret(
	instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildTlsSecretForPortieris(instance)
	return r.createOrUpdateCertSecret(instance, expected)
}

/**********************************************

				Deployment

***********************************************/

func (r *PortierisReconciler) createOrUpdateDeployment(instance *apiv1alpha1.Portieris, expected *appsv1.Deployment) (ctrl.Result, error) {
	ctx := context.Background()
	found := &appsv1.Deployment{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"Deployment.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: instance.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	} else if !res.EqualDeployments(expected, found) {
		// If spec is incorrect, update it and requeue
		reqLogger.Info("Try to update the resource...")
		found.ObjectMeta.Labels = expected.ObjectMeta.Labels
		found.Spec = expected.Spec
		err = r.Update(ctx, found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Namespace", instance.Namespace, "Name", found.Name)
			return ctrl.Result{}, err
		}
		reqLogger.Info("Updating Portieris Controller Deployment", "Deployment.Name", found.Name)
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

func (r *PortierisReconciler) createOrUpdateWebhookDeployment(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildDeploymentForPortieris(instance)
	return r.createOrUpdateDeployment(instance, expected)
}

/**********************************************

				Service

***********************************************/

func (r *PortierisReconciler) createOrUpdateService(instance *apiv1alpha1.Portieris, expected *corev1.Service) (ctrl.Result, error) {
	ctx := context.Background()
	found := &corev1.Service{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"Instance.Spec.ServiceName", instance.GetWebhookServiceName(),
		"Service.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: instance.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil
}

func (r *PortierisReconciler) createOrUpdateWebhookService(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	expected := res.BuildServiceForPortieris(instance)
	return r.createOrUpdateService(instance, expected)
}

/**********************************************

				Webhook

***********************************************/

func (r *PortierisReconciler) createOrUpdateWebhook(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	expected := res.BuildMutatingWebhookConfigurationForPortieris(instance)
	found := &admv1.MutatingWebhookConfiguration{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"MutatingWebhookConfiguration.Name", expected.Name)

	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(instance, expected, r.Scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return ctrl.Result{}, err
	}

	// If PodSecurityPolicy does not exist, create it and requeue
	err = r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new resource")
		// locad cabundle
		secret := &corev1.Secret{}
		err = r.Get(ctx, types.NamespacedName{Name: instance.GetWebhookServerTlsSecretName(), Namespace: instance.Namespace}, secret)
		if err != nil {
			reqLogger.Error(err, "Fail to load CABundle from Secret")
		}
		cabundle, ok := secret.Data["ca.crt"]
		if ok {
			expected.Webhooks[0].ClientConfig.CABundle = cabundle
		}

		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip reconcile: resource already exists")
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new resource")
			return ctrl.Result{}, err
		}
		// Created successfully - return and requeue

		reqLogger.Info("Webhook has been created.", "Name", instance.Name)
		evtName := fmt.Sprintf("portieris-webhook-reconciled")
		_ = r.createOrUpdateWebhookEvent(instance, evtName, expected.Name)

		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// No extra validation

	// No reconcile was necessary
	return ctrl.Result{}, nil

}

// delete webhookconfiguration
func (r *PortierisReconciler) deleteWebhook(instance *apiv1alpha1.Portieris) (ctrl.Result, error) {
	ctx := context.Background()
	expected := res.BuildMutatingWebhookConfigurationForPortieris(instance)
	found := &admv1.MutatingWebhookConfiguration{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"MutatingWebhookConfiguration.Name", expected.Name)

	err := r.Get(ctx, types.NamespacedName{Name: expected.Name}, found)

	if err == nil {
		reqLogger.Info("Deleting the Portieris webhook")
		err = r.Delete(ctx, found)
		if err != nil {
			reqLogger.Error(err, "Failed to delete the Portieris Webhook")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else if errors.IsNotFound(err) {
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
	} else {
		return ctrl.Result{}, err
	}
}

// wait function
func (r *PortierisReconciler) isDeploymentAvailable(instance *apiv1alpha1.Portieris) bool {
	ctx := context.Background()
	found := &appsv1.Deployment{}

	// If Deployment does not exist, return false
	err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return false
	} else if err != nil {
		return false
	}

	// return true only if deployment is available
	if found.Status.AvailableReplicas > 0 {
		return true
	}

	return false
}

func (r *PortierisReconciler) createOrUpdateWebhookEvent(instance *apiv1alpha1.Portieris, evtName, webhookName string) error {
	ctx := context.Background()
	evtNamespace := instance.Namespace
	involvedObject := v1.ObjectReference{
		Namespace:  evtNamespace,
		APIVersion: instance.APIVersion,
		Kind:       instance.Kind,
		Name:       instance.Name,
	}
	now := time.Now()
	evtSourceName := "Portieris"
	reason := "webhook-reconciled"
	msg := fmt.Sprintf("[PortierisEvent] Portieris reconciled MutatingWebhookConfiguration \"%s\"", webhookName)
	expected := &v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      evtName,
			Namespace: evtNamespace,
		},
		InvolvedObject:      involvedObject,
		Type:                evtSourceName,
		Source:              v1.EventSource{Component: evtSourceName},
		ReportingController: evtSourceName,
		ReportingInstance:   evtName,
		Action:              evtName,
		FirstTimestamp:      metav1.NewTime(now),
		LastTimestamp:       metav1.NewTime(now),
		EventTime:           metav1.NewMicroTime(now),
		Message:             msg,
		Reason:              reason,
		Count:               1,
	}
	found := &v1.Event{}

	reqLogger := r.Log.WithValues(
		"Instance.Name", instance.Name,
		"Event.Name", expected.Name)

	// If Event does not exist, create it and requeue
	err := r.Get(ctx, types.NamespacedName{Name: expected.Name, Namespace: expected.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new event")
		err = r.Create(ctx, expected)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			reqLogger.Info("Skip creating event: resource already exists")
			return nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new event")
			return err
		}
		// Created successfully - return and requeue
		return nil
	} else if err != nil {
		return err
	} else {
		// Update Event
		found.Count = found.Count + 1
		found.EventTime = metav1.NewMicroTime(now)
		found.LastTimestamp = metav1.NewTime(now)
		found.Message = msg
		found.Reason = msg
		found.ReportingController = evtSourceName
		found.ReportingInstance = evtName

		err = r.Update(ctx, found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Event", "Namespace", instance.Namespace, "Name", found.Name)
			return err
		}
		reqLogger.Info("Updated Event", "Deployment.Name", found.Name)
		// Spec updated - return and requeue
		return nil
	}
}

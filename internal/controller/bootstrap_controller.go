/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"

	robotopsv1 "cloudhub.cz/robotops/api/v1"

	corev1 "k8s.io/api/core/v1"                   // For core Kubernetes resources (e.g., Pods, Services)
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // For object metadata
	"k8s.io/apimachinery/pkg/runtime"             // For working with Kubernetes runtime objects
	"sigs.k8s.io/controller-runtime/pkg/client"   // For interacting with Kubernetes objects
	"sigs.k8s.io/controller-runtime/pkg/log"      // For logging in controllers
)

// BootstrapReconciler reconciles a Bootstrap object
type BootstrapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=robotops.cloudhub.cz,resources=bootstraps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=robotops.cloudhub.cz,resources=bootstraps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=robotops.cloudhub.cz,resources=bootstraps/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=create;get;list;watch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bootstrap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *BootstrapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx).WithName("MyReconciler").WithValues("resource", req.NamespacedName)

	log.Info("Starting reconciliation loop")

	// TODO(user): your logic here
	// Fetch the Bootstrap instance
	var bootstrap robotopsv1.Bootstrap
	if err := r.Get(ctx, req.NamespacedName, &bootstrap); err != nil {
		log.Error(err, "unable to fetch Bootstrap")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Bootstrap object found", "name", bootstrap.Name)

	// Define a new Pod object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bootstrap.Spec.Pod.Name,
			Namespace: req.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  bootstrap.Spec.Pod.Name,
				Image: bootstrap.Spec.Pod.Image,
			}},
		},
	}
	if err := r.Create(ctx, pod); err != nil {
		log.Error(err, "failed to create new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BootstrapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&robotopsv1.Bootstrap{}).
		Complete(r)
}

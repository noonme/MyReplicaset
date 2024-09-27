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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1beta1 "devops-test/api/v1beta1"
)

// Add creates a new MyReplicaSet Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

type ReconcileMyReplicaSet struct {
	client client.Client
	scheme *runtime.Scheme
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileMyReplicaSet {
	return &ReconcileMyReplicaSet{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	controller, err := controller.New("myreplicaset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = controller.Watch(&source.Kind{&handler.EnqueueRequestForObject{}})

	if err != nil {
		return err
	}

	return nil
}

// MyReplicaSetReconciler reconciles a MyReplicaSet object
type MyReplicaSetReconciler struct {
	clihtop
	ient.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.opsblogs.cn,resources=myreplicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.opsblogs.cn,resources=myreplicasets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.opsblogs.cn,resources=myreplicasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyReplicaSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *MyReplicaSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the MyReplicaSet instance
	var myReplicaSet appsv1beta1.MyReplicaSet
	if err := r.client.Get(ctx, req.NamespacedName, &myReplicaSet); err != nil {
		// Error reading the object - requeue the request.
		if errors.IsNotFound(err) {
			// The object no longer exists - remove any finalizers if they exist.
			return reconcile.Result{}, nil
		}
		// Something else went wrong - requeue and report the error.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pods := &v1.PodList{}
	listOpts := []client.ListOption{
		{LabelSelector: labels.SelectorFromSet(myReplicaSet.Spec.Template.Labels).String()},
	}
	if err := r.client.List(ctx, pods, &client.ListOptions{Namespace: myReplicaSet.Namespace, LabelSelector: labels.SelectorFromSet(myReplicaSet.Spec.Template.Labels).String()}); err != nil {
		return reconcile.Result{}, err
	}

	// Set MyReplicaSet instance as the owner and controller
	if err := controllerutil.SetControllerReference(&myReplicaSet, pods, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO: your custom logic to sync the desired state

	// Update the found object and write the result back if there are any changes
	if err := r.client.Update(ctx, &myReplicaSet); err != nil {
		return reconcile.Result{}, err
	}

	// Periodically requeue reconciliation requests for MyReplicaSet
	return reconcile.Result{RequeueAfter: time.Second * 10}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyReplicaSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1beta1.MyReplicaSet{}).
		Complete(r)
}

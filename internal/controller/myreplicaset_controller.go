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
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	"time"

	appsv1beta1 "devops-test/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//// Add creates a new MyReplicaSet Controller and adds it to the Manager. The Manager will set fields on the Controller
//// and Start it when the Manager is Started.
//func Add(mgr manager.Manager) error {
//	return add(mgr, newReconciler(mgr))
//}
//
//type ReconcileMyReplicaSet struct {
//	client client.Client
//	scheme *runtime.Scheme
//}
//
//// newReconciler returns a new reconcile.Reconciler
//func newReconciler(mgr manager.Manager) *ReconcileMyReplicaSet {
//	return &ReconcileMyReplicaSet{client: mgr.GetClient(), scheme: mgr.GetScheme()}
//}
//
//// add 函数将一个新的控制器添加到管理器中，并使用给定的协调器进行协调
//func add(mgr manager.Manager, r reconcile.Reconciler) error {
//	// 创建一个新的控制器
//	controller, err := controller.New("myreplicaset-controller", mgr, controller.Options{Reconciler: r})
//	if err != nil {
//		// 如果创建控制器失败，返回错误
//		return fmt.Errorf("failed to create controller: %v", err)
//	}
//
//	// 监听 MyReplicaSet 资源的创建、更新和删除事件
//	err = controller.Watch(source.Kind(mgr.GetCache(), &appsv1beta1.MyReplicaSet{}), &handler.EnqueueRequestForObject{Client: mgr.GetClient()})
//	if err != nil {
//		// 如果监听失败，返回错误
//		return fmt.Errorf("failed to watch MyReplicaSet resource: %v", err)
//	}
//
//	// 如果一切正常，返回 nil
//	return nil
//}

// MyReplicaSetReconciler reconciles a MyReplicaSet object
type MyReplicaSetReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
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
	if err := r.Client.Get(ctx, req.NamespacedName, &myReplicaSet); err != nil {
		// Error reading the object - requeue the request.
		if errors.IsNotFound(err) {
			// The object no longer exists - remove any finalizers if they exist.
			return reconcile.Result{}, nil
		}
		// Something else went wrong - requeue and report the error.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	podList := &v1.PodList{}
	labelSelector := myReplicaSet.Spec.Template.ObjectMeta.Labels
	if err := r.Client.List(ctx, podList, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector),
		Namespace:     myReplicaSet.Namespace,
	}); err != nil {
		return reconcile.Result{}, err
	} else {
		fmt.Printf("Found %d Pods in namespace %s\n", len(podList.Items), myReplicaSet.Namespace)

		dep, err := r.podForMyReplicaset(&myReplicaSet, myReplicaSet.Spec)
		for _, y := range dep {
			if err = r.Create(ctx, y); err != nil {
				klog.Error(err, "Failed to create new Deployment",
					"Deployment.Namespace", y.Namespace, "Deployment.Name", y.Name)
				return ctrl.Result{}, err
			}
		}

		// Deployment created successfully
		// We will requeue the reconciliation so that we can ensure the state
		// and move forward for the next operations
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Set MyReplicaSet instance as the owner and controller
	if err := controllerutil.SetControllerReference(&myReplicaSet, &v1.Pod{}, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// TODO: your custom logic to sync the desired state

	// Update the found object and write the result back if there are any changes
	if err := r.Client.Update(ctx, &myReplicaSet); err != nil {
		return reconcile.Result{}, err
	}

	// Periodically requeue reconciliation requests for MyReplicaSet
	return reconcile.Result{RequeueAfter: time.Second * 10}, nil
}

func labelsForMyReplicaset(myreplicaset appsv1beta1.MyReplicaSetSpec) map[string]string {
	//var imageTag string
	//os.Setenv("MEMCACHED_IMAGE", "dockerproxy.cn/memcached:1.4.36-alpine")
	//image, err := imageForMemcached()
	//if err == nil {
	//	imageTag = strings.Split(image, ":")[1]
	//}
	//return map[string]string{"app.kubernetes.io/name": "project",
	//	"app.kubernetes.io/version":    imageTag,
	//	"app.kubernetes.io/managed-by": "MemcachedController",
	//}

	label := myreplicaset.Selector.MatchLabels

	return label
}

func (r *MyReplicaSetReconciler) podForMyReplicaset(
	myreplicaset *appsv1beta1.MyReplicaSet, label appsv1beta1.MyReplicaSetSpec) ([]*v1.Pod, error) {
	ls := labelsForMyReplicaset(label)
	replicas := myreplicaset.Spec.Replicas

	dep := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myreplicaset.Name,
			Namespace: myreplicaset.Namespace,
			Labels:    ls,
		},
		Spec: v1.PodSpec{
			// TODO(user): Uncomment the following code to configure the nodeAffinity expression
			// according to the platforms which are supported by your solution. It is considered
			// best practice to support multiple architectures. build your manager image using the
			// makefile target docker-buildx. Also, you can use docker manifest inspect <image>
			// to check what are the platforms supported.
			// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity
			SecurityContext: &v1.PodSecurityContext{
				RunAsNonRoot: &[]bool{true}[0],
				// IMPORTANT: seccomProfile was introduced with Kubernetes 1.19
				// If you are looking for to produce solutions to be supported
				// on lower versions you must remove this option.
				SeccompProfile: &v1.SeccompProfile{
					Type: v1.SeccompProfileTypeRuntimeDefault,
				},
			},
			Containers: []v1.Container{{
				Image:           label.Template.Spec.Containers[0].Image,
				Name:            label.Template.Spec.Containers[0].Name,
				ImagePullPolicy: v1.PullIfNotPresent,
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						"cpu": resource.MustParse("1"),

						"memory": resource.MustParse("512Mi"),
					},
					Requests: v1.ResourceList{

						"cpu":    resource.MustParse("1"),
						"memory": resource.MustParse("512Mi"),
					},
				},
				// Ensure restrictive context for the container
				// More info: https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
				SecurityContext: &v1.SecurityContext{
					RunAsNonRoot:             &[]bool{true}[0],
					RunAsUser:                &[]int64{1001}[0],
					AllowPrivilegeEscalation: &[]bool{false}[0],
					Capabilities: &v1.Capabilities{
						Drop: []v1.Capability{
							"ALL",
						},
					},
				},
			}},
		},
	}

	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/

	var deps []*v1.Pod
	for i := 0; i < int(replicas); i++ {
		if err := ctrl.SetControllerReference(myreplicaset, dep, r.Scheme); err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyReplicaSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1beta1.MyReplicaSet{}).
		Complete(r)
}

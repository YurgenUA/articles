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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demov1 "operator-demo/api/v1"
)

// MemoryAdjusterReconciler reconciles a MemoryAdjuster object
type MemoryAdjusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=demo.operator.k8s.yfenyuk.io,resources=memoryadjusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=demo.operator.k8s.yfenyuk.io,resources=memoryadjusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=demo.operator.k8s.yfenyuk.io,resources=memoryadjusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MemoryAdjuster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *MemoryAdjusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling MemoryAdjuster...", "name", req.NamespacedName)

	// Fetch the OOMAdjuster instance
	var adjuster demov1.MemoryAdjuster
	if err := r.Get(ctx, req.NamespacedName, &adjuster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Current CR has values", "targetPodLabel", adjuster.Spec.TargetPodLabel, "memoryIncrement", adjuster.Spec.MemoryIncrement)

	if err := r.handleOOMPods(ctx, adjuster); err != nil {
		log.Error(err, "Failed to handle OOM pods")
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *MemoryAdjusterReconciler) handleOOMPods(ctx context.Context, adjuster demov1.MemoryAdjuster) error {
	log := log.FromContext(ctx)
	log.Info("Handling OOM pods")

	// List all pods with the given label in the namespace
	podList := &corev1.PodList{}
	labelSelector := client.MatchingLabels{"app": adjuster.Spec.TargetPodLabel}
	if err := r.List(ctx, podList, client.InNamespace(adjuster.Namespace), labelSelector); err != nil {
		return err
	}

	for _, pod := range podList.Items {
		// Check if any container in the pod has been restarted due to OOM (Exit Code 137)
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.LastTerminationState.Terminated != nil &&
				containerStatus.LastTerminationState.Terminated.ExitCode == 137 &&
				time.Since(containerStatus.LastTerminationState.Terminated.FinishedAt.Time) < time.Minute {

				log.Info("Detected OOMKilled Pod", "Pod", pod.Name)
				err := r.adjustPodResources(ctx, &pod, adjuster)
				if err != nil {
					log.Error(err, "Failed to adjust resources for Pod", "Pod", pod.Name)
				}
				break
			}
		}
	}
	return nil
}

func (r *MemoryAdjusterReconciler) adjustPodResources(ctx context.Context, pod *corev1.Pod, adjuster demov1.MemoryAdjuster) error {
	log := log.FromContext(ctx)
	log.Info("Adjusting resources for Pod", "Pod", pod.Name)

	deployment := &appsv1.Deployment{}
	if err := r.findDeployment(ctx, pod, deployment); err != nil {
		return err
	}

	for i, container := range deployment.Spec.Template.Spec.Containers {
		if err := r.adjustContainerMemory(ctx, deployment, i, &container, adjuster); err != nil {
			return err
		}
	}
	// Update the Deployment with new resource limits
	if err := r.Update(ctx, deployment); err != nil {
		return err
	}

	log.Info("Updated resource limits for Deployment", "Deployment", deployment.Name)
	return nil
}

func (r *MemoryAdjusterReconciler) findDeployment(ctx context.Context, pod *corev1.Pod, deployment *appsv1.Deployment) error {
	// Fetch the parent ReplicaSet
	replicaSet := &appsv1.ReplicaSet{}
	for _, ownerRef := range pod.OwnerReferences {
		if ownerRef.Kind == "ReplicaSet" {
			if err := r.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: ownerRef.Name}, replicaSet); err != nil {
				return err
			}
			break
		}
	}

	if replicaSet == nil {
		return fmt.Errorf("no parent ReplicaSet found for Pod %s", pod.Name)
	}

	// Fetch the parent Deployment
	for _, ownerRef := range replicaSet.OwnerReferences {
		if ownerRef.Kind == "Deployment" {
			if err := r.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: ownerRef.Name}, deployment); err != nil {
				return err
			}
			break
		}
	}

	if deployment == nil {
		return fmt.Errorf("no parent Deployment found for ReplicaSet %s", replicaSet.Name)
	}
	return nil
}

func (r *MemoryAdjusterReconciler) adjustContainerMemory(ctx context.Context, deployment *appsv1.Deployment, containerIndex int, container *corev1.Container, adjuster demov1.MemoryAdjuster) error {
	log := log.FromContext(ctx)
	memoryIncrement, err := resource.ParseQuantity(adjuster.Spec.MemoryIncrement)
	if err != nil {
		return err
	}

	memoryLimit := container.Resources.Limits[corev1.ResourceMemory]
	newMemoryLimit := memoryLimit.DeepCopy()
	newMemoryLimit.Add(memoryIncrement)
	deployment.Spec.Template.Spec.Containers[containerIndex].Resources.Limits[corev1.ResourceMemory] = newMemoryLimit
	log.Info("Container memory limits", "old", memoryLimit, "new", newMemoryLimit)

	// Calculate 10% lower value for GOMEMLIMIT
	gomemlimit := newMemoryLimit.DeepCopy()
	gomemlimit.Sub(resource.MustParse(fmt.Sprintf("%d", gomemlimit.Value()/10)))
	log.Info("GOMEMLIMIT", "new", gomemlimit)

	// Update or add GOMEMLIMIT environment variable
	envUpdated := false
	for j, env := range container.Env {
		if env.Name == "GOMEMLIMIT" {
			deployment.Spec.Template.Spec.Containers[containerIndex].Env[j].Value = gomemlimit.String()
			envUpdated = true
			break
		}
	}
	if !envUpdated {
		deployment.Spec.Template.Spec.Containers[containerIndex].Env = append(deployment.Spec.Template.Spec.Containers[containerIndex].Env, corev1.EnvVar{
			Name:  "GOMEMLIMIT",
			Value: gomemlimit.String(),
		})
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemoryAdjusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.MemoryAdjuster{}).
		Complete(r)
}

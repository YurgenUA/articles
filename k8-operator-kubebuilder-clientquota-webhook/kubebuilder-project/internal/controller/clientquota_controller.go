/*
Copyright 2025.

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
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	quotav1alpha1 "github.com/YurgenUA/k8s-client-quota/api/v1alpha1"
)

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;delete


// ClientQuotaReconciler reconciles a ClientQuota object
type ClientQuotaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ClientQuotaReconciler) readOrInitQuotaMap(ctx context.Context, quota *quotav1alpha1.ClientQuota) (map[string]int, error) {
	cmName := "clientquota-" + quota.Name
	cm := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKey{Name: cmName, Namespace: quota.Namespace}, cm)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create a new ConfigMap with initial quota from spec
			data := make(map[string]string)
			for _, client := range quota.Spec.Clients {
				data[client.Name] = strconv.Itoa(client.QuotaMinutes)
			}

			newCM := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cmName,
					Namespace: quota.Namespace,
					Labels: map[string]string{
						"app": "clientquota",
					},
				},
				Data: data,
			}

			if err := r.Create(ctx, newCM); err != nil {
				return nil, err
			}
			quotaMap := make(map[string]int)
			for k, v := range data {
				quotaMap[k], _ = strconv.Atoi(v)
			}
			return quotaMap, nil
		}
		return nil, err // real error
	}

	// Existing CM found — convert Data to map[string]int
	quotaMap := make(map[string]int)
	for k, v := range cm.Data {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid quota value for client %s: %v", k, err)
		}
		quotaMap[k] = i
	}

	return quotaMap, nil
}

func (r *ClientQuotaReconciler) updateQuotaMap(ctx context.Context, quota *quotav1alpha1.ClientQuota, updated map[string]int) error {
	cmName := "clientquota-" + quota.Name
	cm := &corev1.ConfigMap{}
	if err := r.Get(ctx, client.ObjectKey{Name: cmName, Namespace: quota.Namespace}, cm); err != nil {
		return err
	}

	newData := make(map[string]string)
	for k, v := range updated {
		newData[k] = strconv.Itoa(v)
	}
	cm.Data = newData
	return r.Update(ctx, cm)
}

// +kubebuilder:rbac:groups=quota.operator.k8s.yfenyuk.io,resources=clientquotas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=quota.operator.k8s.yfenyuk.io,resources=clientquotas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=quota.operator.k8s.yfenyuk.io,resources=clientquotas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClientQuota object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *ClientQuotaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling ClientQuota...", "name", req.NamespacedName)

	var quota quotav1alpha1.ClientQuota
	if err := r.Get(ctx, req.NamespacedName, &quota); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	quotaMap, err := r.readOrInitQuotaMap(ctx, &quota)
	if err != nil {
		log.Error(err, "failed to read/create quota map")
		return ctrl.Result{}, err
	}
	log.Info("QuotaMap", "quotaMap", quotaMap)

	var pods corev1.PodList
	if err := r.List(ctx, &pods, client.InNamespace(quota.Spec.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	quota_annotation := "quota.operator.k8s.yfenyuk.io/api-key"
	for _, pod := range pods.Items {
		// If Pod has no quota_annotation, kill it
		if _, ok := pod.Annotations[quota_annotation]; !ok {
			log.Info("Pod has no quota annotation, deleting...", "pod", pod.Name)
			if err := r.Delete(ctx, &pod); err != nil {
				log.Error(err, "Failed to delete Pod", "pod", pod.Name)
				return ctrl.Result{}, err
			}
			continue
		}

		podApiKey := pod.Annotations[quota_annotation]

		// find pod_api_key in quota.Spec.Clients
		var clientSpec quotav1alpha1.ClientSpec
		for _, client := range quota.Spec.Clients {
			if client.APIKey == podApiKey {
				clientSpec = client
				break
			}
		}
		log.Info("ClientSpec", "clientSpec", clientSpec)
		// If pod_api_key not found in quota.Spec.Clients, kill it
		if clientSpec.Name == "" {
			log.Info("Pod API Key not found in quota, deleting...", "pod", pod.Name)
			if err := r.Delete(ctx, &pod); err != nil {
				log.Error(err, "Failed to delete Pod", "pod", pod.Name)
				return ctrl.Result{}, err
			}
			continue
		}

		log.Info("Pod API Key found in quota, checking usages...", "client", clientSpec.Name)
		remainingMinutes, ok := quotaMap[clientSpec.Name]
		if !ok {
			log.Info("Client not found in quotaMap, deleting...", "client", clientSpec.Name)
			if err := r.Delete(ctx, &pod); err != nil {
				log.Error(err, "Failed to delete Pod", "pod", pod.Name)
				return ctrl.Result{}, err
			}
			continue
		}
		if remainingMinutes <= 0 {
			log.Info("Client has no remaining minutes, deleting...", "client", clientSpec.Name)
			if err := r.Delete(ctx, &pod); err != nil {
				log.Error(err, "Failed to delete Pod", "pod", pod.Name)
				return ctrl.Result{}, err
			}
			continue
		} else {
			quotaMap[clientSpec.Name] -= 1
		}

		if err := r.updateQuotaMap(ctx, &quota, quotaMap); err != nil {
			log.Error(err, "failed to update quota map")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClientQuotaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&quotav1alpha1.ClientQuota{}).
		Named("clientquota").
		Complete(r)
}

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

package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	quotav1alpha1 "github.com/YurgenUA/k8s-client-quota/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var clientquotalog = logf.Log.WithName("clientquota-resource")

// SetupClientQuotaWebhookWithManager registers the webhook for ClientQuota in the manager.
func SetupClientQuotaWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&quotav1alpha1.ClientQuota{}).
		WithValidator(&ClientQuotaCustomValidator{}).
		WithDefaulter(&ClientQuotaCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-quota-operator-k8s-yfenyuk-io-v1alpha1-clientquota,mutating=true,failurePolicy=fail,sideEffects=None,groups=quota.operator.k8s.yfenyuk.io,resources=clientquotas,verbs=create;update,versions=v1alpha1,name=mclientquota-v1alpha1.kb.io,admissionReviewVersions=v1

// ClientQuotaCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind ClientQuota when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ClientQuotaCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ClientQuotaCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind ClientQuota.
func (d *ClientQuotaCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	clientquota, ok := obj.(*quotav1alpha1.ClientQuota)

	if !ok {
		return fmt.Errorf("expected an ClientQuota object but got %T", obj)
	}
	clientquotalog.Info("Defaulting for ClientQuota", "name", clientquota.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-quota-operator-k8s-yfenyuk-io-v1alpha1-clientquota,mutating=false,failurePolicy=fail,sideEffects=None,groups=quota.operator.k8s.yfenyuk.io,resources=clientquotas,verbs=create;update,versions=v1alpha1,name=vclientquota-v1alpha1.kb.io,admissionReviewVersions=v1

// ClientQuotaCustomValidator struct is responsible for validating the ClientQuota resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ClientQuotaCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ClientQuotaCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ClientQuota.
func (v *ClientQuotaCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	clientquota, ok := obj.(*quotav1alpha1.ClientQuota)
	if !ok {
		return nil, fmt.Errorf("expected a ClientQuota object but got %T", obj)
	}
	clientquotalog.Info("Validation for ClientQuota upon creation", "name", clientquota.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ClientQuota.
func (v *ClientQuotaCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	clientquota, ok := newObj.(*quotav1alpha1.ClientQuota)
	if !ok {
		return nil, fmt.Errorf("expected a ClientQuota object for the newObj but got %T", newObj)
	}
	clientquotalog.Info("Validation for ClientQuota upon update", "name", clientquota.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ClientQuota.
func (v *ClientQuotaCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	clientquota, ok := obj.(*quotav1alpha1.ClientQuota)
	if !ok {
		return nil, fmt.Errorf("expected a ClientQuota object but got %T", obj)
	}
	clientquotalog.Info("Validation for ClientQuota upon deletion", "name", clientquota.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

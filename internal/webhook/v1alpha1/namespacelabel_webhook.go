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
	"github.com/TalDebi/namespacelabel/internal"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	danav1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
)

// nolint:unused

// SetupNamespaceLabelWebhookWithManager registers the webhook for NamespaceLabel in the manager.
func SetupNamespaceLabelWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&danav1alpha1.NamespaceLabel{}).
		WithValidator(&NamespaceLabelCustomValidator{Client: mgr.GetClient()}).
		Complete()
}

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-dana-dana-io-v1alpha1-namespacelabel,mutating=false,failurePolicy=fail,sideEffects=None,groups=dana.dana.io,resources=namespacelabels,verbs=create;update,versions=v1alpha1,name=vnamespacelabel-v1alpha1.kb.io,admissionReviewVersions=v1

// NamespaceLabelCustomValidator struct is responsible for validating the NamespaceLabel resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type NamespaceLabelCustomValidator struct {
	Client client.Client
}

var _ webhook.CustomValidator = &NamespaceLabelCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type NamespaceLabel.
func (v *NamespaceLabelCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	namespacelabel, err := fetchNamespaceLabel(ctx, obj)
	if err != nil {
		return nil, err
	}

	// Check if there is already a NamespaceLabel in the namespace.
	existingLabels, err := v.fetchNamespaceLabels(ctx, namespacelabel.Namespace)
	if err != nil {
		return nil, fmt.Errorf("could not validate NamespaceLabel: %v", err)
	}

	if len(existingLabels.Items) > 0 {
		return nil, fmt.Errorf("namespace '%s' already has %d NamespaceLabel(s). Only one is allowed per namespace", namespacelabel.Namespace, len(existingLabels.Items))
	}

	// Check if applied labels contain any management labels.
	if err := v.validateLabels(ctx, namespacelabel); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type NamespaceLabel.
func (v *NamespaceLabelCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	namespacelabel, err := fetchNamespaceLabel(ctx, newObj)
	if err != nil {
		return nil, err
	}

	// Check if applied labels contain any management labels.
	if err := v.validateLabels(ctx, namespacelabel); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type NamespaceLabel.
func (v *NamespaceLabelCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// fetchNamespaceLabel fetches NamespaceLabel
func fetchNamespaceLabel(ctx context.Context, obj runtime.Object) (*danav1alpha1.NamespaceLabel, error) {
	logger := log.FromContext(ctx)

	namespacelabel, ok := obj.(*danav1alpha1.NamespaceLabel)
	if !ok {
		return nil, fmt.Errorf("expected a NamespaceLabel object but got %T", obj)
	}

	logger.Info("fetched NamespaceLabel", "name", namespacelabel.GetName())

	return namespacelabel, nil
}

// fetchNamespaceLabels retrieves NamespaceLabel resources in the given namespace.
func (v *NamespaceLabelCustomValidator) fetchNamespaceLabels(ctx context.Context, namespaceName string) (*danav1alpha1.NamespaceLabelList, error) {
	logger := log.FromContext(ctx)

	existingNamespaceLabels := &danav1alpha1.NamespaceLabelList{}

	if err := v.Client.List(ctx, existingNamespaceLabels, client.InNamespace(namespaceName)); err != nil {
		logger.Error(err, "Failed to list NamespaceLabels", "namespaceName", namespaceName)
		return nil, err
	}

	logger.Info("Successfully fetched NamespaceLabels", "namespaceName", namespaceName, "count", len(existingNamespaceLabels.Items))

	return existingNamespaceLabels, nil
}

// validateLabels checks if any of the labels are management labels.
func (v *NamespaceLabelCustomValidator) validateLabels(ctx context.Context, namespacelabel *danav1alpha1.NamespaceLabel) error {
	logger := log.FromContext(ctx)

	for key := range namespacelabel.Spec.Labels {
		if internal.IsManagementLabel(key) {
			return fmt.Errorf("label '%s' is a management label and cannot be used", key)
		}
	}

	logger.Info("validated Namespacelabel's labels", "name", namespacelabel.GetName())

	return nil
}

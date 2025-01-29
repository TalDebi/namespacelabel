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
	"encoding/json"
	"fmt"
	"github.com/TalDebi/namespacelabel/internal"
	admissionv1 "k8s.io/api/admission/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	danav1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
)

// nolint:unused

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

func (v *NamespaceLabelCustomValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	logger := log.FromContext(ctx)
	var namespacelabel danav1alpha1.NamespaceLabel

	if err := json.Unmarshal(req.Object.Raw, &namespacelabel); err != nil {
		logger.Error(err, "Could not unmarshal raw object")
		return admission.Errored(http.StatusBadRequest, err)
	}

	if req.Operation == admissionv1.Create {
		existingLabels, err := v.fetchNamespaceLabels(ctx, namespacelabel.Namespace)
		if err != nil {
			return admission.Errored(http.StatusInternalServerError, fmt.Errorf("could not validate NamespaceLabel: %w", err))
		}

		if len(existingLabels.Items) > 0 {
			return admission.Denied(fmt.Sprintf("namespace '%s' already has %d NamespaceLabel(s). Only one is allowed per namespace", namespacelabel.Namespace, len(existingLabels.Items)))
		}
	}

	if err := v.validateLabels(ctx, &namespacelabel); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	return admission.Allowed("NamespaceLabel validated successfully")
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

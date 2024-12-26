package controller

import (
	"context"
	"fmt"
	"strings"

	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// Management label prefix
const managementLabelPrefix = "kubernetes.io"

// Reconcile the namespace labels based on NamespaceLabel spec
func (r *NamespaceLabelReconciler) reconcileNamespaceLabels(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) error {

	// Track labels to add and remove
	labelsToAdd, labelsToRemove, err := determineLabelChanges(namespaceLabel, ns)

	if err != nil {
		return err
	}

	if err := r.applyLabelsChanges(ctx, ns, labelsToAdd, labelsToRemove); err != nil {
		return err
	}

	return nil
}

// determines labels to add, remove or update
func determineLabelChanges(namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) (map[string]string, map[string]struct{}, error) {
	labelsToAdd := make(map[string]string)
	labelsToRemove := make(map[string]struct{})

	// Collect labels to add or update
	for key, value := range namespaceLabel.Spec.Labels {
		if isManagementLabel(key) {
			return nil, nil, fmt.Errorf("cannot add protected or management label '%s'", key)
		}

		labelsToAdd[key] = value
	}

	// Collect labels to remove
	if ns.Labels != nil {
		for key := range ns.Labels {
			if _, exists := labelsToAdd[key]; !exists && !isManagementLabel(key) {
				labelsToRemove[key] = struct{}{}
			}
		}
	}

	return labelsToAdd, labelsToRemove, nil
}

// check if label is a management label
func isManagementLabel(label string) bool {
	return strings.HasPrefix(label, managementLabelPrefix)
}

// applyLabelsChanges removes and applies labels to the Namespace object
func (r *NamespaceLabelReconciler) applyLabelsChanges(
	ctx context.Context, ns *corev1.Namespace, labelsToAdd map[string]string, labelsToRemove map[string]struct{}) error {

	// Remove labels that are no longer present in NamespaceLabel
	for key := range labelsToRemove {
		delete(ns.Labels, key)
	}

	// Initialize ns.Labels if nil
	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}

	// Apply labels to be added or updated
	for key, value := range labelsToAdd {
		ns.Labels[key] = value
	}

	// Update Namespace with new labels
	if err := r.Update(ctx, ns); err != nil {
		return err
	}

	return nil
}

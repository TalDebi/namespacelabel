package controller

import (
	"context"
	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	"github.com/TalDebi/namespacelabel/internal"
	corev1 "k8s.io/api/core/v1"
)

// reconcileNamespaceLabels reconciles the namespace labels based on NamespaceLabel spec
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

// determineLabelChanges determines labels to add, remove or update
func determineLabelChanges(namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) (map[string]string, map[string]struct{}, error) {
	labelsToAdd := make(map[string]string)
	labelsToRemove := make(map[string]struct{})

	labelsToAdd = collectLabelsToAddOrUpdate(namespaceLabel)

	labelsToRemove = collectLabelsToRemove(ns.Labels, labelsToAdd)

	return labelsToAdd, labelsToRemove, nil
}

// collectLabelsToAddOrUpdate collects the labels from the namespaceLabel and returns a map of labels to add or update.
func collectLabelsToAddOrUpdate(namespaceLabel *danaiov1alpha1.NamespaceLabel) map[string]string {
	labelsToAdd := make(map[string]string)

	for key, value := range namespaceLabel.Spec.Labels {
		labelsToAdd[key] = value
	}

	return labelsToAdd
}

// collectLabelsToRemove identifies labels to be removed from the namespace
func collectLabelsToRemove(nsLabels map[string]string, labelsToAdd map[string]string) map[string]struct{} {
	labelsToRemove := make(map[string]struct{})

	if nsLabels != nil {
		for key := range nsLabels {
			if _, exists := labelsToAdd[key]; !exists && !internal.IsManagementLabel(key) {
				labelsToRemove[key] = struct{}{}
			}
		}
	}

	return labelsToRemove
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

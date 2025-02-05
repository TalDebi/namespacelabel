package controller

import (
	"context"

	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	"github.com/TalDebi/namespacelabel/internal"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// reconcileNamespaceLabels reconciles the namespace labels based on NamespaceLabel spec.
func (r *NamespaceLabelReconciler) reconcileNamespaceLabels(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) error {

	labelsToAdd, labelsToRemove := determineLabelChanges(namespaceLabel, ns)

	if err := r.applyLabelsChanges(ctx, ns, labelsToAdd, labelsToRemove); err != nil {
		return err
	}

	return nil
}

// determineLabelChanges determines labels to add, remove or update.
func determineLabelChanges(namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) (map[string]string, map[string]struct{}) {
	labelsToAdd := collectLabelsToAddOrUpdate(namespaceLabel)

	labelsToRemove := collectLabelsToRemove(ns.Labels, labelsToAdd)

	return labelsToAdd, labelsToRemove
}

// collectLabelsToAddOrUpdate collects the labels from the namespaceLabel and returns a map of labels to add or update.
func collectLabelsToAddOrUpdate(namespaceLabel *danaiov1alpha1.NamespaceLabel) map[string]string {
	labelsToAdd := make(map[string]string)

	for key, value := range namespaceLabel.Spec.Labels {
		labelsToAdd[key] = value
	}

	return labelsToAdd
}

// collectLabelsToRemove identifies labels to be removed from the namespace.
func collectLabelsToRemove(nsLabels map[string]string, labelsToAdd map[string]string) map[string]struct{} {
	labelsToRemove := make(map[string]struct{})

	for key := range nsLabels {
		if _, exists := labelsToAdd[key]; !exists && !internal.IsManagementLabel(key) {
			labelsToRemove[key] = struct{}{}
		}
	}

	return labelsToRemove
}

// applyLabelsChanges removes and applies labels to the Namespace object.
func (r *NamespaceLabelReconciler) applyLabelsChanges(
	ctx context.Context, ns *corev1.Namespace, labelsToAdd map[string]string, labelsToRemove map[string]struct{}) error {
	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}

	for key := range labelsToRemove {
		delete(ns.Labels, key)
	}

	for key, value := range labelsToAdd {
		ns.Labels[key] = value
	}

	if err := r.Update(ctx, ns); err != nil {
		return err
	}

	return nil
}

// removeLabelsFromNamespace removes all the labels not active on the nsl from the current namespace.
func (r *NamespaceLabelReconciler) removeLabelsFromNamespace(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) error {
	logger := log.FromContext(ctx)

	for key := range namespaceLabel.Spec.Labels {
		delete(ns.Labels, key)
	}

	if err := r.Update(ctx, ns); err != nil {
		logger.Error(err, "Failed to update Namespace after removing labels")
		return err
	}

	logger.Info("Labels removed from Namespace successfully", "Namespace", ns.Name)
	return nil
}

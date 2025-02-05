package controller

import (
	"context"

	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	"github.com/TalDebi/namespacelabel/internal"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// handleDeletion handles the process of nsl deletion.
func (r *NamespaceLabelReconciler) handleDeletion(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) error {
	logger := log.FromContext(ctx)

	if err := r.removeLabelsFromNamespace(ctx, namespaceLabel, ns); err != nil {
		logger.Error(err, "Failed to remove labels from Namespace", "Namespace", ns.Name)
		return err
	}

	if err := r.removeFinalizer(ctx, namespaceLabel); err != nil {
		logger.Error(err, "Failed to remove finalizer", "NamespaceLabel", namespaceLabel.Name)
		return err
	}

	logger.Info("Deletion handled successfully", "NamespaceLabel", namespaceLabel.Name)
	return nil
}

// removeFinalizer removes the finalizer from the nsl during the deletion process.
func (r *NamespaceLabelReconciler) removeFinalizer(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel) error {
	logger := log.FromContext(ctx)

	controllerutil.RemoveFinalizer(namespaceLabel, internal.FinalizerName)

	if err := r.Update(ctx, namespaceLabel); err != nil {
		logger.Error(err, "Failed to update NamespaceLabel to remove finalizer")
		return err
	}

	logger.Info("Finalizer removed successfully", "NamespaceLabel", namespaceLabel.Name)
	return nil
}

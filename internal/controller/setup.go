package controller

import (
	"context"
	"github.com/TalDebi/namespacelabel/internal"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&danaiov1alpha1.NamespaceLabel{}).
		Named("namespacelabel").
		Complete(r)
}

// handleDeletion handles the process of nsl deletion
func (r *NamespaceLabelReconciler) handleDeletion(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) error {
	logger := log.FromContext(ctx)

	// Remove labels from the namespace
	if err := r.removeLabelsFromNamespace(ctx, namespaceLabel, ns); err != nil {
		logger.Error(err, "Failed to remove labels from Namespace", "Namespace", ns.Name)
		return err
	}

	// Remove finalizer from the NamespaceLabel
	if err := r.removeFinalizer(ctx, namespaceLabel); err != nil {
		logger.Error(err, "Failed to remove finalizer", "NamespaceLabel", namespaceLabel.Name)
		return err
	}

	logger.Info("Deletion handled successfully", "NamespaceLabel", namespaceLabel.Name)
	return nil
}

// removeLabelsFromNamespace removes all the labels not active on the nsl from the current namespace
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

// removeFinalizer removes the finalizer from the nsl during the deletion process
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

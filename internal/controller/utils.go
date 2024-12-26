package controller

import (
	"context"
	"fmt"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// Fetch NamespaceLabel
func (r *NamespaceLabelReconciler) fetchNamespaceLabel(ctx context.Context, req ctrl.Request) (*danaiov1alpha1.NamespaceLabel, error) {
	namespaceLabel := &danaiov1alpha1.NamespaceLabel{}
	if err := r.Get(ctx, req.NamespacedName, namespaceLabel); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return namespaceLabel, nil
}

// Fetch All NamespaceLabels in the Namespace
func (r *NamespaceLabelReconciler) listNamespaceLabelsInNamespace(ctx context.Context, req ctrl.Request) (*danaiov1alpha1.NamespaceLabelList, error) {
	existingNamespaceLabels := &danaiov1alpha1.NamespaceLabelList{}
	if err := r.List(ctx, existingNamespaceLabels, client.InNamespace(req.Namespace)); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	return existingNamespaceLabels, nil
}

// Fetch Namespace
func (r *NamespaceLabelReconciler) fetchNamespace(ctx context.Context, req ctrl.Request) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	namespaceName := client.ObjectKey{Name: req.Namespace}
	if err := r.Get(ctx, namespaceName, ns); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return ns, nil
}

// Ensure only one NamespaceLabel per namespace
func (r *NamespaceLabelReconciler) ensureSingleNamespaceLabel(ctx context.Context, req ctrl.Request, namespaceLabel *danaiov1alpha1.NamespaceLabel) error {
	existingNamespaceLabels, err := r.listNamespaceLabelsInNamespace(ctx, req)
	if err != nil {
		return err
	}

	if len(existingNamespaceLabels.Items) > 1 {
		return fmt.Errorf("only one NamespaceLabel allowed per namespace")
	}

	return nil
}

func (r *NamespaceLabelReconciler) handleDeletion(
	ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, ns *corev1.Namespace) (ctrl.Result, error) {
	// Remove labels managed by this NamespaceLabel
	for key := range namespaceLabel.Spec.Labels {
		delete(ns.Labels, key)
	}

	if err := r.Update(ctx, ns); err != nil {
		return ctrl.Result{}, err
	}
	controllerutil.RemoveFinalizer(namespaceLabel, finalizerName)
	if err := r.Update(ctx, namespaceLabel); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

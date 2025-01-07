package controller

import (
	"context"
	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// fetchNamespaceLabel fetches the current NamespaceLabel details
func (r *NamespaceLabelReconciler) fetchNamespaceLabel(ctx context.Context, req ctrl.Request) (*danaiov1alpha1.NamespaceLabel, error) {
	namespaceLabel := &danaiov1alpha1.NamespaceLabel{}
	if err := r.Get(ctx, req.NamespacedName, namespaceLabel); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return namespaceLabel, nil
}

// listNamespaceLabelsInNamespace fetches All NamespaceLabels in the Namespace
func (r *NamespaceLabelReconciler) listNamespaceLabelsInNamespace(ctx context.Context, req ctrl.Request) (*danaiov1alpha1.NamespaceLabelList, error) {
	existingNamespaceLabels := &danaiov1alpha1.NamespaceLabelList{}
	if err := r.List(ctx, existingNamespaceLabels, client.InNamespace(req.Namespace)); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	return existingNamespaceLabels, nil
}

// fetchNamespace fetches the current Namespace details
func (r *NamespaceLabelReconciler) fetchNamespace(ctx context.Context, req ctrl.Request) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	namespaceName := client.ObjectKey{Name: req.Namespace}
	if err := r.Get(ctx, namespaceName, ns); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return ns, nil
}

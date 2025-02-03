package internal

import (
	"context"
	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

// IsManagementLabel check if label is a management label.
func IsManagementLabel(label string) bool {
	for _, prefix := range ManagementLabelPrefixes {
		if strings.HasPrefix(label, prefix) {
			return true
		}
	}
	return false
}

// FetchNamespaceLabel fetches the current NamespaceLabel details.
func FetchNamespaceLabel(ctx context.Context, c client.Reader, namespacedName types.NamespacedName) (*danaiov1alpha1.NamespaceLabel, error) {
	namespaceLabel := &danaiov1alpha1.NamespaceLabel{}
	if err := c.Get(ctx, namespacedName, namespaceLabel); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return namespaceLabel, nil
}

// FetchNamespaceLabels retrieves NamespaceLabel resources in the given namespace.
func FetchNamespaceLabels(ctx context.Context, c client.Reader, namespaceName string) (*danaiov1alpha1.NamespaceLabelList, error) {
	logger := log.FromContext(ctx)

	existingNamespaceLabels := &danaiov1alpha1.NamespaceLabelList{}

	if err := c.List(ctx, existingNamespaceLabels, client.InNamespace(namespaceName)); err != nil {
		logger.Error(err, "Failed to list NamespaceLabels", "namespaceName", namespaceName)
		return nil, err
	}

	logger.Info("Successfully fetched NamespaceLabels", "namespaceName", namespaceName, "count", len(existingNamespaceLabels.Items))

	return existingNamespaceLabels, nil
}

// FetchNamespace fetches the current Namespace details.
func FetchNamespace(ctx context.Context, c client.Reader, namespacedName string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	if err := c.Get(ctx, types.NamespacedName{Namespace: namespacedName, Name: namespacedName}, ns); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return ns, nil
}

// ListNamespaceLabelsInNamespace fetches all namespaceLabels in a namespace.
func ListNamespaceLabelsInNamespace(ctx context.Context, c client.Reader, namespaceName string) (*danaiov1alpha1.NamespaceLabelList, error) {
	existingNamespaceLabels := &danaiov1alpha1.NamespaceLabelList{}
	if err := c.List(ctx, existingNamespaceLabels, client.InNamespace(namespaceName)); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	return existingNamespaceLabels, nil
}

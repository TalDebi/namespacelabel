package controller

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
)

// updateConditionsStatus updates the conditions variable
func (r *NamespaceLabelReconciler) updateConditionsStatus(ctx context.Context, namespaceLabel *danaiov1alpha1.NamespaceLabel, conditionType string, status metav1.ConditionStatus, reason, message string) error {
	condition := metav1.Condition{
		Type:               conditionType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}

	// Update or append condition
	namespaceLabel.Status.Conditions = updateNewCondition(namespaceLabel.Status.Conditions, condition)

	// Update status
	if err := r.Status().Update(ctx, namespaceLabel); err != nil {
		r.Log.Error(err, "Failed to update NamespaceLabel status", "NamespaceLabel", namespaceLabel.Name)
		return err
	}

	return nil
}

// updateNewCondition appends a new condition or updates an existing one in the slice of conditions
func updateNewCondition(conditions []metav1.Condition, newCondition metav1.Condition) []metav1.Condition {
	for index := range conditions {
		if conditions[index].Type == newCondition.Type {
			conditions[index] = newCondition
			return conditions
		}
	}

	return append(conditions, newCondition)
}

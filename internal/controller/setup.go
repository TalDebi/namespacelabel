package controller

import (
	danaiov1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&danaiov1alpha1.NamespaceLabel{}).
		Named("namespacelabel").
		Complete(r)
}

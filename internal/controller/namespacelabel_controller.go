/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/TalDebi/namespacelabel/internal"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var finalizerName = internal.FinalizerName

// NamespaceLabelReconciler reconciles a NamespaceLabel object
type NamespaceLabelReconciler struct {
	Scheme *runtime.Scheme
	client.Client
	Log logr.Logger
}

// +kubebuilder:rbac:groups=dana.dana.io,resources=namespacelabels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dana.dana.io,resources=namespacelabels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dana.dana.io,resources=namespacelabels/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=list;watch;get;update

func (r *NamespaceLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	logger := log.FromContext(ctx)

	logger.Info("Starting reconciliation for NamespaceLabel", "Namespace", req.Namespace, "Name", req.Name)

	// Fetch the NamespaceLabel instance
	namespaceLabel, err := r.fetchNamespaceLabel(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	if namespaceLabel == nil {
		return ctrl.Result{}, nil
	}

	logger.Info("Fetched NamespaceLabel", "NamespaceLabel", namespaceLabel)

	// Fetch the Namespace instance
	ns, err := r.fetchNamespace(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}
	if ns == nil {
		return ctrl.Result{}, nil
	}

	logger.Info("Fetched Namespace", "NamespaceLabel", ns)

	// Handle finalizer
	if namespaceLabel.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(namespaceLabel, finalizerName) {
			controllerutil.AddFinalizer(namespaceLabel, finalizerName)
			if err := r.Update(ctx, namespaceLabel); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(namespaceLabel, finalizerName) {
			if err := r.handleDeletion(ctx, namespaceLabel, ns); err != nil {
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	logger.Info("Creating nsl")

	// Reconcile the namespace labels
	if err := r.reconcileNamespaceLabels(ctx, namespaceLabel, ns); err != nil {
		r.updateConditionsStatus(ctx, namespaceLabel, "UpdateLabelsFailed", metav1.ConditionFalse, "UpdateError", err.Error())
		return ctrl.Result{}, err
	}

	r.updateConditionsStatus(ctx, namespaceLabel, "LabelsApplied", metav1.ConditionTrue, "Success", "Namespace labels have been successfully updated")
	logger.Info("nsl Created")

	return ctrl.Result{}, nil
}

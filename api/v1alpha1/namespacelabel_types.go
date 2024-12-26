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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceLabelSpec defines the desired state of NamespaceLabel.
type NamespaceLabelSpec struct {
	// Labels to be added to the Namespace
	Labels map[string]string `json:"labels,omitempty" yaml:"labels"`
}

// NamespaceLabelStatus defines the observed state of NamespaceLabel.
type NamespaceLabelStatus struct {
	// AppliedLabels shows the labels that have been successfully applied
	AppliedLabels map[string]string `json:"appliedLabels,omitempty" yaml:"appliedLabels"`
	// Conditions represents the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// NamespaceLabel is the Schema for the namespacelabels API.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=nsl
// +kubebuilder:printcolumn:name="Labels",type=string,JSONPath=".spec.labels",description="Labels applied to the Namespace", priority=10
// +kubebuilder:object:generate=true
type NamespaceLabel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceLabelSpec   `json:"spec,omitempty"`
	Status NamespaceLabelStatus `json:"status,omitempty"`
}

// NamespaceLabelList contains a list of NamespaceLabel.
// +kubebuilder:object:root=true
// +kubebuilder:object:generate=true
type NamespaceLabelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceLabel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceLabel{}, &NamespaceLabelList{})
}

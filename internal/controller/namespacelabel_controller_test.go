/*
Copyright 2025.

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
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	danav1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
)

var _ = Describe("NamespaceLabel Controller", func() {
	Context("When reconciling a resource", func() {
		const (
			timeout  = time.Second * 10
			interval = time.Millisecond * 500
		)
		testLabels := map[string]string{"env": "test", "team": "dev"}

		ctx := context.Background()
		var typeNamespacedName types.NamespacedName

		BeforeEach(func() {
			resourceName := fmt.Sprintf("test-resource-%d", time.Now().UnixNano())
			testNamespace := fmt.Sprintf("test-ns-%d", time.Now().UnixNano())
			typeNamespacedName = types.NamespacedName{
				Name:      resourceName,
				Namespace: testNamespace,
			}
			By("creating the test namespace")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: typeNamespacedName.Namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("creating the custom resource for the Kind NamespaceLabel")
			namespacelabel := &danav1alpha1.NamespaceLabel{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: testNamespace,
				},
				Spec: danav1alpha1.NamespaceLabelSpec{
					Labels: testLabels,
				},
			}
			Expect(k8sClient.Create(ctx, namespacelabel)).To(Succeed())

			By("Reconciling the created resource")
			controllerReconciler := &NamespaceLabelReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, e := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(e).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			//By("Cleanup the specific resource instance NamespaceLabel")
			//resource := &danav1alpha1.NamespaceLabel{}
			//err := k8sClient.Get(ctx, typeNamespacedName, resource)
			//
			//if err == nil || !errors.IsNotFound(err) {
			//	// Check if the NamespaceLabel is in a stable state before deletion.
			//	Eventually(func() error {
			//		return k8sClient.Get(ctx, typeNamespacedName, resource)
			//	}, timeout, interval).Should(Succeed())
			//
			//	// Check that the resource isn't being updated concurrently (important!)
			//	Consistently(func() error {
			//		return k8sClient.Get(ctx, typeNamespacedName, resource)
			//	}, interval).Should(Succeed())
			//
			//	Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			//	Eventually(func() error {
			//		return k8sClient.Get(ctx, typeNamespacedName, resource)
			//	}, timeout, interval).Should(Satisfy(errors.IsNotFound))
			//} else if !errors.IsNotFound(err) {
			//	Expect(err).To(Succeed())
			//}
			//
			By("cleaning up the test namespace")
			ns := &corev1.Namespace{}
			nsErr := k8sClient.Get(ctx, types.NamespacedName{Name: typeNamespacedName.Namespace}, ns)

			if nsErr == nil || !errors.IsNotFound(nsErr) {
				Expect(k8sClient.Delete(ctx, ns)).To(Succeed())
				Eventually(func() error {
					return k8sClient.Get(ctx, types.NamespacedName{Name: typeNamespacedName.Namespace}, &corev1.Namespace{})
				}, timeout, interval).Should(Satisfy(errors.IsNotFound))
			} else if !errors.IsNotFound(nsErr) {
				Expect(nsErr).To(Succeed())
			}
		})
		It("should successfully reconcile the resource", func() {
			By("Verifying that labels are applied to the namespace")
			updatedNamespace := &corev1.Namespace{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: typeNamespacedName.Namespace}, updatedNamespace)).To(Succeed())
			for key, value := range testLabels {
				Expect(updatedNamespace.Labels).To(HaveKeyWithValue(key, value))
			}

			By("Verifying that NamespaceLabelStatus is updated")
			updatedNamespaceLabel := &danav1alpha1.NamespaceLabel{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, updatedNamespaceLabel)).To(Succeed())

			By("Checking if the status conditions are updated")
			Eventually(func() []metav1.Condition {
				_ = k8sClient.Get(ctx, typeNamespacedName, updatedNamespaceLabel)
				return updatedNamespaceLabel.Status.Conditions
			}, timeout, interval).ShouldNot(BeEmpty())
		})

		It("should update labels in an existing NamespaceLabel", func() {
			By("Updating the NamespaceLabel labels")
			updatedLabels := map[string]string{"env": "prod", "owner": "admin"}
			existingNamespaceLabel := &danav1alpha1.NamespaceLabel{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, existingNamespaceLabel)).To(Succeed())

			existingNamespaceLabel.Spec.Labels = updatedLabels
			Expect(k8sClient.Update(ctx, existingNamespaceLabel)).To(Succeed())

			By("Reconciling the updated resource")
			controllerReconciler := &NamespaceLabelReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying that labels are updated in the namespace")
			updatedNamespace := &corev1.Namespace{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: typeNamespacedName.Namespace}, updatedNamespace)).To(Succeed())

			for key, value := range updatedLabels {
				Expect(updatedNamespace.Labels).To(HaveKeyWithValue(key, value))
			}
		})

		It("should remove labels when NamespaceLabel is deleted", func() {
			By("Deleting the NamespaceLabel resource")
			existingNamespaceLabel := &danav1alpha1.NamespaceLabel{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, existingNamespaceLabel)).To(Succeed())
			Expect(k8sClient.Delete(ctx, existingNamespaceLabel)).To(Succeed())

			By("Reconciling after deletion")
			controllerReconciler := &NamespaceLabelReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying that labels are removed from the namespace")
			updatedNamespace := &corev1.Namespace{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: typeNamespacedName.Namespace}, updatedNamespace)).To(Succeed())

			for key := range testLabels {
				Expect(updatedNamespace.Labels).NotTo(HaveKey(key))
			}
		})
	})
})

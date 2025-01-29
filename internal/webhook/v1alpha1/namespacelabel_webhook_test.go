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

package v1alpha1

import (
	"context"
	"fmt"
	"github.com/TalDebi/namespacelabel/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"time"

	danav1alpha1 "github.com/TalDebi/namespacelabel/api/v1alpha1"
)

var _ = Describe("NamespaceLabel Webhook", func() {
	var (
		obj       *danav1alpha1.NamespaceLabel
		validator NamespaceLabelCustomValidator
		ctx       context.Context
		scheme    *runtime.Scheme
	)

	testLabels := map[string]string{"team": "devops"}

	BeforeEach(func() {
		ctx = context.Background()
		scheme = runtime.NewScheme()
		Expect(danav1alpha1.AddToScheme(scheme)).To(Succeed())
		resourceName := fmt.Sprintf("test-resource-%d", time.Now().UnixNano())
		testNamespace := fmt.Sprintf("test-ns-%d", time.Now().UnixNano())
		obj = &danav1alpha1.NamespaceLabel{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: testNamespace,
			},
			Spec: danav1alpha1.NamespaceLabelSpec{
				Labels: testLabels,
			},
		}
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
		validator = NamespaceLabelCustomValidator{Client: fakeClient}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		Expect(obj).NotTo(BeNil(), "Expected obj to be initialized")
	})

	Context("When creating NamespaceLabel", func() {
		It("Should allow creation if no existing NamespaceLabel is in the namespace", func() {
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should reject creation if a NamespaceLabel already exists in the namespace", func() {
			validator.Client = fake.NewClientBuilder().WithScheme(scheme).WithObjects(obj).Build()

			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
		})

		It("Should reject creation if a management label is used", func() {
			obj.Spec.Labels[internal.ManagementLabelPrefix] = "test"
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("When updating NamespaceLabel", func() {
		It("Should allow updating when labels are valid", func() {
			newObj := obj.DeepCopy()
			newObj.Spec.Labels = testLabels
			_, err := validator.ValidateUpdate(ctx, obj, newObj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should reject update if it introduces a management label", func() {
			existingObj := obj.DeepCopy()
			obj.Spec.Labels[internal.ManagementLabelPrefix] = "test"
			_, err := validator.ValidateUpdate(ctx, existingObj, obj)
			Expect(err).To(HaveOccurred())
		})
	})

})

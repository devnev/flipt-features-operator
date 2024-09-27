/*
Copyright (C) 2024.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fliptv1alpha1 "github.com/devnev/flipt-features-operator/api/v1alpha1"
)

var _ = Describe("FeaturesTarget Controller", func() {
	Context("When reconciling a resource", func() {
		ctx := context.Background()
		resourceNamespacedName := types.NamespacedName{
			Name:      "test-resource",
			Namespace: "default",
		}

		BeforeEach(func() {
			By("clearing any debris from previous tests")
			err := k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(resourceNamespacedName.Namespace))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &fliptv1alpha1.FeaturesTarget{}, client.InNamespace(resourceNamespacedName.Namespace))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &fliptv1alpha1.Features{}, client.InNamespace(resourceNamespacedName.Namespace))
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			By("clearing any debris from the test")
			err := k8sClient.DeleteAllOf(ctx, &corev1.ConfigMap{}, client.InNamespace(resourceNamespacedName.Namespace))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &fliptv1alpha1.FeaturesTarget{}, client.InNamespace(resourceNamespacedName.Namespace))
			Expect(err).NotTo(HaveOccurred())
			err = k8sClient.DeleteAllOf(ctx, &fliptv1alpha1.Features{}, client.InNamespace(resourceNamespacedName.Namespace))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully reconcile a resource", func() {
			const configMapName = "flipt-local"
			By("creating a resource with no sources")
			createFeaturesTargetWithNoSources(resourceNamespacedName, configMapName)

			By("reconciling the created resource")
			reconcileFeaturesTarget(k8sClient, resourceNamespacedName)

			By("fetching the target configmap")
			var cm corev1.ConfigMap
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: resourceNamespacedName.Namespace, Name: configMapName}, &cm)
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).To(BeEmpty())
		})

		It("should reconcile a FeaturesTarget with no source and ignore any Features resources", func() {
			const configMapName = "flipt-local"
			By("creating a FeaturesTarget with no sources")
			createFeaturesTargetWithNoSources(resourceNamespacedName, configMapName)
			By("creating a Features resource")
			createFeaturesEmptyFeatures(resourceNamespacedName, "default")

			By("reconciling the created resource")
			reconcileFeaturesTarget(k8sClient, resourceNamespacedName)

			By("fetching the target configmap")
			var cm corev1.ConfigMap
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: resourceNamespacedName.Namespace, Name: configMapName}, &cm)
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).To(BeEmpty())
		})

		It("should reconcile a FeaturesTarget with missing selector to include an existing Feature", func() {
			const configMapName = "flipt-local"
			By("creating a FeaturesTarget with a nothing selector source")
			createFeaturesTargetWithNilSelectorSource(resourceNamespacedName, configMapName)
			By("creating a Features resource")
			createFeaturesEmptyFeatures(resourceNamespacedName, "default")

			By("reconciling the created resource")
			reconcileFeaturesTarget(k8sClient, resourceNamespacedName)

			By("fetching the target configmap")
			var cm corev1.ConfigMap
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: resourceNamespacedName.Namespace, Name: configMapName}, &cm)
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).NotTo(BeEmpty())
			Expect(cm.Data).To(HaveKey("default_test-resource.yml"))
			Expect(cm.Data).To(HaveLen(1))
		})

		It("should reconcile a FeaturesTarget with an empty selector to include an existing Feature", func() {
			const configMapName = "flipt-local"
			By("creating a FeaturesTarget with an everything selector source")
			createFeaturesTargetWithEmptySelectorSource(resourceNamespacedName, configMapName)
			By("creating a Features resource")
			createFeaturesEmptyFeatures(resourceNamespacedName, "default")

			By("reconciling the created resource")
			reconcileFeaturesTarget(k8sClient, resourceNamespacedName)

			By("fetching the target configmap")
			var cm corev1.ConfigMap
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: resourceNamespacedName.Namespace, Name: configMapName}, &cm)
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).NotTo(BeEmpty())
			Expect(cm.Data).To(HaveKey("default_test-resource.yml"))
			Expect(cm.Data).To(HaveLen(1))
		})

		It("should reconcile a FeaturesTarget and matching Features to a configmap value", func() {
			const configMapName = "flipt-local"
			By("creating a FeaturesTarget with an everything selector source")
			createFeaturesTargetWithEmptySelectorSource(resourceNamespacedName, configMapName)
			By("creating a Features resource")
			createFeaturesFullyPopulated(resourceNamespacedName, "default")

			By("reconciling the created resource")
			reconcileFeaturesTarget(k8sClient, resourceNamespacedName)

			By("fetching the target configmap")
			var cm corev1.ConfigMap
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: resourceNamespacedName.Namespace, Name: configMapName}, &cm)
			Expect(err).NotTo(HaveOccurred())
			Expect(cm.Data).NotTo(BeEmpty())
			Expect(cm.Data).To(HaveKey("default_test-resource.yml"))
			Expect(cm.Data).To(HaveLen(1))
			Expect(cm.Data["default_test-resource.yml"]).To(MatchYAML(`
namespace:
  key: default
flags:
- key: bool-flag
  name: Boolean
  type: boolean
  description: A boolean Flag
  enabled: true
- key: variant-flag
  name: Variant
  description: A boolean Flag
  enabled: false
  variants:
  - key: var-key
    name: Variant Name
    description: A
segments:
- key: segment
  name: Segment
  description: A Segment
  constraints:
  - {}`))
		})
	})
})

func reconcileFeaturesTarget(k8sClient client.Client, namespacedName types.NamespacedName) {
	controllerReconciler := &FeaturesTargetReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),
	}
	_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
		NamespacedName: namespacedName,
	})
	Expect(err).NotTo(HaveOccurred())
}

func createFeaturesTargetWithNoSources(namespacedName types.NamespacedName, configMapName string) {
	resource := &fliptv1alpha1.FeaturesTarget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: fliptv1alpha1.FeaturesTargetSpec{
			ConfigMap: &fliptv1alpha1.FeaturesTargetSpecConfigMap{
				Name: configMapName,
			},
		},
	}
	Expect(k8sClient.Create(ctx, resource)).To(Succeed())
}

func createFeaturesTargetWithNilSelectorSource(namespacedName types.NamespacedName, configMapName string) {
	resource := &fliptv1alpha1.FeaturesTarget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: fliptv1alpha1.FeaturesTargetSpec{
			ConfigMap: &fliptv1alpha1.FeaturesTargetSpecConfigMap{
				Name: configMapName,
			},
			Sources: []fliptv1alpha1.FeaturesTargetSpecSource{{
				Selector: nil,
			}},
		},
	}
	Expect(k8sClient.Create(ctx, resource)).To(Succeed())
}

func createFeaturesTargetWithEmptySelectorSource(namespacedName types.NamespacedName, configMapName string) {
	resource := &fliptv1alpha1.FeaturesTarget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: fliptv1alpha1.FeaturesTargetSpec{
			ConfigMap: &fliptv1alpha1.FeaturesTargetSpecConfigMap{
				Name: configMapName,
			},
			Sources: []fliptv1alpha1.FeaturesTargetSpecSource{{
				Selector: &metav1.LabelSelector{},
			}},
		},
	}
	Expect(k8sClient.Create(ctx, resource)).To(Succeed())
}

func createFeaturesEmptyFeatures(namespacedName types.NamespacedName, featuresNamespace string) {
	resource := &fliptv1alpha1.Features{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: fliptv1alpha1.FeaturesSpec{
			Features: fliptv1alpha1.Document{
				Namespace: &fliptv1alpha1.Namespace{
					Key: featuresNamespace,
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, resource)).To(Succeed())
}

func createFeaturesFullyPopulated(namespacedName types.NamespacedName, featuresNamespace string) {
	resource := &fliptv1alpha1.Features{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: fliptv1alpha1.FeaturesSpec{
			Features: fliptv1alpha1.Document{
				Namespace: &fliptv1alpha1.Namespace{
					Key: featuresNamespace,
				},
				// TODO: rules and rollouts
				Flags: []*fliptv1alpha1.Flag{{
					Key:         "bool-flag",
					Name:        "Boolean",
					Description: "A boolean Flag",
					Enabled:     true,
					Type:        "boolean",
					// FIXME: arbitrary metadata doesn't serialize yet
					// Metadata:    map[string]json.RawMessage{"md": []byte("1234")},
				}, {
					Key:         "variant-flag",
					Name:        "Variant",
					Description: "A boolean Flag",
					Enabled:     false,
					// FIXME: arbitrary metadata doesn't serialize yet
					// Metadata:    map[string]json.RawMessage{"md": []byte("1234")},
					Variants: []*fliptv1alpha1.Variant{{Default: false, Key: "var-key", Name: "Variant Name", Description: "A"}},
				}},
				// TODO: other segment types
				Segments: []*fliptv1alpha1.Segment{{
					Key:         "segment",
					Name:        "Segment",
					Description: "A Segment",
					Constraints: []*fliptv1alpha1.Constraint{{}},
				}},
			},
		},
	}
	Expect(k8sClient.Create(ctx, resource)).To(Succeed())
}

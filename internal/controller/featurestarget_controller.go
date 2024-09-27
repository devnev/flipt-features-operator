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
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	yaml "sigs.k8s.io/yaml/goyaml.v2"

	fliptv1alpha1 "github.com/devnev/flipt-features-operator/api/v1alpha1"
)

// FeaturesTargetReconciler reconciles a FeaturesTarget object
type FeaturesTargetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=flipt.nev.dev,resources=featurestargets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flipt.nev.dev,resources=featurestargets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=flipt.nev.dev,resources=featurestargets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FeaturesTarget object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *FeaturesTargetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var featuresTarget fliptv1alpha1.FeaturesTarget
	if err := r.Get(ctx, req.NamespacedName, &featuresTarget); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	configmap := &corev1.ConfigMap{}
	if err := controllerutil.SetControllerReference(&featuresTarget, configmap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	configmap.Data = make(map[string]string)

	for _, source := range featuresTarget.Spec.Sources {
		var listOpts client.ListOptions
		if sel, err := metav1.LabelSelectorAsSelector(source.Selector); err != nil {
			return ctrl.Result{}, err
		} else {
			listOpts.LabelSelector = sel
		}

		if len(source.Namespaces) > 0 {
			for _, ns := range source.Namespaces {
				listOpts.Namespace = ns
				var featuresList fliptv1alpha1.FeaturesList
				if err := r.List(ctx, &featuresList, &listOpts); err != nil {
					return ctrl.Result{}, err
				}
				for _, features := range featuresList.Items {
					r.addFeatures(configmap, &features, source.NamespaceMapping)
				}
			}
		} else {
			var featuresList fliptv1alpha1.FeaturesList
			if err := r.List(ctx, &featuresList, &listOpts); err != nil {
				return ctrl.Result{}, err
			}
			for _, features := range featuresList.Items {
				r.addFeatures(configmap, &features, source.NamespaceMapping)
			}
		}
	}

	existing := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKeyFromObject(configmap), existing)
	if err != nil && !errors.IsNotFound(err) {
		return ctrl.Result{}, err
	} else if errors.IsNotFound(err) {
		r.Create(ctx, configmap)
	} else {
		for key := range existing.Data {
			if _, exists := configmap.Data[key]; !exists {
				delete(existing.Data, key)
			}
		}
		if existing.Data == nil {
			existing.Data = make(map[string]string)
		}
		for key, value := range configmap.Data {
			// New value invalid, don't update
			if value == "" {
				continue
			}
			existing.Data[key] = value
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FeaturesTargetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fliptv1alpha1.FeaturesTarget{}).
		Owns(&corev1.ConfigMap{}).
		Watches(&fliptv1alpha1.Features{}, handler.EnqueueRequestsFromMapFunc(r.findTargetsForFeatures), builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Complete(r)
}

func (r *FeaturesTargetReconciler) addFeatures(cm *corev1.ConfigMap, features *fliptv1alpha1.Features, NamespaceMapping fliptv1alpha1.FeaturesTargetSpecNamespaceMapping) {
	name := fmt.Sprintf("%s/%s/features.yml", features.Namespace, features.Name)
	// Prevent deletion of the old value in case of error
	cm.Data[name] = ""

	namespaceValid := true
	fsf := features.Spec.Features
	switch NamespaceMapping {
	case "":
	case fliptv1alpha1.PreserveNamespace:
	case fliptv1alpha1.MustMatchNamespace:
		if fsf.Namespace == nil {
			namespaceValid = false
		} else if fsf.Namespace.Name != features.GetNamespace() {
			namespaceValid = false
		}
	case fliptv1alpha1.OverrideNamespace:
		if fsf.Namespace == nil {
			fsf.Namespace = &fliptv1alpha1.Namespace{}
		}
		fsf.Namespace.Name = features.GetNamespace()
	case fliptv1alpha1.RequireNamespace:
		if fsf.Namespace == nil || fsf.Namespace.Name == "" {
			namespaceValid = false
		}
	}
	if !namespaceValid {
		return
	}

	data, err := yaml.Marshal(features.Spec.Features)
	if err != nil {
		return
	}
	cm.Data[name] = string(data)
}

func (r *FeaturesTargetReconciler) findTargetsForFeatures(ctx context.Context, features client.Object) []reconcile.Request {
	var targets fliptv1alpha1.FeaturesTargetList
	if err := r.List(ctx, &targets); err != nil {
		logger := log.FromContext(ctx, "Features", types.NamespacedName{Namespace: features.GetNamespace(), Name: features.GetName()})
		logger.Error(err, "failed to list FeaturesTarget objects")
		return nil
	}

	var requests []reconcile.Request

Targets:
	for _, target := range targets.Items {
		for _, source := range target.Spec.Sources {
			if len(source.Namespaces) > 0 || !slices.Contains(source.Namespaces, features.GetNamespace()) {
				continue
			}
			sel, _ := metav1.LabelSelectorAsSelector(source.Selector)
			if sel.Matches(labels.Set(features.GetLabels())) {
				requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Namespace: target.GetNamespace(), Name: target.GetName()}})
				continue Targets
			}
		}
	}

	return requests
}

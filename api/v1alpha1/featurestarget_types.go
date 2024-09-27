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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FeaturesTargetSpec defines the desired state of FeaturesTarget
type FeaturesTargetSpec struct {
	// Sepcifies the ConfigMap resource that is updated with the compiled features configuration
	ConfigMap *FeaturesTargetSpecConfigMap `json:"configmap,omitempty"`

	Sources []FeaturesTargetSpecSource `json:"sources,omitempty"`
}

type FeaturesTargetSpecConfigMap struct {
	// Name of the ConfigMap to update. Required.
	//+kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
}

type FeaturesTargetSpecSource struct {
	// Namespaces in which to look for feature sources. If empty, all namespaces are searched.
	Namespaces []string `json:"namespaces,omitempty"`
	// Label selector to filter which feature sources to consider.
	Selector metav1.LabelSelector `json:"selector,omitempty"`
	// How to handle the namespace set in the source features config.
	//  - override: The features namespace is replaced with the kubernetes namespace of the source
	//  - mustmatch: The features namespace must match the kubernetes namespace of the source
	//  - require: The features namespace must be non-empty
	//  - preserve: Leave the feature namespace as-is
	//+kubebuilder:validation:Enum=override;mustmatch;require;preserve
	NamespaceMapping string `json:"namespaceMapping,omitempty"`
}

// FeaturesTargetStatus defines the observed state of FeaturesTarget
type FeaturesTargetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FeaturesTarget is the Schema for the featurestargets API
type FeaturesTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeaturesTargetSpec   `json:"spec,omitempty"`
	Status FeaturesTargetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FeaturesTargetList contains a list of FeaturesTarget
type FeaturesTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FeaturesTarget `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FeaturesTarget{}, &FeaturesTargetList{})
}

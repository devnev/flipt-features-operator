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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FeaturesTargetSpec defines the desired state of FeaturesTarget
type FeaturesTargetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of FeaturesTarget. Edit featurestarget_types.go to remove/update
	Foo string `json:"foo,omitempty"`
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

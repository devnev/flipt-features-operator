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

// FeaturesSpec defines the desired state of Features
type FeaturesSpec struct {
	Features Document `json:"features,omitempty"`
}

// FeaturesStatus defines the observed state of Features
type FeaturesStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Features is the Schema for the features API
type Features struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeaturesSpec   `json:"spec,omitempty"`
	Status FeaturesStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FeaturesList contains a list of Features
type FeaturesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Features `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Features{}, &FeaturesList{})
}

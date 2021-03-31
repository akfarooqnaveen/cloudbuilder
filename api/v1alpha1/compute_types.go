/*


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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ComputeSpec defines the desired state of Compute
type ComputeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Compute. Edit Compute_types.go to remove/update
	CloudProviderName string `json:"cloudprovidername"`
	ComputeName       string `json:"computename"`
	OSImage           string `json:"osimage"`
	Shape             string `json:"shape"`
	Region            string `json:"region"`
	Zone              string `json:"zone"`
	Network           string `json:"network"`
	Subnet            string `json:"subnet"`
}

// ComputeStatus defines the observed state of Compute
type ComputeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Compute is the Schema for the computes API
type Compute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComputeSpec   `json:"spec,omitempty"`
	Status ComputeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ComputeList contains a list of Compute
type ComputeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Compute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Compute{}, &ComputeList{})
}

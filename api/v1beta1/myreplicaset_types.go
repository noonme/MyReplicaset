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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyReplicaSetSpec defines the desired state of MyReplicaSet
type MyReplicaSetSpec struct {
	// Replicas is the number of desired replicas
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
	// Template is the object that describes the pod that will be created for the replicas
	// +optional
	Template v1.PodTemplateSpec `json:"template,omitempty"`
}

// MyReplicaSetStatus defines the observed state of MyReplicaSet
type MyReplicaSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Replicas int32 `json:"replicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MyReplicaSet is the Schema for the myreplicasets API
type MyReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyReplicaSetSpec   `json:"spec,omitempty"`
	Status MyReplicaSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MyReplicaSetList contains a list of MyReplicaSet
type MyReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyReplicaSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyReplicaSet{}, &MyReplicaSetList{})
}

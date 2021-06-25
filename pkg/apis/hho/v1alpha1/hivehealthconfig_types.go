package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HiveHealthConfigSpec defines the desired state of HiveHealthConfig
type HiveHealthConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// HiveHealthConfigStatus defines the observed state of HiveHealthConfig
type HiveHealthConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HiveHealthConfig is the Schema for the hivehealthconfigs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=hivehealthconfigs,scope=Namespaced
type HiveHealthConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HiveHealthConfigSpec   `json:"spec,omitempty"`
	Status HiveHealthConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HiveHealthConfigList contains a list of HiveHealthConfig
type HiveHealthConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HiveHealthConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HiveHealthConfig{}, &HiveHealthConfigList{})
}

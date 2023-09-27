package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PolicySpec `json:"spec"`
}

type PolicySpec struct {
	// Rules is a list of Rule instances. A Policy contains multiple rules and
	// each rule can validate, mutate, or generate resources.
	Rules []Rule `json:"rules,omitempty"`
}

type Rule struct {
	// Name is a label to identify the rule, It must be unique within the policy.
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`

	// MatchResources defines when this policy rule should be applied. The match
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the user name or role.
	// At least one kind is required.
	MatchResources MatchResources `json:"match,omitempty"`

	// ExcludeResources defines when this policy rule should not be applied. The exclude
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the name or role.
	// +optional
	ExcludeResources MatchResources `json:"exclude,omitempty"`

	// Validation is used to validate matching resources.
	// +optional
	Validation *Validation `json:"validate,omitempty"`
}

type MatchResources struct {
	// Any allows specifying resources which will be ORed
	// +optional
	Any ResourceFilters `json:"any,omitempty"`

	// All allows specifying resources which will be ANDed
	// +optional
	All ResourceFilters `json:"all,omitempty"`
}

// ResourceFilters is a slice of ResourceFilter
type ResourceFilters []ResourceFilter

// ResourceFilter allow users to "AND" or "OR" between resources
type ResourceFilter struct {
	// ResourceDescription contains information about the resource being created or modified.
	// +kubebuilder:validation:Type=object
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Resource map[string]interface{} `json:"resources,omitempty"`
}

// Validation defines checks to be performed on matching resources.
type Validation struct {
	// Message specifies a custom message to be displayed on failure.
	// +optional
	Message string `json:"message,omitempty"`

	// Pattern specifies an overlay-style pattern used to check resources.
	// +kubebuilder:validation:Type=object
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Pattern map[string]interface{} `json:"pattern,omitempty"`
}

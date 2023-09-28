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

	// Context defines variables and data sources that can be used during rule execution.
	// +optional
	Context []ContextEntry `json:"context,omitempty"`

	// MatchResources defines when this policy rule should be applied. The match
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the user name or role.
	// At least one kind is required.
	MatchResources *MatchResources `json:"match,omitempty"`

	// ExcludeResources defines when this policy rule should not be applied. The exclude
	// criteria can include resource information (e.g. kind, name, namespace, labels)
	// and admission review request information like the name or role.
	// +optional
	ExcludeResources *MatchResources `json:"exclude,omitempty"`

	// Validation is used to validate matching resources.
	// +optional
	Validation *Validation `json:"validate,omitempty"`
}

// ContextEntry adds variables and data sources to a rule Context. Either a
// ConfigMap reference or a APILookup must be provided.
type ContextEntry struct {
	// Name is the variable name.
	Name string `json:"name"`

	// Variable defines an arbitrary JMESPath context variable that can be defined inline.
	Variable *Variable `json:"variable,omitempty"`
}

// Variable defines an arbitrary JMESPath context variable that can be defined inline.
type Variable struct {
	// Value is any arbitrary JSON object representable in YAML or JSON form.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value interface{} `json:"value,omitempty"`

	// JMESPath is an optional JMESPath Expression that can be used to transform the variable.
	JMESPath string `json:"jmesPath,omitempty"`

	// Default is an optional arbitrary JSON object that the variable may take if the JMESPath expression evaluates to nil
	// +kubebuilder:validation:Type=object
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Default interface{} `json:"default,omitempty"`
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
	Resource map[string]interface{} `json:"resource,omitempty"`
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

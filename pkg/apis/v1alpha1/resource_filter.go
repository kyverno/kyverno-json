package v1alpha1

// ResourceFilter allow users to "AND" or "OR" between resources
type ResourceFilter struct {
	// ResourceDescription contains information about the resource being created or modified.
	// +kubebuilder:validation:Type=object
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Resource map[string]interface{} `json:"resource,omitempty"`
}

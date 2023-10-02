package v1alpha1

// Validation defines checks to be performed on matching resources.
type Validation struct {
	// Message specifies a custom message to be displayed on failure.
	Message string `json:"message,omitempty"`

	// Pattern specifies an overlay-style pattern used to check resources.
	// +kubebuilder:validation:Type=object
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Pattern map[string]interface{} `json:"pattern,omitempty"`
}

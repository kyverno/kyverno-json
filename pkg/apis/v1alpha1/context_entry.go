package v1alpha1

// ContextEntry adds variables and data sources to a rule context.
type ContextEntry struct {
	// Name is the entry name.
	Name string `json:"name"`

	// Variable defines an arbitrary variable.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +optional
	Variable Any `json:"variable,omitempty"`
}

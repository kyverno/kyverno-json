package v1alpha1

// Variable defines an arbitrary JMESPath context variable that can be defined inline.
type Variable struct {
	// Value is any arbitrary object.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value Any `json:"value,omitempty"`
}

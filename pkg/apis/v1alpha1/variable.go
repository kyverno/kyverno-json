package v1alpha1

// Variable defines an arbitrary JMESPath context variable that can be defined inline.
type Variable struct {
	// Value is any arbitrary JSON object representable in YAML or JSON form.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value interface{} `json:"value,omitempty"`
}

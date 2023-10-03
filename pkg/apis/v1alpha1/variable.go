package v1alpha1

// Variable defines an arbitrary JMESPath context variable that can be defined inline.
// +k8s:deepcopy-gen=false
type Variable struct {
	// Value is any arbitrary object.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value interface{} `json:"value,omitempty"`
}

func (in *Variable) DeepCopy() *Variable {
	// TODO
	return nil
}

package v1alpha1

import (
	"github.com/jinzhu/copier"
)

// Variable defines an arbitrary JMESPath context variable that can be defined inline.
// +k8s:deepcopy-gen=false
type Variable struct {
	// Value is any arbitrary object.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value interface{} `json:"value,omitempty"`
}

func (in *Variable) DeepCopy() *Variable {
	if in == nil {
		return nil
	}
	out := &Variable{}
	if err := copier.Copy(out, in); err != nil {
		panic("deep copy failed")
	}
	return out
}

package v1alpha1

import (
	"github.com/jinzhu/copier"
)

type Value = interface{}

// +k8s:deepcopy-gen=false
type Any struct {
	// TODO: this is needed to workaround a bug in api machinery code
	// https://kubernetes.slack.com/archives/C0EG7JC6T/p1696331287543159
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value `json:",inline"`
}

func (in *Any) DeepCopyInto(out *Any) {
	if err := copier.Copy(out, in); err != nil {
		panic("deep copy failed")
	}
}

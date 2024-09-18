package v1alpha1

import "k8s.io/apimachinery/pkg/util/json"

// +kubebuilder:validation:Type:=string
type Message struct {
	_template string
}

func (a *Message) Template() string {
	return a._template
}

func (a *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(a._template)
}

func (a *Message) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	a._template = v
	return nil
}

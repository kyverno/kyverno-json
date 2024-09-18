package v1alpha1

import (
	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/util/json"
)

// AssertionTree represents an assertion tree.
// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
type AssertionTree struct {
	tree any `json:"-"`
}

func (t *AssertionTree) Raw() any {
	return t.tree
}

func (a *AssertionTree) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.tree)
}

func (a *AssertionTree) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	a.tree = v
	return nil
}

func (in *AssertionTree) DeepCopyInto(out *AssertionTree) {
	if err := copier.CopyWithOption(out, in, copier.Option{DeepCopy: true}); err != nil {
		panic("deep copy failed")
	}
}

// Assertion contains an assertion tree associated with a message.
type Assertion struct {
	// Message is the message associated message.
	// +optional
	Message string `json:"message,omitempty"`

	// Check is the assertion check definition.
	Check AssertionTree `json:"check"`
}

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/util/json"
)

// +k8s:deepcopy-gen=false
// +kubebuilder:validation:XPreserveUnknownFields
// +kubebuilder:validation:Type:=""
// AssertionTree represents an assertion tree.
type AssertionTree struct {
	// +optional
	tree any `json:"-"`
}

func NewAssertionTree(value any) AssertionTree {
	return AssertionTree{value}
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
	out.tree = deepCopy(in.tree)
}

// Assertion contains an assertion tree associated with a message.
type Assertion struct {
	// Message is the message associated message.
	// +optional
	Message string `json:"message,omitempty"`

	// Check is the assertion check definition.
	Check AssertionTree `json:"check"`
}

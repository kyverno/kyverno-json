package v1alpha1

type Assert struct {
	// Any allows specifying resources which will be ORed.
	Any []Assertion `json:"any,omitempty"`

	// All allows specifying resources which will be ANDed.
	All []Assertion `json:"all,omitempty"`
}

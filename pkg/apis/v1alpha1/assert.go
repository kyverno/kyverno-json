package v1alpha1

// Assert defines collections of assertions.
type Assert struct {
	// Any allows specifying assertions which will be ORed.
	// +optional
	Any []Assertion `json:"any,omitempty"`

	// All allows specifying assertions which will be ANDed.
	// +optional
	All []Assertion `json:"all,omitempty"`
}

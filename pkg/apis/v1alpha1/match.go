package v1alpha1

// Match defines collections of assertion trees.
type Match struct {
	// Any allows specifying assertion trees which will be ORed.
	// +optional
	Any []Any `json:"any,omitempty"`

	// All allows specifying assertion trees which will be ANDed.
	// +optional
	All []Any `json:"all,omitempty"`
}

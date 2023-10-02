package v1alpha1

type MatchResources struct {
	// Any allows specifying resources which will be ORed.
	// +optional
	Any ResourceFilters `json:"any,omitempty"`

	// All allows specifying resources which will be ANDed.
	// +optional
	All ResourceFilters `json:"all,omitempty"`
}

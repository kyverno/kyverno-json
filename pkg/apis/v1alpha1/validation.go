package v1alpha1

// Validation defines checks to be performed on matching resources.
type Validation struct {
	// Assert specifies an overlay-style pattern used to check resources.
	Assert *Assert `json:"assert,omitempty"`
}

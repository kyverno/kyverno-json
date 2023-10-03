package v1alpha1

// Validation defines checks to be performed on matching resources.
type Validation struct {
	// Message specifies a custom message to be displayed on failure.
	Message string `json:"message,omitempty"`

	// Assert specifies an overlay-style pattern used to check resources.
	Assert *Match `json:"assert,omitempty"`
}

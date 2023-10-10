package v1alpha1

type Assertion struct {
	// Message is the variable associated message.
	Message string `json:"message,omitempty"`

	// Check is the assertion check definition.
	Check Any `json:"check"`
}

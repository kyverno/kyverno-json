package v1alpha1

// Assertion contains an assertion tree associated with a message.
type Assertion struct {
	// Message is the message associated message.
	// +optional
	Message string `json:"message,omitempty"`

	// Engine defines the default engine to use when evaluating expressions.
	Engine *Engine `json:"engine,omitempty"`

	// Check is the assertion check definition.
	Check AssertionTree `json:"check"`
}

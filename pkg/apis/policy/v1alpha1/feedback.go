package v1alpha1

// Feedback contains a feedback entry.
type Feedback struct {
	// Name is the feedback entry name.
	Name string `json:"name"`

	// Value is the feedback entry value (a JMESPath expression).
	Value string `json:"value"`
}

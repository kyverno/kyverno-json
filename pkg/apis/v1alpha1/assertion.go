package v1alpha1

type Value = interface{}

type Assertion struct {
	// TODO: this is needed to workaround a bug in api machinery code
	// https://kubernetes.slack.com/archives/C0EG7JC6T/p1696331287543159
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Value `json:",inline"`
}

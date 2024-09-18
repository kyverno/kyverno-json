package v1alpha1

// Engine defines the engine to use when evaluating expressions.
// +kubebuilder:validation:Enum:=jp;cel
type Engine string

const (
	EngineJP  Engine = "jp"
	EngineCEL Engine = "cel"
)

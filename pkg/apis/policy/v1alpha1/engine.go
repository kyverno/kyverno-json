package v1alpha1

// +kubebuilder:validation:Enum:=jp;cel
type Engine string

const (
	EngineJP  Engine = "jp"
	EngineCEL Engine = "cel"
)

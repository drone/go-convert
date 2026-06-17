package yaml

type StepUpload struct {
	Inputs map[string]*Input `json:"inputs,omitempty" yaml:"inputs,omitempty"`
}

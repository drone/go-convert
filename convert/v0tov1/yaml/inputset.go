package yaml

type (
	// InputSet defines a v1 input set configuration.
	InputSet struct {
		Overlay *Pipeline `json:"overlay,omitempty" yaml:"overlay,omitempty"`
	}

	// InputSetConfig is the root wrapper for v1 input set YAML.
	InputSetConfig struct {
		Inputs *InputSet `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	}
)

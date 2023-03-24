package yaml

type WorkflowDispatchCondition struct {
	Inputs map[string]InputDefinition `yaml:"inputs,omitempty"`
}

type InputDefinition struct {
	Description string      `yaml:"description,omitempty"`
	Required    bool        `yaml:"required,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
	Type        string      `yaml:"type,omitempty"`
	Options     interface{} `yaml:"options,omitempty"`
}

type Inputs struct {
	LogLevel    string `yaml:"logLevel"`
	PrintTags   bool   `yaml:"print_tags"`
	Tags        string `yaml:"tags"`
	Environment string `yaml:"environment"`
}

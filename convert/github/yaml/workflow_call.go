package yaml

type WorkflowCallCondition struct {
	Workflows []string                   `yaml:"workflows,omitempty"`
	Inputs    map[string]interface{}     `yaml:"inputs,omitempty"`
	Outputs   map[string]interface{}     `yaml:"outputs,omitempty"`
	Secrets   map[string]WorkflowSecrets `yaml:"secrets,omitempty"`
}

type WorkflowSecrets struct {
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
}

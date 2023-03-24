package yaml

type WorkflowRunCondition struct {
	Workflows      []string `yaml:"workflows,omitempty"`
	Types          []string `yaml:"types,omitempty"`
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
}

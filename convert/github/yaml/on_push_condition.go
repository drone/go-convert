package yaml

type PushCondition struct {
	Branches       []string `yaml:"branches,omitempty"`
	BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	Paths          []string `yaml:"paths,omitempty"`
	PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
	Tags           []string `yaml:"tags,omitempty"`
	TagsIgnore     []string `yaml:"tags-ignore,omitempty"`
}

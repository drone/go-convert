package yaml

import "fmt"

type Conditions struct {
	Refs       []string `yaml:"refs,omitempty"`
	Variables  []string `yaml:"variables,omitempty"`
	Changes    []string `yaml:"changes,omitempty"`
	Kubernetes string   `yaml:"kubernetes,omitempty"`
}

func (c *Conditions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 []string
	var out2 = struct {
		Refs       []string `yaml:"refs"`
		Variables  []string `yaml:"variables"`
		Changes    []string `yaml:"changes"`
		Kubernetes string   `yaml:"kubernetes"`
	}{}

	// Try to unmarshal to out2 struct
	if err := unmarshal(&out2); err == nil {
		c.Refs = out2.Refs
		c.Variables = out2.Variables
		c.Changes = out2.Changes
		c.Kubernetes = out2.Kubernetes
		return nil
	}

	// Try to unmarshal to out1 slice of strings
	if err := unmarshal(&out1); err == nil {
		c.Refs = out1
		return nil
	}

	return fmt.Errorf("failed to unmarshal conditions")
}

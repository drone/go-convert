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
	var out2 string
	var out3 = struct {
		Refs       []string `yaml:"refs"`
		Variables  []string `yaml:"variables"`
		Changes    []string `yaml:"changes"`
		Kubernetes string   `yaml:"kubernetes"`
	}{}

	// Try to unmarshal to out3 struct
	if err := unmarshal(&out3); err == nil {
		c.Refs = out3.Refs
		c.Variables = out3.Variables
		c.Changes = out3.Changes
		c.Kubernetes = out3.Kubernetes
		return nil
	}

	// Try to unmarshal to out1 slice of strings
	if err := unmarshal(&out1); err == nil {
		c.Refs = out1
		return nil
	}

	// Try to unmarshal to out2 single string
	if err := unmarshal(&out2); err == nil {
		c.Refs = []string{out2}
		return nil
	}

	return fmt.Errorf("failed to unmarshal conditions")
}

func (c *Conditions) MarshalYAML() (interface{}, error) {
	if len(c.Refs) > 0 && len(c.Variables) == 0 && len(c.Changes) == 0 && c.Kubernetes == "" {
		return c.Refs, nil
	} else {
		output := make(map[string]interface{})
		if len(c.Refs) > 0 {
			output["refs"] = c.Refs
		}
		if len(c.Variables) > 0 {
			output["variables"] = c.Variables
		}
		if len(c.Changes) > 0 {
			output["changes"] = c.Changes
		}
		if c.Kubernetes != "" {
			output["kubernetes"] = c.Kubernetes
		}
		return output, nil
	}
}

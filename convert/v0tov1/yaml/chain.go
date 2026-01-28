package yaml

import "github.com/drone/go-convert/internal/flexible"

type (
	Chain struct {
		Uses string     `json:"uses,omitempty" yaml:"uses,omitempty"`
		With *ChainWith `json:"with,omitempty" yaml:"with,omitempty"`
	}

	ChainWith struct {
		InputSets *flexible.Field[[]string]         `json:"input-sets,omitempty" yaml:"input-sets,omitempty"`
		Outputs   []*ChainOutput         `json:"outputs,omitempty" yaml:"outputs,omitempty"`
		Inputs    map[string]interface{} `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	}

	ChainOutput struct {
		Name  string `json:"name,omitempty" yaml:"name,omitempty"`
		Value string `json:"value,omitempty" yaml:"value,omitempty"`
	}
)

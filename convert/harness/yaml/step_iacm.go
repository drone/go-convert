package yaml

import (
	"github.com/drone/go-convert/internal/flexible"
)

type (
	StepIACMTerraformPlugin struct {
		CommonStepSpec
		Command string `json:"command,omitempty" yaml:"command,omitempty"`
		Target *flexible.Field[[]string] `json:"target,omitempty" yaml:"target,omitempty"`
		Replace *flexible.Field[[]string] `json:"replace,omitempty" yaml:"replace,omitempty"`
		ImportVars *flexible.Field[[]IACMImportVar] `json:"importVars,omitempty" yaml:"importVars,omitempty"`
	}

	StepIACMOpenTofuPlugin struct {
		CommonStepSpec
		Command string `json:"command,omitempty" yaml:"command,omitempty"`
		Target *flexible.Field[[]string] `json:"target,omitempty" yaml:"target,omitempty"`
		Replace *flexible.Field[[]string] `json:"replace,omitempty" yaml:"replace,omitempty"`
		ImportVars *flexible.Field[[]IACMImportVar] `json:"importVars,omitempty" yaml:"importVars,omitempty"`
	}

	IACMImportVar struct {
		Id string `json:"id,omitempty" yaml:"id,omitempty"`
		Address string `json:"address,omitempty" yaml:"address,omitempty"`
	}
)
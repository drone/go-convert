package yaml

import "github.com/drone/go-convert/internal/flexible"

type (
	// EnvironmentGroup defines the environment group
	EnvironmentGroup struct {
		EnvGroupRef  string                                `json:"envGroupRef,omitempty"  yaml:"envGroupRef,omitempty"`
		DeployToAll  *flexible.Field[bool]                 `json:"deployToAll,omitempty"  yaml:"deployToAll,omitempty"`
		Environments *flexible.Field[[]*Environment]       `json:"environments,omitempty" yaml:"environments,omitempty"`
		Filters      *flexible.Field[[]*EnvironmentFilter] `json:"filters,omitempty"      yaml:"filters,omitempty"`
		Metadata     *EnvironmentMetadata                  `json:"metadata,omitempty"     yaml:"metadata,omitempty"`
		UseFromStage *UseFromStage                         `json:"useFromStage,omitempty" yaml:"useFromStage,omitempty"`
	}

	Environment struct {
		EnvironmentRef            string                                       `json:"environmentRef,omitempty"            yaml:"environmentRef,omitempty"`
		DeployToAll               *flexible.Field[bool]                        `json:"deployToAll,omitempty"               yaml:"deployToAll,omitempty"`
		EnvironmentInputs         interface{}                                  `json:"environmentInputs,omitempty"        yaml:"environmentInputs,omitempty"`
		GitBranch                 string                                       `json:"gitBranch,omitempty"                yaml:"gitBranch,omitempty"`
		InfrastructureDefinition  *flexible.Field[InfrastructureDefinition]    `json:"infrastructureDefinition,omitempty"  yaml:"infrastructureDefinition,omitempty"`
		InfrastructureDefinitions *flexible.Field[[]*InfrastructureDefinition] `json:"infrastructureDefinitions,omitempty" yaml:"infrastructureDefinitions,omitempty"`
		Filters                   *flexible.Field[[]*EnvironmentFilter]        `json:"filters,omitempty"                  yaml:"filters,omitempty"`
		ServiceOverrideInputs     interface{}                                  `json:"serviceOverrideInputs,omitempty"    yaml:"serviceOverrideInputs,omitempty"`
		UseFromStage              *UseFromStage                                `json:"useFromStage,omitempty"             yaml:"useFromStage,omitempty"`
	}

	Environments struct {
		Metadata     *EnvironmentMetadata                  `json:"metadata,omitempty"     yaml:"metadata,omitempty"`
		Values       *flexible.Field[[]*Environment]       `json:"values,omitempty"       yaml:"values,omitempty"`
		Filters      *flexible.Field[[]*EnvironmentFilter] `json:"filters,omitempty"      yaml:"filters,omitempty"`
		UseFromStage *UseFromStage                         `json:"useFromStage,omitempty" yaml:"useFromStage,omitempty"`
	}

	// EnvironmentMetadata defines environment metadata
	EnvironmentMetadata struct {
		Parallel *flexible.Field[bool] `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	InfrastructureDefinition struct {
		Identifier string      `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		Inputs     interface{} `json:"inputs,omitempty"     yaml:"inputs,omitempty"`
		Metadata   string      `json:"metadata,omitempty"   yaml:"metadata,omitempty"`
	}

	EnvironmentFilter struct {
		Identifier string                    `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		Type       string                    `json:"type,omitempty"       yaml:"type,omitempty"`
		Entities   *flexible.Field[[]string] `json:"entities,omitempty"   yaml:"entities,omitempty"`
		Spec       *EnvironmentFilterSpec    `json:"spec,omitempty"       yaml:"spec,omitempty"`
	}

	EnvironmentFilterSpec struct {
		Tags      *flexible.Field[map[string]string] `json:"tags,omitempty"      yaml:"tags,omitempty"`
		MatchType string                             `json:"matchType,omitempty" yaml:"matchType,omitempty"`
	}
)
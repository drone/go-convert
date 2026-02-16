package yaml

import "github.com/drone/go-convert/internal/flexible"

type (
	// EnvironmentGroup defines the environment group
	EnvironmentGroup struct {
		EnvGroupRef  string               `json:"envGroupRef,omitempty"  yaml:"envGroupRef,omitempty"`
		Metadata     *EnvironmentMetadata `json:"metadata,omitempty"     yaml:"metadata,omitempty"`
		DeployToAll  interface{}          `json:"deployToAll,omitempty"  yaml:"deployToAll,omitempty"`
		Environments interface{}          `json:"environments,omitempty" yaml:"environments,omitempty"`
	}

	Environment struct {
		EnvironmentRef string `json:"environmentRef,omitempty" yaml:"environmentRef,omitempty"`
		DeployToAll    bool   `json:"deployToAll,omitempty"    yaml:"deployToAll,omitempty"`
		InfrastructureDefinitions flexible.Field[[]*InfrastructureDefinition] `json:"infrastructureDefinitions,omitempty" yaml:"infrastructureDefinitions,omitempty"`
	}

	Environments struct {
		Metadata *EnvironmentMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
		Values []*Environment `json:"values,omitempty" yaml:"values,omitempty"`
	}

	// EnvironmentMetadata defines environment metadata
	EnvironmentMetadata struct {
		Parallel *flexible.Field[bool] `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	InfrastructureDefinition struct {
		Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	}
)
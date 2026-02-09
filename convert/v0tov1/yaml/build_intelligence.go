package yaml

import "github.com/drone/go-convert/internal/flexible"

// BuildIntelligence defines pipeline build intelligence behavior.
type BuildIntelligence struct {
	Enabled *flexible.Field[bool] `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Port    string                `json:"port,omitempty" yaml:"port,omitempty"`
	MavenUrl string                `json:"maven-url,omitempty" yaml:"maven-url,omitempty"`	
}

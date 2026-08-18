package yaml

import "github.com/drone/go-convert/internal/flexible"

// BuildIntelligence defines pipeline build intelligence behavior.
type BuildIntelligence struct {
	Enabled        *flexible.Field[bool] `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Port           string                `json:"port,omitempty" yaml:"port,omitempty"`
	MavenUrl       string                `json:"maven-url,omitempty" yaml:"maven-url,omitempty"`
	Connector      string                `json:"connector,omitempty" yaml:"connector,omitempty"`
	Region         string                `json:"region,omitempty" yaml:"region,omitempty"`
	BucketName     string                `json:"bucket_name,omitempty" yaml:"bucket_name,omitempty"`
	ContainerName  string                `json:"container_name,omitempty" yaml:"container_name,omitempty"`
	StorageAccount string                `json:"storage_account,omitempty" yaml:"storage_account,omitempty"`
	User           *flexible.Field[int]  `json:"user,omitempty" yaml:"user,omitempty"`
	Resources      *ContainerResources   `json:"resources,omitempty" yaml:"resources,omitempty"`
}

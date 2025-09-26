package yaml
// TODO: handle serviceInputs and use flexiblefields where applicable
type (
	DeploymentServices struct {
		Values   []*DeploymentService `json:"values,omitempty" yaml:"values,omitempty"`
		Metadata *ServicesMetadata    `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	}

	// DeploymentService defines the service configuration for deployment
	DeploymentService struct {
		ServiceRef    string                       `json:"serviceRef,omitempty"    yaml:"serviceRef,omitempty"`
		// ServiceInputs *FlexibleField[ServiceInputs] `json:"serviceInputs,omitempty" yaml:"serviceInputs,omitempty"`
	}

	// ServicesMetadata defines the services metadata
	ServicesMetadata struct {
		Parallel bool `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	// ServiceInputs defines the service inputs for deployment
	ServiceInputs struct {
		ServiceDefinition *ServiceDefinition `json:"serviceDefinition,omitempty" yaml:"serviceDefinition,omitempty"`
	}

	// ServiceDefinition defines the service definition
	ServiceDefinition struct {
		Type string                 `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *ServiceDefinitionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// ServiceDefinitionSpec defines the service definition specification
	ServiceDefinitionSpec struct {
		Manifests              []*ManifestConfig       `json:"manifests,omitempty"              yaml:"manifests,omitempty"`
		Artifacts              *ArtifactConfig         `json:"artifacts,omitempty"              yaml:"artifacts,omitempty"`
		ManifestConfigurations *ManifestConfigurations `json:"manifestConfigurations,omitempty" yaml:"manifestConfigurations,omitempty"`
		ConfigFiles            []*ConfigFile           `json:"configFiles,omitempty"            yaml:"configFiles,omitempty"`
	}

	// ManifestConfig defines a manifest configuration
	ManifestConfig struct {
		Manifest *Manifest `json:"manifest,omitempty" yaml:"manifest,omitempty"`
	}

	// Manifest defines a manifest
	Manifest struct {
		Identifier string        `json:"identifier,omitempty"   yaml:"identifier,omitempty"`
		Type       string        `json:"type,omitempty"         yaml:"type,omitempty"`
		Spec       *ManifestSpec `json:"spec,omitempty"         yaml:"spec,omitempty"`
	}

	// ManifestSpec defines the manifest specification
	ManifestSpec struct {
		Store        *Store `json:"store,omitempty"        yaml:"store,omitempty"`
		SubChartPath string `json:"subChartPath,omitempty" yaml:"subChartPath,omitempty"`
	}

	// Store defines a store configuration
	Store struct {
		Type string     `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *StoreSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// StoreSpec defines the store specification
	StoreSpec struct {
		FolderPath   string   `json:"folderPath,omitempty"    yaml:"folderPath,omitempty"`
		Branch       string   `json:"branch,omitempty"        yaml:"branch,omitempty"`
		ConnectorRef string   `json:"connectorRef,omitempty"  yaml:"connectorRef,omitempty"`
		RepoName     string   `json:"repoName,omitempty"      yaml:"repoName,omitempty"`
		Paths        []string `json:"paths,omitempty"       yaml:"paths,omitempty"`
		GitFetchType string   `json:"gitFetchType,omitempty"  yaml:"gitFetchType,omitempty"`
		CommitId     string   `json:"commitId,omitempty"      yaml:"commitId,omitempty"`
		ValuesPaths  []string `json:"valuesPaths,omitempty" yaml:"valuesPaths,omitempty"`
	}

	// ArtifactConfig defines artifact configuration
	ArtifactConfig struct {
		Primary *PrimaryArtifact `json:"primary,omitempty" yaml:"primary,omitempty"`
	}

	// PrimaryArtifact defines the primary artifact
	PrimaryArtifact struct {
		PrimaryArtifactRef string            `json:"primaryArtifactRef,omitempty" yaml:"primaryArtifactRef,omitempty"`
		Sources            []*ArtifactSource `json:"sources,omitempty"            yaml:"sources,omitempty"`
	}

	// ArtifactSource defines an artifact source
	ArtifactSource struct {
		Identifier string              `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		Type       string              `json:"type,omitempty"       yaml:"type,omitempty"`
		Spec       *ArtifactSourceSpec `json:"spec,omitempty"       yaml:"spec,omitempty"`
	}

	// ArtifactSourceSpec defines the artifact source specification
	ArtifactSourceSpec struct {
		ConnectorRef string `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		Tag          string `json:"tag,omitempty"          yaml:"tag,omitempty"`
	}

	// ManifestConfigurations defines manifest configurations
	ManifestConfigurations struct {
		PrimaryManifestRef string `json:"primaryManifestRef,omitempty" yaml:"primaryManifestRef,omitempty"`
	}

	// ConfigFile defines a configuration file
	ConfigFile struct {
		ConfigFile *ConfigFileSpec `json:"configFile,omitempty" yaml:"configFile,omitempty"`
	}

	// ConfigFileSpec defines the configuration file specification
	ConfigFileSpec struct {
		Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		Spec       *Store `json:"spec,omitempty"       yaml:"spec,omitempty"`
	}
)
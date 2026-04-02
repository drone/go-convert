package yaml

type (
	// InputSet defines a v0 input set configuration.
	InputSet struct {
		Name     string            `json:"name,omitempty"              yaml:"name,omitempty"`
		ID       string            `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Org      string            `json:"orgIdentifier,omitempty"     yaml:"orgIdentifier,omitempty"`
		Project  string            `json:"projectIdentifier,omitempty" yaml:"projectIdentifier,omitempty"`
		Pipeline *Pipeline         `json:"pipeline,omitempty"          yaml:"pipeline,omitempty"`
	}
)

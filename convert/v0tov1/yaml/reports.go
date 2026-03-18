package yaml

import "github.com/drone/go-convert/internal/flexible"

type Reports struct {
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	Paths *flexible.Field[[]string] `json:"paths,omitempty" yaml:"paths,omitempty"`
}
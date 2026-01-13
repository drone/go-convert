package yaml

import "github.com/drone/go-convert/internal/flexible"

type VolumeSecret struct {
	Name      string                `json:"name,omitempty"     yaml:"name,omitempty"`
	Optional  *flexible.Field[bool] `json:"optional,omitempty" yaml:"optional,omitempty"`
	MountPath string                `json:"mount-path,omitempty" yaml:"mount-path,omitempty"`
}
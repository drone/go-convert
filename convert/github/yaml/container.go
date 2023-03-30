package yaml

import "errors"

type Container struct {
	Image       string            `yaml:"image,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Options     string            `yaml:"options,omitempty"`
	Credentials *Credentials      `yaml:"credentials,omitempty"`
}

func (v *Container) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	if err := unmarshal(&out1); err == nil {
		v.Image = out1
		return nil
	}
	var out2 struct {
		Image       string            `yaml:"image,omitempty"`
		Env         map[string]string `yaml:"env,omitempty"`
		Ports       []string          `yaml:"ports,omitempty"`
		Volumes     []string          `yaml:"volumes,omitempty"`
		Options     string            `yaml:"options,omitempty"`
		Credentials *Credentials      `yaml:"credentials,omitempty"`
	}
	if err := unmarshal(&out2); err == nil {
		*v = out2
		return nil
	}
	return errors.New("failed to unmarshal container")
}

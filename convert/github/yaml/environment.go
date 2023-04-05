package yaml

import "errors"

type Environment struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url,omitempty"`
}

func (v *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}
	var out2 struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url,omitempty"`
	}
	if err := unmarshal(&out2); err == nil {
		*v = out2
		return nil
	}
	return errors.New("failed to unmarshal environment")
}

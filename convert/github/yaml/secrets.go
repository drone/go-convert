package yaml

import (
	"errors"
)

type Secrets struct {
	Inherit bool
	Values  map[string]string
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Secrets) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 map[string]string

	if err := unmarshal(&out1); err == nil {
		if out1 == "inherit" {
			v.Inherit = true
			return nil
		} else {
			return errors.New("invalid string value for secrets")
		}
	}
	if err := unmarshal(&out2); err == nil {
		v.Values = out2
		return nil
	}
	return errors.New("failed to unmarshal secrets")
}

// MarshalYAML implements the marshal interface.
func (v *Secrets) MarshalYAML() (interface{}, error) {
	if v.Inherit {
		return map[string]bool{"inherit": true}, nil
	}
	return v.Values, nil
}

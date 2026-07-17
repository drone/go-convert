package yaml

import "encoding/json"

type (
	// InputSet defines a v1 input set configuration.
	InputSet struct {
		Overlay   *Pipeline              `json:"-" yaml:"-"`
		Variables map[string]interface{} `json:"-" yaml:"-"`
	}

	// InputSetConfig is the root wrapper for v1 input set YAML.
	InputSetConfig struct {
		Inputs *InputSet `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	}
)

// MarshalJSON flattens the overlay and any lifted pipeline-level variables
// into a single map so they render as siblings under the top-level "inputs"
// key: the overlay under "overlay" and each variable as a scalar key/value.
func (i *InputSet) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	if i.Overlay != nil {
		m["overlay"] = i.Overlay
	}
	for k, v := range i.Variables {
		m[k] = v
	}
	return json.Marshal(m)
}

package yaml

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// MarshalPipeline marshals the given Pipeline into YAML with a top-level
// 'pipeline:' key, producing output in the form:
//
//	pipeline:
//	  id: ...
//	  name: ...
//	  stages: ...
//	  ...
//
// This matches the expected Harness v1 YAML shape.
func MarshalPipeline(p *Pipeline) ([]byte, error) {
	wrapper := struct {
		Pipeline *Pipeline `json:"pipeline"`
	}{
		Pipeline: p,
	}

	// First marshal to JSON
	jsonBytes, err := json.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	jsonData, err := jsonToInterface(jsonBytes)
	if err != nil {
		return nil, err
	}

	// Marshal to YAML
	return yaml.Marshal(jsonData)
}

// WritePipelineFile writes the Pipeline to the given file path in the same
// top-level 'pipeline:' YAML shape as MarshalPipeline.
func WritePipelineFile(path string, p *Pipeline) error {
	b, err := MarshalPipeline(p)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// MarshalInputSet marshals the given InputSet into YAML with a top-level
// 'inputs:' key, producing output in the form:
//
//	inputs:
//	  overlay:
//	    ...
func MarshalInputSet(i *InputSet) ([]byte, error) {
	wrapper := &InputSetConfig{
		Inputs: i,
	}

	// First marshal to JSON
	jsonBytes, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	jsonData, err := jsonToInterface(jsonBytes)
	if err != nil {
		return nil, err
	}

	// Marshal to YAML
	return yaml.Marshal(jsonData)
}

// WriteInputSetFile writes the InputSet to the given file path.
func WriteInputSetFile(path string, i *InputSet) error {
	b, err := MarshalInputSet(i)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// MarshalTemplate marshals the given Template into YAML with a top-level
//
//	template:
//	  step: ...
//	  or
//	  stage: ...
//	  or
//	  pipeline: ...
//
// This matches the expected Harness v1 YAML shape.
func MarshalTemplate(t *Template) ([]byte, error) {
	wrapper := struct {
		Template *Template `json:"template"`
	}{
		Template: t,
	}

	// First marshal to JSON
	jsonBytes, err := json.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	jsonData, err := jsonToInterface(jsonBytes)
	if err != nil {
		return nil, err
	}

	// Marshal to YAML
	return yaml.Marshal(jsonData)
}

// WriteTemplateFile writes the Template to the given file path.
func WriteTemplateFile(path string, t *Template) error {
	b, err := MarshalTemplate(t)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// MarshalTrigger marshals the given Trigger into YAML with a top-level
// 'trigger:' key, producing output in the form:
//
//	trigger:
//	  name: ...
//	  identifier: ...
//	  source: ...
//	  inputYaml: ...
//
// This matches the expected Harness v1 YAML shape.
func MarshalTrigger(t *Trigger) ([]byte, error) {
	wrapper := &TriggerConfig{
		Trigger: t,
	}

	// First marshal to JSON
	jsonBytes, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	jsonData, err := jsonToInterface(jsonBytes)
	if err != nil {
		return nil, err
	}

	// Marshal to YAML
	return yaml.Marshal(jsonData)
}

// jsonToInterface decodes JSON bytes into an interface{} tree, using
// json.Decoder with UseNumber() to preserve number representations.
// The resulting tree is then walked to convert json.Number values
// to int64 (for integers) or float64 (for floats), so that YAML
// marshaling outputs integers as e.g. 2147483647 instead of 2.147483647e+09.
func jsonToInterface(data []byte) (interface{}, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return convertJSONNumbers(v), nil
}

// convertJSONNumbers recursively walks an interface{} tree and converts
// json.Number values to int64 (if the number is a valid integer) or
// float64 (otherwise). This prevents large integers from being
// represented in scientific notation when marshaled to YAML.
// It also detects Go strings that look like YAML numbers (e.g. "345e2256")
// and wraps them so that the YAML emitter produces a quoted scalar.
func convertJSONNumbers(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			val[k] = convertJSONNumbers(child)
		}
		return val
	case []interface{}:
		for i, child := range val {
			val[i] = convertJSONNumbers(child)
		}
		return val
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return i
		}
		if f, err := val.Float64(); err == nil {
			return f
		}
		return val.String()
	case string:
		if yamlNumberLike.MatchString(val) {
			return quotedString(val)
		}
		// Multi-line strings (e.g. run scripts) should render as YAML literal
		// block scalars (|-) for readability. yaml.v3 refuses literal style and
		// falls back to a double-quoted flow scalar when any line has trailing
		// whitespace, so trim it per line first.
		if strings.Contains(val, "\n") {
			return literalString(trimTrailingSpacePerLine(val))
		}
		return val
	default:
		return v
	}
}

// yamlNumberLike matches strings that YAML parsers may resolve as floats
// in scientific notation (e.g. "345e2256", "1.5e10", ".5E-3").
// These must be quoted in YAML output to preserve their string type.
var yamlNumberLike = regexp.MustCompile(
	`^[-+]?([0-9]+(\.[0-9]*)?|\.[0-9]+)[eE][-+]?[0-9]+$`,
)

// quotedString is a string wrapper whose MarshalYAML method forces the
// YAML emitter to use double-quoted style, preventing the value from
// being misinterpreted as a number.
type quotedString string

func (q quotedString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: string(q),
		Style: yaml.DoubleQuotedStyle,
		Tag:   "!!str",
	}, nil
}

// literalString is a string wrapper whose MarshalYAML method forces the YAML
// emitter to use literal block style (|-) instead of a double-quoted flow
// scalar for multi-line values.
type literalString string

func (l literalString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: string(l),
		Style: yaml.LiteralStyle,
		Tag:   "!!str",
	}, nil
}

// trimTrailingSpacePerLine removes trailing spaces and tabs from each line of
// s. Trailing whitespace is meaningless for shell scripts but prevents yaml.v3
// from emitting literal block scalars, so stripping it lets multi-line values
// render as |- blocks.
func trimTrailingSpacePerLine(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

// WriteTriggerFile writes the Trigger to the given file path.
func WriteTriggerFile(path string, t *Trigger) error {
	b, err := MarshalTrigger(t)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

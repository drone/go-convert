package converter

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// v0InputSetDoc is used to parse a Harness v0 input set YAML document.
type v0InputSetDoc struct {
	InputSet v0InputSetFields `yaml:"inputSet"`
}

type v0InputSetFields struct {
	Name       string                 `yaml:"name"`
	Identifier string                 `yaml:"identifier"`
	Org        string                 `yaml:"orgIdentifier"`
	Project    string                 `yaml:"projectIdentifier"`
	Pipeline   map[string]interface{} `yaml:"pipeline"`
}

// v1InputSetOutput is the serialization shape for a v1 input set YAML document.
type v1InputSetOutput struct {
	Inputset v1InputSetFields `yaml:"inputset"`
}

type v1InputSetFields struct {
	Name     string                 `yaml:"name,omitempty"`
	Pipeline string                 `yaml:"pipeline,omitempty"`
	Inputs   map[string]interface{} `yaml:"inputs,omitempty"`
}

// InputSet converts a Harness v0 input set YAML string into v1 YAML bytes.
// The input must have a top-level "inputSet:" key.
// If refMapping is provided, template references in the output will be replaced.
//
// Conversion strategy:
//   - pipeline.identifier  → inputset.pipeline  (reference to the target pipeline)
//   - all other fields in the pipeline fragment are flattened into dotted key paths
//     under inputset.inputs (e.g. "properties.ci.codebase.build.type": "PR")
func InputSet(yamlStr string, refMapping map[string]string) ([]byte, error) {
	if err := validateTopLevelKey(yamlStr, "inputSet"); err != nil {
		return nil, err
	}

	var doc v0InputSetDoc
	if err := yaml.Unmarshal([]byte(yamlStr), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse v0 input set: %w", err)
	}
	if doc.InputSet.Pipeline == nil {
		return nil, fmt.Errorf("input set has no 'pipeline' section")
	}

	pipelineRef, _ := doc.InputSet.Pipeline["identifier"].(string)

	// Flatten the pipeline fragment (excluding identifier) into dot-path keys.
	inputs := make(map[string]interface{})
	for k, v := range doc.InputSet.Pipeline {
		if k == "identifier" {
			continue
		}
		flattenInto(k, v, inputs)
	}

	out := v1InputSetOutput{
		Inputset: v1InputSetFields{
			Name:     doc.InputSet.Name,
			Pipeline: pipelineRef,
		},
	}
	if len(inputs) > 0 {
		out.Inputset.Inputs = inputs
	}
	yamlBytes, err := yaml.Marshal(out)
	if err != nil {
		return nil, err
	}

	// Apply template reference replacements if mapping is provided
	return ReplaceTemplateRefs(yamlBytes, refMapping)
}

// flattenInto recursively flattens nested map values into dot-separated keys
// inside out. Slices and scalar values are stored directly at their path.
func flattenInto(prefix string, val interface{}, out map[string]interface{}) {
	m, ok := val.(map[string]interface{})
	if !ok {
		// Leaf value (string, int, bool, slice, etc.) — store directly.
		out[prefix] = val
		return
	}
	for k, child := range m {
		flattenInto(prefix+"."+k, child, out)
	}
}

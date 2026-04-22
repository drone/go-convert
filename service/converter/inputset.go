package converter

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// InputSet converts a Harness v0 input set YAML string into v1 YAML bytes.
// The input must have a top-level "inputSet:" key.
// If refMapping is provided, template references in the output will be replaced.
//
// Conversion strategy:
//   - The v0 inputSet.pipeline is converted to v1 inputs.overlay using the pipeline converter
func InputSet(yamlStr string, refMapping map[string]string) ([]byte, error) {
	if err := validateTopLevelKey(yamlStr, "inputSet"); err != nil {
		return nil, err
	}

	// Parse using the v0 Config struct which handles inputset parsing
	v0Config, err := v0.ParseString(yamlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse v0 input set: %w", err)
	}

	if v0Config.InputSet == nil {
		return nil, fmt.Errorf("input set parsing returned nil")
	}

	if v0Config.InputSet.Pipeline == nil {
		return nil, fmt.Errorf("input set has no 'pipeline' section")
	}

	// Use the pipeline converter to convert the input set
	c := pipelineconverter.NewPipelineConverter()
	v1InputSet := c.ConvertInputSet(v0Config.InputSet)
	if v1InputSet == nil {
		return nil, fmt.Errorf("input set conversion returned nil")
	}

	// Marshal using the v1 MarshalInputSet function
	yamlBytes, err := v1.MarshalInputSet(v1InputSet)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 input set: %w", err)
	}

	// Apply template reference replacements if mapping is provided
	return ReplaceTemplateRefs(yamlBytes, refMapping)
}

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
// contextPipelineYAML is an optional v0 pipeline YAML used purely as
// expression-postprocess context (see buildContextFromPipelineYAML). Pass ""
// to run postprocess without FQN context.
//
// Conversion strategy:
//   - The v0 inputSet.pipeline is converted to v1 inputs.overlay using the pipeline converter
func InputSet(yamlStr string, refMapping map[string]string, contextPipelineYAML string) (*Result, error) {
	if err := validateTopLevelKey(yamlStr, "inputSet"); err != nil {
		return nil, err
	}

	done := beginAPIConversion()
	defer done()

	v0Config, unknownFields, err := v0.ParseStringWithUnknownFields(yamlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse v0 input set: %w", err)
	}

	if v0Config.InputSet == nil {
		return nil, fmt.Errorf("input set parsing returned nil")
	}

	if v0Config.InputSet.Pipeline == nil {
		return nil, fmt.Errorf("input set has no 'pipeline' section")
	}

	c := pipelineconverter.NewPipelineConverter()
	v1InputSet := c.ConvertInputSet(v0Config.InputSet)
	if v1InputSet == nil {
		return nil, fmt.Errorf("input set conversion returned nil")
	}

	// Single-pass expression post-process. If the caller supplied a context
	// pipeline_yaml, derive a step-type map and walk in FQN mode; otherwise
	// fall back to nil context (no FQN).
	stepTypeMap, useFQN := buildContextFromPipelineYAML(contextPipelineYAML)
	pipelineconverter.PostProcessExpressions(v1InputSet, stepTypeMap, useFQN)

	yamlBytes, err := v1.MarshalInputSet(v1InputSet)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 input set: %w", err)
	}

	yamlBytes, err = ReplaceTemplateRefs(yamlBytes, refMapping)
	if err != nil {
		return nil, err
	}
	return &Result{
		YAML:          yamlBytes,
		UnknownFields: unknownFields,
		Summary:       buildAPISummary(unknownFields),
	}, nil
}

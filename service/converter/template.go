package converter

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// Template converts a Harness v0 template YAML string into v1 YAML bytes.
// Supported template types are Pipeline, Stage, Step, and StepGroup.
// The input must have a top-level "template:" key.
// If refMapping is provided, template references in the output will be replaced.
// contextPipelineYAML is an optional v0 pipeline YAML used purely as
// expression-postprocess context (see buildContextFromPipelineYAML). Pass ""
// to run postprocess without FQN context.
func Template(yamlStr string, refMapping map[string]string, contextPipelineYAML string) (*Result, error) {
	if err := validateTopLevelKey(yamlStr, "template"); err != nil {
		return nil, err
	}

	done := beginAPIConversion()
	defer done()

	v0Config, unknownFields, err := v0.ParseStringWithUnknownFields(yamlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse v0 template: %w", err)
	}

	if v0Config.Template == nil {
		return nil, fmt.Errorf("template parsing returned nil")
	}

	if v0Config.Template.Type == "" {
		return nil, fmt.Errorf("template 'type' field is required (Pipeline, Stage, Step, or StepGroup)")
	}

	if v0Config.Template.Spec == nil {
		return nil, fmt.Errorf("template 'spec' field is required")
	}

	c := pipelineconverter.NewPipelineConverter()
	v1Template := c.ConvertTemplate(v0Config.Template)
	if v1Template == nil {
		return nil, fmt.Errorf("template conversion returned nil")
	}

	// Single-pass expression post-process. If the caller supplied a context
	// pipeline_yaml, derive a step-type map and walk in FQN mode; otherwise
	// fall back to nil context (no FQN).
	stepTypeMap, useFQN := buildContextFromPipelineYAML(contextPipelineYAML)
	pipelineconverter.PostProcessExpressions(v1Template, stepTypeMap, useFQN)

	yamlBytes, err := v1.MarshalTemplate(v1Template)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 template: %w", err)
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

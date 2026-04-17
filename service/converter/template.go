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
func Template(yamlStr string, refMapping map[string]string) ([]byte, error) {
	if err := validateTopLevelKey(yamlStr, "template"); err != nil {
		return nil, err
	}

	// Parse using the v0 Config struct which handles template parsing
	v0Config, err := v0.ParseString(yamlStr)
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

	// Use the pipeline converter to convert the template
	c := pipelineconverter.NewPipelineConverter()
	v1Template := c.ConvertTemplate(v0Config.Template)
	if v1Template == nil {
		return nil, fmt.Errorf("template conversion returned nil")
	}

	// Marshal using the v1 MarshalTemplate function
	yamlBytes, err := v1.MarshalTemplate(v1Template)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 template: %w", err)
	}

	// Apply template reference replacements if mapping is provided
	return ReplaceTemplateRefs(yamlBytes, refMapping)
}

package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	"gopkg.in/yaml.v3"
)

// v0TemplateDoc is used to parse a Harness v0 template YAML document.
type v0TemplateDoc struct {
	Template v0TemplateFields `yaml:"template"`
}

type v0TemplateFields struct {
	Name       string                 `yaml:"name"`
	Identifier string                 `yaml:"identifier"`
	Org        string                 `yaml:"orgIdentifier"`
	Project    string                 `yaml:"projectIdentifier"`
	Type       string                 `yaml:"type"` // Pipeline | Stage | Step
	Spec       map[string]interface{} `yaml:"spec"`
}

// v1TemplateOutput is the serialization shape for a v1 template YAML document.
type v1TemplateOutput struct {
	Template v1TemplateFields `yaml:"template"`
}

type v1TemplateFields struct {
	Name       string      `yaml:"name,omitempty"`
	Identifier string      `yaml:"identifier,omitempty"`
	Org        string      `yaml:"orgIdentifier,omitempty"`
	Project    string      `yaml:"projectIdentifier,omitempty"`
	Type       string      `yaml:"type,omitempty"`
	Spec       interface{} `yaml:"spec,omitempty"`
}

// Template converts a Harness v0 template YAML string into v1 YAML bytes.
// Supported template types are Pipeline, Stage, and Step.
// The input must have a top-level "template:" key.
// If refMapping is provided, template references in the output will be replaced.
func Template(yamlStr string, refMapping map[string]string) ([]byte, error) {
	if err := validateTopLevelKey(yamlStr, "template"); err != nil {
		return nil, err
	}

	var doc v0TemplateDoc
	if err := yaml.Unmarshal([]byte(yamlStr), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse v0 template: %w", err)
	}
	if doc.Template.Type == "" {
		return nil, fmt.Errorf("template 'type' field is required (Pipeline, Stage, or Step)")
	}
	if doc.Template.Spec == nil {
		return nil, fmt.Errorf("template 'spec' field is required")
	}

	var (
		v1Spec interface{}
		err    error
	)
	switch strings.ToLower(doc.Template.Type) {
	case "pipeline":
		v1Spec, err = convertPipelineTemplateSpec(doc.Template.Spec)
	case "stage":
		v1Spec, err = convertStageTemplateSpec(doc.Template.Spec)
	case "step":
		v1Spec, err = convertStepTemplateSpec(doc.Template.Spec)
	default:
		return nil, fmt.Errorf("unsupported template type %q (must be Pipeline, Stage, or Step)", doc.Template.Type)
	}
	if err != nil {
		return nil, err
	}

	out := v1TemplateOutput{
		Template: v1TemplateFields{
			Name:       doc.Template.Name,
			Identifier: doc.Template.Identifier,
			Org:        doc.Template.Org,
			Project:    doc.Template.Project,
			Type:       doc.Template.Type,
			Spec:       v1Spec,
		},
	}
	yamlBytes, err := yaml.Marshal(out)
	if err != nil {
		return nil, err
	}

	// Apply template reference replacements if mapping is provided
	return ReplaceTemplateRefs(yamlBytes, refMapping)
}

// convertPipelineTemplateSpec converts the spec of a Pipeline-type template.
// spec is the content that lives under template.spec in the v0 YAML.
func convertPipelineTemplateSpec(spec map[string]interface{}) (interface{}, error) {
	// Wrap the spec map under a "pipeline:" key so the v0 parser recognises it.
	wrapped := map[string]interface{}{"pipeline": spec}
	wrappedYAML, err := yaml.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pipeline template spec: %w", err)
	}

	v0Config, err := v0.ParseBytes(wrappedYAML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pipeline template spec: %w", err)
	}

	c := pipelineconverter.NewPipelineConverter()
	v1Pipeline := c.ConvertPipeline(&v0Config.Pipeline)
	if v1Pipeline == nil {
		return nil, fmt.Errorf("pipeline template spec conversion returned nil")
	}
	return marshalToInterface(v1Pipeline)
}

// convertStageTemplateSpec converts the spec of a Stage-type template by wrapping it
// in a synthetic one-stage pipeline, converting, then extracting that stage.
func convertStageTemplateSpec(spec map[string]interface{}) (interface{}, error) {
	synthetic := map[string]interface{}{
		"pipeline": map[string]interface{}{
			"identifier": "_template",
			"name":       "_template",
			"stages": []interface{}{
				map[string]interface{}{"stage": spec},
			},
		},
	}
	return convertAndExtract(synthetic, func(stages []interface{}) (interface{}, error) {
		if len(stages) == 0 {
			return nil, fmt.Errorf("stage template conversion produced no stages")
		}
		return stages[0], nil
	})
}

// convertStepTemplateSpec converts the spec of a Step-type template by wrapping it
// in a synthetic CI pipeline/stage, converting, then extracting that step.
func convertStepTemplateSpec(spec map[string]interface{}) (interface{}, error) {
	synthetic := map[string]interface{}{
		"pipeline": map[string]interface{}{
			"identifier": "_template",
			"name":       "_template",
			"stages": []interface{}{
				map[string]interface{}{
					"stage": map[string]interface{}{
						"identifier": "_template_stage",
						"name":       "_template_stage",
						"type":       "CI",
						"spec": map[string]interface{}{
							"execution": map[string]interface{}{
								"steps": []interface{}{
									map[string]interface{}{"step": spec},
								},
							},
						},
					},
				},
			},
		},
	}
	return convertAndExtract(synthetic, func(stages []interface{}) (interface{}, error) {
		if len(stages) == 0 {
			return nil, fmt.Errorf("step template conversion produced no stages")
		}
		stageMap, ok := stages[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type for stage in step template conversion")
		}
		steps, _ := stageMap["steps"].([]interface{})
		if len(steps) == 0 {
			return nil, fmt.Errorf("step template conversion produced no steps")
		}
		return steps[0], nil
	})
}

// convertAndExtract marshals a synthetic pipeline map, runs it through the v0→v1
// converter, and calls extract with the resulting v1 stages slice.
func convertAndExtract(
	syntheticPipeline map[string]interface{},
	extract func(stages []interface{}) (interface{}, error),
) (interface{}, error) {
	syntheticYAML, err := yaml.Marshal(syntheticPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal synthetic pipeline: %w", err)
	}

	v0Config, err := v0.ParseBytes(syntheticYAML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse synthetic pipeline: %w", err)
	}

	c := pipelineconverter.NewPipelineConverter()
	v1Pipeline := c.ConvertPipeline(&v0Config.Pipeline)
	if v1Pipeline == nil {
		return nil, fmt.Errorf("synthetic pipeline conversion returned nil")
	}

	// Convert to interface{} via JSON round-trip to honour json struct tags.
	pipelineIface, err := marshalToInterface(v1Pipeline)
	if err != nil {
		return nil, err
	}

	pipelineMap, ok := pipelineIface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for converted pipeline")
	}

	stages, _ := pipelineMap["stages"].([]interface{})
	return extract(stages)
}

// marshalToInterface converts a struct with json tags to map[string]interface{}
// via a JSON round-trip so that yaml.Marshal uses the json tag names as keys.
func marshalToInterface(v interface{}) (interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %w", err)
	}
	var result interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal to interface failed: %w", err)
	}
	return result, nil
}


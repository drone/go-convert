package converter

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// Trigger converts a Harness v0 trigger YAML string into v1 YAML bytes.
// The input must have a top-level "trigger:" key.
// templateRefMapping rewrites template references in the output;
// pipelineRefMapping rewrites pipeline identifiers (including the
// trigger's pipelineIdentifier and any chain.uses inside the embedded
// inputYaml). Either or both may be nil/empty.
// contextPipelineYAML is an optional v0 pipeline YAML used purely as
// expression-postprocess context (see buildContextFromPipelineYAML) for the
// trigger wrapper. Pass "" to run postprocess without FQN context. The
// trigger's embedded inputYaml is always post-processed in FQN mode using
// its own inner pipeline as context (independent of this argument).
//
// Conversion strategy:
//   - The trigger structure remains mostly unchanged
//   - Only the inputYaml content is converted to v1 format (similar to input set conversion)
func Trigger(yamlStr string, templateRefMapping, pipelineRefMapping map[string]string, contextPipelineYAML string) (*Result, error) {
	if err := validateTopLevelKey(yamlStr, "trigger"); err != nil {
		return nil, err
	}

	done := beginAPIConversion()
	defer done()

	v0Config, unknownFields, err := v0.ParseStringWithUnknownFields(yamlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse v0 trigger: %w", err)
	}

	if v0Config.Trigger == nil {
		return nil, fmt.Errorf("trigger parsing returned nil")
	}


	// Single-pass expression post-process on the wrapper. The embedded
	// inputYaml string is post-processed inside ConvertTrigger before
	// marshalling (the wrapper walk skips that field). If the caller
	// supplied a context pipeline_yaml, derive a step-type map and walk in
	// FQN mode; otherwise fall back to nil context (no FQN).
	stepTypeMap, useFQN := buildContextFromPipelineYAML(contextPipelineYAML)
	c := pipelineconverter.NewPipelineConverter()
	
	v1Trigger := c.ConvertTrigger(v0Config.Trigger, stepTypeMap, useFQN)
	if v1Trigger == nil {
		return nil, fmt.Errorf("trigger conversion returned nil")
	}

	yamlBytes, err := v1.MarshalTrigger(v1Trigger)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 trigger: %w", err)
	}

	yamlBytes, err = ApplyRefMappings(yamlBytes, templateRefMapping, pipelineRefMapping)
	if err != nil {
		return nil, err
	}
	return &Result{
		YAML:          yamlBytes,
		UnknownFields: unknownFields,
		Summary:       buildAPISummary(unknownFields),
	}, nil
}

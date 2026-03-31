package converter

import (
	"fmt"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// Pipeline converts a Harness v0 pipeline YAML string into v1 YAML bytes.
// The input must have a top-level "pipeline:" key.
func Pipeline(yamlStr string) ([]byte, error) {
	if err := validateTopLevelKey(yamlStr, "pipeline"); err != nil {
		return nil, err
	}

	v0Config, err := v0.ParseString(yamlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse v0 pipeline: %w", err)
	}
	if v0Config == nil {
		return nil, fmt.Errorf("failed to parse v0 pipeline: result is nil")
	}

	c := pipelineconverter.NewPipelineConverter()
	v1Pipeline := c.ConvertPipeline(&v0Config.Pipeline)
	if v1Pipeline == nil {
		return nil, fmt.Errorf("conversion returned nil: check that the pipeline has at least one supported stage")
	}

	out, err := v1.MarshalPipeline(v1Pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 pipeline: %w", err)
	}
	return out, nil
}

// validateTopLevelKey returns an error when yamlStr does not start with "key:".
// It tolerates a leading YAML document separator (---).
func validateTopLevelKey(yamlStr, key string) error {
	s := strings.TrimSpace(yamlStr)
	if strings.HasPrefix(s, "---") {
		s = strings.TrimSpace(s[3:])
	}
	if !strings.HasPrefix(s, key+":") {
		return fmt.Errorf("expected top-level '%s:' key", key)
	}
	return nil
}

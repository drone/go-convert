package yaml

import (
	"encoding/json"
	"os"

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

	// Convert JSON to interface{} for YAML marshaling
	var jsonData interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
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

package yaml

import (
	"encoding/json"
	"fmt"
)

type (
	Template struct {
		ID           string      `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Name         string      `json:"name,omitempty"              yaml:"name,omitempty"`
		Type         string      `json:"type,omitempty"              yaml:"type,omitempty"`
		VersionLabel string      `json:"versionLabel,omitempty"      yaml:"versionLabel,omitempty"`
		Project      string      `json:"projectIdentifier,omitempty" yaml:"projectIdentifier,omitempty"`
		Org          string      `json:"orgIdentifier,omitempty"     yaml:"orgIdentifier,omitempty"`
		Spec         interface{} `json:"spec,omitempty"              yaml:"spec,omitempty"`
	}
)

func (t *Template) UnmarshalJSON(data []byte) error {
	// First unmarshal into a temporary struct to get the type
	type Alias Template
	aux := &struct {
		Spec json.RawMessage `json:"spec,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Unmarshal Spec based on type
	switch t.Type {
	case "Step":
		var spec Step
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Step template spec: %w", err)
		}
		t.Spec = &spec
	case "Stage":
		var spec Stage
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Stage template spec: %w", err)
		}
		t.Spec = &spec
	case "StepGroup":
		var spec StepGroup
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal StepGroup template spec: %w", err)
		}
		t.Spec = &spec
	case "Pipeline":
		var spec Pipeline
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Pipeline template spec: %w", err)
		}
		t.Spec = &spec
	default:
		return fmt.Errorf("unsupported template type: %s (only Step, Stage, StepGroup, and Pipeline are supported)", t.Type)
	}

	return nil
}

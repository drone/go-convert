// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pipelineconverter

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertTrigger converts a v0 Trigger to v1 Trigger format.
// The trigger structure remains mostly unchanged, only the inputYaml
// content is converted from v0 to v1 format.
func (c *PipelineConverter) ConvertTrigger(src *v0.Trigger, stepTypeMap map[string]*StepInfo, useFQN bool) *v1.Trigger {
	if src == nil {
		return nil
	}

	dst := &v1.Trigger{
		Name:               src.Name,
		ID:                 src.ID,
		Enabled:            src.Enabled,
		StagesToExecute:    src.StagesToExecute,
		Description:        src.Description,
		Tags:               src.Tags,
		Org:                src.Org,
		Project:            src.Project,
		PipelineIdentifier: src.PipelineIdentifier,
		InputSetBranchName: src.InputSetBranchName,
		InputSetReferences: src.InputSetReferences,
	}

	// Convert source
	if src.Source != nil {
		dst.Source = convertTriggerSource(src.Source)
	}

	// Convert inputYaml if present
	// The inputYaml contains a v0 pipeline YAML that needs to be converted to v1 format
	if src.InputYaml != "" {
		convertedInputYaml := c.convertInputYaml(src.InputYaml, stepTypeMap, useFQN)
		dst.InputYaml = convertedInputYaml
	}

	// NOTE: Wrapper-level expression post-processing is performed by the
	// top-level caller in a single pass. The embedded inputYaml string is
	// post-processed locally inside convertInputYaml before marshaling
	// (because the wrapper walk skips the InputYaml field).
	return dst
}

// convertTriggerSource converts v0 TriggerSource to v1 TriggerSource.
func convertTriggerSource(src *v0.TriggerSource) *v1.TriggerSource {
	if src == nil {
		return nil
	}

	dst := &v1.TriggerSource{
		Type: src.Type,
	}

	if src.Spec != nil {
		dst.Spec = &v1.TriggerSourceSpec{
			Type: src.Spec.Type,
			Spec: src.Spec.Spec,
		}
	}

	return dst
}

// convertInputYaml converts the inputYaml string from v0 to v1 format.
// The inputYaml contains a pipeline structure similar to inputSet conversion.
// v0 format: pipeline: { identifier: ..., stages: [...], variables: [...] }
// v1 format: inputs: { overlay: { stages: [...], inputs: [...] } }
//
// On any failure (parse, conversion, marshal) the original v0 inputYaml is
// returned unchanged AND a structured ERROR message is logged via the
// MessageLogger so the surface area shows up in summaries and API reports.
func (c *PipelineConverter) convertInputYaml(inputYaml string, stepTypeMap map[string]*StepInfo, useFQN bool) string {
	if inputYaml == "" {
		return ""
	}

	// Parse the inputYaml as a v0 pipeline
	// The inputYaml has "pipeline:" at the top level
	v0Config, _, err := v0.ParseStringWithUnknownFields(inputYaml)
	if err != nil {
		GetMessageLogger().LogError(
			"TRIGGER_INPUT_YAML_PARSE_FAILED",
			"failed to parse trigger inputYaml as v0 pipeline; emitting v0 fragment unchanged",
			WithContext(map[string]string{"error": err.Error()}),
		)
		return inputYaml
	}

	// Use a fresh PipelineConverter so the inner pipeline conversion does
	// not leak step context into / out of the outer trigger conversion.
	inner := NewPipelineConverter()
	v1Pipeline := inner.ConvertPipeline(&v0Config.Pipeline)
	if v1Pipeline == nil {
		GetMessageLogger().LogError(
			"TRIGGER_INPUT_YAML_CONVERT_NIL",
			"trigger inputYaml conversion returned nil; emitting v0 fragment unchanged",
		)
		return inputYaml
	}

	PostProcessExpressions(v1Pipeline, stepTypeMap, useFQN)

	// Create the v1 InputSet structure with the converted pipeline as overlay
	v1InputSet := &v1.InputSet{
		Overlay: v1Pipeline,
	}

	// Marshal to YAML
	yamlBytes, err := v1.MarshalInputSet(v1InputSet)
	if err != nil {
		GetMessageLogger().LogError(
			"TRIGGER_INPUT_YAML_MARSHAL_FAILED",
			"failed to marshal converted trigger inputYaml; emitting v0 fragment unchanged",
			WithContext(map[string]string{"error": err.Error()}),
		)
		return inputYaml
	}

	// Ensure the output ends with a newline for proper YAML formatting
	result := string(yamlBytes)
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}

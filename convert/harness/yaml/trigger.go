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

package yaml

type (
	// Trigger defines a v0 trigger configuration.
	Trigger struct {
		Name               string                 `json:"name,omitempty"               yaml:"name,omitempty"`
		ID                 string                 `json:"identifier,omitempty"         yaml:"identifier,omitempty"`
		Enabled            bool                   `json:"enabled,omitempty"            yaml:"enabled,omitempty"`
		StagesToExecute    []string               `json:"stagesToExecute,omitempty"    yaml:"stagesToExecute,omitempty"`
		Description        string                 `json:"description,omitempty"        yaml:"description,omitempty"`
		Tags               map[string]string      `json:"tags,omitempty"               yaml:"tags,omitempty"`
		Org                string                 `json:"orgIdentifier,omitempty"      yaml:"orgIdentifier,omitempty"`
		Project            string                 `json:"projectIdentifier,omitempty"  yaml:"projectIdentifier,omitempty"`
		PipelineIdentifier string                 `json:"pipelineIdentifier,omitempty" yaml:"pipelineIdentifier,omitempty"`
		Source             *TriggerSource         `json:"source,omitempty"             yaml:"source,omitempty"`
		InputSetBranchName string                 `json:"inputSetBranchName,omitempty" yaml:"inputSetBranchName,omitempty"`
		InputYaml          string                 `json:"inputYaml,omitempty"          yaml:"inputYaml,omitempty"`
		InputSetReferences []string               `json:"inputSetRefs,omitempty" yaml:"inputSetRefs,omitempty"`
	}

	// TriggerSource defines the trigger source configuration.
	TriggerSource struct {
		Type string             `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *TriggerSourceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// TriggerSourceSpec defines the trigger source spec.
	TriggerSourceSpec struct {
		Type string      `json:"type,omitempty" yaml:"type,omitempty"`
		Spec interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	}
)

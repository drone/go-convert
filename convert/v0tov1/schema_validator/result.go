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

package schema_validator

// ValidationResult holds the outcome of validating a single v1 YAML file
// against the JSON Schema. It is JSON-serialized into schema_validation.json.
type ValidationResult struct {
	Valid        bool             `json:"valid"`
	EntityType   string           `json:"entity_type"`
	FilePath     string           `json:"file_path"`
	ErrorMessage string           `json:"error_message"`
	SchemaErrors []SchemaErrorDTO `json:"schema_errors"`
}

// SchemaErrorDTO is a single schema validation error with location context.
type SchemaErrorDTO struct {
	Message        string         `json:"message"`
	MessageWithFQN string         `json:"message_with_fqn"`
	FQN            string         `json:"fqn"`
	StageInfo      *NodeErrorInfo `json:"stage_info"`
	StepInfo       *NodeErrorInfo `json:"step_info"`
}

// NodeErrorInfo describes the nearest stage or step ancestor of the error.
type NodeErrorInfo struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	FQN        string `json:"fqn"`
}

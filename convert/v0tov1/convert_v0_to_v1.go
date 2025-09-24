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

package v0tov1

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// convertStages converts a list of v0 Stages to v1 Stages.
func convertStages(src []*v0.Stages) []*v1.Stage {
	dst := make([]*v1.Stage, 0)
	for _, stages := range src {
		if stages.Stage != nil {
			dst = append(dst, convertStage(stages.Stage))
		}
		// TODO: handle stages.Parallel recursively when needed
	}
	return dst
}

// convertStage converts a v0 Stage to a v1 Stage.
func convertStage(src *v0.Stage) *v1.Stage {
	if src == nil {
		return nil
	}

	var steps []*v1.Step
	var rollback []*v1.Step
	var service *v1.ServiceRef
	var environment *v1.EnvironmentRef
	var onFailure []*v1.FailureStrategy

	switch spec := src.Spec.(type) {
	case *v0.StageCI:
		steps = convert_helpers.ConvertSteps(spec.Execution.Steps)
	case *v0.StageDeployment:
		// Convert deployment steps
		if spec.Execution != nil {
			steps = convert_helpers.ConvertSteps(spec.Execution.Steps)
			rollback = convert_helpers.ConvertSteps(spec.Execution.RollbackSteps)
		}

		// Convert service configuration
		if spec.Service != nil {
			service = convert_helpers.ConvertDeploymentService(spec.Service)
		} else if spec.Services != nil {
			service = convert_helpers.ConvertDeploymentServices(spec.Services)
		}

		// Convert environment configuration - check all possible sources
		if spec.Environment != nil {
			environment = convert_helpers.ConvertEnvironment(spec.Environment)
		} else if spec.Environments != nil {
			environment = convert_helpers.ConvertEnvironments(spec.Environments)
		} else if spec.EnvironmentGroup != nil {
			environment = convert_helpers.ConvertEnvironmentGroup(spec.EnvironmentGroup)
		}
	default:
		fmt.Println("stage type: " + src.Type + " is not yet supported!")
		// non-CI/Deployment stages currently not converted
	}
	onFailure = convert_helpers.ConvertFailureStrategies(src.FailureStrategies)
	delegate := &v1.Delegate{
		Filter: src.DelegateSelectors,
	}
	return &v1.Stage{
		Id:          src.ID,
		Name:        src.Name,
		Steps:       steps,
		Rollback:    rollback,
		Service:     service,
		Environment: environment,
		OnFailure:   onFailure,
		Inputs:      convertVariables(src.Vars),
		Delegate:    delegate,
	}
}

// convertBarriers converts a list of v0 Barriers to v1 Barriers.
func convertBarriers(src []*v0.Barrier) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, 0, len(src))
	for _, barrier := range src {
		if barrier == nil {
			continue
		}
		dst = append(dst, barrier.Name)
	}
	return dst
}

// convertVariables converts a list of v0 Variables to v1 Inputs.
func convertVariables(src []*v0.Variable) map[string]*v1.Input {
	if len(src) == 0 {
		return nil
	}

	dst := make(map[string]*v1.Input)
	for _, variable := range src {
		if variable == nil || variable.Name == "" {
			continue
		}

		input := &v1.Input{
			Type:     convertVariableType(variable.Type),
			Required: true, // Default to required, can be adjusted based on requirements
		}

		// Set default value if provided
		if variable.Value != "" {
			input.Default = variable.Value
		}

		// Set mask to true for secret types
		if variable.Type == "Secret" {
			input.Mask = true
		}

		dst[variable.Name] = input
	}

	return dst
}

// convertVariableType converts v0 variable type to v1 input type.
func convertVariableType(v0Type string) string {
	switch v0Type {
	case "Secret":
		return "string" // Secrets are still strings in v1, but with mask=true
	case "Text":
		return "string"
	default:
		return "string" // Default to string type
	}
}

// ConvertPipeline converts a v0 Pipeline to a v1 Pipeline.
func ConvertPipeline(src *v0.Pipeline) *v1.Pipeline {
	if src == nil {
		return nil
	}
	var barriers []string
	if src.FlowControl != nil {
		barriers = convertBarriers(src.FlowControl.Barriers)
	}
	dst := &v1.Pipeline{
		Id:       src.ID,
		Name:     src.Name,
		Inputs:   convertVariables(src.Variables),
		Stages:   convertStages(src.Stages),
		Barriers: barriers,
	}

	return dst
}

func Main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run convert.go <input-file> [output-file]")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Determine output file
	var outputFile string
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	} else {
		ext := filepath.Ext(inputFile)
		base := strings.TrimSuffix(inputFile, ext)
		outputFile = base + ".v1" + ext
	}

	// Read v0 YAML file and fit to v0 structs
	v0Config, err := v0.ParseFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to parse v0 pipeline file: %v", err)
	}

	// Convert to v1 structs
	v1Pipeline := ConvertPipeline(&v0Config.Pipeline)
	if v1Pipeline == nil {
		log.Fatal("Failed to convert pipeline to v1 format")
	}

	// Write v1 YAML file with top-level 'pipeline:' key using writer
	if err := v1.WritePipelineFile(outputFile, v1Pipeline); err != nil {
		log.Fatalf("Failed to write v1 pipeline YAML: %v", err)
	}

	fmt.Printf("Converted %s -> %s\n", inputFile, outputFile)
}

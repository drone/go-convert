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

// Normalizes normalizes the pipeline configuration by
// grouping all top-level steps into stages.
func Normalize(config *Config) {
	var steps []*Steps
	var stage *Steps
	for _, step := range config.Pipelines.Default {
		// if steps are already grouped in a stage,
		// append to our consolidated list.
		if step.Stage != nil {
			stage = nil // reset the stage
			steps = append(steps, step)
			continue
		}

		// create a new stage to group subsequent
		// steps.
		if stage == nil {
			stage = new(Steps)
			stage.Stage = new(Stage)
			steps = append(steps, stage)
		}

		// append the step to the stage
		stage.Stage.Steps = append(stage.Stage.Steps, step)
	}

	// replace original steps with normalized steps,
	// grouped by stages.
	config.Pipelines.Default = steps
}

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

package expand

import (
	"encoding/json"

	"github.com/bradrydzewski/spec/yaml"
	"github.com/bradrydzewski/spec/yaml/utils/expand/matrix"
	"github.com/bradrydzewski/spec/yaml/utils/walk"
)

// Expand expands the matrix strategies.
func Expand(in *yaml.Schema) error {
	walk.Walk(in, func(v interface{}) error {
		switch vv := v.(type) {
		case *yaml.Stage:
			ExpandStage(vv)
		case *yaml.Step:
			ExpandStep(vv)
		}
		return nil
	})

	return nil
}

// ExpandStage expands the matrix strategies for the stage.
func ExpandStage(v *yaml.Stage) {
	// exit if there is not strategy
	// or matrix defined.
	if v.Strategy == nil || v.Strategy.Matrix == nil {
		return
	}

	// exit if there's no axis
	if len(v.Strategy.Matrix.Axis) == 0 {
		return
	}

	// calculate the matrix permutations.
	perms := matrix.Calc(v.Strategy.Matrix.Axis)

	var stages []*yaml.Stage
	// create a new stage for each item in the
	// matrix and update the strategy to include
	// ony the relevant matrix includes for this
	// specific permutation.
	for _, perm := range perms {

		// we need to make a deep copy to prevent
		// multiple stages from sharing the same
		// child objects
		stage := deepCopyStage(v)

		// unset the stage strategy
		stage.Strategy = nil

		// append the matrx variables to
		// the context.
		if stage.Context == nil {
			stage.Context = new(yaml.Context)
		}
		stage.Context.Matrix = perm

		// append to the stage list
		stages = append(stages, stage)
	}

	// unset the name and identifier
	v.Id = ""
	v.Name = ""

	// unset the stage types
	v.Strategy = nil
	v.Approval = nil
	v.Group = nil
	v.Parallel = nil
	v.Template = nil
	v.Steps = nil

	// replace with the parallel stage, where
	// each stage in the group is a permutation
	// in the matrix.
	v.Parallel = &yaml.StageGroup{
		Stages: stages,
	}
}

// ExpandStage expands the matrix strategies for the step.
func ExpandStep(v *yaml.Step) {
	// exit if there is not strategy
	// or matrix defined.
	if v.Strategy == nil || v.Strategy.Matrix == nil {
		return
	}

	// exit if there's no axis
	if len(v.Strategy.Matrix.Axis) == 0 {
		return
	}

	// calculate the matrix permutations.
	perms := matrix.Calc(v.Strategy.Matrix.Axis)

	var steps []*yaml.Step
	// create a new stage for each item in the
	// matrix and update the strategy to include
	// ony the relevant matrix includes for this
	// specific permutation.
	for _, perm := range perms {

		// we need to make a deep copy to prevent
		// multiple steps from sharing the same
		// child objects
		step := deepCopyStep(v)

		// unset the step strategy
		step.Strategy = nil

		// append the matrx variables to
		// the context.
		if step.Context == nil {
			step.Context = new(yaml.Context)
		}
		step.Context.Matrix = perm

		// append to the step list
		steps = append(steps, step)
	}

	// unset the name and identifier
	v.Id = ""
	v.Name = ""

	// unset the matrix
	v.Strategy = nil

	// unset the step types
	v.Strategy = nil
	v.Action = nil
	v.Approval = nil
	v.Background = nil
	v.Group = nil
	v.Parallel = nil
	v.Run = nil
	v.RunTest = nil
	v.Template = nil

	// replace with the parallel step, where
	// each step in the group is a permutation
	// in the matrix.
	v.Parallel = &yaml.StepGroup{
		Steps: steps,
	}
}

// helper function creates a deep copy of a stage
func deepCopyStage(in *yaml.Stage) *yaml.Stage {
	out := new(yaml.Stage)
	raw, _ := json.Marshal(in)   // assumes no errors
	_ = json.Unmarshal(raw, out) // assumes no errors
	return out
}

// helper function creates a deep copy of a stage
func deepCopyStep(in *yaml.Step) *yaml.Step {
	out := new(yaml.Step)
	raw, _ := json.Marshal(in)   // assumes no errors
	_ = json.Unmarshal(raw, out) // assumes no errors
	return out
}

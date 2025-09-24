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

package normalize

import (
	"github.com/bradrydzewski/spec/yaml"
	"github.com/bradrydzewski/spec/yaml/utils/walk"
)

// Normalize normalizes the yaml to ensure all stages and
// steps have unique identifiers.
func Normalize(in *yaml.Schema) error {
	gen := newGenerator()

	// check to see if the yaml is using github actions
	// syntax and convert as needed.
	if len(in.Jobs) > 0 {
		if in.Pipeline == nil {
			in.Pipeline = new(yaml.Pipeline)
		}
		for name, stage := range in.Jobs {
			stage.Name = name
			in.Pipeline.Stages = append(in.Pipeline.Stages, stage)
		}
		// unset jobs once converted to stages.
		in.Jobs = nil
	}

	// walk the yaml and ensure all stages and steps
	// have unique identifiers.
	walk.Walk(in, func(i interface{}) error {
		switch v := i.(type) {
		case *yaml.Stage:
			v.Id = gen.generate(v.Id, v.Name, stageType(v))
		case *yaml.Step:
			v.Id = gen.generate(v.Id, v.Name, stepType(v))
		}

		return nil
	})

	return nil
}

// helper function returns the step type.
func stepType(step *yaml.Step) string {
	switch {
	case step.Run != nil:
		return "run"
	case step.Group != nil:
		return "group"
	case step.Parallel != nil:
		return "parallel"
	case step.Template != nil:
		return "template"
	case step.Action != nil:
		return "action"
	default:
		return "step"
	}
}

// helper function returns the stage type.
func stageType(step *yaml.Stage) string {
	switch {
	case step.Group != nil:
		return "group"
	case step.Parallel != nil:
		return "parallel"
	case step.Template != nil:
		return "template"
	default:
		return "stage"
	}
}

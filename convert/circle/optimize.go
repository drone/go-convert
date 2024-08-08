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

package circle

import (
	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// this is a poor attempt at an optimization to remove
// individual cache steps and replace with cache intelligence,
// if possible. The criteris is that there can only be one
// save step, and zero or one restore steps. If both a save
// and restore step exist, they must have the same cache key.
func optimizeCache(stage *harness.StageCI) {
	var save *harness.StepPlugin
	var restore *harness.StepPlugin
	var steps []*harness.Step

	for _, step := range stage.Steps {
		spec, ok := step.Spec.(*harness.StepPlugin)
		if ok == false || spec.With == nil {
			steps = append(steps, step)
			continue
		}
		if spec.With["rebuild"] == "true" {
			if save != nil {
				return // exit if multiple save
			} else {
				save = spec
			}
		} else if spec.With["restore"] == "true" {
			if restore != nil {
				return // exit if multiple restore
			} else {
				restore = spec
			}
		} else {
			steps = append(steps, step)
		}
	}
	if save == nil {
		return
	}
	if restore != nil &&
		save.With["cache_key"] != restore.With["cache_key"] {
		return
	}

	stage.Steps = steps
	stage.Cache = &harness.Cache{
		Enabled: true,
		Key:     save.With["cache_key"].(string),
		Paths:   save.With["mount"].(circle.Stringorslice),
	}
}

// this is a helper function that optimizes stages that
// have a single group step.
func optimizeGroup(stage *harness.StageCI) {
	if len(stage.Steps) != 1 {
		return
	}
	step := stage.Steps[0]
	if step.Spec == nil {
		return // should never happen
	}
	if group, ok := step.Spec.(*harness.StepGroup); ok {
		stage.Steps = group.Steps
	}
}

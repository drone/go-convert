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

package bitbucket

import (
	"sort"

	bitbucket "github.com/drone/go-convert/convert/bitbucket/yaml"
)

//
// this file contains helper functions that help
// extract data from the yaml.
//

// helper function that recursively returns all stage steps
// including parallel steps.
func extractSteps(stage *bitbucket.Stage) []*bitbucket.Step {
	var steps []*bitbucket.Step
	// iterate through steps the he stage
	for _, step := range stage.Steps {
		if step.Step != nil {
			steps = append(steps, step.Step)
		}
		// iterate through parallel steps
		if step.Parallel != nil {
			for _, step2 := range step.Parallel.Steps {
				if step2.Step != nil {
					steps = append(steps, step2.Step)
				}
			}
		}
	}
	return steps
}

// helper function that recursively returnsall stage steps,
// including parallel steps.
func extractAllSteps(stages []*bitbucket.Steps) []*bitbucket.Step {
	var steps []*bitbucket.Step
	for _, stage := range stages {
		if stage.Stage != nil {
			steps = append(steps, extractSteps(stage.Stage)...)
		}
	}
	return steps
}

// helper function that returns a recommended machine size
// based on the gloal machine size, as well as any step-level
// overrides.
func extractSize(opts *bitbucket.Options, stage *bitbucket.Stage) bitbucket.Size {
	var size bitbucket.Size

	// start with the global size, if set
	if opts != nil {
		size = opts.Size
	}

	// loop through the steps and if a step
	// defines a size greater than the global
	// size, us the step size instead.
	for _, step := range extractSteps(stage) {
		if step.Size > size {
			size = step.Size
		}
	}
	return size
}

// helper function that returns runs-on tags for routing
// a stage to a specific set of workers. In Bitbucket, each
// step can be routed to a different machine, which is not
// compatible with Harness. We emulate this behavior by
// aggregating all run-on tags and applying them at the
// stage level in Harness.
func extractRunsOn(stage *bitbucket.Stage) []string {
	set := map[string]struct{}{}

	// loop through the steps and if a step
	// defines a size greater than the global
	// size, us the step size instead.
	for _, step := range extractSteps(stage) {
		for _, s := range step.RunsOn {
			set[s] = struct{}{}
		}
	}

	// convert the map to a slice.
	var unique []string
	for k := range set {
		unique = append(unique, k)
	}

	// sort for deterministic unit testing
	sort.Strings(unique)

	return unique
}

// helper function that returns a list of all caches used
// by all steps in a given stage.
func extractCache(stage *bitbucket.Stage) []string {
	set := map[string]struct{}{}

	// loop through the steps and if a step
	// defines cache directories
	for _, step := range extractSteps(stage) {
		for _, s := range step.Caches {
			set[s] = struct{}{}
		}
	}

	// convert the map to a slice.
	var unique []string
	for k := range set {
		unique = append(unique, k)
	}

	// sort for deterministic unit testing
	sort.Strings(unique)

	return unique
}

// helper function that returns a list of all services used
// by all steps in a given stage.
func extractServices(stage *bitbucket.Stage) []string {
	set := map[string]struct{}{}

	// loop through the steps and if a step
	// defines cache directories
	for _, step := range extractSteps(stage) {
		for _, s := range step.Services {
			set[s] = struct{}{}
		}
	}

	// convert the map to a slice.
	var unique []string
	for k := range set {
		unique = append(unique, k)
	}

	// sort for deterministic unit testing
	sort.Strings(unique)

	return unique
}

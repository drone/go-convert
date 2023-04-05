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

package node

import (
	"fmt"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb step to a Harness step.
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "install", "install-packages", "install-yarn":
		return convertInstall(step)
	case "test":
		return convertTest(step)
	default:
		return convertTest(step)
	}
}

// helper function converts a node/install-packages
// orb to a run step.
func convertInstall(step *circle.Custom) *harness.Step {
	if step.Params["install-yarn"] == true ||
		step.Params["pkg-manager"] == "yarn" {
		return &harness.Step{
			Name: "install_packages",
			Type: "script",
			Spec: &harness.StepExec{
				Run: "yarn install",
			},
		}
	}
	return &harness.Step{
		Name: "install_packages",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "npm install",
		},
	}
}

// helper function converts a node/test
// orb to a run step.
//
// https://github.com/CircleCI-Public/node-orb/blob/master/src/jobs/test.yml
func convertTest(step *circle.Custom) *harness.Step {
	cmd := "test"
	run := "npm %s"

	switch step.Params["test-results-for"] {
	case "mocha":
		run = `npm run %s -- --reporter mocha-multi --reporter-options spec=-,mocha-junit-reporter=-`
	case "jest":
		run = `npm run %s -- --reporters=default --reporters=jest-junit`
	}

	if p, ok := step.Params["run-command"].(string); ok {
		cmd = p
	}

	return &harness.Step{
		Name: "test",
		Type: "script",
		Spec: &harness.StepExec{
			Run: fmt.Sprintf(run, cmd),
		},
	}
}

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

// Convert converts an Orb to a Harness step.
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "install", "install-packages", "install-yarn":
		return convertInstallCmd(step)
	case "run":
		return convertRunJob(step)
	case "test":
		return convertTestJob(step)
	default:
		return nil
	}
}

// helper function converts a node/test job
func convertTestJob(step *circle.Custom) *harness.Step {
	var steps []*harness.Step
	steps = append(steps, convertInstallCmd(step))
	steps = append(steps, convertTestCmd(step))
	return &harness.Step{
		Name: "test",
		Type: "group",
		Spec: &harness.StepGroup{
			Steps: steps,
		},
	}
}

// helper function converts a node/test
// command to a run step.
//
// https://github.com/CircleCI-Public/node-orb/blob/master/src/jobs/test.yml
func convertTestCmd(step *circle.Custom) *harness.Step {
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
		Name: "run_tests",
		Type: "script",
		Spec: &harness.StepExec{
			Run: fmt.Sprintf(run, cmd),
		},
	}
}

// helper function converts a node/test job
func convertRunJob(step *circle.Custom) *harness.Step {
	var steps []*harness.Step
	steps = append(steps, convertInstallCmd(step))
	steps = append(steps, convertRunCmd(step))
	return &harness.Step{
		Name: "test",
		Type: "group",
		Spec: &harness.StepGroup{
			Steps: steps,
		},
	}
}

// helper function converts a node/run
// orb to a run step.
func convertRunCmd(step *circle.Custom) *harness.Step {
	manager := "npm"
	command := "ci"
	if s, _ := step.Params["yarn-run"].(string); s != "" {
		manager = "yarn"
		command = s
	}
	if s, _ := step.Params["npm-run"].(string); s != "" {
		manager = "npm"
		command = s
	}
	return &harness.Step{
		Name: fmt.Sprintf("%s_run", manager),
		Type: "script",
		Spec: &harness.StepExec{
			Run: fmt.Sprintf("%s run %s", manager, command),
		},
	}
}

// helper function converts a node/install-packages
// command to a run step.
func convertInstallCmd(step *circle.Custom) *harness.Step {
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

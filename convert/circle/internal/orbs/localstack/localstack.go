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

package localstack

import (
	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb to a Harness step.
// https://circleci.com/developer/orbs/orb/localstack/platform
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "start", "startup":
		return convertStart(step)
	case "wait":
		return convertWait(step)
	default:
		return nil
	}
}

func convertStart(step *circle.Custom) *harness.Step {
	return &harness.Step{
		Name: "localstack",
		Type: "background",
		Spec: &harness.StepBackground{
			Image: "localstack/localstack",
			Ports: []string{"4566"},
		},
	}
}

func convertWait(step *circle.Custom) *harness.Step {
	return &harness.Step{
		Name: "localstack_wait",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "sleep 10", //  sleep for enough time to localstack to start
		},
	}
}

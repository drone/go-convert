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

package golang

import (
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

//
// https://circleci.com/developer/orbs/orb/circleci/go
//

// Convert converts an Orb to a Harness step.
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "install":
		return convertInstall(step)
	case "test":
		return convertTest(step)
	default:
		return nil
	}
}

func convertInstall(step *circle.Custom) *harness.Step {
	return &harness.Step{
		Name: "go_install",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "go install ./...",
		},
	}
}

func convertTest(step *circle.Custom) *harness.Step {
	parts := []string{
		"go test -cover",
	}

	if v, ok := step.Params["covermode"]; ok {
		parts = append(parts, fmt.Sprintf("-covermode %v", v))
	}

	if v, ok := step.Params["coverpkg"]; ok {
		parts = append(parts, fmt.Sprintf("-covermode %v", v))
	}

	if v, ok := step.Params["count"]; ok {
		parts = append(parts, fmt.Sprintf("-count %v", v))
	}

	if v, ok := step.Params["parallel"]; ok {
		parts = append(parts, fmt.Sprintf("-parallel %v", v))
	}

	if v, _ := step.Params["verbose"].(bool); v {
		parts = append(parts, "-v")
	}

	if v, _ := step.Params["race"].(bool); v {
		parts = append(parts, "-race")
	}

	if v, _ := step.Params["short"].(bool); v {
		parts = append(parts, "-short")
	}

	if v, ok := step.Params["packages"]; ok {
		parts = append(parts, fmt.Sprint(v))
	} else {
		parts = append(parts, "./...")
	}

	return &harness.Step{
		Name: "go_test",
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(parts, " "),
		},
	}
}

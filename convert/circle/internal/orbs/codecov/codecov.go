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

package codecov

import (
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb to a Harness step.
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "upload":
		return convertUpload(step)
	default:
		return nil
	}
}

func convertUpload(step *circle.Custom) *harness.Step {
	token := ""
	if s, _ := step.Params["token"].(string); s != "" {
		token = s
	}
	name := "$DRONE_BUILD_NUMBER"
	if s, _ := step.Params["upload_name"].(string); s != "" {
		name = s
	}

	parts := []string{
		"curl -Os https://uploader.codecov.io/latest/linux/codecov",
		"chmod +x codecov",
		fmt.Sprintf("./codecov -t %s -n %s", token, name),
	}

	return &harness.Step{
		Name: "codecov",
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(parts, "\n"),
		},
	}
}

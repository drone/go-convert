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

package coveralls

import (
	"bytes"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb to a Harness step.
// https://circleci.com/developer/orbs/orb/coveralls/coveralls
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "upload":
		return convertUpload(step)
	default:
		return nil
	}
}

func convertUpload(step *circle.Custom) *harness.Step {
	var buf bytes.Buffer
	buf.WriteString("./coveralls")

	if b, _ := step.Params["verbose"].(bool); b == true {
		buf.WriteString("  --debug")
	}

	if b, _ := step.Params["dry_run"].(bool); b == true {
		buf.WriteString("  --dry-run")
	}

	if s, _ := step.Params["base_path"].(string); s != "" {
		buf.WriteString("  --base-path " + s)
	}

	if s, _ := step.Params["coverage_format"].(string); s != "" {
		buf.WriteString("  --format " + s)
	}

	if s, _ := step.Params["coverage_file"].(string); s != "" {
		buf.WriteString("  --file " + s)
	}

	parts := []string{
		"curl -sLO https://github.com/coverallsapp/coverage-reporter/releases/latest/download/coveralls-linux.tar.gz",
		"tar -xzf coveralls-linux.tar.gz",
		buf.String(),
	}

	var token string
	if s, _ := step.Params["token"].(string); s != "" {
		token = s
	}

	return &harness.Step{
		Name: "coveralls",
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(parts, "\n"),
			Envs: map[string]string{
				"COVERALLS_REPO_TOKEN": token,
			},
		},
	}
}

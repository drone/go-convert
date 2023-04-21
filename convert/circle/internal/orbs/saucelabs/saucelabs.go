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

package saucelabs

import (
	"bytes"
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb to a Harness step.
// https://circleci.com/developer/orbs/orb/saucelabs/saucectl-run
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "saucectl-run":
		return convertRun(step)
	default:
		return nil
	}
}

func convertRun(step *circle.Custom) *harness.Step {
	var buf bytes.Buffer
	buf.WriteString("saucectl")

	if s, _ := step.Params["config-file"].(string); s != "" {
		buf.WriteString(" -c " + s)
	}

	if s, _ := step.Params["region"].(string); s != "" {
		buf.WriteString(" --region " + s)
	}

	if s, _ := step.Params["select-suite"].(string); s != "" {
		buf.WriteString(" --select-suite" + s)
	}

	if s, _ := step.Params["sauceignore"].(string); s != "" {
		buf.WriteString(" --sauceignore" + s)
	}

	if s, _ := step.Params["tunnel-name"].(string); s != "" {
		buf.WriteString(" --tunnel-name" + s)
	}

	if s, _ := step.Params["tunnel-owner"].(string); s != "" {
		buf.WriteString(" --tunnel-owner" + s)
	}

	if b, _ := step.Params["show-console-log"].(bool); b {
		buf.WriteString(" --show-console-log")
	}

	if b, _ := step.Params["ccy"].(bool); b {
		buf.WriteString(" --ccy")
	}

	if b, _ := step.Params["test-env-silent"].(bool); b {
		buf.WriteString(" --test-env-silent")
	}

	if n, ok := step.Params["timeout"]; ok {
		buf.WriteString(" --timeout " + fmt.Sprint(n))
	}

	if n, ok := step.Params["retries"]; ok {
		buf.WriteString(" --retries " + fmt.Sprint(n))
	}

	cmds := []string{
		"curl -L https://saucelabs.github.io/saucectl/install | bash -s -- -b /usr/local/bin",
		buf.String(),
	}

	// add username and access key as environment
	// variables.
	envs := map[string]string{}
	if s, _ := step.Params["sauce-username"].(string); s != "" {
		envs["SAUCE_USERNAME"] = s
	}
	if s, _ := step.Params["sauce-access-key"].(string); s != "" {
		envs["SAUCE_ACCESS_KEY"] = s
	}

	// convert user-defined environment variable string
	// to an environment map.
	if s, _ := step.Params["env"].(string); s != "" {
		for _, env := range strings.Split(s, "\n") {
			if parts := strings.SplitN(env, "=", 2); len(parts) == 2 {
				a := parts[0]
				b := parts[1]
				envs[a] = b
			}
		}
	}

	// currently harness does not provide an option to
	// change the step working directory, so we need to
	// prepend the cd command to the script.
	if s, _ := step.Params["working-directory"].(string); s != "" {
		cmds = append([]string{"cd " + s}, cmds...)
	}

	return &harness.Step{
		Name: "saucelabs",
		Type: "background",
		Spec: &harness.StepBackground{
			Run:  strings.Join(cmds, "\n"),
			Envs: envs,
		},
	}
}

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

package ruby

import (
	"bytes"
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb to a Harness step.
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "install":
		return nil
	case "install-deps":
		return convertInstallDeps(step)
	case "rspec-test":
		return convertRspecTest(step)
	case "rubocop-check":
		return convertRubocopCheck(step)
	default:
		return nil
	}
}

// helper function converts a ruby/rspec-test
// orb to a run step.
func convertRspecTest(step *circle.Custom) *harness.Step {
	outpath := "/tmp/test-results/rspec"
	if s, _ := step.Params["path"].(string); s != "" {
		outpath = s
	}

	var script []string
	script = append(script, fmt.Sprintf(`mkdir -p %q`, outpath))
	script = append(script, fmt.Sprintf(
		`bundle exec rspec --profile 10 --format RspecJunitFormatter --out %s/results.xml --format progress`,
		outpath,
	),
	)

	return &harness.Step{
		Name: "rspec_test",
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(script, "\n"),
		},
	}
}

// helper function converts a ruby/install-deps
// orb to a run step.
func convertInstallDeps(step *circle.Custom) *harness.Step {
	path := "./vendor/bundle"
	clean := false
	gemfile := "Gemfile"
	if s, _ := step.Params["path"].(string); s != "" {
		path = s
	}
	if s, _ := step.Params["gemfile"].(string); s != "" {
		gemfile = s
	}
	if s, _ := step.Params["clean-bundle"].(bool); s {
		clean = s
	}

	var script []string
	if path == "./vendor/bundle" {
		script = append(script, `bundle config deployment 'true'`)
	}

	script = append(script, fmt.Sprintf(`bundle config gemfile %q`, gemfile))
	script = append(script, fmt.Sprintf(`bundle config path %q`, path))

	if clean {
		script = append(script, `bundle check || (bundle install && bundle clean --force)`)
	} else {
		script = append(script, `bundle check || bundle install`)
	}

	return &harness.Step{
		Name: "install_deps",
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(script, "\n"),
		},
	}
}

// helper function converts a ruby/rubocop-check
// orb to a run step.
func convertRubocopCheck(step *circle.Custom) *harness.Step {
	checkpath := "."
	format := "progress"
	outpath := "/tmp/rubocop-results"
	parallel := false

	if s, _ := step.Params["check-path"].(string); s != "" {
		checkpath = s
	}
	if s, _ := step.Params["out-path"].(string); s != "" {
		outpath = s
	}
	if s, _ := step.Params["parallel"].(bool); s {
		parallel = s
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`mkdir -p %q`, outpath))
	buf.WriteString("\n")

	if !parallel {
		buf.WriteString(
			fmt.Sprintf(
				`bundle exec rubocop %s --out %s/check-results.xml --format %s`,
				checkpath,
				outpath,
				format,
			),
		)
	} else {
		buf.WriteString(
			fmt.Sprintf(
				`bundle exec rubocop %s --out %s/check-results.xml --format %s --parallel`,
				checkpath,
				outpath,
				format,
			),
		)
	}

	return &harness.Step{
		Name: "rubocop_check",
		Type: "script",
		Spec: &harness.StepExec{
			Run: buf.String(),
		},
	}
}

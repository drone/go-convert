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
	"strings"
	"testing"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/google/go-cmp/cmp"
)

func TestUpload(t *testing.T) {
	in := &circle.Custom{
		Params: map[string]interface{}{
			"token": "SECRET",
		},
	}

	got := Convert("upload", in)
	want := &harness.Step{
		Name: "coveralls",
		Type: "script",
		Spec: &harness.StepExec{
			Envs: map[string]string{"COVERALLS_REPO_TOKEN": "SECRET"},
			Run: strings.Join(
				[]string{
					"curl -sLO https://github.com/coverallsapp/coverage-reporter/releases/latest/download/coveralls-linux.tar.gz",
					"tar -xzf coveralls-linux.tar.gz",
					"./coveralls",
				}, "\n",
			),
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected orb conversion")
		t.Log(diff)
	}
}

func TestUnknownCommand(t *testing.T) {
	if Convert("unknown", nil) != nil {
		t.Errorf("Expect unknown command returns nil step")
	}
}

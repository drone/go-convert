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
	"testing"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/google/go-cmp/cmp"
)

func TestConvertTest(t *testing.T) {
	in := &circle.Custom{
		Params: map[string]interface{}{},
	}

	got := Convert("test", in)
	want := &harness.Step{
		Name: "go_test",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "go test -cover ./...",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected orb conversion")
		t.Log(diff)
	}
}

func TestConvertTestParams(t *testing.T) {
	in := &circle.Custom{
		Params: map[string]interface{}{
			"verbose":  true,
			"race":     true,
			"short":    true,
			"parallel": 5,
		},
	}

	got := Convert("test", in)
	want := &harness.Step{
		Name: "go_test",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "go test -cover -parallel 5 -v -race -short ./...",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected orb conversion")
		t.Log(diff)
	}
}

func TestConvertInstall(t *testing.T) {
	in := &circle.Custom{
		Params: map[string]interface{}{},
	}

	got := Convert("install", in)
	want := &harness.Step{
		Name: "go_install",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "go install ./...",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected orb conversion")
		t.Log(diff)
	}
}

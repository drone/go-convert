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

package yaml

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestInclude(t *testing.T) {
	tests := []struct {
		yaml string
		want Include
	}{
		{
			yaml: `".gitlab-ci-production.yml"`,
			want: Include{
				Local: ".gitlab-ci-production.yml",
			},
		},
		{
			yaml: `{ "project": "my-group/my-project", "file": "/templates/.gitlab-ci-template.yml" }`,
			want: Include{
				Project: "my-group/my-project",
				File:    []string{"/templates/.gitlab-ci-template.yml"},
			},
		},
		{
			yaml: `{ "project": "my-group/my-project", "file": [ "/templates/.builds.yml", "/templates/.tests.yml" ] }`,
			want: Include{
				Project: "my-group/my-project",
				File:    []string{"/templates/.builds.yml", "/templates/.tests.yml"},
			},
		},
		{
			yaml: `{ "remote": "https://gitlab.com/example-project/-/raw/main/.gitlab-ci.yml" }`,
			want: Include{
				Remote: "https://gitlab.com/example-project/-/raw/main/.gitlab-ci.yml",
			},
		},
		{
			yaml: `{ "template": "Auto-DevOps.gitlab-ci.yml" }`,
			want: Include{
				Template: "Auto-DevOps.gitlab-ci.yml",
			},
		},
		{
			yaml: `{ "project": "my-group/my-project", "ref": "main", "file": "/templates/.gitlab-ci-template.yml" }`,
			want: Include{
				Project: "my-group/my-project",
				Ref:     "main",
				File:    []string{"/templates/.gitlab-ci-template.yml"},
			},
		},
		{
			yaml: `{ "project": "my-group/my-project", "ref": "v1.0.0", "file": "/templates/.gitlab-ci-template.yml" }`,
			want: Include{
				Project: "my-group/my-project",
				Ref:     "v1.0.0",
				File:    []string{"/templates/.gitlab-ci-template.yml"},
			},
		},
		{
			yaml: `{ "project": "my-group/my-project", "ref": "787123b47f14b552955ca2786bc9542ae66fee5b", "file": "/templates/.gitlab-ci-template.yml" }`,
			want: Include{
				Project: "my-group/my-project",
				Ref:     "787123b47f14b552955ca2786bc9542ae66fee5b",
				File:    []string{"/templates/.gitlab-ci-template.yml"},
			},
		},
	}

	for i, test := range tests {
		got := new(Include)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestInclude_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[]"), new(Include))
	if err == nil || err.Error() != "failed to unmarshal include" {
		t.Errorf("Expect error, got %s", err)
	}
}

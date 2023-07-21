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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestArtifacts(t *testing.T) {
	tests := []struct {
		yaml string
		want Artifacts
	}{
		{
			yaml: `
  paths:
    - binaries/
    - .config`,
			want: Artifacts{
				Paths: Stringorslice{
					"binaries/",
					".config",
				},
			},
		},
		{
			yaml: `
  paths:
    - test.txt`,
			want: Artifacts{
				Paths: Stringorslice{
					"test.txt",
				},
			},
		},
		{
			yaml: `
  paths: []`,
			want: Artifacts{},
		},
		{
			yaml: `
  paths:
    - binaries/
  exclude:
    - binaries/**/*.o`,
			want: Artifacts{
				Paths: Stringorslice{
					"binaries/",
				},
				Exclude: Stringorslice{
					"binaries/**/*.o",
				},
			},
		},
		{
			yaml: `
  paths:
    - test.txt
  exclude:
    - test.log`,
			want: Artifacts{
				Paths: Stringorslice{
					"test.txt",
				},
				Exclude: Stringorslice{
					"test.log",
				},
			},
		},
		{
			yaml: `
  paths:
    - test.txt
  exclude: []`,
			want: Artifacts{
				Paths: Stringorslice{
					"test.txt",
				},
			},
		},
		{
			yaml: `
  expire_in: '42'`,
			want: Artifacts{
				ExpireIn: "42",
			},
		},
		{
			yaml: `
  expire_in: '2 hrs 20 min'`,
			want: Artifacts{
				ExpireIn: "2 hrs 20 min",
			},
		},
		{
			yaml: `
  expire_in: '6 mos 1 day'`,
			want: Artifacts{
				ExpireIn: "6 mos 1 day",
			},
		},
		{
			yaml: `
  expire_in: 'never'`,
			want: Artifacts{
				ExpireIn: "never",
			},
		},
		{
			yaml: `
  expose_as: 'artifact 1'
  paths: ['file.txt']`,
			want: Artifacts{
				ExposeAs: "artifact 1",
				Paths:    Stringorslice{"file.txt"},
			},
		},
		{
			yaml: `
  name: "job1-artifacts-file"
  paths: ["binaries/"]`,
			want: Artifacts{
				Name:  "job1-artifacts-file",
				Paths: Stringorslice{"binaries/"},
			},
		},
		{
			yaml: `
  public: false`,
			want: Artifacts{
				Public: pointerToBool(false),
			},
		},
		{
			yaml: `
  public: true`,
			want: Artifacts{
				Public: pointerToBool(true),
			},
		},
		{
			yaml: `
  reports:
    junit: rspec.xml`,
			want: Artifacts{
				Reports: map[string]interface{}{
					"junit": "rspec.xml",
				},
			},
		},
		{
			yaml: `
  reports:
    type1: output1.xml
    type2: output2.xml`,
			want: Artifacts{
				Reports: map[string]interface{}{
					"type1": "output1.xml",
					"type2": "output2.xml",
				},
			},
		},
		{
			yaml: `
  untracked: true`,
			want: Artifacts{
				Untracked: true,
			},
		},
		{
			yaml: `
  untracked: false`,
			want: Artifacts{
				Untracked: false,
			},
		},
		{
			yaml: `
  paths:
  - binaries/`,
			want: Artifacts{
				Paths:     []string{"binaries/"},
				Untracked: false, // default value if not defined
			},
		},
		{
			yaml: `
  when: on_failure`,
			want: Artifacts{
				When: "on_failure",
			},
		},
		{
			yaml: `
  when: always`,
			want: Artifacts{
				When: "always",
			},
		},
	}

	for i, test := range tests {
		got := new(Artifacts)
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

func TestArtifacts_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[]"), new(Artifacts))
	if err == nil || !strings.Contains(err.Error(), "cannot unmarshal") {
		t.Errorf("Expect error, got %s", err)
	}
}

func pointerToBool(b bool) *bool {
	return &b
}

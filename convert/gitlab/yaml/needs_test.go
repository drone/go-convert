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

func TestNeeds(t *testing.T) {
	tests := []struct {
		yaml string
		want Needs
	}{
		{
			yaml: `"linux:build"`,
			want: Needs{
				Job: "linux:build",
			},
		},
		{
			yaml: `{ "job": "linux:build", "ref": "main" }`,
			want: Needs{
				Job: "linux:build",
				Ref: "main",
			},
		},
		{
			yaml: `{ "job": "linux:build", "ref": "main", "artifacts": true, "optional": true }`,
			want: Needs{
				Job:       "linux:build",
				Ref:       "main",
				Artifacts: true,
				Optional:  true,
			},
		},
	}

	for i, test := range tests {
		got := new(Needs)
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

func TestNeeds_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[]"), new(Needs))
	if err == nil || err.Error() != "failed to unmarshal needs" {
		t.Errorf("Expect error, got %s", err)
	}
}

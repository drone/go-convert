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

func TestStages(t *testing.T) {
	tests := []struct {
		yaml string
		want Stages
	}{
		// test string value
		{
			yaml: `"unit tests"`,
			want: Stages{
				Items: []*Stage{
					{Name: "unit tests"},
				},
			},
		},
		// test string slice value
		{
			yaml: `[ "unit tests" ]`,
			want: Stages{
				Items: []*Stage{
					{Name: "unit tests"},
				},
			},
		},
		// test map value
		{
			yaml: `{ name: "unit tests", "if": "branch = master" }`,
			want: Stages{
				Items: []*Stage{
					{Name: "unit tests", If: "branch = master"},
				},
			},
		},
		// test map slice value
		{
			yaml: `[{ name: "unit tests", "if": "branch = master" }]`,
			want: Stages{
				Items: []*Stage{
					{Name: "unit tests", If: "branch = master"},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Stages)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Log(test.yaml)
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Log(test.yaml)
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestStage_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Stage))
	if err == nil || err.Error() != "failed to unmarshal stage" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestStages_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Stages))
	if err == nil || err.Error() != "failed to unmarshal stages" {
		t.Errorf("Expect error, got %s", err)
	}
}

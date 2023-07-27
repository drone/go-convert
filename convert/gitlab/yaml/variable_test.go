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

func TestVariable(t *testing.T) {
	tests := []struct {
		yaml string
		want Variable
	}{
		{
			yaml: `"https://example.com/"`,
			want: Variable{
				Value:  "https://example.com/",
				Expand: toPointerBool(true),
			},
		},
		{
			yaml: `{ "description": "The deployment note." }`,
			want: Variable{
				Desc:   "The deployment note.",
				Expand: toPointerBool(true),
			},
		},
		{
			yaml: `{ "description": "The deployment target.", "value": "staging" }`,
			want: Variable{
				Value:  "staging",
				Desc:   "The deployment target.",
				Expand: toPointerBool(true),
			},
		},
		{
			yaml: `{ "options": [ "production", "staging", "development" ], "value": "staging" }`,
			want: Variable{
				Value:   "staging",
				Options: []string{"production", "staging", "development"},
				Expand:  toPointerBool(true),
			},
		},
		{
			yaml: `{ "value": "value3 $VAR1", "expand": false }`,
			want: Variable{
				Value:  "value3 $VAR1",
				Expand: toPointerBool(false),
			},
		},
	}

	for i, test := range tests {
		got := new(Variable)
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

func TestVariable_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[]"), new(Variable))
	if err == nil || err.Error() != "failed to unmarshal variable" {
		t.Errorf("Expect error, got %s", err)
	}
}

func toPointerBool(b bool) *bool {
	return &b
}

// variables:
//   DEPLOY_NOTE:
//     description: "The deployment note. Explain the reason for this deployment."

// variables:
//   DEPLOY_ENVIRONMENT:
//     value: "staging"
//     description: "The deployment target. Change this variable to 'canary' or 'production' if needed."

// variables:
//   DEPLOY_ENVIRONMENT:
//     value: "staging"
//     options:
//       - "production"
//       - "staging"
//       - "canary"
//     description: "The deployment target. Set to 'staging' by default."

// variables:
//   VAR1: value1
//   VAR2: value2 $VAR1
//   VAR3:
//     value: value3 $VAR1
//     expand: false

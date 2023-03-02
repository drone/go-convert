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

func TestRetry(t *testing.T) {
	tests := []struct {
		yaml string
		want Retry
	}{
		{
			yaml: `3`,
			want: Retry{
				Max: 3,
			},
		},
		{
			yaml: `{ "max": 3 }`,
			want: Retry{
				Max: 3,
			},
		},
		{
			yaml: `{ "max": 3, "when": "always" }`,
			want: Retry{
				Max:  3,
				When: []string{"always"},
			},
		},
		{
			yaml: `{ "max": 3, "when": [ "unknown_failure", "script_failure" ] }`,
			want: Retry{
				Max:  3,
				When: []string{"unknown_failure", "script_failure"},
			},
		},
	}

	for i, test := range tests {
		got := new(Retry)
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

func TestRetry_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[]"), new(Retry))
	if err == nil || err.Error() != "failed to unmarshal retry" {
		t.Errorf("Expect error, got %s", err)
	}
}

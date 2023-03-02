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

func TestEnv(t *testing.T) {
	tests := []struct {
		yaml string
		want Env
	}{
		// test string in key=value format
		{
			yaml: `"FOO=foo"`,
			want: Env{
				Global: []map[string]string{
					{"FOO": "foo"},
				},
			},
		},
		// test map
		{
			yaml: `{ FOO: foo }`,
			want: Env{
				Global: []map[string]string{
					{"FOO": "foo"},
				},
			},
		},
		// test map array
		{
			yaml: `[ { FOO: foo } ]`,
			want: Env{
				Global: []map[string]string{
					{"FOO": "foo"},
				},
			},
		},
		// test slice array of key=value items
		{
			yaml: `[ FOO=foo, BAR=bar ]`,
			want: Env{
				Global: []map[string]string{
					{"FOO": "foo", "BAR": "bar"},
				},
			},
		},
		// test env array with different value types
		// first value type is a string in key=value format
		// second value type is slice of map values
		{
			yaml: `[ "FOO=foo", { BAR: bar } ]`,
			want: Env{
				Global: []map[string]string{
					{"FOO": "foo"},
					{"BAR": "bar"},
				},
			},
		},
		// test env struct
		{
			yaml: `{ global: [ FOO: foo ], jobs: [ BAR: bar ] }`,
			want: Env{
				Global: []map[string]string{
					{"FOO": "foo"},
				},
				Jobs: []map[string]string{
					{"BAR": "bar"},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Env)
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

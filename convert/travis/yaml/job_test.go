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

func TestJobs(t *testing.T) {
	tests := []struct {
		yaml string
		want Jobs
	}{
		// test map value
		{
			yaml: `{ language: ruby, os: linux, dist: trusty }`,
			want: Jobs{
				Include: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
			},
		},
		// test map slice value
		{
			yaml: `[{ language: ruby, os: linux, dist: trusty }]`,
			want: Jobs{
				Include: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
			},
		},
		// test advanced
		{
			yaml: `{ include: [{ language: ruby, os: linux, dist: trusty }], exclude: [{ language: ruby, os: linux, dist: trusty }], allow_failures: [{ language: ruby, os: linux, dist: trusty }] }`,
			want: Jobs{
				Include: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
				Exclude: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
				AllowFailures: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
			},
		},
		// test fast_finish
		{
			yaml: `{ include: [{ language: ruby, os: linux, dist: trusty }], fast_finish: true }`,
			want: Jobs{
				Include: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
				FastFinish: true,
			},
		},
		// test fast_finish alias
		{
			yaml: `{ include: [{ language: ruby, os: linux, dist: trusty }], fast_failure: true }`,
			want: Jobs{
				Include: []map[string]string{
					{"language": "ruby", "os": "linux", "dist": "trusty"},
				},
				FastFinish: true,
			},
		},
	}

	for i, test := range tests {
		got := new(Jobs)
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

func TestJobs_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Jobs))
	if err == nil || err.Error() != "failed to unmarshal jobs" {
		t.Errorf("Expect error, got %s", err)
	}
}

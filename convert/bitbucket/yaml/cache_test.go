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

	"gopkg.in/yaml.v3"
)

func TestCache_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Cache))
	if err == nil || err.Error() != "failed to unmarshal cache" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestCache_Marshal(t *testing.T) {
	tests := []struct {
		before Cache
		after  string
	}{
		{
			before: Cache{
				Path: "node_modules",
			},
			after: "node_modules\n",
		},
		{
			before: Cache{
				Path: "node_modules",
				Key: &CacheKey{
					Files: []string{"package.json"},
				},
			},
			after: "key:\n    files:\n        - package.json\npath: node_modules\n",
		},
	}

	for _, test := range tests {
		after, err := yaml.Marshal(&test.before)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := string(after), test.after; got != want {
			t.Errorf("want yaml %q, got %q", want, got)
		}
	}
}

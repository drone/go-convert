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

package downgrader

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestConvert(t *testing.T) {
	tests, err := filepath.Glob("testdata/*.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// convert the yaml file from github to harness
			downgrader := New()
			tmp1, err := downgrader.DowngradeFile(test)
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the converted yaml file to a map
			got := map[string]interface{}{}
			if err := yaml.Unmarshal(tmp1, &got); err != nil {
				t.Error(err)
				return
			}

			got = normalizeMap(got)

			// parse the golden yaml file
			data, err := os.ReadFile(test + ".golden")
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the golden yaml file to a map
			want := map[string]interface{}{}
			if err := yaml.Unmarshal(data, &want); err != nil {
				t.Error(err)
				return
			}

			want = normalizeMap(want)

			// compare the converted yaml to the golden file
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Unexpected conversion result")
				t.Log(diff)
			}
		})
	}
}

func normalizeMap(m map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(m))
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]
		switch t := v.(type) {
		case map[string]interface{}:
			normalized[k] = normalizeMap(t)
		case map[interface{}]interface{}:
			normalizedMap := make(map[string]interface{}, len(t))
			for k, v := range t {
				normalizedMap[fmt.Sprintf("%v", k)] = v
			}
			normalized[k] = normalizeMap(normalizedMap)
		case []interface{}:
			normalized[k] = normalizeSlice(t)
		default:
			normalized[k] = v
		}
	}

	return normalized
}

func normalizeSlice(s []interface{}) []interface{} {
	for i, v := range s {
		if m, ok := v.(map[string]interface{}); ok {
			s[i] = normalizeMap(m)
		}
	}
	sort.SliceStable(s, func(i, j int) bool {
		mi, oki := s[i].(map[string]interface{})
		mj, okj := s[j].(map[string]interface{})
		if !oki || !okj {
			// At least one of the values is not a map, so don't attempt to sort
			return false
		}
		ni, inj := mi["name"].(string)
		nj, inj := mj["name"].(string)
		if inj {
			// Both maps have a string name field, so sort by these
			return ni < nj
		}
		// At least one map doesn't have a string name field, so don't attempt to sort
		return false
	})
	return s
}

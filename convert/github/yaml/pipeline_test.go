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
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestPipeline(t *testing.T) {
	tests, err := filepath.Glob("testdata/*/*.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {

		switch test {
		case "testdata/matrix/example-3.yaml",
			"testdata/matrix/example-8.yaml":
			// skip these tests due to unsupported syntax
			// TODO these should be eventually re-enabled
			continue
		}

		t.Run(test, func(t *testing.T) {

			// parse the yaml file
			tmp1, err := ParseFile(test)
			if err != nil {
				t.Error(err)
				return
			}

			// marshal the yaml file
			tmp2, err := yaml.Marshal(tmp1)
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the yaml file to a map
			got := map[string]interface{}{}
			if err := yaml.Unmarshal(tmp2, &got); err != nil {
				t.Error(err)
				return
			}

			// parse the golden yaml file and unmarshal
			data, err := ioutil.ReadFile(test + ".golden")
			if err != nil {
				// skip tests with no golden files
				// TODO these should be re-enabled
				return
			}

			// unmarshal the golden yaml file
			want := map[string]interface{}{}
			if err := yaml.Unmarshal(data, &want); err != nil {
				t.Error(err)
				return
			}

			// compare the parsed yaml to the golden file
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Unexpected parsing result")
				t.Log(diff)
			}
		})
	}
}

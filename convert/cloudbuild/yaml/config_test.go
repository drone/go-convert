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
	tests, err := filepath.Glob("testdata/*.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			temp1, err := ioutil.ReadFile(test)
			if err != nil {
				t.Error(err)
				return
			}

			// parse the yaml file
			temp2, err := ParseBytes(temp1)
			if err != nil {
				t.Error(err)
				return
			}

			// marshal the yaml file
			temp3, err := yaml.Marshal(temp2)
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the yaml file to a map
			got := map[string]interface{}{}
			if err := yaml.Unmarshal(temp3, &got); err != nil {
				t.Error(err)
				return
			}

			// unmarshal the json file
			want := map[string]interface{}{}
			if err := yaml.Unmarshal(temp1, &want); err != nil {
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

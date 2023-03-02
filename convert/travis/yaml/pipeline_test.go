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
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestPipeline(t *testing.T) {
	tests := []string{
		"testdata/addons/browserstack",
		"testdata/addons/sonar",
		"testdata/golang",
		"testdata/rust",
		"testdata/scala",
		"testdata/smalltalk",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// parse the yaml file
			tmp1, err := ParseFile(test + ".yaml")
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

			// parse the golden json file and unmarshal
			data, err := ioutil.ReadFile(test + ".json")
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the json file
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

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

package azure

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestConvert(t *testing.T) {
	tests, err := filepath.Glob("testdata/*/*.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// convert the yaml file from azure to harness
			converter := New()
			tmp1, err := converter.ConvertFile(test)
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

			// parse the golden yaml file
			data, err := ioutil.ReadFile(test + ".golden")
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

			// compare the converted yaml to the golden file
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Unexpected conversion result")
				t.Log(diff)
			}
		})
	}
}

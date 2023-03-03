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

package bitbucket

import (
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestConvert(t *testing.T) {
	// tests, err := filepath.Glob("testdata/*/*.yaml")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// TODO use glob once we have more test cases complete
	tests := []string{
		"testdata/clone/example1.yaml",
		"testdata/clone/example2.yaml",
		"testdata/clone/example3.yaml", // TODO LFS
		"testdata/clone/example4.yaml", // TODO convert BITBUCKET_ variables
		"testdata/clone/example5.yaml",
		"testdata/clone/example6.yaml",
		"testdata/clone/example7.yaml",
		"testdata/clone/example8.yaml", // TODO LFS
		"testdata/clone/example9.yaml", // TODO LFS
		"testdata/clone/example10.yaml",

		"testdata/definitions/example1.yaml",
		// "testdata/definitions/example2.yaml", // TODO handle non-default pipelines
		"testdata/definitions/example3.yaml",
		"testdata/definitions/example4.yaml",
		"testdata/definitions/example5.yaml",
		"testdata/definitions/example6.yaml",
		"testdata/definitions/example7.yaml",
		"testdata/definitions/example8.yaml",
		"testdata/definitions/example9.yaml",

		"testdata/global/example1.yaml",
		"testdata/global/example2.yaml",
		"testdata/global/example3.yaml",
		"testdata/global/example4.yaml",
		"testdata/global/example5.yaml",

		"testdata/image/example3.yaml",
		// "testdata/image/example4.yaml", // username, password
		"testdata/image/example5.yaml",
		"testdata/image/example6.yaml",
		"testdata/image/example7.yaml",
		// "testdata/image/example8.yaml", // username, password
		// "testdata/image/example9.yaml", // username, password
		// "testdata/image/example10.yaml", // username, password
		// "testdata/image/example11.yaml", // username, password
		// "testdata/image/example12.yaml", // services, username, password
		"testdata/image/example13.yaml",
		"testdata/image/example14.yaml",
		"testdata/image/example15.yaml",
		// "testdata/image/example16.yaml", // aws
		// "testdata/image/example17.yaml", // aws
		// "testdata/image/example18.yaml", // aws
		// "testdata/image/example19.yaml", // aws
		// "testdata/image/example20.yaml", // aws, services

		"testdata/parallel/example1.yaml",
		"testdata/parallel/example2.yaml",
		"testdata/parallel/example3.yaml", // TODO fail-fast
		"testdata/parallel/example4.yaml", // TODO fail-fast

		"testdata/stages/example1.yaml",
		"testdata/stages/example2.yaml",
		"testdata/stages/example3.yaml", // TODO changesets
		"testdata/stages/example4.yaml", // TODO changesets
		"testdata/stages/example5.yaml", // TODO deploy, trigger
		"testdata/stages/example6.yaml",
		"testdata/stages/example7.yaml", // TODO trigger

		"testdata/steps/example1.yaml",
		"testdata/steps/example2.yaml",
		"testdata/steps/example3.yaml",
		"testdata/steps/example4.yaml",
		"testdata/steps/example5.yaml",
		"testdata/steps/example6.yaml",
		"testdata/steps/example7.yaml",
		"testdata/steps/example8.yaml",
		"testdata/steps/example9.yaml",  // TODO trigger, deploy
		"testdata/steps/example10.yaml", // TODO oidc
		"testdata/steps/example11.yaml", // TODO trigger
		"testdata/steps/example12.yaml", // TODO artifacts
		"testdata/steps/example13.yaml", // TODO changeset
		"testdata/steps/example14.yaml", // TODO changeset
		"testdata/steps/example15.yaml",
		"testdata/steps/example16.yaml",
		"testdata/steps/example17.yaml",
		"testdata/steps/example18.yaml", // TODO artifacts
		"testdata/steps/example19.yaml", // TODO artifacts
		"testdata/steps/example20.yaml",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// convert the yaml file from bitbucket to harness
			// tmp1, err := FromFile(test)
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }

			// convert the yaml file from bitbucket to harness
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

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

func TestHomebrew(t *testing.T) {
	tests := []struct {
		yaml string
		want Homebrew
	}{
		// string value
		{
			yaml: `true`,
			want: Homebrew{
				Update: true,
			},
		},
		// string value
		{
			yaml: `"beanstalk"`,
			want: Homebrew{
				Update:   true,
				Packages: []string{"beanstalk"},
			},
		},
		// string slice value
		{
			yaml: `["beanstalk"]`,
			want: Homebrew{
				Update:   true,
				Packages: []string{"beanstalk"},
			},
		},
		// struct value
		{
			yaml: `{packages: ["beanstalk"], taps: homebrew/cask-versions, casks: java8 }`,
			want: Homebrew{
				Update:   false,
				Packages: []string{"beanstalk"},
				Taps:     []string{"homebrew/cask-versions"},
				Casks:    []string{"java8"},
			},
		},
		// struct value with brewfile bool
		{
			yaml: `{packages: ["beanstalk"], brewfile: true }`,
			want: Homebrew{
				Update:   false,
				Packages: []string{"beanstalk"},
				Brewfile: "Brewfile",
			},
		},
		// struct value with brewfile path
		{
			yaml: `{packages: ["beanstalk"], brewfile: Brewfile.travis }`,
			want: Homebrew{
				Update:   false,
				Packages: []string{"beanstalk"},
				Brewfile: "Brewfile.travis",
			},
		},
	}

	for i, test := range tests {
		got := new(Homebrew)
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

func TestHomebrew_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Homebrew))
	if err == nil || err.Error() != "failed to unmarshal homebrew" {
		t.Errorf("Expect error, got %s", err)
	}
}

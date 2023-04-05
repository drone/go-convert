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

func TestCache(t *testing.T) {
	tests := []struct {
		yaml string
		want Cache
	}{
		{
			yaml: `true`,
			want: Cache{
				Timeout: 3,
			},
		},
		{
			yaml: `{}`,
			want: Cache{
				Timeout: 3,
			},
		},
		{
			yaml: `{timeout: 42, directories: [ /go, /node ]}`,
			want: Cache{
				Timeout:     42,
				Directories: Stringorslice{"/go", "/node"},
			},
		},
		{
			yaml: `"apt"`,
			want: Cache{
				Timeout: 3,
				Apt:     true,
			},
		},
		{
			yaml: `"bundler"`,
			want: Cache{
				Timeout: 3,
				Bundler: true,
			},
		},
		{
			yaml: `"cargo"`,
			want: Cache{
				Timeout: 3,
				Cargo:   true,
			},
		},
		{
			yaml: `"ccache"`,
			want: Cache{
				Timeout: 3,
				Ccache:  true,
			},
		},
		{
			yaml: `"cocoapods"`,
			want: Cache{
				Timeout:   3,
				Cocoapods: true,
			},
		},
		{
			yaml: `"npm"`,
			want: Cache{
				Timeout: 3,
				Npm:     true,
			},
		},
		{
			yaml: `"packages"`,
			want: Cache{
				Timeout:  3,
				Packages: true,
			},
		},
		{
			yaml: `"pip"`,
			want: Cache{
				Timeout: 3,
				Pip:     true,
			},
		},
		{
			yaml: `"yarn"`,
			want: Cache{
				Timeout: 3,
				Yarn:    true,
			},
		},
		{
			yaml: `"edge"`,
			want: Cache{
				Timeout: 3,
				Edge:    true,
			},
		},
	}

	for i, test := range tests {
		got := new(Cache)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestCache_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Cache))
	if err == nil || err.Error() != "failed to unmarshal cache" {
		t.Errorf("Expect error, got %s", err)
	}
}

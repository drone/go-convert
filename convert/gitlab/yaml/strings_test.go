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

func TestStrings(t *testing.T) {
	tests := []struct {
		yaml string
		want Stringorslice
	}{
		{
			yaml: `"hello"`,
			want: Stringorslice{"hello"},
		},
		{
			yaml: `[ "hello", "world" ]`,
			want: Stringorslice{"hello", "world"},
		},
		{
			yaml: `[ ]`,
			want: nil,
		},
	}

	for i, test := range tests {
		got := new(Stringorslice)
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

func TestStrings_Error(t *testing.T) {
	if err := yaml.Unmarshal([]byte("{}"), new(Stringorslice)); err == nil {
		t.Errorf("Expect error when unmarshaling a into a string slice")
	}

	if err := yaml.Unmarshal([]byte("[1]"), new(Stringorslice)); err == nil {
		t.Errorf("Expect error when unmarshaling a float slice into a string slice")
	}
}

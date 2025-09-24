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

package normalize

import "testing"

func TestNamer(t *testing.T) {
	tests := []struct {
		inputs []string
		output string
	}{
		{
			inputs: []string{"", "", "run"},
			output: "run",
		},
		{
			inputs: []string{"", "Unit Test", "run"},
			output: "unittest",
		},
		{
			inputs: []string{"test", "Unit Test", "run"},
			output: "test",
		},
	}

	for _, test := range tests {
		gen := newGenerator()
		if got, want := gen.generate(test.inputs...), test.output; got != want {
			t.Errorf("got name %s, want %s", got, want)
		}
	}
}

func TestNamerUnique(t *testing.T) {
	tests := []struct {
		inputs []string
		output string
	}{
		{
			inputs: []string{"run"},
			output: "run",
		},
		{
			inputs: []string{"run"},
			output: "run1",
		},
		{
			inputs: []string{"run"},
			output: "run2",
		},
	}

	gen := newGenerator()
	for _, test := range tests {
		if got, want := gen.generate(test.inputs...), test.output; got != want {
			t.Errorf("got name %s, want %s", got, want)
		}
	}
}

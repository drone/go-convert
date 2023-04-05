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

import "testing"

func TestName(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		{"Foo", "Foo"},
		{"Foo 1", "Foo 1"},
		{"Foo ", "Foo"},
		{"Foo-Bar", "Foo-Bar"},
		{"Foo_Bar", "Foo_Bar"},
		{"Foo~Bar", "FooBar"},
		{"_FooBar", "_FooBar"},
		{"-FooBar", "FooBar"},
		{" FooBar", "FooBar"},
		{"~FooBar", "FooBar"},
	}

	for _, test := range tests {
		t.Run(test.before, func(t *testing.T) {
			got, want := convertName(test.before), test.after
			if got != want {
				t.Errorf("Want name %q, got %q", want, got)
			}
		})
	}
}

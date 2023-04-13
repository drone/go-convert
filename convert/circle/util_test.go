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

package circle

import "testing"

func TestSplitOrb(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		command string
	}{
		{
			name:    "node/install",
			alias:   "node",
			command: "install",
		},
		{
			name:    "node",
			alias:   "node",
			command: "",
		},
	}
	for _, test := range tests {
		alias, command := splitOrb(test.name)
		if got, want := alias, test.alias; got != want {
			t.Errorf("Got alias %s want %s", got, want)
		}
		if got, want := command, test.command; got != want {
			t.Errorf("Got command %s want %s", got, want)
		}
	}
}

func TestSplitOrbVersion(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		version string
	}{
		{
			name:    "node@1.0.0",
			alias:   "node",
			version: "1.0.0",
		},
		{
			name:    "node",
			alias:   "node",
			version: "",
		},
	}
	for _, test := range tests {
		alias, command := splitOrbVersion(test.name)
		if got, want := alias, test.alias; got != want {
			t.Errorf("Got alias %s want %s", got, want)
		}
		if got, want := command, test.version; got != want {
			t.Errorf("Got version %s want %s", got, want)
		}
	}
}

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
	"fmt"

	bitbucket "github.com/drone/go-convert/convert/bitbucket/yaml"
)

type state struct {
	config *bitbucket.Config
	stage  *bitbucket.Stage
	steps  *bitbucket.Steps
	step   *bitbucket.Step
	script *bitbucket.Script

	names map[string]struct{}
}

// reset resets the converter.
func (s *state) reset() {
	s.config = nil
	s.stage = nil
	s.steps = nil
	s.step = nil
	s.script = nil
	s.names = map[string]struct{}{}
}

// helper function to generate a unique name.
func (s *state) generateName(name, kind string) string {
	// if the name is empty use the type of step
	// as the name.
	if name == "" {
		name = kind
	}

	// if the name is unused we can return the
	// name as-is.
	if _, ok := s.names[name]; !ok {
		s.names[name] = struct{}{}
		return name
	}

	// iterate on the name until we find a
	//unique name.
	for i := 0; ; i++ {
		// append a sequential int to the suffix of the
		// name for uniqueness.
		temp := name + fmt.Sprint(i)
		// if the name is not in the set, add the name
		// to the set and return to the caller.
		if _, ok := s.names[temp]; !ok {
			s.names[temp] = struct{}{}
			return temp
		}
	}
}

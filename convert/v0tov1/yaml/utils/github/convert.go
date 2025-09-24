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

package github

import (
	"github.com/bradrydzewski/spec/yaml"
)

// Convert converts relevant GitHub structures in the yaml
// to the equivalent Harness structures, if applicable.
func Convert(in *yaml.Schema) error {

	// check to see if the yaml is using github actions
	// syntax and convert as needed.
	if len(in.Jobs) > 0 {
		if in.Pipeline == nil {
			in.Pipeline = new(yaml.Pipeline)
		}
		for name, stage := range in.Jobs {
			stage.Name = name
			in.Pipeline.Stages = append(in.Pipeline.Stages, stage)
		}
		// unset jobs once converted to stages.
		in.Jobs = nil
	}

	return nil
}

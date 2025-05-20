// Copyright 2023 Harness, Inc.
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

package json

import (
	"log"

	harness "github.com/drone/spec/dist/go"
)

// ConvertReadTrusted converts a Jenkins readTrusted step to a Harness step
func ConvertReadTrusted(node Node) *harness.Step {
	// Get the path parameter which is required
	path, ok := node.ParameterMap["path"].(string)
	if !ok {
		log.Printf("readTrusted: missing required parameter 'path'")
		return nil
	}

	// Create the step with required parameters
	step := &harness.Step{
		Name: "read_trusted",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/read-trusted",
			With: map[string]interface{}{
				"file_path":      path,
				"trusted_branch": "<+input>",
				"git_pat":        "<+input>",
			},
		},
	}

	return step
}

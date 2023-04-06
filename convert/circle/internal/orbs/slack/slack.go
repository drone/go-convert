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

package slack

import (
	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// Convert converts an Orb step to a Harness step.
func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "":
		return nil // not supported
	case "notify":
		return convertNotify(step)
	default:
		return convertNotify(step)
	}
}

// helper function converts a slack/notify
// orb to a run step.
func convertNotify(step *circle.Custom) *harness.Step {
	return &harness.Step{
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/slack",
			With:  map[string]interface{}{
				// convert step.Params here
			},
		},
	}
}

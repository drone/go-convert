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
	case "on-hold":
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
	customMessage, _ := step.Params["custom"].(string)
	channel, _ := step.Params["channel"].(string)
	accessToken, _ := step.Params["access_token"].(string)
	mentions, _ := step.Params["mentions"].(string)
	template, _ := step.Params["template"].(string)
	webhook, _ := step.Params["webhook"].(string)
	message, _ := step.Params["message"].(string)
	color, _ := step.Params["color"].(string)

	withMap := map[string]interface{}{}

	if channel != "" {
		withMap["channel"] = channel
	}
	if color != "" {
		withMap["color"] = color
	}
	if message != "" {
		withMap["message"] = message
	}
	if accessToken != "" {
		withMap["access.token"] = accessToken
	}
	if webhook != "" {
		withMap["webhook"] = webhook
	}
	if customMessage != "" {
		withMap["custom.block"] = customMessage
	}
	if mentions != "" {
		withMap["recipient"] = mentions
	}
	if template != "" {
		withMap["template"] = template
	}

	return &harness.Step{
		Type: "plugin",
		Name: "notify_slack",
		Spec: &harness.StepPlugin{
			Image: "plugins/slack",
			With:  withMap,
		},
	}
}

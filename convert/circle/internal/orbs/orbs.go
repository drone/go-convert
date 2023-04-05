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

package orbs

import (
	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/convert/circle/internal/orbs/node"
	"github.com/drone/go-convert/convert/circle/internal/orbs/slack"
)

// Convert converts an Orb step to a Harness step.
func Convert(name, command string, step *circle.Custom) *harness.Step {
	switch name {
	case "circleci/slack":
		return slack.Convert(command, step)
	case "circleci/node":
		return node.Convert(command, step)
	default:
		return nil
	}
}

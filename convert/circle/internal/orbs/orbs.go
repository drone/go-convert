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
	browser_tools "github.com/drone/go-convert/convert/circle/internal/orbs/browser-tools"
	"github.com/drone/go-convert/convert/circle/internal/orbs/datadog"
	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/convert/circle/internal/orbs/codecov"
	"github.com/drone/go-convert/convert/circle/internal/orbs/coveralls"
	"github.com/drone/go-convert/convert/circle/internal/orbs/golang"
	"github.com/drone/go-convert/convert/circle/internal/orbs/localstack"
	"github.com/drone/go-convert/convert/circle/internal/orbs/node"
	"github.com/drone/go-convert/convert/circle/internal/orbs/ruby"
	"github.com/drone/go-convert/convert/circle/internal/orbs/saucelabs"
	"github.com/drone/go-convert/convert/circle/internal/orbs/slack"
)

// Convert converts an Orb step to a Harness step.
func Convert(name, command, version string, step *circle.Custom) *harness.Step {
	switch name {
	case "codecov/codecov":
		return codecov.Convert(command, step)
	case "coveralls/coveralls":
		return coveralls.Convert(command, step)
	case "circleci/browser-tools":
		return browser_tools.Convert(command, version, step)
	case "circleci/go":
		return golang.Convert(command, step)
	case "circleci/node":
		return node.Convert(command, step)
	case "circleci/ruby":
		return ruby.Convert(command, step)
	case "circleci/slack":
		return slack.Convert(command, step)
	case "datadog/agent":
		return datadog.Convert(command, step)
	case "localstack/platform":
		return localstack.Convert(command, step)
	case "saucelabs/saucectl-run":
		return saucelabs.Convert(command, step)
	default:
		return nil
	}
}

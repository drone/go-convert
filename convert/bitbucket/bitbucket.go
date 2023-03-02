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

// Package bitbucket converts Bitbucket pipelines to Harness pipelines.
package bitbucket

import (
	"bytes"
	"io"
	"os"

	bitbucket "github.com/drone/go-convert/convert/bitbucket/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

// From converts the legacy drone yaml format to the
// unified yaml format.
func From(r io.Reader) ([]byte, error) {
	// unmarshal the bitbucket yaml
	src, err := bitbucket.Parse(r)
	if err != nil {
		return nil, err
	}

	// normalize the yaml and ensure
	// all root-level steps are grouped
	// by stage to simplify conversion.
	bitbucket.Normalize(src)

	// TODO Clone
	// TODO Caches
	// TODO Artifacts
	// TODO Services
	// TODO Conditions
	// TODO FastFail
	// TODO Deployment
	// TODO Trigger

	// create the harness stage spec
	// spec := &harness.StageCI{
	// 	// TODO Clone
	// 	// TODO Repository
	// 	// TODO Delegate
	// 	// TODO Platform
	// 	// TODO Runtime
	// 	// TODO Envs
	// }

	// // convert the global instance size
	// if src.Options != nil && src.Options.Size != "" {
	// 	var size string
	// 	switch src.Options.Size {
	// 	case "2x": // 8GB
	// 		size = "large"
	// 	case "4x": // 16GB
	// 		size = "xlarge"
	// 	case "8x": // 32GB
	// 		size = "xxlarge"
	// 	default: // 4GB
	// 		size = "standard"
	// 	}
	// 	spec.Runtime = &harness.Runtime{
	// 		Type: "cloud",
	// 		Spec: &harness.RuntimeCloud{
	// 			Size: size,
	// 		},
	// 	}
	// }

	// // create the harness stage.
	// stage := &harness.Stage{
	// 	Name: "build",
	// 	Type: "ci",
	// 	Spec: spec,
	// 	// TODO When
	// 	// TODO On
	// }

	// create the harness pipeline
	pipeline := &harness.Pipeline{
		Version: 1,
		Default: convertDefault(src),
		// Stages:  []*harness.Stage{stage},
	}

	// create the converter state
	state := new(state)
	state.names = map[string]struct{}{}
	state.config = src // push the config to the state

	for _, steps := range src.Pipelines.Default {
		// if steps.Parallel != nil {
		// 	// TODO parallel steps
		// 	// TODO fast fail
		// }
		// if steps.Step != nil {
		// 	state.step = steps.Step // push the step to the state
		// 	step := convertSteps(state)
		// 	spec.Steps = append(spec.Steps, step)
		// }
		if steps.Stage != nil {
			// TODO stage
			// TODO fast fail
			state.stage = steps.Stage // push the stage to the state
			stage := convertStage(state)
			pipeline.Stages = append(pipeline.Stages, stage)
		}
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// FromBytes converts the legacy drone yaml format to the
// unified yaml format.
func FromBytes(b []byte) ([]byte, error) {
	return From(
		bytes.NewBuffer(b),
	)
}

// FromString converts the legacy drone yaml format to the
// unified yaml format.
func FromString(s string) ([]byte, error) {
	return FromBytes(
		[]byte(s),
	)
}

// FromFile converts the legacy drone yaml format to the
// unified yaml format.
func FromFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return From(f)
}

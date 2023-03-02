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

// Package gitlab converts Gitlab pipelines to Harness pipelines.
package gitlab

import (
	"bytes"
	"io"
	"os"
	"strings"

	harness "github.com/drone/spec/dist/go"
	gitlab "github.com/drone/spec/dist/go/convert/gitlab/yaml"

	"github.com/ghodss/yaml"
)

// From converts the legacy drone yaml format to the
// unified yaml format.
func From(r io.Reader) ([]byte, error) {
	src, err := gitlab.Parse(r)
	if err != nil {
		return nil, err
	}

	dst := new(harness.Pipeline)

	// TODO handle includes
	// src.Include

	// TODO handle stages
	// src.Stages

	// TODO handle variables
	// src.Variables

	// TODO handle workflow
	// src.Workflow

	// create the harness stage.
	dstStage := &harness.Stage{
		Name: "build",
		Type: "ci",
		// When: convertCond(from.Trigger),
		Spec: &harness.StageCI{
			// Delegate: convertNode(from.Node),
			// Envs:     copyenv(from.Environment),
			// Platform: convertPlatform(from.Platform),
			// Runtime:  convertRuntime(from),
			// Steps:    convertSteps(from),
		},
	}
	dst.Stages = append(dst.Stages, dstStage)

	// iterage through named stages
	for _, stagename := range src.Stages {

		// children steps converted from gitlab to harness.
		var steps []*harness.Step

		// iterate through jobs and find jobs assigned to
		// the stage. skip other stages.
		for jobname, job := range src.Jobs {
			if job == nil {
				continue
			}
			if job.Stage != stagename {
				continue
			}

			// gitlab only supports run steps.
			spec := new(harness.StepExec)

			if job.Image != nil {
				spec.Image = job.Image.Name
				spec.Pull = job.Image.PullPolicy
			} else if src.Default != nil && src.Default.Image != nil {
				spec.Image = src.Default.Image.Name
				spec.Pull = src.Default.Image.PullPolicy
			} else if src.Image != nil {
				spec.Image = src.Image.Name
				spec.Pull = src.Image.PullPolicy
			}

			// aggregate all scripts
			script := append(job.Before)
			script = append(script, job.Script...)
			script = append(script, job.After...)

			// and then combine as a single string
			spec.Run = strings.Join(script, "\n")

			// job.Cache
			// job.Retry
			// job.Services
			// job.Timeout
			// job.Tags
			// job.Secrets

			steps = append(steps, &harness.Step{
				Name: jobname,
				Type: "script",
				Spec: spec,
			})
		}

		// if not steps converted, move to next stage
		if len(steps) == 0 {
			continue
		}

		// if there is a single step, append to the stage.
		if len(steps) == 1 {
			dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, steps[0]) // HACK
			continue
		}

		// else if there are multiple steps, wrap with a parallel
		// step to mirror gitlab behavior.
		group := &harness.Step{
			Name: stagename,
			Type: "parallel",
			Spec: &harness.StepParallel{
				Steps: steps,
			},
		}
		dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, group) // HACK
	}

	out, err := yaml.Marshal(dst)
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

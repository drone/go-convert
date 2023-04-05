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

	gitlab "github.com/drone/go-convert/convert/gitlab/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	config *gitlab.Pipeline
	job    *gitlab.Job
}

// Converter converts a Gitlab pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers

	// config *gitlab.Pipeline
	// job    *gitlab.Job
}

// New creates a new Converter that converts a GitLab
// pipeline to a Harness v1 pipeline.
func New(options ...Option) *Converter {
	d := new(Converter)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = store.New()

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default kubernetes namespace.
	if d.kubeNamespace == "" {
		d.kubeNamespace = "default"
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	src, err := gitlab.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(&context{
		config: src,
	})
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.Convert(
		bytes.NewBuffer(b),
	)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.Convert(
		bytes.NewBufferString(s),
	)
}

// ConvertFile downgrades a v1 pipeline.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

// converts converts a GitLab pipeline to a Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {

	// create the harness pipeline
	dst := &harness.Pipeline{
		Version: 1,
		// Default: convertDefault(d.config),
	}

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
	for _, stagename := range ctx.config.Stages {

		// children steps converted from gitlab to harness.
		var steps []*harness.Step

		// iterate through jobs and find jobs assigned to
		// the stage. skip other stages.
		for jobname, job := range ctx.config.Jobs {
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
			} else if ctx.config.Default != nil && ctx.config.Default.Image != nil {
				spec.Image = ctx.config.Default.Image.Name
				spec.Pull = ctx.config.Default.Image.PullPolicy
			} else if ctx.config.Image != nil {
				spec.Image = ctx.config.Image.Name
				spec.Pull = ctx.config.Image.PullPolicy
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

	// marshal the harness yaml
	out, err := yaml.Marshal(dst)
	if err != nil {
		return nil, err
	}

	return out, nil
}

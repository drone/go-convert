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
	"sort"
	"strconv"
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
			Envs: convertVariables(ctx.config.Variables),
			// Platform: convertPlatform(from.Platform),
			// Runtime:  convertRuntime(from),
			// Steps:    convertSteps(from),
		},
	}
	dst.Stages = append(dst.Stages, dstStage)
	var jobKeys []string
	for jobKey := range ctx.config.Jobs {
		jobKeys = append(jobKeys, jobKey)
	}
	sort.Strings(jobKeys)

	stages := ctx.config.Stages
	if len(stages) == 0 {
		stages = []string{".pre", "build", "test", "deploy", ".post"} // stages don't have to be declared for valid yaml. Default to test
	}

	for name, job := range ctx.config.Jobs {
		if job.Stage == "" {
			job.Stage = "test" // default stage
		} // required for ordering
		switch name {
		case "before_script":
			job.Stage = ".pre"
		case "after_script":
			job.Stage = ".post"
		case "":
			job.Stage = "build"
		}
	}

	for _, stageName := range stages {
		stepGroup := &harness.StepGroup{
			Steps: []*harness.Step{},
		}

		// iterate through jobs and find jobs assigned to
		// the stage. skip other stages.
		for _, jobName := range jobKeys {
			job := ctx.config.Jobs[jobName] // maintaining order here
			if job == nil || job.Stage != stageName {
				continue
			}

			// Convert each job to a step
			step := convertJobToStep(ctx, jobName, job)
			stepGroup.Steps = append(stepGroup.Steps, step...)
		}

		// if not steps converted, move to next stage
		if len(stepGroup.Steps) == 0 {
			continue
		}

		// if there is a single step, append to the stage.
		if len(stepGroup.Steps) == 1 {
			dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, stepGroup.Steps[0])
			continue
		}

		// else if there are multiple steps, wrap with a parallel
		// step to mirror gitlab behavior.
		group := &harness.Step{
			Name: stageName,
			Type: "parallel",
			Spec: &harness.StepParallel{
				Steps: stepGroup.Steps,
			},
		}
		dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, group)
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(dst)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func convertJobToStep(ctx *context, jobName string, job *gitlab.Job) []*harness.Step {
	var steps []*harness.Step
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

	beforeScripts := job.Before
	if len(beforeScripts) == 0 && ctx.config.Default != nil {
		beforeScripts = ctx.config.Default.Before
	}

	afterScripts := job.After
	if len(afterScripts) == 0 && ctx.config.Default != nil {
		afterScripts = ctx.config.Default.After
	}

	// Convert all scripts into a single step
	script := append(beforeScripts)
	script = append(script, job.Script...)
	script = append(script, job.After...)
	script = append(script, afterScripts...)

	spec.Run = strings.Join(script, "\n")

	step := &harness.Step{
		Name: jobName,
		Type: "script",
		Spec: spec,
		On:   convertAllowFailure(job),
	}

	steps = append(steps, step)

	// job.Cache
	// job.Retry
	// job.Services
	// job.Timeout
	// job.Tags
	// job.Secrets

	return steps
}

func convertAllowFailure(job *gitlab.Job) *harness.On {
	if job.AllowFailure != nil && job.AllowFailure.Value {
		var exitCodesStr []string
		for _, code := range job.AllowFailure.ExitCodes {
			exitCodesStr = append(exitCodesStr, strconv.Itoa(code))
		}
		// Sort the slice to maintain order
		sort.Strings(exitCodesStr)

		return &harness.On{
			Failure: &harness.Failure{
				Type:      "ignore",
				ExitCodes: exitCodesStr,
			},
		}
	}
	return nil
}

func convertVariables(variables map[string]*gitlab.Variable) map[string]string {
	result := make(map[string]string)
	var keys []string
	for key := range variables {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		variable := variables[key]
		if variable != nil {
			result[key] = variable.Value
		}
	}

	return result
}

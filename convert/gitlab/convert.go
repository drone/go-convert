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

// ConvertBytes downgrades a v1 pipeline.
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
	}
	cacheFound := false

	// TODO handle includes
	// src.Include

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
			Steps: make([]*harness.Step, 0), // Initialize the Steps slice
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

	for _, stageName := range stages {
		// iterate through jobs and find jobs assigned to the stage. Skip other stages.
		for i, jobName := range jobKeys {
			job := ctx.config.Jobs[jobName] // maintaining order here
			if job.Before != nil {
				job.Stage = ".pre"
			}
			if job.After != nil {
				job.Stage = ".post"
			}
			if !cacheFound && job.Cache != nil {
				dstStage.Spec.(*harness.StageCI).Cache = convertCache(job.Cache) // Update cache if it's defined in the job
				cacheFound = true
			}
			if job == nil || job.Stage != stageName {
				continue
			}

			if job.Before != nil {
				beforeScriptStep := convertScriptToStep(job.Before, "before_script", "", false)
				dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, beforeScriptStep)
			}

			// Convert each job to a step
			steps := convertJobToStep(ctx, jobName, job)

			// Add all steps from the job to the stage
			for _, step := range steps {
				dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, step)
			}

			// If job has an after_script, append it to the steps
			if job.After != nil {
				afterScriptStep := convertScriptToStep(job.After, "after_script", "5m", true)
				dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, afterScriptStep)
			}

			// If first job in the stage, prepend the before_script
			if ctx.config.Default != nil && ctx.config.Default.Before != nil && i == 0 {
				beforeScriptStep := convertScriptToStep(ctx.config.Default.Before, "before_script", "", false)
				dstStage.Spec.(*harness.StageCI).Steps = append([]*harness.Step{beforeScriptStep}, dstStage.Spec.(*harness.StageCI).Steps...)
			}

			// If last job in the stage, append the after_script
			if ctx.config.Default != nil && ctx.config.Default.After != nil && i == len(jobKeys)-1 {
				afterScriptStep := convertScriptToStep(ctx.config.Default.After, "after_script", "5m", true)
				dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, afterScriptStep)
			}
		}
	}
	// marshal the harness yaml
	out, err := yaml.Marshal(dst)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// convertCache converts a GitLab cache to a Harness cache.
func convertCache(cache *gitlab.Cache) *harness.Cache {
	if cache == nil {
		return nil
	}

	return &harness.Cache{
		Enabled: true,
		Key:     cache.Key.Value,
		Paths:   cache.Paths,
		Policy:  cache.Policy,
	}
}

// convertScriptToStep converts a GitLab script to a Harness step.
func convertScriptToStep(script []string, name, timeout string, onFailureIgnore bool) *harness.Step {
	spec := new(harness.StepExec)
	spec.Run = strings.Join(script, "\n")

	step := &harness.Step{
		Name: name,
		Type: "script",
		Spec: spec,
	}
	if timeout != "" {
		step.Timeout = timeout
	}
	if onFailureIgnore {
		step.On = &harness.On{
			Failure: &harness.Failure{
				Type: "ignore",
			},
		}
	}

	return step
}

// convertJobToStep converts a GitLab job to a Harness step.
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

	// Convert all scripts into a single step
	script := append(job.Script)

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

// convertAllowFailure converts a GitLab job's allow_failure to a Harness step's on.failure.
func convertAllowFailure(job *gitlab.Job) *harness.On {
	if job.AllowFailure != nil && job.AllowFailure.Value {
		var exitCodesStr []string
		for _, code := range job.AllowFailure.ExitCodes {
			exitCodesStr = append(exitCodesStr, strconv.Itoa(code))
		}
		// Sort the slice to maintain order
		sort.Strings(exitCodesStr)

		on := &harness.On{
			Failure: &harness.Failure{
				Type: "ignore",
			},
		}
		if len(exitCodesStr) > 0 {
			on.Failure.ExitCodes = exitCodesStr
		}
	}
	return nil
}

// convertVariables converts a GitLab variables map to a Harness variables map.
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

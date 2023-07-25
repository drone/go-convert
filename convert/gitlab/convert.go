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
	stagesLength := len(stages)
	if stagesLength == 0 {
		stages = []string{".pre", "build", "test", "deploy", ".post"} // stages don't have to be declared for valid yaml. Default to test
	}

	for _, stageName := range stages {
		stageSteps := make([]*harness.Step, 0)

		// maintain stage name if set
		if stagesLength != 0 {
			dstStage.Name = stageName
		}

		// iterate through jobs and find jobs assigned to the stage. Skip other stages.
		for _, jobName := range jobKeys {
			job := ctx.config.Jobs[jobName] // maintaining order here
			if job.Before != nil {
				job.Stage = ".pre"
			}
			if job.After != nil {
				job.Stage = ".post"
			}
			if job.Stage == "" {
				job.Stage = "test" // default
			}
			if !cacheFound && job.Cache != nil {
				dstStage.Spec.(*harness.StageCI).Cache = convertCache(job.Cache) // Update cache if it's defined in the job
				cacheFound = true
			}

			if len(job.Extends) > 0 {
				for _, extend := range job.Extends {
					if templateJob, ok := ctx.config.TemplateJobs[extend]; ok {
						// Perform deep merge of the template job into the current job.
						var err error
						job = mergeJobConfiguration(templateJob, job)
						if err != nil {
							return nil, err
						}
					}
				}
			}

			if job == nil || job.Stage != stageName {
				continue
			}
			// Convert each job to a step
			steps := convertJobToStep(ctx, jobName, job)

			for _, step := range steps {
				// Prepend the pipeline-level before_script
				if ctx.config.BeforeScript != nil {
					prependScript := convertScriptToStep(ctx.config.BeforeScript, "", "", false)
					step.Spec.(*harness.StepExec).Run = prependScript.Spec.(*harness.StepExec).Run + "\n" + step.Spec.(*harness.StepExec).Run
				}

				// Prepend the job-specific before_script
				if job.Before != nil {
					prependScript := convertScriptToStep(job.Before, "", "", false)
					step.Spec.(*harness.StepExec).Run = prependScript.Spec.(*harness.StepExec).Run + "\n" + step.Spec.(*harness.StepExec).Run
				}
				stageSteps = append(stageSteps, step)
			}

			if job.Inherit != nil && job.Inherit.Variables != nil {
				dstStage.Spec.(*harness.StageCI).Envs = convertInheritedVariables(job, dstStage.Spec.(*harness.StageCI).Envs)
			}
		}
		// If there are multiple steps, wrap them with a parallel group to mirror gitlab behavior
		if len(stageSteps) > 1 {
			group := &harness.Step{
				Type: "parallel",
				Spec: &harness.StepParallel{
					Steps: stageSteps,
				},
			}
			dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, group)
		} else if len(stageSteps) == 1 {
			// If there's a single step, append it to the stage directly
			dstStage.Spec.(*harness.StageCI).Steps = append(dstStage.Spec.(*harness.StageCI).Steps, stageSteps[0])
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
		spec.Image, spec.Pull = convertImageAndPullPolicy(job.Image)
	} else if job.Inherit == nil || job.Inherit.Default == nil || !job.Inherit.Default.All {
		if ctx.config.Default != nil && ctx.config.Default.Image != nil {
			spec.Image, spec.Pull = convertImageAndPullPolicy(ctx.config.Default.Image)
		} else if ctx.config.Image != nil {
			spec.Image, spec.Pull = convertImageAndPullPolicy(ctx.config.Image)
		}
	}

	if job.Inherit == nil || job.Inherit.Default == nil || job.Inherit.Default.All {
		convertInheritDefaultFields(spec, ctx.config.Default, nil)
	} else {
		convertInheritDefaultFields(spec, ctx.config.Default, job.Inherit.Default.Keys)
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

// convertInheritDefaultFields converts the default fields from the default job into the current job.
func convertInheritDefaultFields(spec *harness.StepExec, defaultJob *gitlab.Default, keys []string) {
	if defaultJob == nil {
		return
	}
	if keys == nil {
		keys = []string{
			"after_script", "before_script", "artifacts", "cache", "image",
			"interruptible", "retry", "services", "tags", "duration",
		}
	}
	for _, key := range keys {
		switch key {
		case "after_script":
			if len(defaultJob.After) > 0 {
				spec.Run = strings.Join(defaultJob.After, "\n")
			}
		case "before_script":
			if len(defaultJob.Before) > 0 {
				spec.Run = strings.Join(defaultJob.Before, "\n") + "\n" + spec.Run
			}
		case "artifacts":
			if defaultJob.Artifacts != nil {
				//TODO no supported
			}
		case "cache":
			if defaultJob.Cache != nil {
				//TODO
			}
		case "image":
			if defaultJob.Image != nil {
				spec.Image, spec.Pull = convertImageAndPullPolicy(defaultJob.Image)
			}
		case "interruptible":
			//TODO not supported
		case "retry":
			if defaultJob.Retry != nil {
				//TODO not supported
			}
		case "services":
			if len(defaultJob.Services) > 0 {
				//TODO
			}
		case "tags":
			if len(defaultJob.Tags) > 0 {
				//spec.Tags = strings.Join(defaultJob.Tags, ", ") //TODO
			}
		case "duration":
			//spec.Timeout = defaultJob.Timeout //TODO
		}
	}
}

// convertInherit converts the inherit fields from the default job into the current job.
func convertInheritedVariables(job *gitlab.Job, stageEnvs map[string]string) map[string]string {
	if job.Inherit == nil || job.Inherit.Variables == nil {
		return stageEnvs
	}

	if job.Inherit.Variables.All {
		return stageEnvs
	}

	// If inherit.variables is an array, only the variables in the array are inherited.
	if job.Inherit.Variables.Keys != nil {
		newEnvs := make(map[string]string)
		for _, key := range job.Inherit.Variables.Keys {
			if value, ok := stageEnvs[key]; ok {
				newEnvs[key] = value
			}
		}
		return newEnvs
	}

	return stageEnvs
}

// convertImageAndPullPolicy converts a GitLab image to a Harness image and pull policy.
func convertImageAndPullPolicy(image *gitlab.Image) (string, string) {
	var name string
	var pullPolicy string

	if image != nil {
		name = image.Name

		if len(image.PullPolicy) == 1 {
			pullPolicyMapping := map[string]string{
				"always":         "always",
				"never":          "never",
				"if-not-present": "if-not-exists",
			}

			pullPolicy = pullPolicyMapping[image.PullPolicy[0]]
		}
	}

	return name, pullPolicy
}

// mergeJobConfiguration merges the child job configuration into the parent job configuration.
func mergeJobConfiguration(child *gitlab.Job, parent *gitlab.Job) *gitlab.Job {
	mergedJob := &gitlab.Job{}

	mergedJob.After = child.After
	if len(mergedJob.After) == 0 {
		mergedJob.After = parent.After
	}

	mergedJob.Artifacts = child.Artifacts
	if mergedJob.Artifacts == nil {
		mergedJob.Artifacts = parent.Artifacts
	}

	mergedJob.AllowFailure = child.AllowFailure
	if mergedJob.AllowFailure == nil {
		mergedJob.AllowFailure = parent.AllowFailure
	}

	mergedJob.Before = child.Before
	if len(mergedJob.Before) == 0 {
		mergedJob.Before = parent.Before
	}

	mergedJob.Cache = child.Cache
	if mergedJob.Cache == nil {
		mergedJob.Cache = parent.Cache
	}

	mergedJob.Coverage = child.Coverage
	if mergedJob.Coverage == "" {
		mergedJob.Coverage = parent.Coverage
	}

	mergedJob.DASTConfiguration = child.DASTConfiguration
	if mergedJob.DASTConfiguration == nil {
		mergedJob.DASTConfiguration = parent.DASTConfiguration
	}

	mergedJob.Dependencies = child.Dependencies
	if len(mergedJob.Dependencies) == 0 {
		mergedJob.Dependencies = parent.Dependencies
	}

	mergedJob.Environment = child.Environment
	if mergedJob.Environment == nil {
		mergedJob.Environment = parent.Environment
	}

	mergedJob.Extends = child.Extends
	if len(mergedJob.Extends) == 0 {
		mergedJob.Extends = parent.Extends
	}

	mergedJob.Image = child.Image
	if mergedJob.Image == nil {
		mergedJob.Image = parent.Image
	}

	mergedJob.Inherit = child.Inherit
	if mergedJob.Inherit == nil {
		mergedJob.Inherit = parent.Inherit
	}

	mergedJob.Interruptible = child.Interruptible
	if !mergedJob.Interruptible {
		mergedJob.Interruptible = parent.Interruptible
	}

	mergedJob.Needs = child.Needs
	if mergedJob.Needs == nil {
		mergedJob.Needs = parent.Needs
	}

	mergedJob.Only = child.Only
	if mergedJob.Only == nil {
		mergedJob.Only = parent.Only
	}

	mergedJob.Pages = child.Pages
	if mergedJob.Pages == nil {
		mergedJob.Pages = parent.Pages
	}

	mergedJob.Parallel = child.Parallel
	if mergedJob.Parallel == nil {
		mergedJob.Parallel = parent.Parallel
	}

	mergedJob.Release = child.Release
	if mergedJob.Release == nil {
		mergedJob.Release = parent.Release
	}

	mergedJob.ResourceGroup = child.ResourceGroup
	if mergedJob.ResourceGroup == "" {
		mergedJob.ResourceGroup = parent.ResourceGroup
	}

	mergedJob.Retry = child.Retry
	if mergedJob.Retry == nil {
		mergedJob.Retry = parent.Retry
	}

	mergedJob.Rules = child.Rules
	if mergedJob.Rules == nil {
		mergedJob.Rules = parent.Rules
	}

	mergedJob.Script = child.Script
	if len(mergedJob.Script) == 0 {
		mergedJob.Script = parent.Script
	}

	mergedJob.Secrets = child.Secrets
	if mergedJob.Secrets == nil {
		mergedJob.Secrets = parent.Secrets
	}

	mergedJob.Services = child.Services
	if mergedJob.Services == nil {
		mergedJob.Services = parent.Services
	}

	mergedJob.Stage = child.Stage
	if mergedJob.Stage == "" {
		mergedJob.Stage = parent.Stage
	}

	mergedJob.Tags = child.Tags
	if len(mergedJob.Tags) == 0 {
		mergedJob.Tags = parent.Tags
	}

	mergedJob.Timeout = child.Timeout
	if mergedJob.Timeout == "" {
		mergedJob.Timeout = parent.Timeout
	}

	mergedJob.Trigger = child.Trigger
	if mergedJob.Trigger == nil {
		mergedJob.Trigger = parent.Trigger
	}

	mergedJob.Variables = child.Variables
	if mergedJob.Variables == nil {
		mergedJob.Variables = parent.Variables
	}

	mergedJob.When = child.When
	if mergedJob.When == "" {
		mergedJob.When = parent.When
	}

	return mergedJob
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

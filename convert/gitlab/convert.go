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
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	gitlab "github.com/drone/go-convert/convert/gitlab/yaml"
	"github.com/drone/go-convert/internal/store"
	harness "github.com/drone/spec/dist/go"

	"dario.cat/mergo"
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

	// create the harness pipeline spec
	dst := &harness.Pipeline{}

	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    dst,
	}

	cacheFound := false

	// TODO handle includes
	// src.Include

	if ctx.config.Workflow != nil {
		// TODO pipeline.name removed from spec
		// dst.Name = ctx.config.Workflow.Name
	}

	// create the harness stage.
	dstStage := &harness.Stage{
		Name: "test",
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
	} else {
		dstStage.Name = ctx.config.Stages[0]
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
						mergedJob, err := mergeJobConfiguration(job, templateJob)
						if err != nil {
							return nil, err
						}
						job = mergedJob
					}
				}
			}

			if job == nil || job.Stage != stageName {
				continue
			}

			if job.Parallel != nil {
				if job.Parallel.Matrix != nil {
					for i, matrix := range job.Parallel.Matrix {
						steps := convertJobToStep(ctx, fmt.Sprintf("%s-%d", jobName, i), job, matrix)
						stageSteps = append(stageSteps, steps...)
					}
				}
			} else {
				// Convert each job to a step
				steps := convertJobToStep(ctx, jobName, job, nil)
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
	out, err := yaml.Marshal(config)
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
		step.Failure = &harness.FailureList{
			Items: []*harness.Failure{
				{
					Action: &harness.FailureAction{
						Type: "ignore",
					},
				},
			},
		}
	}

	return step
}

// convertJobToStep converts a GitLab job to a Harness step.
func convertJobToStep(ctx *context, jobName string, job *gitlab.Job, matrix map[string][]string) []*harness.Step {
	var steps []*harness.Step
	spec := new(harness.StepExec)

	if imageProvided(job.Image) {
		spec = convertImage(job.Image)
	} else if useDefaultImage(job, ctx) {
		spec = convertImage(ctx.config.Default.Image)
	} else if imageProvided(ctx.config.Image) {
		spec = convertImage(ctx.config.Image)
	}

	if job.Inherit == nil || job.Inherit.Default == nil || job.Inherit.Default.All {
		convertInheritDefaultFields(spec, ctx.config.Default, nil)
	} else {
		convertInheritDefaultFields(spec, ctx.config.Default, job.Inherit.Default.Keys)
	}

	// Convert all scripts into a single step
	script := append(job.Script)

	spec.Run = strings.Join(script, "\n")

	var on *harness.FailureList
	if job.Retry != nil {
		on = convertRetry(job)
	} else if job.AllowFailure != nil {
		on = convertAllowFailure(job)
	}

	// set step environment variables
	if job.Variables != nil || job.Secrets != nil || matrix != nil {
		spec.Envs = make(map[string]string)

		// job variables become step variables
		if job.Variables != nil {
			envVariables := convertVariables(job.Variables)
			for key := range envVariables {
				spec.Envs[key] = envVariables[key]
			}
		}

		// job secrets become step variables that reference Harness secrets
		if job.Secrets != nil {
			envSecrets := convertSecrets(job.Secrets)
			for key := range envSecrets {
				spec.Envs[key] = envSecrets[key]
			}
		}

		// job matrix axes become step variables that reference Harness matrix values
		if matrix != nil {
			envMatrix := convertVariablesMatrix(matrix)
			for key := range envMatrix {
				spec.Envs[key] = envMatrix[key]
			}
		}
	}

	var strategy *harness.Strategy
	if matrix != nil {
		strategy = convertStrategy(matrix)
	}

	step := &harness.Step{
		Name: jobName,
		Type: "script",
		Spec: spec,
	}
	// map on if exists
	if on != nil {
		step.Failure = on
	}
	// map strategy if exists
	if strategy != nil {
		step.Strategy = strategy
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

func imageProvided(image *gitlab.Image) bool {
	return image != nil
}

func isInheritAll(job *gitlab.Job) bool {
	return job.Inherit != nil && job.Inherit.Default != nil && job.Inherit.Default.All
}

func useDefaultImage(job *gitlab.Job, ctx *context) bool {
	return !isInheritAll(job) && ctx.config.Default != nil && imageProvided(ctx.config.Default.Image)
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
				spec = convertImage(defaultJob.Image)
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

// convertImage extracts the image name, pull policy, and entrypoint from a GitLab image.
func convertImage(image *gitlab.Image) *harness.StepExec {
	spec := &harness.StepExec{}

	if image != nil {
		spec.Image = image.Name
		if len(image.PullPolicy) == 1 {
			pullPolicyMapping := map[string]string{
				"always":         "always",
				"never":          "never",
				"if-not-present": "if-not-exists",
			}

			spec.Pull = pullPolicyMapping[image.PullPolicy[0]]
		}
		if len(image.Entrypoint) > 0 {
			spec.Entrypoint = image.Entrypoint[0]
			if len(image.Entrypoint) > 1 {
				spec.Args = image.Entrypoint[1:]
			}
		}
	}

	return spec
}

func mergeJobConfiguration(child *gitlab.Job, parent *gitlab.Job) (*gitlab.Job, error) {
	mergedJob := &gitlab.Job{}

	// Copy all fields from the parent job into mergedJob.
	if err := mergo.Merge(mergedJob, parent, mergo.WithOverride); err != nil {
		return nil, err
	}

	// Then, copy all non-empty fields from the child job into mergedJob.
	if err := mergo.Merge(mergedJob, child, mergo.WithOverride); err != nil {
		return nil, err
	}

	return mergedJob, nil
}

// convertAllowFailure converts a GitLab job's allow_failure to a Harness step's on.failure.
func convertAllowFailure(job *gitlab.Job) *harness.FailureList {
	if job.AllowFailure != nil && job.AllowFailure.Value {
		var exitCodesStr []string
		for _, code := range job.AllowFailure.ExitCodes {
			exitCodesStr = append(exitCodesStr, strconv.Itoa(code))
		}
		// Sort the slice to maintain order
		sort.Strings(exitCodesStr)

		on := &harness.FailureList{
			Items: []*harness.Failure{
				{
					Errors: []string{"all"},
					Action: &harness.FailureAction{
						Type: "ignore",
					},
				},
			},
		}
		if len(exitCodesStr) > 0 {
			// TODO exit_code needs to be re-added to spec
			// on.Failure.ExitCodes = exitCodesStr
		}
		return on
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

// convertVariablesMatrix converts a matrix axis map to a Harness variables map.
func convertVariablesMatrix(axis map[string][]string) map[string]string {
	result := make(map[string]string)

	var keys []string
	for k := range axis {
		keys = append(keys, k)
	}
	sort.Strings(keys) // to maintain order

	for axisName := range axis {
		result[axisName] = fmt.Sprintf("<+matrix.%s>", axisName)
	}

	return result
}

// convertRetry converts a GitLab job's retry to a Harness step's on.failure.retry.
func convertRetry(job *gitlab.Job) *harness.FailureList {
	if job.Retry == nil {
		return nil
	}

	return &harness.FailureList{
		Items: []*harness.Failure{
			{
				Action: &harness.FailureAction{
					Type: "retry",
					Spec: harness.Retry{
						Attempts: int64(job.Retry.Max),
					},
				},
			},
		},
	}
}

// convertSecrets converts a GitLab secrets map to a Harness secrets map.
func convertSecrets(secrets map[string]*gitlab.Secret) map[string]string {
	result := make(map[string]string)

	var keys []string
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys) // to maintain order

	for secretName := range secrets {
		result[secretName] = fmt.Sprintf("<+secrets.getValue(\"%s\")>", secretName)
	}

	return result
}

func convertStrategy(axis map[string][]string) *harness.Strategy {
	return &harness.Strategy{
		Type: "matrix",
		Spec: &harness.Matrix{
			Axis: axis,
		},
	}
}

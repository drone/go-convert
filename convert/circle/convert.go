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

// Package circle converts Circle pipelines to Harness pipelines.
package circle

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/drone/go-convert/convert/circle/internal/orbs"
	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// Converter converts a Circle pipeline to a harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	gcsBucket     string
	gcsToken      string
	gcsEnabled    bool
	s3Enabled     bool
	s3Bucket      string
	s3Region      string
	s3AccessKey   string
	s3SecretKey   string
	dockerhubConn string
	identifiers   *store.Identifiers
}

// New creates a new Converter that converts a Circle
// pipeline to a harness v1 pipeline.
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

	// set the storage engine to s3 if configured.
	if d.s3Bucket != "" {
		d.s3Enabled = true
	}

	// set the storage engine to gcs if configured.
	if d.gcsBucket != "" {
		d.gcsEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	src, err := circle.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(src)
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

// converts converts a circle pipeline pipeline.
func (d *Converter) convert(config *circle.Config) ([]byte, error) {

	// create the harness pipeline spec
	pipeline := &harness.Pipeline{}

	// convert pipeline and job parameters to inputs
	if params := extractParameters(config); len(params) != 0 {
		pipeline.Inputs = convertParameters(params)
	}

	// require a minimum of 1 workflows
	if config.Workflows == nil || len(config.Workflows.Items) == 0 {
		return nil, errors.New("no workflows defined")
	}

	// choose the first workflow in the list, for now
	var pipelines []*harness.Pipeline
	for _, workflow := range config.Workflows.Items {
		pipelines = append(pipelines, d.convertPipeline(workflow, config))
	}

	var buf bytes.Buffer
	for i, pipeline := range pipelines {
		// marshal the harness yaml
		out, err := yaml.Marshal(&harness.Config{
			Version: 1,
			Kind:    "pipeline",
			Spec:    pipeline,
		})
		if err != nil {
			return nil, err
		}

		// replace circle parameters with harness parameters
		out = replaceParams(out, params)

		// write the pipeline to the buffer. if there are
		// multiple yaml files they need to be separated
		// by the document separator.
		if i > 0 {
			buf.WriteString("\n")
			buf.WriteString("---")
			buf.WriteString("\n")
		}
		buf.Write(out)
	}

	return buf.Bytes(), nil
}

// converts converts a circle pipeline pipeline.
func (d *Converter) convertPipeline(workflow *circle.Workflow, config *circle.Config) *harness.Pipeline {

	// create the harness pipeline spec
	pipeline := &harness.Pipeline{}

	// convert pipeline and job parameters to inputs
	if params := extractParameters(config); len(params) != 0 {
		pipeline.Inputs = convertParameters(params)
	}

	// loop through workflow jobs and convert each
	// job to a stage.
	for _, workflowjob := range workflow.Jobs {
		// snapshot the config
		config_ := config

		// lookup the named job
		job, ok := config_.Jobs[workflowjob.Name]
		if !ok {
			// if the job does not exist, check to
			// see if the job is an orb.
			alias, command := splitOrb(workflowjob.Name)

			// lookup the orb and silently skip the
			// job if not found
			orb, ok := config_.Orbs[alias]
			if !ok {
				continue
			}

			// HACK (bradrydzewski) this is a temporary
			// hack to create the configuration for an
			// orb referenced directly in the workflow.
			if orb.Inline == nil {
				// config_ = new(circle.Config)
				// config_.Orbs = map[string]*circle.Orb{
				// 	orb.Name: {},
				// }
				job = &circle.Job{
					Steps: []*circle.Step{
						{
							Custom: &circle.Custom{
								Name:   workflowjob.Name,
								Params: workflowjob.Params,
							},
						},
					},
				}
			} else {
				// lookup the orb command and silently skip
				// the job if not found
				job, ok = orb.Inline.Jobs[command]
				if !ok {
					continue
				}

				// replace the config_ with the orb
				config_ = orb.Inline
			}
		}

		// this section replaces circle matrix expressions
		// with harness circle matrix expressions.
		//
		// before: << parameters.foo >>
		// after: << matrix.foo >>
		replaceParamsMatrix(job, workflowjob.Matrix)

		// convert the circle job to a stage and silently
		// skip any stages that cannot be converted.
		stage := d.convertStage(job, config_)
		if stage == nil {
			continue
		}

		stage.Name = workflowjob.Name

		if v := workflowjob.Matrix; v != nil {
			stage.Strategy = convertMatrix(job, v)
		}

		// TODO workflows.[*].triggers
		// TODO workflows.[*].unless
		// TODO workflows.[*].when
		// TODO workflows.[*].jobs[*].context
		// TODO workflows.[*].jobs[*].filters
		// TODO workflows.[*].jobs[*].type
		// TODO workflows.[*].jobs[*].requires

		// append the converted stage to the pipeline.
		pipeline.Stages = append(pipeline.Stages, stage)
	}

	return pipeline
}

// helper function converts Circle job to a Harness stage.
func (d *Converter) convertStage(job *circle.Job, config *circle.Config) *harness.Stage {

	// create stage spec
	spec := &harness.StageCI{
		Envs:     job.Environment,
		Platform: convertPlatform(job, config),
		Runtime:  convertRuntime(job, config),
		Steps: append(
			defaultBackgroundSteps(job, config),
			d.convertSteps(job.Steps, job, config)...,
		),
	}

	// TODO executor.machine
	// TODO executor.shell
	// TODO executor.working_directory

	// if there are no steps in the stage we
	// can skip adding the stage to the pipeline.
	if len(spec.Steps) == 0 {
		return nil
	}

	// TODO job.branches
	// TODO job.parallelism
	// TODO job.parameters

	optimizeCache(spec)
	optimizeGroup(spec)

	// create the stage
	stage := &harness.Stage{}
	stage.Type = "ci"
	stage.Spec = spec
	return stage
}

// helper function converts Circle steps to Harness steps.
func (d *Converter) convertSteps(steps []*circle.Step, job *circle.Job, config *circle.Config) []*harness.Step {
	var out []*harness.Step
	for _, src := range steps {
		if dst := d.convertStep(src, job, config); dst != nil {
			out = append(out, dst)
		}
	}
	return out
}

// helper function converts a Circle step to a Harness step.
func (d *Converter) convertStep(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	switch {
	case step.AddSSHKeys != nil:
		return d.convertAddSSHKeys(step)
	case step.AttachWorkspace != nil:
		return nil // not supported
	case step.Checkout != nil:
		return nil // ignore
	case step.PersistToWorkspace != nil:
		return nil // not supported
	case step.RestoreCache != nil:
		return d.convertRestoreCache(step)
	case step.Run != nil:
		return d.convertRun(step, job, config)
	case step.SaveCache != nil:
		return d.convertSaveCache(step)
	case step.SetupRemoteDocker != nil:
		return nil // not supported
	case step.StoreArtifacts != nil:
		return d.convertStoreArtifacts(step)
	case step.StoreTestResults != nil:
		return d.convertStoreTestResults(step)
	case step.Unless != nil:
		return d.convertUnlessStep(step, job, config)
	case step.When != nil:
		return d.convertWhenStep(step, job, config)
	case step.Custom != nil:
		return d.convertCustom(step, job, config)
	default:
		return nil
	}
}

//
// Step Types
//

// helper function converts a Circle Run step.
func (d *Converter) convertRun(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	// TODO run.shell
	// TODO run.when
	// TODO run.working_directory
	// TODO docker.auth.username
	// TODO docker.auth.password
	// TODO docker.aws_auth.aws_access_key_id
	// TODO docker.aws_auth.aws_secret_access_key

	var image string
	var entrypoint string
	var args []string
	var user string
	var envs map[string]string
	var shell string

	if docker := extractDocker(job, config); docker != nil {
		image = docker.Image
		entrypoint = "" // TODO needs a Harness v1 spec change
		args = docker.Command
		user = docker.User
		envs = docker.Environment
	}

	runCommand := step.Run.Command
	if job.Shell != "" {
		shellOptions := strings.Split(job.Shell, " ")[1:] // split the shell options from the shell binary
		if len(shellOptions) > 0 {
			shellOptionStr := strings.Join(shellOptions, " ")                  // join the shell options back into a single string
			runCommand = fmt.Sprintf("set %s\n%s", shellOptionStr, runCommand) // prepend the shell options to the run command
		}
		shell = strings.Split(job.Shell, " ")[0]
		shell = strings.Split(shell, "/")[len(strings.Split(shell, "/"))-1]
	} else { // default shell
		shell = "bash"
		runCommand = "set -eo pipefail\n" + runCommand
	}

	if step.Run.Background {
		return &harness.Step{
			Name: step.Run.Name,
			Type: "background",
			Spec: &harness.StepBackground{
				Run:        runCommand,
				Envs:       combineEnvs(step.Run.Environment, envs),
				Image:      image,
				Entrypoint: entrypoint,
				Args:       args,
				User:       user,
				Shell:      shell,
			},
		}
	} else {
		return &harness.Step{
			Name: step.Run.Name,
			Type: "script",
			Spec: &harness.StepExec{
				Run:        runCommand,
				Envs:       combineEnvs(step.Run.Environment, envs),
				Image:      image,
				Entrypoint: entrypoint,
				Args:       args,
				User:       user,
				Shell:      shell,
			},
		}
	}
}

// helper function converts a Circle Restore Cache step.
func (d *Converter) convertRestoreCache(step *circle.Step) *harness.Step {
	// TODO support restore_cache.keys (plural)
	return &harness.Step{
		Name: d.identifiers.Generate(step.RestoreCache.Name, "restore_cache"),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/cache",
			With: map[string]interface{}{
				"bucket":                          `<+ secrets.getValue("aws_bucket") >`,
				"region":                          `<+ secrets.getValue("aws_region") >`,
				"access_key":                      `<+ secrets.getValue("aws_access_key_id") >`,
				"secret_key":                      `<+ secrets.getValue("aws_secret_access_key") >`,
				"cache_key":                       step.RestoreCache.Key,
				"restore":                         "true",
				"exit_code":                       "true",
				"archive_format":                  "tar",
				"backend":                         "s3",
				"backend_operation_timeout":       "1800s",
				"fail_restore_if_key_not_present": "false",
			},
		},
	}
}

// helper function converts a Save Cache step.
func (d *Converter) convertSaveCache(step *circle.Step) *harness.Step {
	// TODO support save_cache.when
	return &harness.Step{
		Name: d.identifiers.Generate(step.SaveCache.Name, "save_cache"),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/cache",
			With: map[string]interface{}{
				"bucket":                          `<+ secrets.getValue("aws_bucket") >`,
				"region":                          `<+ secrets.getValue("aws_region") >`,
				"access_key":                      `<+ secrets.getValue("aws_access_key_id") >`,
				"secret_key":                      `<+ secrets.getValue("aws_secret_access_key") >`,
				"cache_key":                       step.SaveCache.Key,
				"rebuild":                         "true",
				"mount":                           step.SaveCache.Paths,
				"exit_code":                       "true",
				"archive_format":                  "tar",
				"backend":                         "s3",
				"backend_operation_timeout":       "1800s",
				"fail_restore_if_key_not_present": "false",
			},
		},
	}
}

// helper function converts a Add SSH Keys step.
func (d *Converter) convertAddSSHKeys(step *circle.Step) *harness.Step {
	// TODO step.AddSSHKeys.Fingerprints
	return &harness.Step{
		Name: d.identifiers.Generate(step.AddSSHKeys.Name, "add_ssh_keys"),
		Type: "script",
		Spec: &harness.StepExec{
			Run: "echo unable to convert add_ssh_keys step",
		},
	}
}

// helper function converts a Store Artifacts step.
func (d *Converter) convertStoreArtifacts(step *circle.Step) *harness.Step {
	src := step.StoreArtifacts.Path
	dst := step.StoreArtifacts.Destination
	if dst == "" {
		dst = "/"
	}
	return &harness.Step{
		Name: d.identifiers.Generate(step.StoreArtifacts.Name, "store_artifacts"),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/s3",
			With: map[string]interface{}{
				"bucket":     `<+ secrets.getValue("aws_bucket") >`,
				"region":     `<+ secrets.getValue("aws_region") >`,
				"access_key": `<+ secrets.getValue("aws_access_key_id") >`,
				"secret_key": `<+ secrets.getValue("aws_secret_access_key") >`,
				"source":     src,
				"target":     dst,
			},
		},
	}
}

// helper function converts a Test Results step.
func (d *Converter) convertStoreTestResults(step *circle.Step) *harness.Step {
	return &harness.Step{
		Name: d.identifiers.Generate(step.StoreTestResults.Name, "store_test_results"),
		Type: "script",
		Spec: &harness.StepExec{
			Run: "echo upload unit test results",
			Reports: []*harness.Report{
				{
					Path: []string{step.StoreTestResults.Path},
					Type: "junit",
				},
			},
		},
	}
}

// helper function converts a When step.
func (d *Converter) convertWhenStep(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	steps := d.convertSteps(step.When.Steps, job, config)
	if len(steps) == 0 {
		return nil
	}
	// TODO step.When.Condition
	return &harness.Step{
		Type: "group",
		Spec: &harness.StepGroup{
			Steps: steps,
		},
	}
}

// helper function converts an Unless step.
func (d *Converter) convertUnlessStep(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	steps := d.convertSteps(step.Unless.Steps, job, config)
	if len(steps) == 0 {
		return nil
	}
	// TODO step.Unless.Condition
	return &harness.Step{
		Type: "group",
		Spec: &harness.StepGroup{
			Steps: steps,
		},
	}
}

// helper function converts a Custom step.
func (d *Converter) convertCustom(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	// check to see if the step is a re-usable command.
	if _, ok := config.Commands[step.Custom.Name]; ok {
		return d.convertCommand(step, job, config)
	}
	// else convert the orb
	return d.convertOrb(step, job, config)
}

// helper function converts a Command step.
func (d *Converter) convertCommand(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	// extract the command
	command, ok := config.Commands[step.Custom.Name]
	if !ok {
		return nil
	}

	// find and replace parameters
	// https://circleci.com/docs/reusing-config/#using-the-parameters-declaration
	expandParamsCommand(command, step)

	// convert the circle steps to harness steps
	steps := d.convertSteps(command.Steps, job, config)
	if len(steps) == 0 {
		return nil
	}

	// If there is only one step, return it directly instead of creating a group
	if len(steps) == 1 {
		return steps[0]
	}

	// return a step group
	return &harness.Step{
		Type: "group",
		Spec: &harness.StepGroup{
			Steps: steps,
		},
	}
}

// helper function converts an Orb step.
func (d *Converter) convertOrb(step *circle.Step, job *circle.Job, config *circle.Config) *harness.Step {
	// get the orb alias and command
	alias, command := splitOrb(step.Custom.Name)

	// get the orb from the configuration
	orb, ok := config.Orbs[alias]
	if !ok {
		return nil
	}

	// convert inline orbs
	if orb.Inline != nil {
		// use the command to get the job name
		// if the action does not exist, silently
		// ignore the orb.
		job, ok := orb.Inline.Jobs[command]
		if !ok {
			return nil
		}
		// convert the orb steps to harness steps
		// if not steps are returned, silently ignore
		// the orb.
		steps := d.convertSteps(job.Steps, job, orb.Inline)
		if len(steps) == 0 {
			return nil
		}
		// return a step group
		return &harness.Step{
			Type: "group",
			Spec: &harness.StepGroup{
				Steps: steps,
			},
		}
	}

	name, version := splitOrbVersion(orb.Name)

	// convert the orb
	out := orbs.Convert(name, command, version, step.Custom)
	if out != nil {
		return out
	}

	return &harness.Step{
		Name: d.identifiers.Generate(name),
		Type: "script",
		Spec: &harness.StepExec{
			Run: fmt.Sprintf("echo unable to convert orb %s/%s", name, command),
		},
	}
}

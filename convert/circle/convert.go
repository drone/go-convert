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

	// create the harness pipeline
	pipeline := &harness.Pipeline{
		Version: 1,
	}

	// TODO .commands
	// TODO .orbs

	// convert pipeline and job parameters to inputs
	if params := extractParameters(config); len(params) != 0 {
		pipeline.Inputs = convertParameters(params)
	}

	// require a minimum of 1 workflows
	if config.Workflows == nil || len(config.Workflows.Items) == 0 {
		return nil, errors.New("no workflows defined")
	}

	// choose the first workflow in the list, for now
	// TODO convert multiple workflows
	var workflow *circle.Workflow
	for _, item := range config.Workflows.Items {
		workflow = item
		break
	}

	// loop through workflow jobs and convert each
	// job to a stage.
	for _, workflowjob := range workflow.Jobs {

		// TODO workflows.[*].triggers
		// TODO workflows.[*].unless
		// TODO workflows.[*].when
		// TODO workflows.[*].jobs[*].context
		// TODO workflows.[*].jobs[*].filters
		// TODO workflows.[*].jobs[*].matrix
		// TODO workflows.[*].jobs[*].type
		// TODO workflows.[*].jobs[*].requires

		// loop through jobs
		for name, job := range config.Jobs {
			// skip jobs that do not match
			if workflowjob.Name != name {
				continue
			}

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

			// TODO executor.resource_class
			// TODO executor.machine
			// TODO executor.shell
			// TODO executor.working_directory
			// TODO executor.environment

			// if there are no steps in the stage we
			// can skip adding the stage to the pipeline.
			if len(spec.Steps) == 0 {
				continue
			}

			// TODO jobs.[*].branches
			// TODO jobs.[*].parallelism
			// TODO jobs.[*].parameters

			// create the stage
			stage := &harness.Stage{}
			stage.Name = workflowjob.Name
			stage.Type = "ci"
			stage.Spec = spec

			// append the stage to the pipeline
			pipeline.Stages = append(pipeline.Stages, stage)
		}
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	// replace circle parameters with harness parameters
	out = replaceParams(out)

	return out, nil
}

func convertPlatform(job *circle.Job, config *circle.Config) *harness.Platform {
	executor := extractExecutor(job, config)
	if executor == nil {
		return nil
	}
	if executor.Windows != nil {
		return &harness.Platform{
			Os:   harness.OSWindows,
			Arch: harness.ArchAmd64,
		}
	}
	if executor.Machine != nil {
		if strings.Contains(executor.Machine.Image, "win") ||
			strings.Contains(executor.ResourceClass, "win") {
			return &harness.Platform{
				Os:   harness.OSWindows,
				Arch: harness.ArchAmd64,
			}
		}
		if strings.Contains(executor.Machine.Image, "arm") ||
			strings.Contains(executor.ResourceClass, "arm") {
			return &harness.Platform{
				Os:   harness.OSLinux,
				Arch: harness.ArchArm64,
			}
		}
	}
	if executor.Macos != nil {
		return &harness.Platform{
			Os:   harness.OSMacos,
			Arch: harness.ArchArm64,
		}
	}
	return &harness.Platform{
		Os:   harness.OSLinux,
		Arch: harness.ArchAmd64,
	}
}

func convertRuntime(job *circle.Job, config *circle.Config) *harness.Runtime {
	return &harness.Runtime{
		Type: "cloud",
		Spec: &harness.RuntimeCloud{},
	}
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
		return d.convertCustomStep(step)
	default:
		return nil
	}
}

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
	if docker := extractDocker(job, config); docker != nil {
		image = docker.Image
		entrypoint = "" // TODO needs a Harness v1 spec change
		args = docker.Command
		user = docker.User
		envs = docker.Environment
	}

	if step.Run.Background {
		return &harness.Step{
			Name: step.Run.Name,
			Type: "background",
			Spec: &harness.StepBackground{
				Run:        step.Run.Command,
				Envs:       conbineEnvs(step.Run.Environment, envs),
				Image:      image,
				Entrypoint: entrypoint,
				Args:       args,
				User:       user,
			},
		}
	} else {
		return &harness.Step{
			Name: step.Run.Name,
			Type: "script",
			Spec: &harness.StepExec{
				Run:        step.Run.Command,
				Envs:       conbineEnvs(step.Run.Environment, envs),
				Image:      image,
				Entrypoint: entrypoint,
				Args:       args,
				User:       user,
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
func (d *Converter) convertCustomStep(step *circle.Step) *harness.Step {
	return nil
}

func defaultBackgroundSteps(job *circle.Job, config *circle.Config) []*harness.Step {
	var steps []*harness.Step

	executor := extractExecutor(job, config)
	// exit if the executor is nil or if there are
	// no background containers defined.
	if executor == nil || len(executor.Docker) < 1 {
		return nil
	}
	for i := 1; i < len(executor.Docker); i++ {
		docker := executor.Docker[i]
		steps = append(steps, &harness.Step{
			Type: "background",
			Spec: &harness.StepBackground{
				Envs:  docker.Environment,
				Image: docker.Image,
				// TODO entrypoint
				// Entrypoint: docker.Entrypoint,
				Args: docker.Command,
				User: docker.User,
			},
		})
	}
	return steps
}

// func extractBackgroundSteps(job *circle.Job, config *circle.Config) []*circle.Docker {
// 	executor := extractExecutor(job, config)
// 	// exit if the executor is nil or if there are
// 	// no background containers defined.
// 	if executor == nil || len(executor.Docker) < 1 {
// 		return nil
// 	}
// 	return executor.Docker[1:]
// }

func extractDocker(job *circle.Job, config *circle.Config) *circle.Docker {
	executor := extractExecutor(job, config)
	// if the executor defines a docker environment,
	// use the first docker container as the execution
	// container.
	if executor != nil && len(executor.Docker) != 0 {
		return job.Docker[0]
	}
	return nil
}

func extractParameters(config *circle.Config) map[string]*circle.Parameter {
	params := map[string]*circle.Parameter{}

	// extract the parameters from the jobs.
	for _, job := range config.Jobs {
		for k, v := range job.Parameters {
			params[k] = v
		}
	}
	// extract the parameters from the pipeline.
	// these will override job parameters by design.
	for k, v := range config.Parameters {
		params[k] = v
	}
	return params
}

func extractExecutor(job *circle.Job, config *circle.Config) *circle.Executor {
	// return the named executor for the job
	if job.Executor != nil {
		// loop through the global executors.
		for name, executor := range config.Executors {
			if name == job.Executor.Name {
				return executor
			}
		}
	}
	// else create an executor based on the job
	// configuration. we do this because it is easier to
	// work with an executor struct, than both an executor
	// and a job struct.
	return &circle.Executor{
		Docker:        job.Docker,
		ResourceClass: job.ResourceClass,
		Machine:       job.Machine,
		Macos:         job.Macos,
		Windows:       nil,
		Shell:         job.Shell,
		WorkingDir:    job.WorkingDir,
		Environment:   job.Environment,
	}
}

func convertParameters(in map[string]*circle.Parameter) map[string]*harness.Input {
	out := map[string]*harness.Input{}
	for name, param := range in {
		t := param.Type
		switch t {
		case "integer":
			t = "number"
		case "string", "enum", "env_var_name":
			t = "string"
		case "boolean":
			t = "boolean"
		case "executor", "steps":
			// TODO parameter.type execution not supported
			// TODO parameter.type steps not supported
			continue // skip
		}
		var d string
		if param.Default != nil {
			d = fmt.Sprint(param.Default)
		}
		out[name] = &harness.Input{
			Type:        t,
			Default:     d,
			Description: param.Description,
		}
	}
	return out
}

// helper function combines environment variables.
func conbineEnvs(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}

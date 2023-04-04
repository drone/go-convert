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
	"io"
	"os"

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
		// Default: convertDefault(d.config),
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
				Steps:    d.convertSteps(job.Steps, job, config),
			}

			// if there are no steps in the stage we
			// can skip adding the stage to the pipeline.
			if len(spec.Steps) == 0 {
				continue
			}

			// TODO job.Branches
			// TODO job.Docker
			// TODO job.Executor
			// TODO job.IPRanges
			// TODO job.Machine
			// TODO job.Macos
			// TODO job.Parallelism
			// TODO job.Parameters
			// TODO job.ResourceClass
			// TODO job.Shell
			// TODO job.WorkingDir

			// create the stage
			stage := &harness.Stage{}
			stage.Name = workflowjob.Name
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

	return out, nil
}

func convertPlatform(job *circle.Job, config *circle.Config) *harness.Platform {
	if job.Macos != nil {
		return &harness.Platform{
			Os:   harness.OSMacos,
			Arch: harness.ArchArm64,
		}
	}
	// else if the job uses a global executor.
	if job.Executor != nil {
		// loop through the global executors.
		for name, executor := range config.Executors {
			// find the matching execturo.
			if name != job.Executor.Name {
				continue
			}
			if executor.Macos != nil {
				return &harness.Platform{
					Os:   harness.OSMacos,
					Arch: harness.ArchArm64,
				}
			}
		}
	}
	return nil
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
	return &harness.Step{
		Name: step.RestoreCache.Name,
		Type: "script",
		Spec: &harness.StepPlugin{
			// TODO
		},
	}
}

// helper function converts a Save Cache step.
func (d *Converter) convertSaveCache(step *circle.Step) *harness.Step {
	return &harness.Step{
		Name: *&step.SaveCache.Name,
		Type: "plugin",
		Spec: &harness.StepPlugin{
			// TODO
		},
	}
}

// helper function converts a Add SSH Keys step.
func (d *Converter) convertAddSSHKeys(step *circle.Step) *harness.Step {
	// TODO step.AddSSHKeys.Fingerprints
	return &harness.Step{
		Name: step.AddSSHKeys.Name,
		Type: "script",
		Spec: &harness.StepExec{
			Run: "",
		},
	}
}

// helper function converts a Store Artifacts step.
func (d *Converter) convertStoreArtifacts(step *circle.Step) *harness.Step {
	// TODO step.StoreArtifacts.Destination
	// TODO step.StoreArtifacts.Path
	return &harness.Step{
		Name: step.StoreArtifacts.Name,
		Type: "script",
		Spec: &harness.StepExec{
			Run: "",
		},
	}
}

// helper function converts a Test Results step.
func (d *Converter) convertStoreTestResults(step *circle.Step) *harness.Step {
	// TODO step.StoreTestResults.Name
	// TODO step.StoreTestResults.Path
	return &harness.Step{
		Name: step.StoreTestResults.Name,
		Type: "script",
		Spec: &harness.StepExec{
			Run: "",
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

func extractDocker(job *circle.Job, config *circle.Config) *circle.Docker {
	// if the job defines a docker executor
	// we can extract the default docker image.
	if len(job.Docker) != 0 {
		return job.Docker[0]
	}
	// else if the job uses a global executor.
	if job.Executor != nil {
		// loop through the global executors.
		for name, executor := range config.Executors {
			// find the matching execturo.
			if name != job.Executor.Name {
				continue
			}
			// if the matching executor defines
			// a docker execution environment, return
			// the first container in the list.
			if len(executor.Docker) != 0 {
				return executor.Docker[0]
			}
		}
	}
	return nil
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

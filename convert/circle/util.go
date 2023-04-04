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

package circle

import (
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

// helper function converts docker containers from the
// docker executor to background steps.
func defaultBackgroundSteps(job *circle.Job, config *circle.Config) []*harness.Step {
	var steps []*harness.Step

	executor := extractExecutor(job, config)
	// execit if the executor is nil, or if there is
	// less than 1 docker container defined. The first
	// container is used for execution, and subsequent
	// containers are used as background steps.
	if executor == nil || len(executor.Docker) < 1 {
		return nil
	}
	// loop through and convert the docker containers
	// to background steps.
	for i, docker := range executor.Docker {
		// skip the first docker container in the list,
		// since this is used as the run step execution
		// container only.
		if i == 0 {
			continue
		}
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

// helper function extracts the docker configuration
// from a job.
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

// helper function extrats an executor from a job.
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

// helper function extracts and aggregates the circle
// input parameters from the circle pipeline and job.
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

// helper function converts circle parameters to
// harness inputs.
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

// helper function converts circle executor to a
// harness platform.
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

// helper function converts circle executor to a
// harness runtime.
func convertRuntime(job *circle.Job, config *circle.Config) *harness.Runtime {
	return &harness.Runtime{
		Type: "cloud",
		Spec: &harness.RuntimeCloud{},
	}
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

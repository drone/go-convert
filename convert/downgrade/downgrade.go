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

package downgrade

import (
	"fmt"
	"time"

	"github.com/docker/go-units"
	"github.com/drone/go-convert/convert/internal/slug"

	v0 "github.com/drone/go-convert/convert/downgrade/yaml"
	v1 "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

func Downgrade(src *v1.Pipeline, args Args) ([]byte, error) {
	dst := new(v0.Pipeline)
	dst.ID = args.ID

	if args.ID == "" {
		dst.ID = slug.Create(src.Name)
	}

	dst.Name = src.Name
	dst.Org = args.Organization
	dst.Project = args.Project
	dst.Props.CI.Codebase = v0.Codebase{
		Name:  args.Codebase.Repo,
		Conn:  args.Codebase.Connector,
		Build: args.Codebase.Build,
	}

	if dst.Props.CI.Codebase.Build == "" {
		dst.Props.CI.Codebase.Build = "<+input>"
	}

	// convert stages
	for _, stage := range src.Stages {
		// skip nil stages. this is un-necessary but we have
		// this logic in place just to be safe.
		if stage == nil {
			continue
		}

		// skip stages that are not CI stages, for now
		if _, ok := stage.Spec.(*v1.StageCI); !ok {
			continue
		}

		// convert the stage and add to the list
		dst.Stages = append(dst.Stages, v0.Stages{
			Stage: convertStage(stage, args),
		})
	}

	return yaml.Marshal(dst)
}

// helper function converts a drone pipeline stage to a
// harness stage.
//
// TODO env variables to vars (stage-level)
// TODO delegate selectors
// TODO tags
// TODO when
// TODO infrastructure (cloud vs kubernetes)
// TODO failure strategy
// TODO matrix strategy
func convertStage(stage *v1.Stage, args Args) v0.Stage {
	// extract the spec from the v1 stage
	spec := stage.Spec.(*v1.StageCI)

	var steps []v0.Steps
	// convert each drone step to a harness step.
	for _, v := range spec.Steps {
		// the v0 yaml does not have the concept of
		// a group step, so we append all steps in
		// the group directly to the stage to emulate
		// this behavior.
		if _, ok := v.Spec.(*v1.StepGroup); ok {
			steps = append(steps, convertStepGroup(v, args)...)
		} else {
			// else convert the step and append to
			// the stage.
			steps = append(steps, convertStep(v, args))
		}
	}

	// enable clone by default
	enableClone := true
	if spec.Clone != nil && spec.Clone.Disabled == true {
		enableClone = false
	}

	// convert the drone stage to a harness stage.
	return v0.Stage{
		ID:   slug.Create(stage.Name),
		Name: stage.Name,
		Type: v0.StageTypeCI,
		Vars: convertVariables(spec.Envs),
		Spec: v0.StageCI{
			Cache: convertCache(spec.Cache),
			Clone: enableClone,
			Infrastructure: v0.Infrastructure{
				Type: v0.InfraTypeKubernetesDirect,
				Spec: v0.InfraSpec{
					Namespace: args.Kubernetes.Namespace,
					Conn:      args.Kubernetes.Connector,
				},
			},
			Execution: v0.Execution{
				Steps: steps,
			},
		},
	}
}

// helper function converts a drone pipeline step to a
// harness step.
//
// TODO unique identifier
// TODO failure strategy
// TODO matrix strategy
// TODO when
func convertStep(src *v1.Step, args Args) v0.Steps {
	switch src.Spec.(type) {
	case *v1.StepExec:
		return v0.Steps{Step: convertStepRun(src, args)}
	case *v1.StepPlugin:
		return v0.Steps{Step: convertStepPlugin(src, args)}
	case *v1.StepAction:
		return v0.Steps{Step: convertStepAction(src, args)}
	case *v1.StepBitrise:
		return v0.Steps{Step: convertStepBitrise(src, args)}
	case *v1.StepParallel:
		return v0.Steps{Parallel: convertStepParallel(src, args)}
	case *v1.StepBackground:
		return v0.Steps{Step: convertStepBackground(src, args)}
	default:
		return v0.Steps{} // should not happen
	}
}

// helper function to convert a Group step from the v1
// structure to a list of steps. The v0 yaml does not have
// an equivalent to the group step.
func convertStepGroup(src *v1.Step, args Args) []v0.Steps {
	spec_ := src.Spec.(*v1.StepGroup)

	var steps []v0.Steps
	for _, step := range spec_.Steps {
		dst := convertStep(step, args)
		steps = append(steps, v0.Steps{Step: dst.Step})
	}
	return steps
}

// helper function to convert a Parallel step from the v1
// structure to the v0 harness structure.
func convertStepParallel(src *v1.Step, args Args) []v0.Step {
	spec_ := src.Spec.(*v1.StepParallel)

	var steps []v0.Step
	for _, step := range spec_.Steps {
		dst := convertStep(step, args)
		steps = append(steps, dst.Step)
	}
	return steps
}

// helper function to convert a Run step from the v1
// structure to the v0 harness structure.
//
// TODO convert outputs
// TODO convert resources
// TODO convert reports
func convertStepRun(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepExec)
	return v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeRun,
		Timeout: convertTimeout(src.Timeout),
		Spec: v0.StepRun{
			Env:             spec_.Envs,
			Command:         spec_.Run,
			ConnRef:         args.Docker.Connector,
			Image:           spec_.Image,
			ImagePullPolicy: convertImagePull(spec_.Pull),
			Privileged:      spec_.Privileged,
			RunAsUser:       spec_.User,
		},
	}
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
// TODO convert ports
func convertStepBackground(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepBackground)
	return v0.Step{
		ID:   slug.Create(src.Name),
		Name: src.Name,
		Type: v0.StepTypeBackground,
		Spec: v0.StepBackground{
			Command:         spec_.Run,
			ConnRef:         args.Docker.Connector,
			Entrypoint:      []string{spec_.Entrypoint},
			Env:             spec_.Envs,
			Image:           spec_.Image,
			ImagePullPolicy: convertImagePull(spec_.Pull),
			Privileged:      spec_.Privileged,
			RunAsUser:       spec_.User,
		},
	}
}

// helper function to convert a Plugin step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
// TODO convert reports
func convertStepPlugin(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)
	return v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeRun,
		Timeout: convertTimeout(src.Timeout),
		Spec: v0.StepPlugin{
			ConnRef:         args.Docker.Connector,
			Image:           spec_.Image,
			ImagePullPolicy: convertImagePull(spec_.Pull),
			Settings:        convertSettings(spec_.With),
			Privileged:      spec_.Privileged,
			RunAsUser:       spec_.User,
		},
	}
}

// helper function to convert an Action step from the v1
// structure to the v0 harness structure.
func convertStepAction(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepAction)
	return v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeAction,
		Timeout: convertTimeout(src.Timeout),
		Spec: v0.StepAction{
			Uses: args.Docker.Connector,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
	}
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
func convertStepBitrise(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepBitrise)
	return v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeBitrise,
		Timeout: convertTimeout(src.Timeout),
		Spec: v0.StepBitrise{
			Uses: args.Docker.Connector,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
	}
}

func convertCache(src *v1.Cache) v0.Cache {
	var cache v0.Cache
	if src != nil {
		cache = v0.Cache{
			Enabled: src.Enabled,
			Key:     src.Key,
			Paths:   src.Paths,
		}
	}
	return cache
}

func convertVariables(src map[string]string) []v0.Variable {
	var vars []v0.Variable
	for k, v := range src {
		vars = append(vars, v0.Variable{
			Name:  k,
			Value: v,
			Type:  "String",
		})
	}
	return vars
}

func convertSettings(src map[string]interface{}) map[string]string {
	dst := map[string]string{}
	for k, v := range src {
		dst[k] = fmt.Sprint(v)
	}
	return dst
}

func convertTimeout(s string) v0.Duration {
	i, _ := units.FromHumanSize(s)
	return v0.Duration{Duration: time.Duration(i)}
}

func convertImagePull(v string) (s string) {
	switch v {
	case "always":
		return v0.ImagePullAlways
	case "never":
		return v0.ImagePullNever
	case "if-not-exists":
		return v0.ImagePullIfNotPresent
	default:
		return ""
	}
}

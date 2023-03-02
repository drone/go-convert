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

	v0 "github.com/drone/go-convert/convert/downgrade/yaml"
	v1 "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

func Downgrade(src *v1.Pipeline, args Args) ([]byte, error) {
	dst := new(v0.Pipeline)

	dst.ID = src.Name   // TODO ensure not-null, unique
	dst.Name = src.Name // TODO ensure not-null, unique
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
func convertStage(stage *v1.Stage, args Args) v0.Stage {
	// extract the spec from the v1 stage
	spec := stage.Spec.(*v1.StageCI)

	var steps []v0.Steps
	// convert each drone step to a harness step.
	for _, v := range spec.Steps {
		if _, ok := v.Spec.(*v1.StepGroup); ok {
			steps = append(steps,
				convertStepGroup(v, args)...,
			)
		} else {
			step := convertStep(v, args)
			steps = append(steps, step)
		}
	}

	// enable clone by default
	enableClone := true
	if spec.Clone != nil && spec.Clone.Disabled == true {
		enableClone = false
	}

	// convert the drone stage to a harness stage.
	return v0.Stage{
		ID:   stage.Name,
		Name: stage.Name,
		Type: v0.StageTypeCI,
		Spec: v0.StageCI{
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

func convertStepGroup(src *v1.Step, args Args) []v0.Steps {
	spec_ := src.Spec.(*v1.StepGroup)

	var steps []v0.Steps
	for _, step := range spec_.Steps {
		dst := convertStep(step, args)
		steps = append(steps, v0.Steps{Step: dst.Step})
	}
	return steps
}

func convertStepParallel(src *v1.Step, args Args) []v0.Step {
	spec_ := src.Spec.(*v1.StepParallel)

	var steps []v0.Step
	for _, step := range spec_.Steps {
		dst := convertStep(step, args)
		steps = append(steps, dst.Step)
	}
	return steps
}

func convertStepRun(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepExec)

	spec := v0.StepRun{
		Env:             spec_.Envs,
		Command:         spec_.Run,
		ConnRef:         args.Docker.Connector,
		Image:           spec_.Image,
		ImagePullPolicy: convertImagePull(spec_.Pull),
		Privileged:      spec_.Privileged,
		RunAsUser:       spec_.User,
		// TODO convert resources
		// TODO convert reports
	}

	step := v0.Step{
		ID:   src.Name,
		Name: src.Name,
		Type: v0.StepTypeRun,
		Spec: spec,
	}

	return step
}

func convertStepBackground(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepBackground)

	spec := v0.StepBackground{
		Env:             spec_.Envs,
		Entrypoint:      []string{spec_.Entrypoint},
		Command:         spec_.Run,
		ConnRef:         args.Docker.Connector,
		Image:           spec_.Image,
		ImagePullPolicy: convertImagePull(spec_.Pull),
		Privileged:      spec_.Privileged,
		RunAsUser:       spec_.User,
		// TODO convert resources
		// TODO convert ports
	}

	step := v0.Step{
		ID:   src.Name,
		Name: src.Name,
		Type: v0.StepTypeBackground,
		Spec: spec,
	}

	return step
}

func convertStepPlugin(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)

	spec := v0.StepPlugin{
		ConnRef:         args.Docker.Connector,
		Image:           spec_.Image,
		ImagePullPolicy: convertImagePull(spec_.Pull),
		Settings:        convertSettings(spec_.With),
		Privileged:      spec_.Privileged,
		RunAsUser:       spec_.User,
		// TODO convert resources
		// TODO convert reports
	}

	step := v0.Step{
		ID:   src.Name,
		Name: src.Name,
		Type: v0.StepTypeRun,
		Spec: spec,
	}

	return step
}

func convertStepAction(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepAction)
	return v0.Step{
		ID:   src.Name,
		Name: src.Name,
		Type: v0.StepTypeAction,
		Spec: v0.StepAction{
			Uses: args.Docker.Connector,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
	}
}

func convertStepBitrise(src *v1.Step, args Args) v0.Step {
	spec_ := src.Spec.(*v1.StepBitrise)
	return v0.Step{
		ID:   src.Name,
		Name: src.Name,
		Type: v0.StepTypeBitrise,
		Spec: v0.StepBitrise{
			Uses: args.Docker.Connector,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
	}
}

func convertSettings(src map[string]interface{}) map[string]string {
	dst := map[string]string{}
	for k, v := range src {
		dst[k] = fmt.Sprint(v)
	}
	return dst
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

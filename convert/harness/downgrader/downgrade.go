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

package downgrader

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/docker/go-units"
	"github.com/drone/go-convert/internal/slug"
	"github.com/drone/go-convert/internal/store"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

// Downgrader downgrades pipelines from the v0 harness
// configuration format to the v1 configuration format.
type Downgrader struct {
	codebaseName  string
	codebaseConn  string
	dockerhubConn string
	kubeConnector string
	kubeNamespace string
	kubeEnabled   bool
	pipelineId    string
	pipelineName  string
	pipelineOrg   string
	pipelineProj  string
	identifiers   *store.Identifiers
}

// New creates a new Downgrader that downgrades pipelines
// from the v0 harness configuration format to the v1
// configuration format.
func New(options ...Option) *Downgrader {
	d := new(Downgrader)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = new(store.Identifiers)

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default pipeline name.
	if d.pipelineName == "" {
		d.pipelineName = "default"
	}

	// set the default pipeline id.
	if d.pipelineId == "" {
		d.pipelineId = slug.Create(d.pipelineName)
	}

	// set the default pipeline org.
	if d.pipelineOrg == "" {
		d.pipelineOrg = "default"
	}

	// set the default pipeline org.
	if d.pipelineProj == "" {
		d.pipelineProj = "default"
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

// Downgrade downgrades a v1 pipeline.
func (d *Downgrader) Downgrade(b []byte) ([]byte, error) {
	src, err := v1.ParseBytes(b)
	if err != nil {
		return nil, err
	}
	return d.downgrade(src)
}

// DowngradeString downgrades a v1 pipeline.
func (d *Downgrader) DowngradeString(s string) ([]byte, error) {
	return d.Downgrade([]byte(s))
}

// DowngradeString downgrades a v1 pipeline.
func (d *Downgrader) DowngradeFile(path string) ([]byte, error) {
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return d.Downgrade(out)
}

// downgrade downgrades a v1 pipeline.
func (d *Downgrader) downgrade(src *v1.Pipeline) ([]byte, error) {
	dst := new(v0.Pipeline)

	dst.ID = d.pipelineId
	dst.Name = d.pipelineName
	dst.Org = d.pipelineOrg
	dst.Project = d.pipelineProj
	dst.Props.CI.Codebase = v0.Codebase{
		Name:  d.codebaseName,
		Conn:  d.codebaseConn,
		Build: "<+input>",
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
		dst.Stages = append(dst.Stages, &v0.Stages{
			Stage: d.convertStage(stage),
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
// TODO runtime / kubernetes / cloud
// TODO platform / os / arch
func (d *Downgrader) convertStage(stage *v1.Stage) *v0.Stage {
	// extract the spec from the v1 stage
	spec := stage.Spec.(*v1.StageCI)

	var steps []*v0.Steps
	// convert each drone step to a harness step.
	for _, v := range spec.Steps {
		// the v0 yaml does not have the concept of
		// a group step, so we append all steps in
		// the group directly to the stage to emulate
		// this behavior.
		if _, ok := v.Spec.(*v1.StepGroup); ok {
			steps = append(steps, d.convertStepGroup(v)...)
		} else {
			// else convert the step and append to
			// the stage.
			steps = append(steps, d.convertStep(v))
		}
	}

	// enable clone by default
	enableClone := true
	if spec.Clone != nil && spec.Clone.Disabled == true {
		enableClone = false
	}

	// convert the drone stage to a harness stage.
	return &v0.Stage{
		ID:   slug.Create(stage.Name),
		Name: stage.Name,
		Type: v0.StageTypeCI,
		Vars: convertVariables(spec.Envs),
		Spec: v0.StageCI{
			Cache: convertCache(spec.Cache),
			Clone: enableClone,
			Infrastructure: &v0.Infrastructure{
				Type: v0.InfraTypeKubernetesDirect,
				Spec: &v0.InfraSpec{
					Namespace: d.kubeNamespace,
					Conn:      d.kubeConnector,
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
func (d *Downgrader) convertStep(src *v1.Step) *v0.Steps {
	switch src.Spec.(type) {
	case *v1.StepExec:
		return &v0.Steps{Step: d.convertStepRun(src)}
	case *v1.StepPlugin:
		return &v0.Steps{Step: d.convertStepPlugin(src)}
	case *v1.StepAction:
		return &v0.Steps{Step: d.convertStepAction(src)}
	case *v1.StepBitrise:
		return &v0.Steps{Step: d.convertStepBitrise(src)}
	case *v1.StepParallel:
		return &v0.Steps{Parallel: d.convertStepParallel(src)}
	case *v1.StepBackground:
		return &v0.Steps{Step: d.convertStepBackground(src)}
	default:
		return nil // should not happen
	}
}

// helper function to convert a Group step from the v1
// structure to a list of steps. The v0 yaml does not have
// an equivalent to the group step.
func (d *Downgrader) convertStepGroup(src *v1.Step) []*v0.Steps {
	spec_ := src.Spec.(*v1.StepGroup)

	var steps []*v0.Steps
	for _, step := range spec_.Steps {
		dst := d.convertStep(step)
		steps = append(steps, &v0.Steps{Step: dst.Step})
	}
	return steps
}

// helper function to convert a Parallel step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepParallel(src *v1.Step) []*v0.Step {
	spec_ := src.Spec.(*v1.StepParallel)

	var steps []*v0.Step
	for _, step := range spec_.Steps {
		dst := d.convertStep(step)
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
func (d *Downgrader) convertStepRun(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepExec)
	return &v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeRun,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepRun{
			Env:             spec_.Envs,
			Command:         spec_.Run,
			ConnRef:         d.dockerhubConn,
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
func (d *Downgrader) convertStepBackground(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepBackground)

	// convert the entrypoint string to a slice.
	var entypoint []string
	if spec_.Entrypoint != "" {
		entypoint = []string{spec_.Entrypoint}
	}

	return &v0.Step{
		ID:   slug.Create(src.Name),
		Name: src.Name,
		Type: v0.StepTypeBackground,
		Spec: &v0.StepBackground{
			Command:         spec_.Run,
			ConnRef:         d.dockerhubConn,
			Entrypoint:      entypoint,
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
func (d *Downgrader) convertStepPlugin(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)
	return &v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeRun,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepPlugin{
			ConnRef:         d.dockerhubConn,
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
func (d *Downgrader) convertStepAction(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepAction)
	return &v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeAction,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepAction{
			Uses: spec_.Uses,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
	}
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepBitrise(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepBitrise)
	return &v0.Step{
		ID:      slug.Create(src.Name),
		Name:    src.Name,
		Type:    v0.StepTypeBitrise,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepBitrise{
			Uses: spec_.Uses,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
	}
}

func convertCache(src *v1.Cache) *v0.Cache {
	if src == nil {
		return nil
	}
	return &v0.Cache{
		Enabled: src.Enabled,
		Key:     src.Key,
		Paths:   src.Paths,
	}
}

func convertVariables(src map[string]string) []*v0.Variable {
	var vars []*v0.Variable
	for k, v := range src {
		vars = append(vars, &v0.Variable{
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
	if i == -1 {
		return v0.Duration{
			Duration: time.Duration(0),
		}
	} else {
		return v0.Duration{
			Duration: time.Duration(i),
		}
	}
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

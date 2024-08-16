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

package bitbucket

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	bitbucket "github.com/drone/go-convert/convert/bitbucket/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// Converter converts a Bitbucket pipeline to a harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers

	// as we walk the yaml, we store a
	// a snapshot of the current node and
	// its parents.
	config *bitbucket.Config
	stage  *bitbucket.Stage
	steps  *bitbucket.Steps
	step   *bitbucket.Step
	script *bitbucket.Script
}

// New creates a new Converter that converts a Bitbucket
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

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	src, err := bitbucket.Parse(r)
	if err != nil {
		return nil, err
	}
	d.config = src // push the bitbucket config to the state
	return d.convert()
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

// converts converts a bitbucket pipeline pipeline.
func (d *Converter) convert() ([]byte, error) {

	// normalize the yaml and ensure
	// all root-level steps are grouped
	// by stage to simplify conversion.
	bitbucket.Normalize(d.config)

	// create the harness pipeline spec
	pipeline := &harness.Pipeline{
		Options: convertDefault(d.config),
	}

	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    pipeline,
	}

	for _, steps := range d.config.Pipelines.Default {
		if steps.Stage != nil {
			// TODO support for fast-fail
			d.stage = steps.Stage // push the stage to the state
			stage := d.convertStage()
			pipeline.Stages = append(pipeline.Stages, stage)
		}
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// helper function converts a bitbucket stage to
// a harness stage.
func (d *Converter) convertStage() *harness.Stage {

	// create the harness stage spec
	spec := &harness.StageCI{
		Clone: convertClone(d.stage),
		// TODO Repository
		// TODO Delegate
		// TODO Platform
		// TODO Runtime
		// TODO Envs
	}

	// find the step with the largest size and use that
	// size. else fallback to the global size.
	if size := extractSize(d.config.Options, d.stage); size != bitbucket.SizeNone {
		spec.Runtime = &harness.Runtime{
			Type: "cloud",
			Spec: &harness.RuntimeCloud{
				Size: convertSize(size),
			},
		}
	}

	// find the unique cache paths used by this
	// stage and setup harness caching
	if paths := extractCache(d.stage); len(paths) != 0 {
		spec.Cache = convertCache(d.config.Definitions, paths)
	}

	// find the unique services used by this stage and
	// setup the relevant background steps
	if services := extractServices(d.stage); len(services) != 0 {
		spec.Steps = append(spec.Steps, d.convertServices(services)...)
	}

	// create the harness stage.
	stage := &harness.Stage{
		Name: "build",
		Type: "ci",
		Spec: spec,
		// TODO When
		// TODO Failure
	}

	// find the unique selectors and append
	// to the stage.
	if runson := extractRunsOn(d.stage); len(runson) != 0 {
		stage.Delegate = runson
	}

	// default docker service (container-based only)
	if d.config.Options != nil && d.config.Options.Docker {
		spec.Steps = append(spec.Steps, &harness.Step{
			Name: d.identifiers.Generate("dind", "service"),
			Type: "background",
			Spec: &harness.StepBackground{
				Image:      "docker:dind",
				Ports:      []string{"2375", "2376"},
				Network:    "host", // TODO host networking for cloud only
				Privileged: true,
			},
		})
	}

	// default services
	// TODO

	for _, steps := range d.stage.Steps {
		if steps.Parallel != nil {
			// TODO parallel steps
			// TODO fast fail
			d.steps = steps // push the parallel step to the state
			step := d.convertParallel()
			spec.Steps = append(spec.Steps, step)
		}
		if steps.Step != nil {
			d.step = steps.Step // push the step to the state
			step := d.convertStep()
			spec.Steps = append(spec.Steps, step)
		}
	}

	// if the stage has a single step, and that step is a
	// group step, we can eliminate the un-necessary group
	// and add the steps directly to the stage.
	if len(spec.Steps) == 1 {
		if group, ok := spec.Steps[0].Spec.(*harness.StepGroup); ok {
			spec.Steps = group.Steps
		}
	}

	return stage
}

// helper function converts global bitbucket services to
// harness background steps. The list of global bitbucket
// services is filtered by the services string slice.
func (d *Converter) convertServices(services []string) []*harness.Step {
	var steps []*harness.Step

	// if no global services defined, exit
	if d.config.Definitions == nil {
		return nil
	}

	// iterate through services and create background steps
	for _, name := range services {
		// lookup the service and skip if not found,
		// or if there is no image definition
		service, ok := d.config.Definitions.Services[name]
		if !ok {
			continue
		} else if service.Image == nil {
			continue
		}

		spec := &harness.StepBackground{
			Image:   service.Image.Name,
			Envs:    service.Variables,
			Network: "host", // TODO host netowrking for cloud only
		}

		// if the service is of type docker, we
		// should open up the default docker ports
		// and also run in privileged mode.
		if service.Type == "docker" {
			spec.Privileged = true
			spec.Ports = []string{"2375", "2376"} // TODO can we remove this?
			spec.Network = "host"                 // TODO host networking for Cloud only
		}

		// if the service specifies a uid then set the
		// step user identifier.
		if v := service.Image.RunAsUser; v != 0 {
			spec.User = fmt.Sprint(v)
		}

		// if the service defines memory set the
		// harness resource limit.
		if v := service.Memory; v != 0 {
			// memory in bitbucket is measured in megabytes
			// so we need to convert to bytes.
			spec.Resources = &harness.Resources{
				Limits: &harness.Resource{
					Memory: harness.MemStringorInt(v * 1000000),
				},
			}
		}

		step := &harness.Step{
			Name: d.identifiers.Generate(name, "service"),
			Type: "background",
			Spec: spec,
		}

		steps = append(steps, step)
	}
	return steps
}

// helper function converts a bitbucket parallel step
// group to a Harness parallel step group.
func (d *Converter) convertParallel() *harness.Step {

	// create the step group spec
	spec := new(harness.StepParallel)

	for _, src := range d.steps.Parallel.Steps {
		if src.Step != nil {
			d.step = src.Step
			step := d.convertStep()
			spec.Steps = append(spec.Steps, step)
		}
	}

	// else create the step group wrapper.
	return &harness.Step{
		Type: "parallel",
		Spec: spec,
		Name: d.identifiers.Generate("parallel", "parallel"), // TODO can we avoid a name here?
	}
}

// helper function converts a bitbucket step
// to a harness run step or plugin step.
func (d *Converter) convertStep() *harness.Step {
	// create the step group spec
	spec := new(harness.StepGroup)

	// loop through each script item
	for _, script := range d.step.Script {
		d.script = script

		// if a pipe step
		if script.Pipe != nil {
			step := d.convertPipeStep()
			spec.Steps = append(spec.Steps, step)
		}

		// else if a script step
		if script.Pipe == nil {
			step := d.convertScriptStep()
			spec.Steps = append(spec.Steps, step)
		}
	}

	// and loop through each after script item
	for _, script := range d.step.ScriptAfter {
		d.script = script

		// if a pipe step
		if script.Pipe != nil {
			step := d.convertPipeStep()
			spec.Steps = append(spec.Steps, step)
		}

		// else if a script step
		if script.Pipe == nil {
			step := d.convertScriptStep()
			spec.Steps = append(spec.Steps, step)
		}
	}

	// if there is only a single step, no need to
	// create a step group.
	if len(spec.Steps) == 1 {
		return spec.Steps[0]
	}

	// else create the step group wrapper.
	return &harness.Step{
		Type: "group",
		Spec: spec,
		Name: d.identifiers.Generate(d.step.Name, "group"),
	}
}

// helper function converts a script step to a
// harness run step.
func (d *Converter) convertScriptStep() *harness.Step {

	// create the run spec
	spec := &harness.StepExec{
		Run: d.script.Text,

		// TODO configure an optional connector
		// TODO configure pull policy
		// TODO configure envs
		// TODO configure volumes
		// TODO configure resources
	}

	// use the global image, if set
	if image := d.config.Image; image != nil {
		spec.Image = strings.TrimPrefix(image.Name, "docker://")
		if image.RunAsUser != 0 {
			spec.User = fmt.Sprint(image.RunAsUser)
		}
	}

	// use the step image, if set (overrides previous)
	if image := d.step.Image; image != nil {
		spec.Image = strings.TrimPrefix(image.Name, "docker://")
		if image.RunAsUser != 0 {
			spec.User = fmt.Sprint(image.RunAsUser)
		}
	}

	// create the run step wrapper
	step := &harness.Step{
		Type: "script",
		Spec: spec,
		Name: d.identifiers.Generate(d.step.Name, "run"),
	}

	// use the global max-time, if set
	if d.config.Options != nil {
		if v := int64(d.config.Options.MaxTime); v != 0 {
			step.Timeout = minuteToDurationString(v)
		}
	}

	// set the timeout
	if v := int64(d.step.MaxTime); v != 0 {
		step.Timeout = minuteToDurationString(v)
	}

	return step
}

// helper function converts a pipe step to a
// harness plugin step.
func (d *Converter) convertPipeStep() *harness.Step {
	pipe := d.script.Pipe

	// create the plugin spec
	spec := &harness.StepPlugin{
		Image: strings.TrimPrefix(pipe.Image, "docker://"),

		// TODO configure an optional connector
		// TODO configure envs
		// TODO configure volumes
	}

	// append the plugin spec variables
	spec.With = map[string]interface{}{}
	for key, val := range pipe.Variables {
		spec.With[key] = val
	}

	// create the plugin step wrapper
	step := &harness.Step{
		Type: "plugin",
		Spec: spec,
		Name: d.identifiers.Generate(d.step.Name, "plugin"),
	}

	// set the timeout
	if v := int64(d.step.MaxTime); v != 0 {
		step.Timeout = minuteToDurationString(v)
	}

	return step
}

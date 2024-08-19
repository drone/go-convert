// Copyright 2024 Harness, Inc.
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

// Package jenkinsxml converts Jenkins XML pipelines to Harness pipelines.
package jenkinsxml

import (
	"bytes"
	"io"
	"os"

	jenkinsxml "github.com/drone/go-convert/convert/jenkinsxml/xml"
	"github.com/drone/go-convert/internal/store"
	harness "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	config *jenkinsxml.Project
}

// Converter converts a Gitlab pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers
}

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
	src, err := jenkinsxml.Parse(r)
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

	// cacheFound := false

	// create the harness stage.
	dstStage := &harness.Stage{
		Type: "ci",
		// When: convertCond(from.Trigger),
		Spec: &harness.StageCI{
			// Delegate: convertNode(from.Node),
			// Envs: convertVariables(ctx.config.Variables),
			// Platform: convertPlatform(from.Platform),
			// Runtime:  convertRuntime(from),
			Steps: make([]*harness.Step, 0), // Initialize the Steps slice
		},
	}
	dst.Stages = append(dst.Stages, dstStage)
	stageSteps := make([]*harness.Step, 0)

	// TODO: support more than just Shell tasks, such as Ant
	tasks := ctx.config.Builders.HudsonShellTasks
	for _, task := range tasks {
		step := convertShellTaskToStep(task.Command)
		stageSteps = append(stageSteps, step)
	}
	dstStage.Spec.(*harness.StageCI).Steps = stageSteps

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// convertShellTaskToStep converts a Jenkins shell task to a Harness step.
func convertShellTaskToStep(command string) *harness.Step {
	spec := new(harness.StepExec)
	spec.Run = command
	step := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: spec,
	}

	return step
}

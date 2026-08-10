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
	"encoding/xml"
	"io"
	"os"
	"strings"

	harnessconv "github.com/drone/go-convert/convert/harness"
	jenkinsxml "github.com/drone/go-convert/convert/jenkinsxml/xml"
	"github.com/drone/go-convert/internal/store"
	harness "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	config *jenkinsxml.Project
}

// Converter converts a Jenkins XML file to a Harness
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

// converts converts a Jenkins XML pipeline to a Harness pipeline.
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

	// emit a git clone step first when the job uses a Git SCM. The
	// downgrader lifts this step's url into the pipeline codebase, so it
	// must be the first step in the stage.
	if step := convertSCMToStep(ctx.config); step != nil {
		stageSteps = append(stageSteps, step)
	}

	// Builders is nil for job types that do not use a freestyle
	// <builders> block (for example maven2-moduleset or scripted
	// flow-definition pipelines). Guard against a nil dereference and
	// emit a pipeline with no steps rather than panicking.
	var tasks []jenkinsxml.Task
	if ctx.config.Builders != nil {
		tasks = ctx.config.Builders.Tasks
	}
	for _, task := range tasks {
		step := &harness.Step{}

		switch taskname := task.XMLName.Local; taskname {
		case "hudson.tasks.Shell":
			step = convertShellTaskToStep(&task)
		case "hudson.tasks.Ant":
			step = convertAntTaskToStep(&task)
		default:
			step = unsupportedTaskToStep(taskname)
		}

		stageSteps = append(stageSteps, step)
	}

	// maven2-moduleset jobs declare their build via top-level <goals>
	// rather than a <builders> block. Emit an mvn step so these jobs
	// convert to a runnable pipeline instead of an empty stage.
	if step := convertMavenGoalsToStep(ctx.config); step != nil {
		stageSteps = append(stageSteps, step)
	}

	dstStage.Spec.(*harness.StageCI).Steps = stageSteps

	// map Jenkins string build parameters to pipeline inputs.
	if inputs := convertParameters(ctx.config); len(inputs) > 0 {
		dst.Inputs = inputs
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// convertSCMToStep converts a Jenkins Git SCM to a Harness git clone
// plugin step. It returns nil when the job has no Git SCM (the only SCM
// type modelled today). The branch ref is stripped of Jenkins' "*/"
// prefix (for example "*/master" becomes "master").
func convertSCMToStep(project *jenkinsxml.Project) *harness.Step {
	if project == nil || project.SCM == nil {
		return nil
	}
	scm := project.SCM
	if scm.Class != "hudson.plugins.git.GitSCM" || len(scm.RemoteURLs) == 0 {
		return nil
	}

	with := map[string]interface{}{
		"git_url": scm.RemoteURLs[0],
	}
	if len(scm.BranchNames) > 0 {
		with["branch"] = strings.TrimPrefix(scm.BranchNames[0], "*/")
	}

	spec := &harness.StepPlugin{
		Image: harnessconv.GitPluginImage,
		With:  with,
	}
	return &harness.Step{
		Name: "clone",
		Type: "plugin",
		Spec: spec,
	}
}

// convertParameters maps Jenkins string build parameters to Harness
// pipeline inputs. It returns nil when the job declares no parameters.
func convertParameters(project *jenkinsxml.Project) map[string]*harness.Input {
	if project == nil || len(project.Parameters) == 0 {
		return nil
	}

	inputs := make(map[string]*harness.Input, len(project.Parameters))
	for _, param := range project.Parameters {
		if param.Name == "" {
			continue
		}
		inputs[param.Name] = &harness.Input{
			Type:        "string",
			Description: param.Description,
			Default:     param.DefaultValue,
		}
	}
	return inputs
}

// convertMavenGoalsToStep converts the top-level <goals> of a
// maven2-moduleset job to a Harness Run step. It returns nil when the
// project declares no Maven goals (for example a freestyle job).
func convertMavenGoalsToStep(project *jenkinsxml.Project) *harness.Step {
	if project == nil || project.Goals == "" {
		return nil
	}

	spec := new(harness.StepExec)
	spec.Run = "mvn " + project.Goals
	step := &harness.Step{
		Name: "maven",
		Type: "script",
		Spec: spec,
	}

	return step
}

// convertAntTaskToStep converts a Jenkins Ant task to a Harness step.
func convertAntTaskToStep(task *jenkinsxml.Task) *harness.Step {
	antTask := &jenkinsxml.HudsonAntTask{}
	// TODO: wrapping task.Content with 'builders' tags is ugly.
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), antTask)
	if err != nil {
		return nil
	}

	spec := new(harness.StepPlugin)
	spec.Image = "harnesscommunitytest/ant-plugin"
	spec.Inputs = map[string]interface{}{
		"goals": antTask.Targets,
	}
	step := &harness.Step{
		Name: "ant",
		Type: "plugin",
		Spec: spec,
	}

	return step
}

// convertShellTaskToStep converts a Jenkins Shell task to a Harness step.
func convertShellTaskToStep(task *jenkinsxml.Task) *harness.Step {
	shellTask := &jenkinsxml.HudsonShellTask{}
	// TODO: wrapping task.Content with 'builders' tags is ugly.
	err := xml.Unmarshal([]byte("<builders>"+task.Content+"</builders>"), shellTask)
	if err != nil {
		return nil
	}

	spec := new(harness.StepExec)
	spec.Run = shellTask.Command
	step := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: spec,
	}

	return step
}

// unsupportedTaskToStep converts an unsupported Jenkins Task to a placeholder
// Harness step.
func unsupportedTaskToStep(task string) *harness.Step {
	spec := new(harness.StepExec)
	spec.Run = "echo Unsupported field " + task
	step := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: spec,
	}

	return step
}

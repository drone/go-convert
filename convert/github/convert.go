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

// Package github converts GitHub pipelines to Harness pipelines.
package github

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	v1 "github.com/drone/go-convert/convert/github/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	pipeline []*v1.Pipeline
	stage    *v1.Pipeline
}

// Converter converts a GitHub pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers

	// // as we walk the yaml, we store a
	// // a snapshot of the current node and
	// // its parents.
	// config *github.Pipeline
	// stage  *github.Stage
}

// New creates a new Converter that converts a GitHub
// pipeline to a Harness v1 pipeline.
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
	src, err := v1.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(&context{
		pipeline: src,
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

// converts a GitHub pipeline to Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {

	// create the harness pipeline
	pipeline := &harness.Pipeline{
		Version: 1,
		Stages:  []*harness.Stage{},
	}

	for _, from := range ctx.pipeline {
		if from == nil {
			continue
		}

		var actionJob *v1.Job
		if from.Jobs != nil {
			for _, job := range from.Jobs {
				actionJob = &job
				break
			}
		}

		var cloneStage *harness.CloneStage
		if actionJob != nil {
			for _, step := range actionJob.Steps {
				cloneStage = convertClone(step)
				if cloneStage != nil {
					break
				}
			}
		}

		pipeline.Stages = append(pipeline.Stages, &harness.Stage{
			Name: from.Name,
			Type: "ci",
			When: convertCond(&from.On),
			Spec: &harness.StageCI{
				Clone:    cloneStage,
				Envs:     copyEnv(from.Environment),
				Platform: convertRunsOn(actionJob.RunsOn),
				Runtime: &harness.Runtime{
					Type: "machine",
					Spec: harness.RuntimeMachine{},
				},
				Steps: convertSteps(actionJob),
				//Volumes:  convertVolumes(from.Volumes),

				// TODO support for delegate.selectors from from.Node
				// TODO support for stage.variables
			},
		})
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func convertClone(src *v1.Step) *harness.CloneStage {
	if src == nil || !isCheckoutAction(src.Uses) {
		return nil
	}
	dst := new(harness.CloneStage)
	if src.With != nil {
		if depth, ok := src.With["fetch-depth"]; ok {
			dst.Depth, _ = toInt64(depth)
		}
	}
	return dst
}

func convertCond(src *v1.WorkflowTriggers) *harness.When {
	if src == nil || isTriggersEmpty(src) {
		return nil
	}

	exprs := map[string]*harness.Expr{}

	for eventName, eventCondition := range getEventConditions(src) {
		if expr := convertEventCondition(eventCondition); expr != nil {
			exprs[eventName] = expr
		}
	}

	dst := new(harness.When)
	dst.Cond = []map[string]*harness.Expr{exprs}
	return dst
}

func getEventConditions(src *v1.WorkflowTriggers) map[string][]string {
	eventConditions := make(map[string][]string)

	if src.Push != nil {
		eventConditions["push"] = src.Push.Branches
	}
	if src.PullRequest != nil {
		eventConditions["pull_request"] = src.PullRequest.Branches
	}
	return eventConditions
}

func convertEventCondition(src []string) *harness.Expr {
	if len(src) != 0 {
		return &harness.Expr{In: src}
	}
	return nil
}

func isTriggersEmpty(src *v1.WorkflowTriggers) bool {
	return (src.Push == nil || len(src.Push.Branches) == 0) &&
		(src.PullRequest == nil || len(src.PullRequest.Branches) == 0)
}

func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		intValue, err := strconv.Atoi(v)
		return int64(intValue), err
	default:
		return 0, fmt.Errorf("unsupported type for conversion to int64")
	}
}

func isCheckoutAction(action string) bool {
	matched, _ := regexp.MatchString(`^actions/checkout@`, action)
	return matched
}

func convertRunsOn(src string) *harness.Platform {
	if src == "" {
		return nil
	}
	dst := new(harness.Platform)
	switch {
	case strings.Contains(src, "windows"), strings.Contains(src, "win"):
		dst.Os = harness.OSWindows
	case strings.Contains(src, "darwin"), strings.Contains(src, "macos"), strings.Contains(src, "mac"):
		dst.Os = harness.OSDarwin
	default:
		dst.Os = harness.OSLinux
	}
	dst.Arch = harness.ArchAmd64 // we assume amd64 for now
	return dst
}

// copyEnv returns a copy of the environment variable map.
func copyEnv(src map[string]string) map[string]string {
	dst := map[string]string{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func convertSteps(src *v1.Job) []*harness.Step {
	var steps []*harness.Step
	for _, step := range src.Steps {
		if isCheckoutAction(step.Uses) {
			continue
		}
		dst := &harness.Step{
			Name: step.Name,
		}

		if step.Uses != "" {
			dst.Name = step.Name
			dst.Spec = convertAction(step)
			dst.Type = "action"
		} else {
			dst.Name = step.Name
			dst.Spec = convertRun(step)
			dst.Type = "script"
		}
		steps = append(steps, dst)
	}
	return steps
}

func convertAction(src *v1.Step) *harness.StepAction {
	dst := &harness.StepAction{
		Uses: src.Uses,
		With: src.With,
		Envs: src.Environment,
	}
	return dst
}

func convertRun(src *v1.Step) *harness.StepExec {
	dst := &harness.StepExec{
		Run:  src.Run,
		Envs: src.Environment,
	}
	return dst
}

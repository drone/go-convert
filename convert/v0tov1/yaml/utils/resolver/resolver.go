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

package resolver

import (
	"errors"

	"github.com/bradrydzewski/spec/yaml"
	"github.com/bradrydzewski/spec/yaml/utils/walk"
)

// Lookup returns a resource by name, kind and type.
type LookupFunc func(name string) (*yaml.Template, error)

// Resolve resolves the pipeline templates.
func Resolve(in *yaml.Schema, fn LookupFunc) error {
	walk.Walk(in, func(v interface{}) error {
		switch vv := v.(type) {
		case *yaml.Stage:
			ResolveStage(vv, fn)
		case *yaml.Step:
			ResolveStep(vv, fn)
		}
		return nil
	})

	return nil
}

// ResolveStage resolves the stage template.
func ResolveStage(stage *yaml.Stage, fn LookupFunc) error {
	// exit if not a template
	if stage.Template == nil {
		return nil
	}

	// lookup the template
	template, err := fn(stage.Template.Uses)
	if err != nil {
		return err
	}

	// template must have a child step
	if template.Stage == nil {
		return errors.New("template does not contain a stage")
	}

	// create the context if not exists
	if stage.Context == nil {
		stage.Context = new(yaml.Context)
	}

	// append the inputs to the context

	// first we append the default inputs defined
	// in the template.
	stage.Context.Inputs = map[string]any{}
	for k, v := range template.Inputs {
		if v.Default != nil {
			stage.Context.Inputs[k] = v.Default
		}
	}
	// then we append the inputs that are passed
	// to the template in the uses section.
	for k, v := range stage.Template.With {
		stage.Context.Inputs[k] = v
	}

	// merge the parent stage with the template values
	if stage.Approval == nil {
		stage.Approval = template.Stage.Approval
	}
	if stage.Cache == nil {
		stage.Cache = template.Stage.Cache
	}
	if stage.Clone == nil {
		stage.Clone = template.Stage.Clone
	}
	if stage.Concurrency == nil {
		stage.Concurrency = template.Stage.Concurrency
	}
	if stage.Delegate == "" {
		stage.Delegate = template.Stage.Delegate
	}
	if len(stage.Env) == 0 { // TODO append
		stage.Env = template.Stage.Env
	}
	if stage.Environment == nil {
		stage.Environment = template.Stage.Environment
	}
	stage.Group = template.Stage.Group
	if stage.Id == "" {
		stage.Id = template.Stage.Id
	}
	if stage.If == "" {
		stage.If = template.Stage.If
	}
	if stage.Name == "" {
		stage.Name = template.Stage.Name
	}
	if len(stage.Needs) == 0 {
		stage.Needs = template.Stage.Needs
	}
	if stage.OnFailure == nil {
		stage.OnFailure = template.Stage.OnFailure
	}
	if stage.Outputs == nil { // TODO append
		stage.Outputs = template.Stage.Outputs
	}
	stage.Parallel = template.Stage.Parallel
	if stage.Permissions == nil {
		stage.Permissions = template.Stage.Permissions
	}
	if stage.Platform == nil {
		stage.Platform = template.Stage.Platform
	}
	if stage.Rollback == nil { // TODO append
		stage.Rollback = template.Stage.Rollback
	}
	if stage.RunsOn == "" {
		stage.RunsOn = template.Stage.RunsOn
	}
	if stage.Runtime == nil {
		stage.Runtime = template.Stage.Runtime
	}
	if stage.Service == nil {
		stage.Service = template.Stage.Service
	}
	if len(stage.Services) == 0 { // TODO append
		stage.Services = template.Stage.Services
	}
	if stage.Status == nil {
		stage.Status = template.Stage.Status
	}
	stage.Steps = template.Stage.Steps
	if stage.Strategy == nil {
		stage.Strategy = template.Stage.Strategy
	}
	stage.Template = template.Stage.Template
	if stage.Volumes == nil { // TODO append
		stage.Volumes = template.Stage.Volumes
	}
	if stage.Workspace == nil {
		stage.Workspace = template.Stage.Workspace
	}

	return nil
}

// ResolveStep resolves the step template.
func ResolveStep(step *yaml.Step, fn LookupFunc) error {
	// exit if not a template
	if step.Template == nil {
		return nil
	}

	// lookup the yaml
	template, err := fn(step.Template.Uses)
	if err != nil {
		return err
	}

	// template must have a child step
	if template.Step == nil {
		return errors.New("template does not contain a step")
	}

	// create the context if not exists
	if step.Context == nil {
		step.Context = new(yaml.Context)
	}

	// append the inputs to the context

	// first we append the default inputs defined
	// in the template.
	step.Context.Inputs = map[string]any{}
	for k, v := range template.Inputs {
		if v.Default != nil {
			step.Context.Inputs[k] = v.Default
		}
	}
	// then we append the inputs that are passed
	// to the template in the uses section.
	for k, v := range step.Template.With {
		step.Context.Inputs[k] = v
	}

	// merge the parent step with the template values
	step.Action = template.Step.Action
	step.Approval = template.Step.Approval
	step.Background = template.Step.Background
	step.Barrier = template.Step.Barrier
	if step.Delegate == nil {
		step.Delegate = template.Step.Delegate
	}
	if len(step.Env) == 0 { // TODO append
		step.Env = template.Step.Env
	}
	step.Group = template.Step.Group
	if step.Id == "" {
		step.Id = template.Step.Id
	}
	if step.If == "" {
		step.If = template.Step.If
	}
	if step.Name == "" {
		step.Name = template.Step.Name
	}
	if len(step.Needs) == 0 {
		step.Needs = template.Step.Needs
	}
	if step.OnFailure == nil {
		step.OnFailure = template.Step.OnFailure
	}
	step.Parallel = template.Step.Parallel
	step.Queue = template.Step.Queue
	step.Run = template.Step.Run
	step.RunTest = template.Step.RunTest
	step.Queue = template.Step.Queue
	if step.Status == nil {
		step.Status = template.Step.Status
	}
	if step.Strategy == nil {
		step.Strategy = template.Step.Strategy
	}
	step.Template = template.Step.Template
	if step.Timeout == 0 {
		step.Timeout = template.Step.Timeout
	}

	return nil
}

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
	"time"

	github "github.com/drone/go-convert/convert/github/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	pipeline *github.Pipeline
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
	src, err := github.Parse(r)
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

	// create the harness pipeline spec
	pipeline := &harness.Pipeline{
		Stages: []*harness.Stage{},
	}

	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    pipeline,
	}

	// TODO pipeline.name removed from spec
	// pipeline.Name = ctx.pipeline.Name

	if ctx.pipeline.Env != nil {
		pipeline.Options = &harness.Default{
			Envs: ctx.pipeline.Env,
		}
	}

	//pipeline.When = convertOn(from.On) //GAP

	if ctx.pipeline.Jobs != nil {
		for name, job := range ctx.pipeline.Jobs {
			// skip nil jobs to avoid nil-pointer
			if job == nil {
				continue
			}

			var cloneStage *harness.CloneStage
			for _, step := range job.Steps {
				cloneStage = convertClone(step)
				if cloneStage != nil {
					break
				}
			}

			pipeline.Stages = append(pipeline.Stages, &harness.Stage{
				Name:     name,
				Type:     "ci",
				Strategy: convertStrategy(job.Strategy),
				When:     convertIf(job.If),
				Spec: &harness.StageCI{
					Clone:    cloneStage,
					Envs:     job.Env,
					Platform: convertRunsOn(job.RunsOn),
					Runtime: &harness.Runtime{
						Type: "cloud",
						Spec: &harness.RuntimeCloud{},
					},
					Steps: convertSteps(job),
					//Volumes:  convertVolumes(from.Volumes),

					// TODO support for delegate.selectors from.Node
					// TODO support for stage.variables
				},
			})
		}
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func convertClone(src *github.Step) *harness.CloneStage {
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

func convertOn(src *github.On) *harness.When {
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

func convertIf(i string) *harness.When {
	if i == "" {
		return nil
	}

	i = githubExprToJexlExpr(i)

	dst := new(harness.When)
	dst.Eval = i
	return dst
}

func githubExprToJexlExpr(githubExpr string) string {
	// Replace functions
	githubExpr = strings.Replace(githubExpr, "!contains(", "!~ ", -1)
	githubExpr = strings.Replace(githubExpr, "contains(", "=~ ", -1)
	githubExpr = strings.Replace(githubExpr, "startsWith(", "=^ ", -1)
	githubExpr = strings.Replace(githubExpr, "endsWith(", "=$ ", -1)

	// Replace variables
	githubExpr = strings.Replace(githubExpr, "github.event_name", "<+trigger.event>", -1)
	githubExpr = strings.Replace(githubExpr, "github.ref", "<+trigger.payload.ref>", -1)
	githubExpr = strings.Replace(githubExpr, "github.head_ref", "<+trigger.sourceBranch>", -1)
	githubExpr = strings.Replace(githubExpr, "github.event.ref", "<+trigger.payload.ref>", -1)
	githubExpr = strings.Replace(githubExpr, "github.base_ref", "<+trigger.targetBranch>", -1)
	githubExpr = strings.Replace(githubExpr, "github.event.number", "<+trigger.prNumber>", -1)
	githubExpr = strings.Replace(githubExpr, "github.event.pull_request.title", "<+trigger.prTitle>", -1)
	githubExpr = strings.Replace(githubExpr, "github.event.pull_request.body", "<+trigger.payload.pull_request.body>", -1)
	githubExpr = strings.Replace(githubExpr, "github.event.pull_request.html_url", "<+trigger.payload.pull_request.html_url>", -1)
	githubExpr = strings.Replace(githubExpr, "github.event.repository.html_url", "<+trigger.repoUrl>", -1)
	githubExpr = strings.Replace(githubExpr, "github.actor", "<+trigger.gitUser>", -1)
	githubExpr = strings.Replace(githubExpr, "github.actor_email", "<+codebase.gitUserEmail>", -1)

	return githubExpr
}

func getEventConditions(src *github.On) map[string][]string {
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

func isTriggersEmpty(src *github.On) bool {
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
		dst.Os = harness.OSWindows.String()
	case strings.Contains(src, "darwin"), strings.Contains(src, "macos"), strings.Contains(src, "mac"):
		dst.Os = harness.OSDarwin.String()
	default:
		dst.Os = harness.OSLinux.String()
	}
	dst.Arch = harness.ArchAmd64.String() // we assume amd64 for now
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

func convertSteps(src *github.Job) []*harness.Step {
	var steps []*harness.Step
	for serviceName, service := range src.Services {
		if service != nil {
			steps = append(steps, convertServices(service, serviceName))
		}
	}
	for _, step := range src.Steps {
		if isCheckoutAction(step.Uses) {
			continue
		}
		dst := &harness.Step{
			Name: step.Name,
		}

		if step.ContinueOnErr {
			dst.Failure = convertContinueOnError(step)
		}

		if step.Timeout != 0 {
			dst.Timeout = convertTimeout(step)
		}

		if step.Uses != "" {
			dst.Name = step.Name
			dst.Spec = convertAction(step)
			dst.Type = "action"
		} else {
			dst.Name = step.Name
			dst.Spec = convertRun(step, src.Container)
			dst.Type = "script"
		}
		steps = append(steps, dst)
	}
	return steps
}

func convertAction(src *github.Step) *harness.StepAction {
	if src == nil {
		return nil
	}
	dst := &harness.StepAction{
		Uses: src.Uses,
		With: make(map[string]interface{}),
		Envs: src.Env,
	}
	for key, value := range src.With {
		switch v := value.(type) {
		case float64:
			dst.With[key] = fmt.Sprintf("%v", v)
		case bool:
			dst.With[key] = value
		case string:
			if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
				dst.With[key] = value
			} else {
				dst.With[key] = fmt.Sprintf("%s", v)
			}
		default:
			dst.With[key] = value
		}
	}
	return dst
}

func convertContinueOnError(src *github.Step) *harness.FailureList {
	if !src.ContinueOnErr {
		return nil
	}

	return &harness.FailureList{
		Items: []*harness.Failure{
			{
				Errors: []string{"all"},
				Action: &harness.FailureAction{
					Type: "ignore",
					Spec: &harness.Ignore{},
				},
			},
		},
	}
}

func convertRun(src *github.Step, container *github.Container) *harness.StepExec {
	if src == nil {
		return nil
	}
	dst := &harness.StepExec{
		Run:  src.Run,
		Envs: src.Env,
	}
	if container != nil {
		dst.Image = container.Image
	}
	return dst
}

func convertServices(service *github.Service, serviceName string) *harness.Step {
	if service == nil {
		return nil
	}
	return &harness.Step{
		Name: serviceName,
		Type: "background",
		Spec: &harness.StepBackground{
			Image: service.Image,
			Envs:  service.Env,
			Mount: convertMounts(service.Volumes),
			Ports: service.Ports,
			Args:  service.Options,
		},
	}
}

func convertTimeout(src *github.Step) string {
	if src == nil || src.Timeout == 0 {
		return "0"
	}
	return fmt.Sprint(time.Duration(src.Timeout * int(time.Minute)))
}

func convertMounts(volumes []string) []*harness.Mount {
	if len(volumes) == 0 {
		return nil
	}
	var dst []*harness.Mount

	for _, volume := range volumes {
		parts := strings.Split(volume, ":")

		var mount harness.Mount
		if len(parts) > 1 {
			mount.Name = parts[0]
			mount.Path = parts[1]
		} else {
			mount.Path = parts[0]
		}

		dst = append(dst, &mount)
	}

	return dst
}

func convertStrategy(src *github.Strategy) *harness.Strategy {
	if src == nil || src.Matrix == nil {
		return nil
	}

	matrix := src.Matrix

	includeMaps := convertInterfaceMapsToStringMaps(matrix.Include)
	excludeMaps := convertInterfaceMapsToStringMaps(matrix.Exclude)
	dst := &harness.Strategy{
		Type: "matrix",
		Spec: &harness.Matrix{
			Axis:    matrix.Matrix,
			Include: includeMaps,
			Exclude: excludeMaps,
		},
	}
	return dst
}

func convertInterfaceMapsToStringMaps(maps []map[string]interface{}) []map[string]string {
	convertedMaps := make([]map[string]string, len(maps))
	for i, originalMap := range maps {
		convertedMap := make(map[string]string)
		for key, value := range originalMap {
			convertedMap[key] = fmt.Sprintf("%v", value)
		}
		convertedMaps[i] = convertedMap
	}
	return convertedMaps
}

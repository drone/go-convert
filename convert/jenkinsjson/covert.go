// Copyright 2023 Harness, Inc.
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

// Package jenkinsjson converts jenkinsjson pipelines to Harness pipelines.
package jenkinsjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	harness "github.com/drone/spec/dist/go"
	jenkinsjson "github.com/jamie-harness/go-convert/convert/jenkinsjson/json"
	"github.com/jamie-harness/go-convert/internal/store"
	"gopkg.in/yaml.v2"
)

// conversion context
type context struct {
	// config *jenkinsjson.Pipeline
	// job    *jenkinsjson.Job
}

// Converter converts a jenkinsjson pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers
}

// New creates a new Converter that converts a jenkinsjson
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
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	var pipelineJson jenkinsjson.Node
	err := json.Unmarshal(buf.Bytes(), &pipelineJson)

	// create the harness pipeline spec
	dst := &harness.Pipeline{}

	stepsInStage := make([]*harness.Step, 0)

	recursiveParseJsonToSteps(pipelineJson, &stepsInStage)

	// create the harness stage.
	dstStage := &harness.Stage{
		Name: pipelineJson.Name,
		Id:   SanitizeForId(pipelineJson.SpanName, pipelineJson.SpanId),
		Type: "ci",
		// When: convertCond(from.Trigger),
		Spec: &harness.StageCI{
			// Delegate: convertNode(from.Node),
			// Envs: convertVariables(ctx.config.Variables),
			// Platform: convertPlatform(from.Platform),
			// Runtime:  convertRuntime(from),
			Steps: stepsInStage, // Initialize the Steps slice
			// Clone: clone,
			// Repository: repo,
		},
	}
	dst.Stages = append(dst.Stages, dstStage)

	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    dst,
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
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

// converts converts a jenkinsjson pipeline to a Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {

	return nil, nil
}

func recursiveParseJsonToSteps(currentNode jenkinsjson.Node, steps *[]*harness.Step) (*harness.CloneStage, *harness.Repository) {
	var clone *harness.CloneStage
	var repo *harness.Repository

	if len(currentNode.AttributesMap) == 0 {
		for _, child := range currentNode.Children {
			clone, repo = recursiveParseJsonToSteps(child, steps)
		}
	}

	switch currentNode.AttributesMap["jenkins.pipeline.step.type"] {
	case "node", "parallel":
		// node and parallel are wrapper layer to hold actual steps
		// parallel should be handled at it's parent
		for _, child := range currentNode.Children {
			clone, repo = recursiveParseJsonToSteps(child, steps)
		}
	case "stage":
		// this is techinically a step group, we treat it as just steps for now
		if len(currentNode.Children) > 1 {
			// handle parallel from parent
			if currentNode.Children[0].AttributesMap["jenkins.pipeline.step.type"] == "parallel" {
				parallelStepItems := make([]*harness.Step, 0)
				for _, child := range currentNode.Children {
					clone, repo = recursiveParseJsonToSteps(child, &parallelStepItems)
				}
				parallelStep := &harness.Step{
					Name: currentNode.SpanName,
					Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
					Type: "parallel",
					Spec: &harness.StepParallel{
						Steps: parallelStepItems,
					},
				}
				*steps = append(*steps, parallelStep)
			} else {
				for _, child := range currentNode.Children {
					clone, repo = recursiveParseJsonToSteps(child, steps)
				}
			}
		} else {
			for _, child := range currentNode.Children {
				clone, repo = recursiveParseJsonToSteps(child, steps)
			}
		}
	case "sh":
		*steps = append(*steps, jenkinsjson.ConvertSh(currentNode))
	case "junit":
		*steps = append(*steps, jenkinsjson.ConvertJunit(currentNode))
	case "git":
		clone, repo = jenkinsjson.ConvertClone(currentNode)
	case "":
	default:
		*steps = append(*steps, &harness.Step{
			Name: currentNode.SpanName,
			Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
			Type: "script",
			Spec: &harness.StepExec{
				Shell: "sh",
				Run:   fmt.Sprintf("echo %q", "This is a place holder for: "+currentNode.AttributesMap["jenkins.pipeline.step.type"]),
			},
			Desc: "This is a place holder for: " + currentNode.AttributesMap["jenkins.pipeline.step.type"],
		})
	}
	return clone, repo
}

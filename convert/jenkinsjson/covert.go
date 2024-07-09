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
	"strings"

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

// Struct used to store the stages so they can be sorted
type StageWithID struct {
	Stage *harness.Stage
	ID    int
}

type ProcessedTools struct {
	ToolProcessed      bool
	MavenPresent       bool
	GradlePresent      bool
	AntPresent         bool
	MavenProcessed     bool
	GradleProcessed    bool
	AntProcessed       bool
	SonarCubeProcessed bool
	SonarCubePresent   bool
}

var mavenGoals string
var gradleGoals string

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

	processedTools := &ProcessedTools{false, false, false, false, false, false, false, false, false}
	var variable map[string]string
	recursiveParseJsonToStages(&pipelineJson, dst, processedTools, variable)
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

// Recursive function to parse JSON nodes into stages and steps
func recursiveParseJsonToStages(jsonNode *jenkinsjson.Node, dst *harness.Pipeline, processedTools *ProcessedTools, variables map[string]string) {
	// Collect all stages with their IDs
	processedSteps := collectStagesWithID(jsonNode, processedTools, variables)

	// Sort the stages based on their IDs
	// sort.Slice(stagesWithID, func(i, j int) bool {
	// 	return stagesWithID[i].ID < stagesWithID[j].ID
	// })

	spec := &harness.StageCI{
		Steps: processedSteps,
	}

	stage := &harness.Stage{
		Name: "build",
		Id:   "build",
		Type: "ci",
		Spec: spec,
	}

	dst.Stages = append(dst.Stages, stage)
}

func collectStagesWithID(jsonNode *jenkinsjson.Node, processedTools *ProcessedTools, variables map[string]string) []*harness.Step {
	stagesWithID := make([]*harness.Step, 0)
	for _, childNode := range jsonNode.Children {
		if childNode.AttributesMap["jenkins.pipeline.step.type"] == "stage" {

			stageName, ok := childNode.ParameterMap["name"].(string)
			if !ok || stageName == "" {
				stageName = "Unnamed Stage"
			}

			// Skip stages with the name "Declarative Tool Install" since this does not contain any step
			if stageName == "Declarative: Tool Install" {
				continue
			}

			searchChildNodesForStageId(&childNode, stageName)

			stepsInStage := make([]*harness.Step, 0)
			recursiveParseJsonToSteps(childNode, &stepsInStage, processedTools, variables)

			// Create the harness stage for each Jenkins stage
			dstStep := &harness.Step{
				Name: stageName,
				Id:   SanitizeForId(childNode.SpanName, childNode.SpanId),
				Type: "group",
				Spec: &harness.StepGroup{
					Steps: stepsInStage,
				},
			}

			// Convert stageID to integer and store it with the stage
			// id, err := strconv.Atoi(stageID)
			// if err != nil {
			// 	fmt.Println("Error converting stage ID to integer:", err)
			// 	continue
			// }

			stagesWithID = append(stagesWithID, dstStep)
		} else {
			// Recursively process the children
			stagesWithID = append(stagesWithID, collectStagesWithID(&childNode, processedTools, variables)...)
		}
	}
	return stagesWithID
}

func searchChildNodesForStageId(node *jenkinsjson.Node, stageName string) string {
	for _, childNode := range node.Children {
		if childNode.AttributesMap["jenkins.pipeline.step.type"] == "stage" {
			name, ok := childNode.AttributesMap["jenkins.pipeline.step.name"]
			if ok && name == stageName {
				return childNode.AttributesMap["jenkins.pipeline.step.id"]
			}
		}
		// Recursively search in the children
		if id := searchChildNodesForStageId(&childNode, stageName); id != "" {
			return id
		}
	}
	return ""
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

func extractToolType(fullToolType string) string {
	knownToolTypes := []string{"MavenInstallation", "GradleInstallation", "AntInstallation", "SonarRunnerInstallation"}
	for _, knownType := range knownToolTypes {
		if strings.HasSuffix(fullToolType, knownType) {
			return knownType
		}
	}
	// Default return if no known type is matched
	return fullToolType
}

func recursiveParseJsonToSteps(currentNode jenkinsjson.Node, steps *[]*harness.Step, processedTools *ProcessedTools, variables map[string]string) (*harness.CloneStage, *harness.Repository) {
	var clone *harness.CloneStage
	var repo *harness.Repository

	if len(currentNode.AttributesMap) == 0 {
		for _, child := range currentNode.Children {
			clone, repo = recursiveParseJsonToSteps(child, steps, processedTools, variables)
		}
	}

	// Handle all tool nodes at the beginning
	handleTool(currentNode, processedTools)

	switch currentNode.AttributesMap["jenkins.pipeline.step.type"] {
	case "node", "parallel", "script":
		// node, parallel, withEnv, and withMaven are wrapper layers to hold actual steps
		// parallel should be handled at its parent
		for _, child := range currentNode.Children {
			clone, repo = recursiveParseJsonToSteps(child, steps, processedTools, variables)
		}
	case "withEnv":
		var1 := ExtractEnvironmentVariables(currentNode)
		if len(var1) > 0 {
			variables = var1
		}
		for _, child := range currentNode.Children {
			clone, repo = recursiveParseJsonToSteps(child, steps, processedTools, variables)
		}

	case "stage":
		// this is technically a step group, we treat it as just steps for now
		if len(currentNode.Children) > 1 {
			// handle parallel from parent
			if currentNode.Children[0].AttributesMap["jenkins.pipeline.step.type"] == "parallel" {
				parallelStepItems := make([]*harness.Step, 0)
				for _, child := range currentNode.Children {
					clone, repo = recursiveParseJsonToSteps(child, &parallelStepItems, processedTools, variables)
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
					clone, repo = recursiveParseJsonToSteps(child, steps, processedTools, variables)
				}
			}
		} else {
			for _, child := range currentNode.Children {
				clone, repo = recursiveParseJsonToSteps(child, steps, processedTools, variables)
			}
		}
	case "sh":
		*steps = append(*steps, jenkinsjson.ConvertSh(currentNode, variables))
	case "checkout":
		*steps = append(*steps, &harness.Step{
			Name: currentNode.SpanName,
			Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
			Type: "plugin",
			Spec: &harness.StepPlugin{
				Image: "checkout_plugin",
				With: map[string]interface{}{
					"platform": currentNode.AttributesMap["peer.service"],
					"git_url":  currentNode.AttributesMap["http.url"],
					"branch":   currentNode.AttributesMap["git.branch"],
					"depth":    currentNode.AttributesMap["git.clone.depth"],
				},
			},
		})
	case "archiveArtifacts":
		*steps = append(*steps, jenkinsjson.ConvertArchive(currentNode)...)
	case "junit":
		*steps = append(*steps, jenkinsjson.ConvertJunit(currentNode, variables))
	case "git":
	case "sleep":
		*steps = append(*steps, jenkinsjson.ConvertSleep(currentNode, variables))
	case "dir":
		dirPath := currentNode.ParameterMap["path"].(string)
		step := &harness.Step{
			Name: "Deletingdir",
			Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
			Type: "script",
			Spec: &harness.StepExec{
				Shell: "sh",
				Run:   fmt.Sprintf("rm -rf %s", dirPath),
			},
		}
		if len(variables) > 0 {
			step.Spec.(*harness.StepExec).Envs = variables
		}
		*steps = append(*steps, step)
	case "deleteDir":
		step := &harness.Step{
			Name: "Deletingdir",
			Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
			Type: "script",
			Spec: &harness.StepExec{
				Shell: "sh",
				Run: `
					dir_to_delete=$(pwd)
					cd ..
					rm -rf $dir_to_delete
					`,
			},
		}
		if len(variables) > 0 {
			step.Spec.(*harness.StepExec).Envs = variables
		}
		*steps = append(*steps, step)
	case "writeFile":
		*steps = append(*steps, jenkinsjson.ConvertWriteFile(currentNode, variables))
	case "readFile":
		*steps = append(*steps, jenkinsjson.ConvertReadFile(currentNode, variables))
	case "synopsys_detect":
		*steps = append(*steps, jenkinsjson.ConvertSynopsysDetect(currentNode, variables))
	case "":
	case "withAnt", "tool", "envVarsForTool":
	case "withMaven":
		clone, repo = recursiveHandleWithTool(currentNode, steps, processedTools, "maven", "maven", "harnesscommunitytest/maven-plugin:latest", variables)
	case "withGradle":
		clone, repo = recursiveHandleWithTool(currentNode, steps, processedTools, "gradle", "gradle", "harnesscommunitytest/gradle-plugin:latest", variables)
	case "wrap":
		if !processedTools.SonarCubeProcessed && !processedTools.SonarCubePresent {
			clone, repo = recursiveHandleSonarCube(currentNode, steps, processedTools, "sonarCube", "sonarCube", "aosapps/drone-sonar-plugin", variables)
		}
		if processedTools.AntPresent {
			clone, repo = recursiveHandleWithTool(currentNode, steps, processedTools, "ant", "ant", "harnesscommunitytest/ant-plugin:latest", variables)
		}

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

func ExtractEnvironmentVariables(node jenkinsjson.Node) map[string]string {
	envVars := make(map[string]string)
	if environment, ok := node.ParameterMap["overrides"].([]interface{}); ok && len(environment) > 0 {
		for _, envVar := range environment {
			if envVarStr, ok := envVar.(string); ok {
				parts := strings.SplitN(envVarStr, "=", 2)
				if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
					envVars[parts[0]] = parts[1]
				}
			}
		}
		return envVars
	}
	return envVars
}

func recursiveHandleWithTool(currentNode jenkinsjson.Node, steps *[]*harness.Step, processedTools *ProcessedTools, toolType string, pluginName string, pluginImage string, variables map[string]string) (*harness.CloneStage, *harness.Repository) {
	var clone *harness.CloneStage
	var repo *harness.Repository
	for _, child := range currentNode.Children {
		// Check if this child contains the type "sh"
		stepType, ok := child.AttributesMap["jenkins.pipeline.step.type"]
		if ok && stepType == "sh" {
			script, scriptOk := child.ParameterMap["script"]
			if scriptOk {
				// Check if the tool is not processed
				switch toolType {
				case "maven":
					if processedTools.MavenProcessed {
						continue
					}
				case "gradle":
					if processedTools.GradleProcessed {
						continue
					}
				case "ant":
					if processedTools.AntProcessed {
						continue
					}
				}

				value := script.(string) // Store the script value in the global variable
				goals := value
				words := strings.Fields(goals)
				goals = strings.Join(words[1:], " ")

				toolStep := &harness.Step{
					Name: pluginName,
					Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
					Type: "plugin",
					Spec: &harness.StepPlugin{
						Connector: "c.docker",
						Image:     pluginImage,
						With:      map[string]interface{}{"Goals": goals},
					},
				}
				if len(variables) > 0 {
					toolStep.Spec.(*harness.StepPlugin).Envs = variables
				}

				*steps = append(*steps, toolStep)

				// Mark the tool as processed
				switch toolType {
				case "maven":
					processedTools.MavenProcessed = true
				case "gradle":
					processedTools.GradleProcessed = true
				case "ant":
					processedTools.AntProcessed = true
				}

				// Stop recursion if the `sh` script is found and processed
				return clone, repo
			}
		}
		clone, repo = recursiveHandleWithTool(child, steps, processedTools, toolType, pluginName, pluginImage, variables)
	}
	return clone, repo
}

func handleTool(currentNode jenkinsjson.Node, processedTools *ProcessedTools) {
	// Check if this node is a "tool" node
	if stepType, ok := currentNode.AttributesMap["jenkins.pipeline.step.type"]; ok && stepType == "tool" {
		if attr, ok := currentNode.AttributesMap["harness-attribute"]; ok {
			toolAttributes := make(map[string]string)
			if err := json.Unmarshal([]byte(attr), &toolAttributes); err == nil {
				fullToolType := toolAttributes["type"]
				toolType := ""
				if parts := strings.Split(fullToolType, "$"); len(parts) > 1 {
					toolType = parts[1]
				} else {
					toolType = extractToolType(fullToolType)
				}

				switch toolType {
				case "MavenInstallation":
					processedTools.MavenPresent = true
				case "GradleInstallation":
					processedTools.GradlePresent = true
				case "AntInstallation":
					processedTools.AntPresent = true
				default:
					processedTools.ToolProcessed = true
				}
			}
		}
	}

	// Recursively call handleTool for all children
	for _, child := range currentNode.Children {
		handleTool(child, processedTools)
	}
}

func recursiveHandleSonarCube(currentNode jenkinsjson.Node, steps *[]*harness.Step, processedTools *ProcessedTools, toolType string, pluginName string, pluginImage string, variables map[string]string) (*harness.CloneStage, *harness.Repository) {
	var clone *harness.CloneStage
	var repo *harness.Repository

	// Check if this node contains the type "wrap" (for withSonarQubeEnv)
	stepType, ok := currentNode.AttributesMap["jenkins.pipeline.step.type"]
	if ok && stepType == "wrap" {
		// Extract delegate information
		delegate, delegateOk := currentNode.ParameterMap["delegate"]
		if delegateOk {
			symbol, symbolOk := delegate.(map[string]interface{})["symbol"].(string)
			if symbolOk && symbol == "withSonarQubeEnv" {
				// Handle withSonarQubeEnv step
				toolStep := &harness.Step{
					Name: pluginName,
					Id:   SanitizeForId(currentNode.SpanName, currentNode.SpanId),
					Type: "plugin",
					Spec: &harness.StepPlugin{
						Connector: "c.docker",
						Image:     pluginImage,
						With:      map[string]interface{}{"sonar_token": "<+input>", "sonar_host": "<+input>"},
					},
				}
				if len(variables) > 0 {
					toolStep.Spec.(*harness.StepPlugin).Envs = variables
				}

				*steps = append(*steps, toolStep)

				// Mark the tool as processed
				processedTools.SonarCubeProcessed = true
				processedTools.SonarCubePresent = true

				// Since we found and processed the withSonarQubeEnv step, return early
				return clone, repo
			}
		}
	}

	// Recursively handle other child nodes
	for _, child := range currentNode.Children {
		clone, repo = recursiveHandleSonarCube(child, steps, processedTools, toolType, pluginName, pluginImage, variables)
		if clone != nil || repo != nil {
			// If we found and processed the step, return
			return clone, repo
		}
	}

	return clone, repo
}

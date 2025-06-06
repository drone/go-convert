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
	"sort"
	"strconv"
	"strings"

	jenkinsjson "github.com/drone/go-convert/convert/jenkinsjson/json"
	"github.com/drone/go-convert/internal/store"
	harness "github.com/drone/spec/dist/go"
	"gopkg.in/yaml.v2"
)

// conversion context
type context struct {
	// config *jenkinsjson.Pipeline
	// job    *jenkinsjson.Job
}

// Struct used to store the stages so they can be sorted
type StepGroupWithID struct {
	Step *harness.Step
	ID   int
}
type StepWithID struct {
	Step *harness.Step
	ID   int
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
	Tags               []string
}

// In a unified pipeline trace, 'sh' types identified as branched/conditional will be pre-fixed with '_unifiedTraceBranch'
// This is done because we don't want these steps to be merged, so that pipeline analysis is easier.
const unifiedBranchedShStep = "sh_unifiedTraceBranch"

var tags = []string{
	"eks", "ec2",
	"gke", "gcp", "gcloud",
	"azure",
	"artifactory", "jfrog", "gcr", "gar",
	"java", "python", "go ", "ruby", "nodejs", "javascript", "typescript", "scala", "kotlin", "groovy", "csharp", "php", "perl", "bash", "powershell",
	"docker ", "git ", "npm ", "yarn ", "node ", "maven", "mvn ", "gradle", "sbt ", "bazel", "pip ", "dotnet ", "msbuild",
	"sfdx ",
	"hadoop", "mariadb", "mysql", "psql", "mongo", "redis", "jdbc",
}

var mavenGoals string
var gradleGoals string

var defaultWindowsImage string = "mcr.microsoft.com/powershell"

// Converter converts a jenkinsjson pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers

	// Infrastructure configuration
	infrastructure string // "cloud", "kubernetes", or "local"
	os             string // "linux", "mac", or "windows"
	arch           string // "amd64" or "arm64"

	useIntelligence bool
	configFile      string
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
	// Validate infrastructure options first
	if err := d.ValidateInfrastructureOptions(); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	var pipelineJson jenkinsjson.Node
	err := json.Unmarshal(buf.Bytes(), &pipelineJson)

	// create the harness pipeline spec
	dst := &harness.Pipeline{}

	processedTools := &ProcessedTools{false, false, false, false, false, false, false, false, false, []string{}}
	var variable map[string]string
	d.recursiveParseJsonToStages(&pipelineJson, dst, processedTools, variable)
	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Name:    jenkinsjson.SanitizeForName(pipelineJson.Name),
		Kind:    "pipeline",
		Type:    strings.Join(processedTools.Tags, ","),
		Spec:    dst,
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Recursive function to parse JSON nodes into stages and steps
func (d *Converter) recursiveParseJsonToStages(jsonNode *jenkinsjson.Node, dst *harness.Pipeline, processedTools *ProcessedTools, variables map[string]string) {
	stepGroupWithID := make([]StepGroupWithID, 0)
	var defaultDockerImage string = "alpine" // Default Docker Image
	// Collect all stages with their IDs
	collectStagesWithID(jsonNode, processedTools, &stepGroupWithID, variables, defaultDockerImage)
	// Sort the stages based on their IDs
	sort.Slice(stepGroupWithID, func(i, j int) bool {
		return stepGroupWithID[i].ID < stepGroupWithID[j].ID
	})
	sortedSteps := make([]*harness.Step, len(stepGroupWithID))
	for i, stepGroup := range stepGroupWithID {
		sortedSteps[i] = stepGroup.Step
	}

	// Create the StageCI spec with the sorted steps
	spec := &harness.StageCI{
		Steps: sortedSteps,
	}

	if d.useIntelligence {
		spec.Cache = &harness.Cache{
			Enabled: true,
		}
	}

	// Only add infrastructure configuration if any of the CLI options are specified
	if d.infrastructure != "" || d.os != "" || d.arch != "" {
		// Add Platform if OS or Arch is specified
		if d.os != "" || d.arch != "" {
			spec.Platform = &harness.Platform{}
			if d.os != "" {
				// Normalize OS names
				switch d.os {
				case "mac", "darwin", "macos":
					spec.Platform.Os = "macos"
				default:
					spec.Platform.Os = strings.ToLower(d.os)
				}
			}
			if d.arch != "" {
				spec.Platform.Arch = strings.ToLower(d.arch)
			}
		}

		// Add Runtime based on infrastructure type
		switch d.infrastructure {
		case "kubernetes", "k8s":
			spec.Runtime = &harness.Runtime{
				Type: "kubernetes",
			}
		case "local", "shell", "docker":
			spec.Runtime = &harness.Runtime{
				Type: "shell",
			}
		default: // "cloud"
			spec.Runtime = &harness.Runtime{
				Type: "cloud",
			}
		}
	}

	stage := &harness.Stage{
		Name: "build",
		Id:   "build",
		Type: "ci",
		Spec: spec,
	}

	dst.Stages = append(dst.Stages, stage)
}

func collectStagesWithID(jsonNode *jenkinsjson.Node, processedTools *ProcessedTools, stepGroupWithId *[]StepGroupWithID, variables map[string]string, dockerImage string) {
	if jsonNode.AttributesMap["jenkins.pipeline.step.type"] == "withDockerContainer" {
		if img, ok := jsonNode.ParameterMap["image"].(string); ok {
			dockerImage = img // Update Docker image from the node
		}
	}

	if jsonNode.AttributesMap["jenkins.pipeline.step.type"] == "stage" || jsonNode.AttributesMap["jenkins.pipeline.step.type"] == "node" {

		stageName := searchChildNodesForStageName(jsonNode)
		if stageName == "" {
			stageName = "Unnamed Stage"
		}

		// Skip stages with the name "Declarative Tool Install" since this does not contain any step
		if stageName == "Declarative: Tool Install" {
			return
		}

		stageID := searchChildNodesForStageId(jsonNode, stageName)

		stepsInStage := make([]*harness.Step, 0)
		for _, childNode := range jsonNode.Children {
			// Recursively process the children
			recursiveParseJsonToSteps(childNode, stepGroupWithId, &stepsInStage, processedTools, variables, dockerImage)
		}

		if len(stepsInStage) > 0 {
			// Create the stepGroup for the new stage
			dstStep := &harness.Step{
				Name: jenkinsjson.SanitizeForName(stageName),
				Id:   jenkinsjson.SanitizeForId(stageName, jsonNode.SpanId),
				Type: "group",
				Spec: &harness.StepGroup{
					Steps: stepsInStage,
				},
			}

			// identify technology tags
			for _, step := range stepsInStage {
				if step == nil {
					continue
				}
				switch step.Spec.(type) {
				case *harness.StepExec:
					exec := step.Spec.(*harness.StepExec)
					for _, tag := range tags {
						tagNoSpace := strings.TrimSpace(tag)
						if strings.Contains(strings.ToLower(exec.Run), tag) || strings.Contains(strings.ToLower(exec.Image), tagNoSpace) {
							processedTools.Tags = append(processedTools.Tags, tagNoSpace)
						}
					}
				}
			}

			// Convert stageID to integer and store it with the stage
			id, err := strconv.Atoi(stageID)
			if err != nil {
				fmt.Printf("Error converting stage ID to integer, spanId=%s: %v\n", jsonNode.SpanId, err)
				id = 0
			}

			*stepGroupWithId = append(*stepGroupWithId, StepGroupWithID{Step: dstStep, ID: id})
		}
	} else {
		for _, childNode := range jsonNode.Children {
			// Recursively process the children
			collectStagesWithID(&childNode, processedTools, stepGroupWithId, variables, dockerImage)
		}
	}
}

func searchChildNodesForStageId(node *jenkinsjson.Node, stageName string) string {
	name, ok := node.AttributesMap["jenkins.pipeline.step.name"]
	if ok && name == stageName {
		return node.AttributesMap["jenkins.pipeline.step.id"]
	}

	for _, childNode := range node.Children {
		// Recursively search in the children
		if id := searchChildNodesForStageId(&childNode, stageName); id != "" {
			return id
		}
	}
	return ""
}

func searchChildNodesForStageName(node *jenkinsjson.Node) string {
	if name, ok := node.ParameterMap["name"].(string); ok {
		return name
	}
	if name, ok := node.AttributesMap["jenkins.pipeline.step.name"]; ok {
		return name
	}
	for _, childNode := range node.Children {
		// Recursively search in the children
		if name := searchChildNodesForStageName(&childNode); name != "" {
			return name
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

func recursiveParseJsonToSteps(currentNode jenkinsjson.Node, stepGroupWithId *[]StepGroupWithID, steps *[]*harness.Step, processedTools *ProcessedTools, variables map[string]string, dockerImage string) (*harness.CloneStage, *harness.Repository) {
	stepWithIDList := make([]StepWithID, 0)

	var timeout string
	// Collect all steps with their IDs
	collectStepsWithID(currentNode, stepGroupWithId, &stepWithIDList, processedTools, variables, timeout, dockerImage)
	// Sort the steps based on their IDs
	sort.Slice(stepWithIDList, func(i, j int) bool {
		return stepWithIDList[i].ID < stepWithIDList[j].ID
	})
	mergeRunSteps(&stepWithIDList)

	sortedSteps := make([]*harness.Step, len(stepWithIDList))

	for i, step := range stepWithIDList {
		sortedSteps[i] = step.Step
	}

	*steps = append(*steps, sortedSteps...)

	return nil, nil
}

func collectStepsWithID(currentNode jenkinsjson.Node, stepGroupWithId *[]StepGroupWithID, stepWithIDList *[]StepWithID, processedTools *ProcessedTools, variables map[string]string, timeout string, dockerImage string) (*harness.CloneStage, *harness.Repository) {
	var clone *harness.CloneStage
	var repo *harness.Repository

	// parameterMap and harness-attribute are now aligned in recent trace files and always contain identical data.
	// However, if older traces can only contain harness-attribute, let's update parameterMap manually as we
	// use it exclusively for lookups below.
	if len(currentNode.ParameterMap) == 0 {
		if harnessAttr, ok := currentNode.AttributesMap["harness-attribute"]; ok {
			var paramMap map[string]interface{}
			err := json.Unmarshal([]byte(harnessAttr), &paramMap)
			if err == nil {
				currentNode.ParameterMap = paramMap
			} else {
				fmt.Println("Error parsing harness-attribute:", err)
			}
		}
	}

	// Handle all tool nodes at the beginning
	handleTool(currentNode, processedTools)

	stepId := currentNode.AttributesMap["jenkins.pipeline.step.id"]
	var id int
	if stepId != "" {
		var err error
		id, err = strconv.Atoi(stepId)
		if err != nil {
			fmt.Println("Error converting step ID to integer:", err)
		}
	}

	switch currentNode.AttributesMap["jenkins.pipeline.step.type"] {
	case "node":
		collectStagesWithID(&currentNode, processedTools, stepGroupWithId, variables, dockerImage)
		return clone, repo
	case "", "parallel", "script", "withAnt", "tool", "envVarsForTool", "ws", "ansiColor", "newBuildInfo", "getArtifactoryServer":
		// for spans where we should ignore its content and just process the children

	case "withEnv":
		var1 := ExtractEnvironmentVariables(currentNode)
		variables = mergeMaps(variables, var1)

	case "withCredentials":
		withCredentialsVars := jenkinsjson.ConvertWithCredentials(currentNode)
		variables = mergeMaps(variables, withCredentialsVars)

	case "withDockerContainer":
		if img, ok := currentNode.ParameterMap["image"].(string); ok {
			dockerImage = img
		}

	case "stage":
		// this is technically a step group, we treat it as just steps for now
		if len(currentNode.Children) > 0 {
			// handle parallel from parent
			var hasParallelStep = false
			for _, child := range currentNode.Children {
				if child.AttributesMap["jenkins.pipeline.step.type"] == "parallel" {
					hasParallelStep = true
					break
				}
			}

			if hasParallelStep {
				parallelStepGroupWithID := make([]StepGroupWithID, 0)
				collectStagesWithID(&currentNode, processedTools, &parallelStepGroupWithID, variables, dockerImage)
				sort.Slice(parallelStepGroupWithID, func(i, j int) bool {
					return parallelStepGroupWithID[i].ID < parallelStepGroupWithID[j].ID
				})

				sortedSteps := make([]*harness.Step, len(parallelStepGroupWithID))
				for i, stepGroup := range parallelStepGroupWithID {
					sortedSteps[i] = stepGroup.Step
				}

				parallelStep := &harness.Step{
					Name: currentNode.SpanName,
					Id:   jenkinsjson.SanitizeForId(currentNode.SpanName, currentNode.SpanId),
					Type: "parallel",
					Spec: &harness.StepParallel{
						Steps: sortedSteps,
					},
				}
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: parallelStep, ID: id})
			} else {
				collectStagesWithID(&currentNode, processedTools, stepGroupWithId, variables, dockerImage)
			}
			return clone, repo
		}

	case "sh":
		step := jenkinsjson.ConvertSh(currentNode, variables, timeout, dockerImage, "")
		if step != nil {
			*stepWithIDList = append(*stepWithIDList, StepWithID{Step: step, ID: id})
		}
	case unifiedBranchedShStep:
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSh(currentNode, variables, timeout, dockerImage, "_unifiedTraceBranch"), ID: id})
	case "bat":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertBat(currentNode, variables, timeout, defaultWindowsImage), ID: id})
	case "pwsh":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertPwsh(currentNode, variables, timeout, defaultWindowsImage), ID: id})
	case "powershell":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertPowerShell(currentNode, variables, timeout, defaultWindowsImage), ID: id})
	case "timeout":
		if len(currentNode.ParameterMap) > 0 {
			var unit string
			unitExists := false

			// Extract unit safely
			if val, ok := currentNode.ParameterMap["unit"].(string); ok {
				unit = val
				unitExists = true
			}

			// Extract and handle time safely
			var time int
			timeExists := false
			if val, ok := currentNode.ParameterMap["time"].(int); ok {
				time = val
				timeExists = true
			} else if val, ok := currentNode.ParameterMap["time"].(float64); ok {
				time = int(val)
				timeExists = true
			} else {
				fmt.Println("time is of a type I don't know how to handle")
			}

			if timeExists {
				if unit == "SECONDS" && time < 10 {
					time = 10
				}
				if !unitExists {
					unit = "MINUTES"
				}
				switch unit {
				case "MINUTES":
					timeout = strconv.Itoa(time) + "m"
				case "SECONDS":
					timeout = strconv.Itoa(time) + "s"
				case "HOURS":
					timeout = strconv.Itoa(time) + "h"
				case "DAYS":
					timeout = strconv.Itoa(time) + "d"
				default:
					timeout = strconv.Itoa(time) + unit
				}

			} else {
				fmt.Println("Time key does not exist in ParameterMap or its value is not a recognized type.")
			}
		}
	case "git", "checkout":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertCheckout(currentNode, variables), ID: id})
	case "archiveArtifacts":
		archiveSteps := jenkinsjson.ConvertArchive(currentNode)
		for _, step := range archiveSteps {
			*stepWithIDList = append(*stepWithIDList, StepWithID{Step: step, ID: id})
		}
	case "emailext":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertEmailext(currentNode, variables, timeout), ID: id})
		processedTools.Tags = append(processedTools.Tags, "email")
	case "junit":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertJunit(currentNode, variables), ID: id})
		processedTools.Tags = append(processedTools.Tags, "junit")
	case "sleep":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSleep(currentNode, variables), ID: id})
	case "dir":
		// dir can have repeated dir steps in its children, these do not include the ParametersMap
		if dirStep, ok := jenkinsjson.ConvertDir(currentNode, variables); ok {
			*stepWithIDList = append(*stepWithIDList, StepWithID{Step: dirStep, ID: id})
		}
	case "deleteDir":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertDeleteDir(currentNode, variables), ID: id})
	case "writeYaml":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertWriteYaml(currentNode, variables), ID: id})
	case "writeFile":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertWriteFile(currentNode, variables), ID: id})
	case "readFile":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertReadFile(currentNode, variables), ID: id})
	case "readJSON":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertReadJson(currentNode, variables), ID: id})
	case "verifySha1":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertVerifySha1(currentNode, variables, dockerImage), ID: id})
	case "readYaml":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertReadYaml(currentNode, variables), ID: id})
	case "sha256":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSHA256(currentNode, variables, dockerImage), ID: id})
	case "readCSV":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertReadCsv(currentNode, variables), ID: id})
	case "writeJSON":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertWriteJSON(currentNode, variables), ID: id})
	case "sha1":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSHA1(currentNode, variables, dockerImage), ID: id})
	case "synopsys_detect":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSynopsysDetect(currentNode, variables), ID: id})
	case "verifySha256":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertVerifySha256(currentNode, variables, dockerImage), ID: id})
	case "artifactoryUpload":
		specStr, ok := currentNode.ParameterMap["spec"].(string)
		if !ok {
			fmt.Println("Invalid or missing 'spec' in ParameterMap for artifactoryUpload")
			break
		}

		var spec jenkinsjson.Spec
		err := json.Unmarshal([]byte(specStr), &spec)
		if err != nil {
			fmt.Println("Error unmarshalling 'spec' for artifactoryUpload:", err)
			break
		}

		// Iterate over each file specification
		for _, fileSpec := range spec.Files {
			// Create a new node with 'pattern' and 'target' in ParameterMap
			newNode := currentNode
			newNode.ParameterMap["pattern"] = fileSpec.Pattern
			newNode.ParameterMap["target"] = fileSpec.Target

			// Convert and append the Artifactory Upload step
			convertedStep := jenkinsjson.ConvertArtifactUploadJfrog(newNode, variables, timeout)
			if convertedStep != nil {
				*stepWithIDList = append(*stepWithIDList, StepWithID{
					Step: convertedStep,
					ID:   id,
				})
			} else {
				fmt.Println("Failed to convert artifactoryUpload step")
			}
		}

	case "anchore":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertAnchore(currentNode, variables), ID: id})

	case "dockerPushStep", "rtDockerPush":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertDockerPushStep(currentNode, variables, timeout), ID: id})

	case "withMaven":
		clone, repo = recursiveHandleWithTool(currentNode, stepWithIDList, processedTools, "maven", "maven", "maven:latest", variables, timeout)
	case "withGradle":
		clone, repo = recursiveHandleWithTool(currentNode, stepWithIDList, processedTools, "gradle", "gradle", "gradle:latest", variables, timeout)
	case "wrap":
		if !processedTools.SonarCubeProcessed && !processedTools.SonarCubePresent {
			clone, repo = recursiveHandleSonarCube(currentNode, stepWithIDList, processedTools, "sonarCube", "sonarCube", "aosapps/drone-sonar-plugin", variables, timeout)
		}
		if processedTools.AntPresent {
			clone, repo = recursiveHandleWithTool(currentNode, stepWithIDList, processedTools, "ant", "ant", "frekele/ant:latest", variables, timeout)
		}
		if symbol, ok := currentNode.ParameterMap["delegate"].(map[string]interface{})["symbol"]; ok && symbol == "nodejs" {
			*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertNodejs(currentNode), ID: id})
		}
		return clone, repo

	case "zip":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertZip(currentNode, variables), ID: id})

	case "unzip":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertUnzip(currentNode, variables), ID: id})

	case "findFiles":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFindFiles(currentNode), ID: id})

	case "tar":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertTar(currentNode, variables), ID: id})

	case "untar":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertUntar(currentNode, variables), ID: id})
	case "fileOperations":
		// Step 1: Extract the 'delegate' map from the 'parameterMap'
		delegate, ok := currentNode.ParameterMap["delegate"].(map[string]interface{})
		if !ok {
			fmt.Println("Missing 'delegate' in parameterMap")
			break
		}

		// Step 2: Extract the 'arguments' map from the 'delegate'
		arguments, ok := delegate["arguments"].(map[string]interface{})
		if !ok {
			fmt.Println("Missing 'arguments' in delegate map")
			break
		}

		// Step 3: Extract the list of anonymous operations
		anonymousOps, ok := arguments["<anonymous>"].([]interface{})
		if !ok {
			fmt.Println("No anonymous operations found in arguments")
			break
		}

		// Step 4: Iterate over each operation and handle based on the 'symbol' type
		for _, op := range anonymousOps {
			// Convert the operation to a map for easy access
			operation, ok := op.(map[string]interface{})
			if !ok {
				fmt.Println("Invalid operation format")
				continue
			}

			// Extract the 'symbol' to determine the type of file operation
			symbol, ok := operation["symbol"].(string)
			if !ok {
				fmt.Println("Operation symbol not found or not a string")
				continue
			}

			// Step 5: Process each operation based on its 'symbol'
			switch symbol {
			case "fileCreateOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileCreate(currentNode, operation), ID: id})
			case "fileCopyOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileCopy(currentNode, operation), ID: id})
			case "fileDeleteOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileDelete(currentNode, operation), ID: id})
			case "fileDownloadOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileDownload(currentNode, operation), ID: id})
			case "fileRenameOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileRename(currentNode, operation), ID: id})
			case "filePropertiesToJsonOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileJson(currentNode, operation), ID: id})
			case "fileJoinOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileJoin(currentNode, operation), ID: id})
			case "fileTransformOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileTranform(currentNode, operation), ID: id})
			case "folderCopyOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFolderCopy(currentNode, operation), ID: id})
			case "folderCreateOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFolderCreate(currentNode, operation), ID: id})
			case "folderDeleteOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFolderDelete(currentNode, operation), ID: id})
			case "folderRenameOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFolderRename(currentNode, operation), ID: id})
			case "fileUnTarOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileUntar(currentNode, operation), ID: id})
			case "fileUnZipOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileUnzip(currentNode, operation), ID: id})
			case "fileZipOperation":
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFileZip(currentNode, operation), ID: id})
			default:
				fmt.Println("Unsupported file operation:", symbol)
			}
		}

	case "publishHTML":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertPublishHtml(currentNode, variables), ID: id})

	case "httpRequest":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertHttpRequest(currentNode, variables), ID: id})

	case "jacoco":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertJacoco(currentNode, variables), ID: id})

	case "cobertura":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertCobertura(currentNode, variables), ID: id})

	case "slackSend":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSlackSend(currentNode, variables), ID: id})

	case "slackUploadFile":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSlackUploadFile(currentNode, variables), ID: id})

	case "flywayrunner":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertFlywayRunner(currentNode, variables), ID: id})

	case "slackUserIdFromEmail":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSlackUserIdFromEmail(currentNode, variables), ID: id})

	case "slackUserIdsFromCommitters":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSlackUserIdsFromCommitters(currentNode, variables), ID: id})

	case "nexusArtifactUploader":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertNexusArtifactUploader(currentNode, variables), ID: id})

	case "rtDownload":
		fallthrough
	case "rtMavenRun":
		fallthrough
	case "rtGradleRun":
		fallthrough
	case "rtPublishBuildInfo":
		fallthrough
	case "rtPromote":
		fallthrough
	case "xrayScan":
		step := jenkinsjson.ConvertArtifactoryRtCommand(currentNode.AttributesMap["jenkins.pipeline.step.type"], currentNode, variables)
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: step, ID: id})

	case "testResultsAggregator":
		*stepWithIDList = append(*stepWithIDList,
			StepWithID{Step: jenkinsjson.ConvertTestResultsAggregator(currentNode, variables), ID: id})

	case "readMavenPom":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertReadMavenPom(currentNode), ID: id})

	case "jiraSendBuildInfo":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertJiraBuildInfo(currentNode, variables), ID: id})

	case "jiraSendDeploymentInfo":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertJiraDeploymentInfo(currentNode, variables), ID: id})

	case "nunit":
		// Step 1: Extract the 'delegate' map from the 'parameterMap'
		delegate, ok := currentNode.ParameterMap["delegate"].(map[string]interface{})
		if !ok {
			fmt.Println("Missing 'delegate' in parameterMap")
			break
		}

		// Step 2: Extract the 'arguments' map from the 'delegate'
		arguments, ok := delegate["arguments"].(map[string]interface{})
		if !ok {
			fmt.Println("Missing 'arguments' in delegate map")
			break
		}
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertNunit(currentNode, arguments), ID: id})

	case "s3Upload":
		entries := jenkinsjson.ExtractEntries(currentNode)
		if entries == nil {
			fmt.Println("No entries exists for s3Upload:collectStepsWithID")
			break
		}
		// Initialize an index counter
		index := 0
		// Iterate over each entry
		for _, entry := range entries {
			gzipFlag, ok := entry["gzipFiles"].(bool)
			if !ok {
				// Set default value to false if the key does not exist or has an invalid value
				gzipFlag = false
			}

			// Call function to handle gzip and upload logic
			if gzipFlag {
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.Converts3Archive(currentNode, entry, index), ID: id})
			}

			*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.Converts3Upload(currentNode, entry, index), ID: id})
			// Increment the index for each entry
			index++
		}

	case "allure":
		steps := jenkinsjson.ConvertAllureSteps(currentNode)
		for _, step := range steps {
			*stepWithIDList = append(*stepWithIDList, StepWithID{Step: step, ID: id})
		}

	case "withKubeConfig":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertKubeCtl(currentNode, currentNode.ParameterMap), ID: id})
	case "mail":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertMailer(currentNode, currentNode.ParameterMap), ID: id})
		processedTools.Tags = append(processedTools.Tags, "email")

	case "pagerdutyChangeEvent":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertPagerDutyChangeEvent(currentNode, currentNode.ParameterMap), ID: id})

	case "pagerduty":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertPagerDuty(currentNode, currentNode.ParameterMap), ID: id})

	case "notifyEndpoints":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertNotification(currentNode, currentNode.ParameterMap), ID: id})

	case "gatlingArchive":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertGatling(currentNode), ID: id})

	case "unarchive":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertUnarchive(currentNode, currentNode.ParameterMap), ID: id})

	case "ansiblePlaybook":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertAnsiblePlaybook(currentNode, currentNode.ParameterMap), ID: id})

	case "ansibleVault":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertAnsibleVault(currentNode, currentNode.ParameterMap), ID: id})

	case "ansibleAdhoc":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertAnsibleAdhoc(currentNode, currentNode.ParameterMap), ID: id})

	case "waitForQualityGate":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertSonarQualityGate(currentNode), ID: id})

	case "testNG":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertTestng(currentNode, currentNode.ParameterMap), ID: id})

	case "cucumber":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertCucumber(currentNode, currentNode.ParameterMap), ID: id})

	case "robot":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertRobot(currentNode, currentNode.ParameterMap), ID: id})

	case "browserstack":
		if _, exists := currentNode.ParameterMap["credentialsId"]; exists {
			*stepWithIDList = append(*stepWithIDList, StepWithID{
				Step: jenkinsjson.ConvertBrowserStack(currentNode),
				ID:   id,
			})
		}

	case "readTrusted":
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: jenkinsjson.ConvertReadTrusted(currentNode), ID: id})

	default:
		placeholderStr := fmt.Sprintf("echo %q", "This is a place holder for: "+currentNode.AttributesMap["jenkins.pipeline.step.type"])
		b, err := json.MarshalIndent(currentNode.ParameterMap, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		placeholderStr += "\n" + prependCommentHashToLines(string(b))
		*stepWithIDList = append(*stepWithIDList, StepWithID{Step: &harness.Step{
			Name: jenkinsjson.SanitizeForName(currentNode.SpanName),
			Id:   jenkinsjson.SanitizeForId(currentNode.SpanName, currentNode.SpanId),
			Type: "script",
			Spec: &harness.StepExec{
				Shell: "sh",
				Run:   placeholderStr,
			},
			Desc: "This is a place holder for: " + currentNode.AttributesMap["jenkins.pipeline.step.type"],
		}, ID: id})
	}

	for _, child := range currentNode.Children {
		clone, repo = collectStepsWithID(child, stepGroupWithId, stepWithIDList, processedTools, variables, timeout, dockerImage)
	}

	return clone, repo
}
func mergeMaps(dest, src map[string]string) map[string]string {
	result := make(map[string]string)
	for key, value := range dest {
		result[key] = value
	}
	for key, value := range src {
		result[key] = value
	}
	return result
}

func prependCommentHashToLines(input string) string {
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		lines[i] = "# " + line
	}

	return "\n# Here is the parameterMap for this function:\n" + strings.Join(lines, "\n")
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

func recursiveHandleWithTool(currentNode jenkinsjson.Node, stepWithIDList *[]StepWithID, processedTools *ProcessedTools, toolType string, buildName string, buildImage string, variables map[string]string, timeout string) (*harness.CloneStage, *harness.Repository) {
	var clone *harness.CloneStage
	var repo *harness.Repository
	for _, child := range currentNode.Children {
		// Check if this child contains the type "sh"
		stepType, ok := child.AttributesMap["jenkins.pipeline.step.type"]
		if ok && (stepType == "sh" || stepType == unifiedBranchedShStep) {
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
				stepId := child.AttributesMap["jenkins.pipeline.step.id"]
				var id int
				if stepId != "" {
					var err error
					id, err = strconv.Atoi(stepId)
					if err != nil {
						fmt.Println("Error converting step ID to integer:", err)
					}
				}

				value := script.(string) // Store the script value in the global variable
				goals := value
				words := strings.Fields(goals)
				goals = strings.Join(words[1:], " ")

				toolStep := &harness.Step{
					Name:    buildName,
					Timeout: timeout,
					Id:      jenkinsjson.SanitizeForId(currentNode.SpanName, currentNode.SpanId),
					Type:    "script",
					Spec: &harness.StepExec{
						Shell: "sh",
						Image: buildImage,
						Run:   value,
					},
				}
				if len(variables) > 0 {
					toolStep.Spec.(*harness.StepExec).Envs = variables
				}
				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: toolStep, ID: id})

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
		clone, repo = recursiveHandleWithTool(child, stepWithIDList, processedTools, toolType, buildName, buildImage, variables, timeout)
	}
	return clone, repo
}

func handleTool(currentNode jenkinsjson.Node, processedTools *ProcessedTools) {
	// Check if this node is a "tool" node
	if stepType, ok := currentNode.AttributesMap["jenkins.pipeline.step.type"]; ok && stepType == "tool" {
		typeVal, typeExists := currentNode.ParameterMap["type"].(string)
		if typeExists {
			toolType := ""
			// Determine the tool type from 'typeVal'
			if parts := strings.Split(typeVal, "$"); len(parts) > 1 {
				toolType = parts[1]
			} else {
				toolType = extractToolType(typeVal)
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

	// Recursively call handleTool for all children
	for _, child := range currentNode.Children {
		handleTool(child, processedTools)
	}
}

func recursiveHandleSonarCube(currentNode jenkinsjson.Node, stepWithIDList *[]StepWithID, processedTools *ProcessedTools, toolType string, pluginName string, pluginImage string, variables map[string]string, timeout string) (*harness.CloneStage, *harness.Repository) {
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
					Name:    pluginName,
					Timeout: timeout,
					Id:      jenkinsjson.SanitizeForId(currentNode.SpanName, currentNode.SpanId),
					Type:    "plugin",
					Spec: &harness.StepPlugin{
						Connector: "c.docker",
						Image:     pluginImage,
						With:      map[string]interface{}{"sonar_token": "<+input>", "sonar_host": "<+input>"},
					},
				}

				stepId := currentNode.AttributesMap["jenkins.pipeline.step.id"]
				var id int
				if stepId != "" {
					var err error
					id, err = strconv.Atoi(stepId)
					if err != nil {
						fmt.Println("Error converting step ID to integer:", err)
					}
				}

				if len(variables) > 0 {
					toolStep.Spec.(*harness.StepPlugin).Envs = variables
				}

				*stepWithIDList = append(*stepWithIDList, StepWithID{Step: toolStep, ID: id})

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
		clone, repo = recursiveHandleSonarCube(child, stepWithIDList, processedTools, toolType, pluginName, pluginImage, variables, timeout)
		if clone != nil || repo != nil {
			// If we found and processed the step, return
			return clone, repo
		}
	}

	return clone, repo
}

func mergeRunSteps(steps *[]StepWithID) {
	if len(*steps) < 2 {
		return
	}

	merged := []StepWithID{}
	cursor := (*steps)[0]
	for i := 1; i < len(*steps); i++ {
		current := (*steps)[i]
		// if can merge, store all current content in cursor
		if canMergeSteps(cursor.Step, current.Step) {
			previousExec := cursor.Step.Spec.(*harness.StepExec)
			currentExec := current.Step.Spec.(*harness.StepExec)
			previousExec.Run += "\n" + currentExec.Run
			cursor.Step.Name = jenkinsjson.SanitizeForName(cursor.Step.Name + "_" + current.Step.Name)
		} else {
			// if not able to merge, push cursor and reset cursor to current one
			merged = append(merged, cursor)
			cursor = current
		}
	}

	merged = append(merged, cursor)
	*steps = merged
}

func canMergeSteps(step1, step2 *harness.Step) bool {
	if step1 == nil || step2 == nil {
		return false
	}
	if step1.Type != "script" || step2.Type != "script" {
		return false
	}

	exec1, ok1 := step1.Spec.(*harness.StepExec)
	exec2, ok2 := step2.Spec.(*harness.StepExec)

	if !ok1 || !ok2 {
		return false
	}

	return exec1.Image == exec2.Image &&
		exec1.Connector == exec2.Connector &&
		exec1.Shell == exec2.Shell &&
		ENVmapsEqual(exec1.Envs, exec2.Envs) &&
		exec1.Entrypoint == exec2.Entrypoint &&
		ARGSslicesEqual(exec1.Args, exec2.Args) &&
		exec1.Privileged == exec2.Privileged &&
		exec1.Network == exec2.Network
}

func ENVmapsEqual(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func ARGSslicesEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// ValidateInfrastructureOptions validates the infrastructure configuration options
func (d *Converter) ValidateInfrastructureOptions() error {
	// Validate infrastructure type
	if d.infrastructure != "" {
		switch d.infrastructure {
		case "cloud", "kubernetes", "k8s", "local":
			// valid values
		default:
			return fmt.Errorf("invalid infrastructure type: %s. Must be one of: cloud, kubernetes, k8s, local", d.infrastructure)
		}
	}

	// Validate architecture
	if d.arch != "" {
		switch d.arch {
		case "amd64", "arm64":
			// valid values
		default:
			return fmt.Errorf("invalid architecture: %s. Must be one of: amd64, arm64", d.arch)
		}
	}

	// Validate operating system
	if d.os != "" {
		switch d.os {
		case "linux", "windows", "mac", "darwin":
			// valid values
		default:
			return fmt.Errorf("invalid operating system: %s. Must be one of: linux, windows, mac, darwin", d.os)
		}
	}

	return nil
}

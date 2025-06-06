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
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/drone/go-convert/convert/jenkinsjson/json"
	"github.com/drone/go-convert/internal/rand"

	"github.com/drone/go-convert/convert/harness"
	"github.com/drone/go-convert/internal/slug"
	"github.com/drone/go-convert/internal/store"

	downgraderYaml "github.com/drone/go-convert/convert/harness/downgrader/yaml"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/spec/dist/go"
	"github.com/ghodss/yaml"
)

// Downgrader downgrades pipelines from the v0 harness
// configuration format to the v1 configuration format.
type Downgrader struct {
	codebaseName    string
	codebaseConn    string
	dockerhubConn   string
	kubeConnector   string
	kubeNamespace   string
	kubeEnabled     bool
	pipelineId      string
	pipelineName    string
	pipelineOrg     string
	pipelineProj    string
	defaultImage    string
	useIntelligence bool
	randomId        bool
	identifiers     *store.Identifiers
}

const MaxDepth = 100

var eventMap = map[string]map[string]string{
	"pull_request": {
		"jexl":            "<+trigger.event>",
		"event":           "PR",
		"operator":        "==",
		"inverseOperator": "!=",
	},
	"push": {
		"jexl":            "<+trigger.event>",
		"event":           "PUSH",
		"operator":        "==",
		"inverseOperator": "!=",
	},
	"tag": {
		"jexl":            "<+trigger.payload.ref>",
		"event":           "refs/tags/",
		"operator":        "=^",
		"inverseOperator": "!^",
	},
}

// New creates a new Downgrader that downgrades pipelines
// from the v0 harness configuration format to the v1
// configuration format.
func New(options ...Option) *Downgrader {
	d := new(Downgrader)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = store.New()

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default pipeline name.
	if d.pipelineName == "" {
		d.pipelineName = harness.DefaultName
	}

	// set the default pipeline id.
	if d.pipelineId == "" {
		if d.randomId {
			d.pipelineId = slug.CreateWithRandom(d.pipelineName)
		} else {
			d.pipelineId = slug.Create(d.pipelineName)
		}
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
	src, err := downgraderYaml.ParseBytes(b)
	if err != nil {
		return nil, err
	}
	return d.DowngradeFrom(src)
}

// DowngradeString downgrades a v1 pipeline.
func (d *Downgrader) DowngradeString(s string) ([]byte, error) {
	return d.Downgrade([]byte(s))
}

// DowngradeFile downgrades a v1 pipeline.
func (d *Downgrader) DowngradeFile(path string) ([]byte, error) {
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return d.Downgrade(out)
}

// DowngradeFrom downgrades a v1 pipeline object.
func (d *Downgrader) DowngradeFrom(src []*v1.Config) ([]byte, error) {
	return d.downgrade(src)
}

// downgrade downgrades a v1 pipeline.
func (d *Downgrader) downgrade(src []*v1.Config) ([]byte, error) {
	var buf bytes.Buffer
	for i, p := range src {
		config := new(v0.Config)

		config.Pipeline.ID = d.pipelineId
		config.Pipeline.Name = d.pipelineName
		if config.Pipeline.Name == harness.DefaultName && p.Name != "" {
			config.Pipeline.Name = p.Name
		}
		if config.Pipeline.ID == harness.DefaultName && p.Name != "" {
			config.Pipeline.ID = slug.Create(p.Name)
		}

		config.Pipeline.Org = d.pipelineOrg
		config.Pipeline.Project = d.pipelineProj
		config.Pipeline.Tags = make(map[string]string)
		for _, tag := range strings.Split(p.Type, ",") {
			if tag != "" {
				config.Pipeline.Tags[tag] = ""
			}
		}
		config.Pipeline.Props.CI.Codebase = v0.Codebase{
			Name:  d.codebaseName,
			Conn:  d.codebaseConn,
			Build: "<+input>",
		}
		// FIXME: this is subject to a nil pointer
		if p.Spec.(*v1.Pipeline).Options != nil {
			config.Pipeline.Variables = convertVariables(p.Spec.(*v1.Pipeline).Options.Envs)
		}

		d.populateGitConnectorIfApplicable(p, config)

		// convert stages
		// FIXME: this is subject to a nil pointer
		for _, stage := range p.Spec.(*v1.Pipeline).Stages {
			// skip nil stages. this is un-necessary, we have
			// this logic in place just to be safe.
			if stage == nil {
				continue
			}

			// skip stages that are not CI stages, for now
			if _, ok := stage.Spec.(*v1.StageCI); !ok {
				continue
			}

			// convert the stage and add to the list
			config.Pipeline.Stages = append(config.Pipeline.Stages, &v0.Stages{
				Stage: d.convertStage(stage),
			})
		}
		out, err := yaml.Marshal(config)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			buf.WriteString("\n---\n")
		}
		buf.Write(out)
	}
	return buf.Bytes(), nil
}

func (d *Downgrader) populateGitConnectorIfApplicable(p *v1.Config, config *v0.Config) {
	if !d.isGitCloneFirstAndOnly(p) {
		return
	}

	for _, stage := range p.Spec.(*v1.Pipeline).Stages {
		spec := stage.Spec.(*v1.StageCI)
		gitConnector, newSteps := d.extractGitConnectorConfig(spec.Steps)
		spec.Steps = newSteps
		if gitConnector != nil {
			//config.GitConnector = gitConnector
			config.Pipeline.Props.CI.Codebase.Name = gitConnector.Spec.Url
		}
	}
}

func (d *Downgrader) isGitCloneFirstAndOnly(p *v1.Config) bool {
	pipeline, ok := p.Spec.(*v1.Pipeline)
	if !ok {
		return false
	}
	if !d.hasExactlyOneGitClone(pipeline) {
		return false
	}

	return d.isGitCloneFirstStep(pipeline)
}

func (d *Downgrader) isGitCloneFirstStep(pipeline *v1.Pipeline) bool {
	if len(pipeline.Stages) == 0 {
		return false
	}

	spec, ok := pipeline.Stages[0].Spec.(*v1.StageCI)
	if !ok {
		return false
	}
	if len(spec.Steps) == 0 {
		return false
	}

	step := spec.Steps[0]
	if step.Type == "group" {
		spec := step.Spec.(*v1.StepGroup)
		if len(spec.Steps) == 0 {
			return false
		}
		step = spec.Steps[0]
	}

	if isGitCloneStep(step) {
		return true
	}
	return false
}

func isGitCloneStep(step *v1.Step) bool {
	if step.Type != "plugin" {
		return false
	}
	pluginSpec, ok := step.Spec.(*v1.StepPlugin)
	return ok && pluginSpec.Image == harness.GitPluginImage
}

func (d *Downgrader) hasExactlyOneGitClone(pipeline *v1.Pipeline) bool {
	gitConnectorCount := 0
	for _, stage := range pipeline.Stages {
		spec, ok := stage.Spec.(*v1.StageCI)
		if !ok {
			continue
		}
		gitConnectorCount += d.countGitConnectors(spec.Steps)
	}

	return gitConnectorCount == 1
}

func (d *Downgrader) countGitConnectors(steps []*v1.Step) int {
	gitConnectorCount := 0
	for _, step := range steps {
		if step.Type == "group" {
			spec := step.Spec.(*v1.StepGroup)
			gitConnectorCount += d.countGitConnectors(spec.Steps)
		} else if isGitCloneStep(step) {
			gitConnectorCount++
		}
	}
	return gitConnectorCount
}

func (d *Downgrader) extractGitConnectorConfig(steps []*v1.Step) (*v0.GitConnector, []*v1.Step) {
	for stepIdx, step := range steps {
		if step.Type == "group" {
			spec := step.Spec.(*v1.StepGroup)
			connector, newSteps := d.extractGitConnectorConfig(spec.Steps)
			spec.Steps = newSteps
			if connector != nil {
				steps = slices.Delete(steps, stepIdx, stepIdx+1)
				return connector, steps
			}
		} else if isGitCloneStep(step) {
			pluginSpec := step.Spec.(*v1.StepPlugin)
			gitUrl, ok := pluginSpec.With["git_url"].(string)
			if !ok {
				break
			}

			r := regexp.MustCompile(".*/([\\w-.]+)(.git)?")
			repoNameMatch := r.FindStringSubmatch(gitUrl)
			if len(repoNameMatch) <= 1 {
				break
			}
			repoName := json.SanitizeForName(repoNameMatch[1])

			username, ok := pluginSpec.With["git.username"].(string)
			gitConnector := &v0.GitConnector{
				ID:      json.SanitizeForId(repoName, rand.Alphanumeric(8)),
				Name:    repoName,
				Org:     d.pipelineOrg,
				Project: d.pipelineProj,
				Type:    "Git",
				Spec: &v0.GitConnectorSpec{
					Url:            gitUrl,
					Type:           "Http",
					ConnectionType: "Repo",
					Spec: &v0.GitConnectionSpec{
						Username:    username,
						PasswordRef: "account." + repoName,
					},
				},
			}
			steps = slices.Delete(steps, stepIdx, stepIdx+1)
			return gitConnector, steps
		}
	}
	return nil, steps
}

// helper function converts a drone pipeline stage to a
// harness stage.
//
// TODO env variables to vars (stage-level)
// TODO delegate selectors
// TODO tags
// TODO when
// TODO failure strategy
// TODO volumes
// TODO if no stage clone, use global clone, if exists
func (d *Downgrader) convertStage(stage *v1.Stage) *v0.Stage {
	// extract the spec from the v1 stage
	spec := stage.Spec.(*v1.StageCI)

	var steps []*v0.Steps
	// convert each drone step to a harness step.
	for _, v := range spec.Steps {
		steps = append(steps, d.convertStep(v))
	}

	// enable clone by default
	enableClone := true
	if spec.Clone != nil && spec.Clone.Disabled == true {
		enableClone = false
	}

	// convert volumes
	if len(spec.Volumes) > 0 {
		// TODO
	}

	//
	// START TODO - refactor this into a helper function
	//

	var infra *v0.Infrastructure
	var runtime *v0.Runtime

	if spec.Runtime != nil {
		// convert kubernetes
		if kube, ok := spec.Runtime.Spec.(*v1.RuntimeKube); ok {
			infra = &v0.Infrastructure{
				Type: v0.InfraTypeKubernetesDirect,
				Spec: &v0.InfraSpec{
					Namespace: kube.Namespace,
					Conn:      kube.Connector,
				},
			}
			if infra.Spec.Namespace == "" {
				kube.Namespace = d.kubeNamespace
			}
			if infra.Spec.Conn == "" {
				kube.Connector = d.kubeConnector
			}
		}

		// convert cloud
		if _, ok := spec.Runtime.Spec.(*v1.RuntimeCloud); ok {
			runtime = &v0.Runtime{
				Type: "Cloud",
				Spec: struct{}{},
			}
		}
	}

	// if neither cloud nor kubernetes are specified
	// we default to cloud.
	if runtime == nil && infra == nil {
		runtime = &v0.Runtime{
			Type: "Cloud",
			Spec: struct{}{},
		}
	}

	// if the user explicitly provides a kubernetes connector,
	// we should override whatever was in the source yaml and
	// force kubernetes.
	if d.kubeConnector != "" {
		runtime = nil
		infra = &v0.Infrastructure{
			Type: v0.InfraTypeKubernetesDirect,
			Spec: &v0.InfraSpec{
				Namespace: d.kubeNamespace,
				Conn:      d.kubeConnector,
			},
		}
	}

	//
	// END TODO
	//

	// convert the stage to a harness stage.
	return &v0.Stage{
		ID: d.identifiers.Generate(
			slug.Create(stage.Id),
			slug.Create(stage.Name),
			slug.Create(stage.Type),
		),
		Name: convertName(stage.Name),
		Type: v0.StageTypeCI,
		Spec: v0.StageCI{
			BuildIntelligence: d.convertUseIntelligence(),
			Cache:             convertCache(spec.Cache),
			Clone:             enableClone,
			Infrastructure:    infra,
			Platform:          convertPlatform(spec.Platform, runtime),
			Runtime:           runtime,
			Execution: v0.Execution{
				Steps: steps,
			},
		},
		When:     convertStageWhen(stage.When, ""),
		Strategy: convertStrategy(stage.Strategy),
		Vars:     convertVariables(spec.Envs),
	}
}

// convertStrategy converts the v1.Strategy to the v0.Strategy
func convertStrategy(v1Strategy *v1.Strategy) *v0.Strategy {
	if v1Strategy == nil {
		return nil
	}
	v0Strategy := v0.Strategy{}
	switch v1Strategy.Type {
	case "matrix":
		v0Matrix := convertMatrix(v1Strategy.Spec.(*v1.Matrix))
		v0Strategy.Matrix = v0Matrix
	default:
	}

	return &v0Strategy
}

// convertMatrix converts the v1.Matrix to the v0.Matrix
func convertMatrix(v1Matrix *v1.Matrix) map[string]interface{} {
	matrix := make(map[string]interface{})

	// Convert axis
	for key, values := range v1Matrix.Axis {
		matrix[key] = values
	}

	// Convert exclusions
	var exclusions []v0.Exclusion
	for _, v1Exclusion := range v1Matrix.Exclude {
		exclusion := make(v0.Exclusion)
		for key, values := range v1Exclusion {
			exclusion[key] = values
		}
		exclusions = append(exclusions, exclusion)
	}
	matrix["exclude"] = exclusions

	// Convert maxConcurrency
	if v1Matrix.Concurrency != 0 {
		matrix["maxConcurrency"] = v1Matrix.Concurrency
	}

	return matrix
}

// helper function converts a drone pipeline step to a
// harness step.
//
// TODO unique identifier
// TODO failure strategy
// TODO matrix strategy
// TODO when
func (d *Downgrader) convertStep(src *v1.Step) *v0.Steps {
	if src == nil {
		return nil
	}
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
	case *v1.StepGroup:
		return &v0.Steps{StepGroup: d.convertStepGroup(src)}
	case *v1.StepBackground:
		return &v0.Steps{Step: d.convertStepBackground(src)}
	default:
		return nil // should not happen
	}
}

// helper function to convert a Group step from the v1
// structure to a list of steps. The v0 yaml does not have
// an equivalent to the group step.
func (d *Downgrader) convertStepGroup(src *v1.Step) *v0.StepGroup {
	spec_ := src.Spec.(*v1.StepGroup)

	var steps []*v0.Steps
	for _, step := range spec_.Steps {
		dst := d.convertStep(step)
		steps = append(steps, dst)
	}
	return &v0.StepGroup{
		ID:      src.Id,
		Name:    convertName(src.Name),
		Timeout: convertTimeout(src.Timeout),
		Steps:   steps,
	}
}

// helper function to convert a Parallel step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepParallel(src *v1.Step) []*v0.Steps {
	spec_ := src.Spec.(*v1.StepParallel)

	var steps []*v0.Steps
	for _, step := range spec_.Steps {
		if dst := d.convertStep(step); dst != nil {
			if dst.Step != nil {
				steps = append(steps, &v0.Steps{Step: dst.Step})
			} else {
				if dstg := d.convertStepGroup(step); dstg != nil {
					steps = append(steps, &v0.Steps{StepGroup: dstg})
				}
			}
		}
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
	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}

	// Convert outputs
	var outputs []*v0.Output
	for _, output := range spec_.Outputs {
		outputs = append(outputs, convertOutput(output))
	}
	connectorRef := d.dockerhubConn
	if spec_.Connector != "" {
		connectorRef = spec_.Connector
	}

	return &v0.Step{
		ID:      id,
		Name:    convertName(src.Name),
		Type:    d.convertStepType(spec_),
		Timeout: convertTimeout(src.Timeout),

		Spec: &v0.StepRun{
			Env:             spec_.Envs,
			Command:         spec_.Run,
			ConnRef:         connectorRef,
			Image:           convertImage(spec_.Image, d.defaultImage),
			ImagePullPolicy: convertImagePull(spec_.Pull),
			Outputs:         outputs, // Add this line
			Privileged:      spec_.Privileged,
			RunAsUser:       spec_.User,
			Reports:         convertReports(spec_.Reports),
			Shell:           strings.Title(spec_.Shell),
		},
		When:     convertStepWhen(src.When, id),
		Strategy: convertStrategy(src.Strategy),
	}
}

func (d *Downgrader) convertStepType(spec_ *v1.StepExec) string {
	stepType := v0.StepTypeRun
	if d.shouldUseTestIntelligence(spec_) {
		stepType = v0.StepTypeTest
	}
	return stepType
}

func (d *Downgrader) shouldUseTestIntelligence(spec_ *v1.StepExec) bool {
	if !d.useIntelligence {
		return false
	}
	if spec_ == nil || spec_.Run == "" {
		return false
	}

	tiTools := []string{"mvn ", "gradle ", "gradlew ", "sbt ", "bazel ", "pytest", "rspec"}
	for _, substr := range tiTools {
		if strings.Contains(spec_.Run, substr) {
			return true
		}
	}
	return false
}

// helper function to convert reports from the v1 to v0
func convertReports(reports []*v1.Report) *v0.Report {
	if reports == nil || len(reports) == 0 {
		return nil
	}

	// Initialize an empty slice to store all paths
	allPaths := []string{}

	// Loop over reports and collect all paths
	for _, report := range reports {
		allPaths = append(allPaths, report.Path...)
	}

	reportJunit := v0.ReportJunit{
		Paths: allPaths,
	}

	v0Report := v0.Report{
		// Assuming all reports have the same type
		Type: "JUnit",
		Spec: &reportJunit,
	}

	return &v0Report
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
func (d *Downgrader) convertStepBackground(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepBackground)
	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}
	// convert the entrypoint string to a slice.
	var entypoint []string
	if spec_.Entrypoint != "" {
		entypoint = []string{spec_.Entrypoint}
	}
	return &v0.Step{
		ID:   id,
		Name: convertName(src.Name),
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
			PortBindings:    convertPorts(spec_.Ports),
		},
		When: convertStepWhen(src.When, id),
	}
}

// helper function to convert a Plugin step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
// TODO convert reports
func (d *Downgrader) convertStepPlugin(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)

	if strings.Contains(spec_.Image, "plugins/kaniko:latest") {
		return d.convertStepPluginToDocker(src)
	}

	if strings.Contains(spec_.Image, "plugins/trivy:latest") {
		return d.convertStepPluginToTrivy(src)
	}

	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}

	switch spec_.Image {
	case harness.GitPluginImage:
		setting := convertSettings(spec_.With)
		return &v0.Step{
			ID:   id,
			Name: "GitClone",
			Type: v0.StepTypeGitClone,

			Spec: &v0.StepGitClone{
				Repository: setting["git_url"].(string),
				BuildType:  "<+input>",
				// TODO this directory should be populated differently for each clone
				CloneDirectory: "./",
			},
			When: convertStepWhen(src.When, id),
		}
	case "plugins/artifactory:latest":
		setting := convertSettings(spec_.With)
		return &v0.Step{
			ID:   id,
			Name: "ArtifactoryUpload",
			Type: v0.StepTypeArtifactoryUpdload,

			Spec: &v0.StepArtifactoryUpload{
				Target:     setting["target"].(string),
				SourcePath: setting["source"].(string),
				ConnRef:    "<+input>",
			},
			When: convertStepWhen(src.When, id),
		}
	default:
		return &v0.Step{
			ID:      id,
			Name:    src.Name,
			Type:    v0.StepTypePlugin,
			Timeout: convertTimeout(src.Timeout),
			Spec: &v0.StepPlugin{
				Env:             spec_.Envs,
				ConnRef:         d.dockerhubConn,
				Image:           spec_.Image,
				ImagePullPolicy: convertImagePull(spec_.Pull),
				Settings:        convertSettings(spec_.With),
				Privileged:      spec_.Privileged,
				RunAsUser:       spec_.User,
			},
			When: convertStepWhen(src.When, id),
		}
	}
}

// helper function to convert a Plugin step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
// TODO convert reports
func (d *Downgrader) convertStepPluginToDocker(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)
	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}

	stepDocker := &v0.StepDocker{
		Caching:      true,
		ConnectorRef: "<+input>",
		Privileged:   spec_.Privileged,
		RunAsUser:    spec_.User,
	}

	// Check if "repo" exists in spec_.With
	if repo, ok := spec_.With["repo"].(string); ok {
		// If "repo" exists, set it in StepDocker
		stepDocker.Repo = repo
	}
	if context, ok := spec_.With["context"].(string); ok {
		// If "context" exists, set it in StepDocker
		stepDocker.Context = context
	}
	if dockerfile, ok := spec_.With["dockerfile"].(string); ok {
		// If "dockerfile" exists, set it in StepDocker
		stepDocker.Dockerfile = dockerfile
	}
	if target, ok := spec_.With["target"].(string); ok {
		// If "target" exists, set it in StepDocker
		stepDocker.Target = target
	}
	// if cacheRepo, ok := spec_.With["cache_repo"].(string); ok {
	// 	// If "cache_repo" exists, set it in StepDocker
	// 	stepDocker.RemoteCacheRepo = cacheRepo
	// }
	if tagsInterface, ok := spec_.With["tags"]; ok {
		stepDocker.Tags = extractStringSlice(tagsInterface)
	}
	if buildArgsInterface, ok := spec_.With["build_args"]; ok {
		stepDocker.BuildsArgs = extractStringMap(buildArgsInterface)
	}
	if labelsInterface, ok := spec_.With["custom_labels"]; ok {
		stepDocker.Labels = extractStringMap(labelsInterface)
	}

	return &v0.Step{
		ID:      id,
		Name:    src.Name,
		Type:    v0.StepTypeBuildAndPushDockerRegistry,
		Timeout: convertTimeout(src.Timeout),
		Spec:    stepDocker,
		When:    convertStepWhen(src.When, id),
		Env:     spec_.Envs,
	}
}

func (d *Downgrader) convertStepPluginToTrivy(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)
	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}

	target := &v0.STOTarget{
		Type:      "container",
		Detection: "auto",
	}
	advanced := &v0.STOAdvanced{
		Log: &v0.STOAdvancedLog{
			Level: "info",
		},
	}

	image := &v0.STOImage{
		Type: "docker_v2",
	}

	if imageInterface, ok := spec_.With["image"].(string); ok {
		image.Name = imageInterface
	}

	if tagsInterface, ok := spec_.With["tag"].(string); ok {
		image.Tag = tagsInterface
	}
	stepTrivy := &v0.StepTrivy{
		Mode:       "orchestration",
		Config:     "default",
		Privileged: true,
		Target:     target,
		Advanced:   advanced,
		Image:      image,
	}

	return &v0.Step{
		ID:      id,
		Name:    src.Name,
		Type:    v0.StepTypeAquaTrivy,
		Timeout: convertTimeout(src.Timeout),
		Spec:    stepTrivy,
		When:    convertStepWhen(src.When, id),
		Env:     spec_.Envs,
	}
}

// helper function to convert an Action step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepAction(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepAction)
	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}
	return &v0.Step{
		ID:      id,
		Name:    convertName(src.Name),
		Type:    v0.StepTypeAction,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepAction{
			Uses: spec_.Uses,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
		When: convertStepWhen(src.When, id),
	}
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepBitrise(src *v1.Step) *v0.Step {
	spec_ := src.Spec.(*v1.StepBitrise)
	var id = d.identifiers.Generate(
		slug.Create(src.Id),
		slug.Create(src.Name),
		slug.Create(src.Type))
	if src.Name == "" {
		src.Name = id
	}
	return &v0.Step{
		ID:      id,
		Name:    convertName(src.Name),
		Type:    v0.StepTypeBitrise,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepBitrise{
			Uses: spec_.Uses,
			With: convertSettings(spec_.With),
			Envs: spec_.Envs,
		},
		When: convertStepWhen(src.When, id),
	}
}

func convertPorts(ports []string) map[string]string {
	if len(ports) == 0 {
		return nil
	}
	bindings := make(map[string]string, len(ports))
	for _, port := range ports {
		split := strings.Split(port, ":")
		if len(split) == 1 {
			bindings[split[0]] = split[0]
		} else if len(split) == 2 {
			bindings[split[0]] = split[1]
		}
	}
	return bindings
}

func (d *Downgrader) convertUseIntelligence() *v0.BuildIntelligence {
	if !d.useIntelligence {
		return nil
	}
	return &v0.BuildIntelligence{
		Enabled: true,
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
	if src == nil {
		return nil
	}
	var vars []*v0.Variable
	var keys []string
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := src[k]
		vars = append(vars, &v0.Variable{
			Name:  k,
			Value: v,
			Type:  "String",
		})
	}
	return vars
}

func convertSettings(src map[string]interface{}) map[string]interface{} {
	dst := map[string]interface{}{}
	for k, v := range src {
		switch v := v.(type) {
		case []interface{}:
			var strList []string
			for _, item := range v {
				strList = append(strList, fmt.Sprint(item))
			}
			dst[k] = strList
		default:
			dst[k] = fmt.Sprint(v)
		}
	}
	return dst
}

func convertTimeout(s string) v0.Duration {
	d, _ := time.ParseDuration(s)
	return v0.Duration{
		Duration: d,
	}
}

func convertImage(s string, defaultImage string) string {
	image := s
	if image == "" {
		image = defaultImage
	}
	return image
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

func convertPlatform(platform *v1.Platform, runtime *v0.Runtime) *v0.Platform {
	if platform != nil {
		var os, arch string

		// convert the OS name
		switch platform.Os {
		case "linux":
			os = "Linux"
		case "windows":
			os = "Windows"
		case "macos", "mac", "darwin":
			os = "MacOS"
		default:
			os = "Linux"
		}

		// convert the Arch name
		switch platform.Arch {
		case "amd64":
			arch = "Amd64"
		case "arm", "arm64":
			arch = "Arm64"
		default:
			// choose the default architecture
			// based on the os.
			switch os {
			case "MacOS":
				arch = "Arm64"
			default:
				arch = "Amd64"
			}
		}

		// ensure supported infra when using harness cloud
		if runtime != nil && runtime.Type == "Cloud" {
			switch os {
			case "MacOS":
				// force amd64 for Mac when using Cloud
				arch = "Arm64"
			case "Windows":
				// force amd64 for Windows when using Cloud
				arch = "Amd64"
			}
		}

		return &v0.Platform{
			OS:   os,
			Arch: arch,
		}
	} else {
		// default to linux amd64
		return &v0.Platform{
			OS:   "Linux",
			Arch: "Amd64",
		}
	}
}

func convertStepWhen(when *v1.When, stepId string) *v0.StepWhen {
	if when == nil {
		return nil
	}

	newWhen := &v0.StepWhen{
		StageStatus: "Success", // default
	}
	var conditions []string

	for _, cond := range when.Cond {
		for k, v := range cond {
			switch k {
			case "event":
				if v.In != nil {
					var eventConditions []string
					for _, event := range v.In {
						eventMapping, ok := eventMap[event]
						if !ok {
							continue
						}
						eventConditions = append(eventConditions, fmt.Sprintf("%s %s %q", eventMapping["jexl"], eventMapping["operator"], eventMapping["event"]))
					}
					if len(eventConditions) > 0 {
						conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(eventConditions, " || ")))
					}
				}
				if v.Not != nil && v.Not.In != nil {
					var notEventConditions []string
					for _, event := range v.Not.In {
						eventMapping, ok := eventMap[event]
						if !ok {
							continue
						}
						notEventConditions = append(notEventConditions, fmt.Sprintf("%s %s %q", eventMapping["jexl"], eventMapping["inverseOperator"], eventMapping["event"]))
					}
					if len(notEventConditions) > 0 {
						conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(notEventConditions, " && ")))
					}
				}
			case "status":
				if v.Eq != "" {
					newWhen.StageStatus = v.Eq
				}
				if v.In != nil {
					var statusConditions []string
					for _, status := range v.In {
						statusConditions = append(statusConditions, fmt.Sprintf("<+execution.steps.%s.status> == %q", stepId, status))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(statusConditions, " || ")))
				}
			case "branch":
				if v.In != nil {
					var branchConditions []string
					for _, branch := range v.In {
						branchConditions = append(branchConditions, fmt.Sprintf("<+trigger.targetBranch> == %q", branch))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(branchConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notBranchConditions []string
					for _, branch := range v.Not.In {
						notBranchConditions = append(notBranchConditions, fmt.Sprintf("<+trigger.targetBranch> != %q", branch))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(notBranchConditions, " && ")))
				}
			case "repo":
				if v.In != nil {
					var repoConditions []string
					for _, repo := range v.In {
						repoConditions = append(repoConditions, fmt.Sprintf("<+trigger.payload.repository.name> == %q", repo))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(repoConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notRepoConditions []string
					for _, repo := range v.Not.In {
						notRepoConditions = append(notRepoConditions, fmt.Sprintf("<+trigger.payload.repository.name> != %q", repo))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(notRepoConditions, " && ")))
				}
			case "ref":
				if v.In != nil {
					var refConditions []string
					for _, ref := range v.In {
						refConditions = append(refConditions, fmt.Sprintf("<+trigger.payload.ref> == %q", ref))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(refConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notRefConditions []string
					for _, ref := range v.Not.In {
						notRefConditions = append(notRefConditions, fmt.Sprintf("<+trigger.payload.ref> != %q", ref))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(notRefConditions, " && ")))
				}
			}
		}
	}

	if len(conditions) > 0 {
		newWhen.Condition = strings.Join(conditions, " && ")
	}

	return newWhen
}

func convertOutput(output string) *v0.Output {
	return &v0.Output{
		Name: output,
	}
}

func convertStageWhen(when *v1.When, stepId string) *v0.StageWhen {
	if when == nil {
		return nil
	}

	newWhen := &v0.StageWhen{
		PipelineStatus: "Success", // default
	}
	var conditions []string

	for _, cond := range when.Cond {
		for k, v := range cond {
			switch k {
			case "event":
				if v.In != nil {
					var eventConditions []string
					for _, event := range v.In {
						eventMapping, ok := eventMap[event]
						if !ok {
							continue
						}
						eventConditions = append(eventConditions, fmt.Sprintf("%s %s %q", eventMapping["jexl"], eventMapping["operator"], eventMapping["event"]))
					}
					if len(eventConditions) > 0 {
						conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(eventConditions, " || ")))
					}
				}
				if v.Not != nil && v.Not.In != nil {
					var notEventConditions []string
					for _, event := range v.Not.In {
						eventMapping, ok := eventMap[event]
						if !ok {
							continue
						}
						notEventConditions = append(notEventConditions, fmt.Sprintf("%s %s %q", eventMapping["jexl"], eventMapping["inverseOperator"], eventMapping["event"]))
					}
					if len(notEventConditions) > 0 {
						conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(notEventConditions, " && ")))
					}
				}
			case "status":
				if v.Eq != "" {
					newWhen.PipelineStatus = v.Eq
				}
				if v.In != nil {
					var statusConditions []string
					for _, status := range v.In {
						statusConditions = append(statusConditions, fmt.Sprintf("<+execution.steps.%s.status> == %q", stepId, status))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(statusConditions, " || ")))
				}
			case "branch":
				if v.In != nil {
					var branchConditions []string
					for _, branch := range v.In {
						branchConditions = append(branchConditions, fmt.Sprintf("<+trigger.targetBranch> == %q", branch))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(branchConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notBranchConditions []string
					for _, branch := range v.Not.In {
						notBranchConditions = append(notBranchConditions, fmt.Sprintf("<+trigger.targetBranch> != %q", branch))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(notBranchConditions, " && ")))
				}
			case "repo":
				if v.In != nil {
					var repoConditions []string
					for _, repo := range v.In {
						repoConditions = append(repoConditions, fmt.Sprintf("<+trigger.payload.repository.name> == %q", repo))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(repoConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notRepoConditions []string
					for _, repo := range v.Not.In {
						notRepoConditions = append(notRepoConditions, fmt.Sprintf("<+trigger.payload.repository.name> != %q", repo))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(notRepoConditions, " && ")))
				}
			case "ref":
				if v.In != nil {
					var refConditions []string
					for _, ref := range v.In {
						refConditions = append(refConditions, fmt.Sprintf("<+trigger.payload.ref> == %q", ref))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(refConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notRefConditions []string
					for _, ref := range v.Not.In {
						notRefConditions = append(notRefConditions, fmt.Sprintf("<+trigger.payload.ref> != %q", ref))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(notRefConditions, " && ")))
				}
			}
		}
	}

	if len(conditions) > 0 {
		newWhen.Condition = strings.Join(conditions, " && ")
	}

	return newWhen
}

func extractStringSlice(input interface{}) []string {
	switch data := input.(type) {
	case []interface{}:
		var result []string
		for _, item := range data {
			result = append(result, strings.SplitN(fmt.Sprintf("%v", item), ",", -1)...)
		}
		return result
	case interface{}:
		return []string{fmt.Sprintf("%v", data)}
	default:
		return nil
	}
}

func extractStringMap(input interface{}) map[string]string {
	switch data := input.(type) {
	case []interface{}:
		result := make(map[string]string)
		for _, item := range data {
			for _, keyValue := range strings.SplitN(fmt.Sprintf("%v", item), ",", -1) {
				pair := strings.SplitN(keyValue, "=", 2)
				if len(pair) == 2 {
					result[pair[0]] = pair[1]
				}
			}
		}
		return result
	case interface{}:
		result := make(map[string]string)
		for _, keyValue := range strings.SplitN(fmt.Sprintf("%v", data), ",", -1) {
			pair := strings.SplitN(keyValue, "=", 2)
			if len(pair) == 2 {
				result[pair[0]] = pair[1]
			}
		}
		return result
	default:
		return nil
	}
}

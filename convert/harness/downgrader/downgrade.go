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
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/drone/go-convert/internal/slug"
	"github.com/drone/go-convert/internal/store"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

// Downgrader downgrades pipelines from the v0 harness
// configuration format to the v1 configuration format.
type Downgrader struct {
	codebaseName  string
	codebaseConn  string
	dockerhubConn string
	kubeConnector string
	kubeNamespace string
	kubeEnabled   bool
	pipelineId    string
	pipelineName  string
	pipelineOrg   string
	pipelineProj  string
	identifiers   *store.Identifiers
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
		d.pipelineName = "default"
	}

	// set the default pipeline id.
	if d.pipelineId == "" {
		d.pipelineId = slug.Create(d.pipelineName)
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
	src, err := v1.ParseBytes(b)
	if err != nil {
		return nil, err
	}
	return d.downgrade(src)
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

// downgrade downgrades a v1 pipeline.
func (d *Downgrader) downgrade(src *v1.Pipeline) ([]byte, error) {
	config := new(v0.Config)
	config.Pipeline.ID = d.pipelineId
	config.Pipeline.Name = d.pipelineName
	config.Pipeline.Org = d.pipelineOrg
	config.Pipeline.Project = d.pipelineProj
	config.Pipeline.Props.CI.Codebase = v0.Codebase{
		Name:  d.codebaseName,
		Conn:  d.codebaseConn,
		Build: "<+input>",
	}
	if src.Options != nil {
		config.Pipeline.Variables = convertVariables(src.Options.Envs)
	}

	// convert stages
	for _, stage := range src.Stages {
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

	return yaml.Marshal(config)
}

// helper function converts a drone pipeline stage to a
// harness stage.
//
// TODO env variables to vars (stage-level)
// TODO delegate selectors
// TODO tags
// TODO when
// TODO failure strategy
// TODO matrix strategy
// TODO volumes
// TODO if no stage clone, use global clone, if exists
func (d *Downgrader) convertStage(stage *v1.Stage) *v0.Stage {
	// extract the spec from the v1 stage
	spec := stage.Spec.(*v1.StageCI)

	var steps []*v0.Steps

	stageEnv := spec.Envs // Convert stage environment variables

	// convert each drone step to a harness step.
	for _, v := range spec.Steps {
		// the v0 yaml does not have the concept of
		// a group step, so we append all steps in
		// the group directly to the stage to emulate
		// this behavior.
		if _, ok := v.Spec.(*v1.StepGroup); ok {
			steps = append(steps, d.convertStepGroup(v, stageEnv, 10)...)

		} else {
			// else convert the step and append to
			// the stage.
			steps = append(steps, d.convertStep(v, stageEnv))
		}
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
			Cache:          convertCache(spec.Cache),
			Clone:          enableClone,
			Infrastructure: infra,
			Platform:       convertPlatform(spec.Platform, runtime),
			Runtime:        runtime,
			Execution: v0.Execution{
				Steps: steps,
			},
		},
		When: convertStageWhen(stage.When, ""),
	}
}

// helper function converts a drone pipeline step to a
// harness step.
//
// TODO unique identifier
// TODO failure strategy
// TODO matrix strategy
// TODO when
func (d *Downgrader) convertStep(src *v1.Step, stageEnv map[string]string) *v0.Steps {
	switch src.Spec.(type) {
	case *v1.StepExec:
		return &v0.Steps{Step: d.convertStepRun(src, stageEnv)}
	case *v1.StepPlugin:
		return &v0.Steps{Step: d.convertStepPlugin(src, stageEnv)}
	case *v1.StepAction:
		return &v0.Steps{Step: d.convertStepAction(src, stageEnv)}
	case *v1.StepBitrise:
		return &v0.Steps{Step: d.convertStepBitrise(src, stageEnv)}
	case *v1.StepParallel:
		return &v0.Steps{Parallel: d.convertStepParallel(src, stageEnv)}
	case *v1.StepBackground:
		return &v0.Steps{Step: d.convertStepBackground(src, stageEnv)}
	default:
		return nil // should not happen
	}
}

// helper function to convert a Group step from the v1
// structure to a list of steps. The v0 yaml does not have
// an equivalent to the group step.
func (d *Downgrader) convertStepGroup(src *v1.Step, stageEnv map[string]string, depth int) []*v0.Steps {
	if depth > MaxDepth {
		return nil // Reached maximum depth. Stop recursion to prevent stack overflow
	}
	spec_ := src.Spec.(*v1.StepGroup)

	var steps []*v0.Steps
	for _, step := range spec_.Steps {
		// If this step is a step group, recursively convert it
		if _, ok := step.Spec.(*v1.StepGroup); ok {
			steps = append(steps, d.convertStepGroup(step, stageEnv, depth+1)...)
		} else {
			// Else, convert the step
			dst := d.convertStep(step, stageEnv)
			steps = append(steps, &v0.Steps{Step: dst.Step})
		}
	}
	return steps
}

// helper function to convert a Parallel step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepParallel(src *v1.Step, stageEnv map[string]string) []*v0.Step {
	spec_ := src.Spec.(*v1.StepParallel)

	var steps []*v0.Step
	for _, step := range spec_.Steps {
		dst := d.convertStep(step, stageEnv)
		steps = append(steps, dst.Step)
	}
	return steps
}

// helper function to convert a Run step from the v1
// structure to the v0 harness structure.
//
// TODO convert outputs
// TODO convert resources
// TODO convert reports
func (d *Downgrader) convertStepRun(src *v1.Step, stageEnv map[string]string) *v0.Step {
	spec_ := src.Spec.(*v1.StepExec)
	var id = d.identifiers.Generate(
		slug.Create(src.Name))
	return &v0.Step{
		ID:      id,
		Name:    convertName(src.Name),
		Type:    v0.StepTypeRun,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepRun{
			Env:             spec_.Envs,
			Command:         spec_.Run,
			ConnRef:         d.dockerhubConn,
			Image:           spec_.Image,
			ImagePullPolicy: convertImagePull(spec_.Pull),
			Privileged:      spec_.Privileged,
			RunAsUser:       spec_.User,
		},
		When: convertStepWhen(src.When, id),
		Env:  stageEnv,
	}
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
func (d *Downgrader) convertStepBackground(src *v1.Step, stageEnv map[string]string) *v0.Step {
	spec_ := src.Spec.(*v1.StepBackground)

	// convert the entrypoint string to a slice.
	var entypoint []string
	if spec_.Entrypoint != "" {
		entypoint = []string{spec_.Entrypoint}
	}
	var id = d.identifiers.Generate(
		slug.Create(src.Name))
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
		Env:  stageEnv,
	}
}

// helper function to convert a Plugin step from the v1
// structure to the v0 harness structure.
//
// TODO convert resources
// TODO convert reports
func (d *Downgrader) convertStepPlugin(src *v1.Step, stageEnv map[string]string) *v0.Step {
	spec_ := src.Spec.(*v1.StepPlugin)
	var id = d.identifiers.Generate(
		slug.Create(src.Name))
	return &v0.Step{
		ID:      id,
		Name:    convertName(src.Name),
		Type:    v0.StepTypePlugin,
		Timeout: convertTimeout(src.Timeout),
		Spec: &v0.StepPlugin{
			ConnRef:         d.dockerhubConn,
			Image:           spec_.Image,
			ImagePullPolicy: convertImagePull(spec_.Pull),
			Settings:        convertSettings(spec_.With),
			Privileged:      spec_.Privileged,
			RunAsUser:       spec_.User,
		},
		When: convertStepWhen(src.When, id),
		Env:  stageEnv,
	}
}

// helper function to convert an Action step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepAction(src *v1.Step, stageEnv map[string]string) *v0.Step {
	spec_ := src.Spec.(*v1.StepAction)
	var id = d.identifiers.Generate(
		slug.Create(src.Name))
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
		Env:  stageEnv,
	}
}

// helper function to convert a Bitrise step from the v1
// structure to the v0 harness structure.
func (d *Downgrader) convertStepBitrise(src *v1.Step, stageEnv map[string]string) *v0.Step {
	spec_ := src.Spec.(*v1.StepBitrise)
	var id = d.identifiers.Generate(
		slug.Create(src.Name))
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
		Env:  stageEnv,
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
	for k, v := range src {
		vars = append(vars, &v0.Variable{
			Name:  k,
			Value: v,
			Type:  "String",
		})
	}
	return vars
}

func convertSettings(src map[string]interface{}) map[string]string {
	dst := map[string]string{}
	for k, v := range src {
		dst[k] = fmt.Sprint(v)
	}
	return dst
}

func convertTimeout(s string) v0.Duration {
	d, _ := time.ParseDuration(s)
	return v0.Duration{
		Duration: d,
	}
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
		switch platform.Os.String() {
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
		switch platform.Arch.String() {
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
						branchConditions = append(branchConditions, fmt.Sprintf("<+trigger.branch> == %q", branch))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(branchConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notBranchConditions []string
					for _, branch := range v.Not.In {
						notBranchConditions = append(notBranchConditions, fmt.Sprintf("<+trigger.branch> != %q", branch))
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
						branchConditions = append(branchConditions, fmt.Sprintf("<+trigger.branch> == %q", branch))
					}
					conditions = append(conditions, fmt.Sprintf("%s", strings.Join(branchConditions, " || ")))
				}
				if v.Not != nil && v.Not.In != nil {
					var notBranchConditions []string
					for _, branch := range v.Not.In {
						notBranchConditions = append(notBranchConditions, fmt.Sprintf("<+trigger.branch> != %q", branch))
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

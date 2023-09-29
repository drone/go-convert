package rio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	harness "github.com/drone/go-convert/convert/harness/yaml"
	rio "github.com/drone/go-convert/convert/rio/yaml"
	"github.com/drone/go-convert/internal/store"
	"gopkg.in/yaml.v3"
)

// as we walk the yaml, we store a
// snapshot of the current node and
// its parents.
type context struct {
	config *rio.Config
}

// Converter converts a Rio pipeline to a Harness
// v0 pipeline.
type Converter struct {
	kubeEnabled     bool
	kubeNamespace   string
	kubeConnector   string
	kubeOs          string
	dockerhubConn   string
	githubConnector string
	identifiers     *store.Identifiers

	pipelineId   string
	pipelineName string
	pipelineOrg  string
	pipelineProj string

	notifyUserGroup string
}

// New creates a new Converter that converts a Rio
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

	// set default kubernetes OS
	if d.kubeOs == "" {
		d.kubeOs = string(harness.InfraOsLinux)
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	// set default docker connector
	if d.dockerhubConn == "" {
		d.dockerhubConn = harness.DefaultDockerConnector
	}

	if d.notifyUserGroup == "" {
		d.notifyUserGroup = "account._account_all_users"
	}
	return d
}

// Convert converts a rio pipeline to v0.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	config, err := rio.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(&context{
		config: config,
	})
}

// ConvertBytes converts a rio pipeline to v0.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.Convert(
		bytes.NewBuffer(b),
	)
}

// ConvertString converts a rio pipeline to v0.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.Convert(
		bytes.NewBufferString(s),
	)
}

// ConvertFile converts a rio pipeline to v0.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

func (d *Converter) convertRunStep(p rio.Pipeline) harness.StepRun {
	step := new(harness.StepRun)
	step.Env = p.Machine.Env
	command := ""
	for _, s := range p.Build.Steps {
		command += fmt.Sprintf("%s\n", s)
	}
	step.Command = command
	step.Shell = harness.ShellPosix
	step.ConnRef = d.dockerhubConn
	step.Image = p.Machine.BaseImage
	return *step
}

func (d *Converter) convertDockerSteps(dockerfile rio.Dockerfile) []harness.StepDocker {
	steps := make([]harness.StepDocker, 0)
	if dockerfile.Version == "" {
		dockerfile.Version = "latest"
	}
	tags := dockerfile.ExtraTags
	tags = append(tags, dockerfile.Version)
	for _, r := range dockerfile.Publish {
		step := new(harness.StepDocker)
		step.Context = dockerfile.Context
		step.Dockerfile = dockerfile.DockerfilePath
		step.Repo = r.Repo
		step.Tags = tags
		step.BuildsArgs = dockerfile.Env
		step.ConnectorRef = d.dockerhubConn
		steps = append(steps, *step)
	}
	return steps
}

func (d *Converter) convertExecution(p rio.Pipeline) harness.Execution {
	execution := harness.Execution{}

	executionSteps := make([]*harness.Steps, 0)
	// Append all commands in build attribute to one run step
	if len(p.Build.Steps) != 0 {
		steps := harness.Steps{
			Step: &harness.Step{
				Name: "Build",
				ID:   "Build",
				Spec: d.convertRunStep(p),
				Type: harness.StepTypeRun,
			},
		}
		executionSteps = append(executionSteps, &steps)
	}

	dockerStepCounter := 1
	if len(p.Pkg.Dockerfile) != 0 {
		if !p.Pkg.Release {
			fmt.Println("# [WARN]: release=false is not supported")
		}

		// Each entry in package.Dockerfile attribute is one build and push step
		for _, dockerfile := range p.Pkg.Dockerfile {
			// each dockerfile entry can have multiple repos
			dockerSteps := d.convertDockerSteps(dockerfile)

			// Append all docker steps
			for _, dStep := range dockerSteps {
				steps := harness.Steps{
					Step: &harness.Step{
						Name: fmt.Sprintf("BuildAndPush_%d", dockerStepCounter),
						ID:   fmt.Sprintf("BuildAndPush_%d", dockerStepCounter),
						Spec: &dStep,
						Type: harness.StepTypeBuildAndPushDockerRegistry,
					},
				}
				executionSteps = append(executionSteps, &steps)
				dockerStepCounter++
			}
		}
	}
	execution.Steps = executionSteps
	return execution
}

func (d *Converter) convertCIStage(p rio.Pipeline, platform string) harness.StageCI {
	stage := harness.StageCI{
		Execution: d.convertExecution(p),
	}
	if d.kubeEnabled {
		infra := harness.Infrastructure{
			Type: harness.InfraTypeKubernetesDirect,
			Spec: &harness.InfraSpec{
				Namespace:             d.kubeNamespace,
				Conn:                  d.kubeConnector,
				AutomountServiceToken: true,
				Os:                    d.kubeOs,
			},
		}
		if platform == "linux/amd64" {
			infra.Spec.NodeSelector = map[string]string{"kubernetes.io/arch": "amd64"}
			infra.Spec.Os = string(harness.InfraOsLinux)
		} else if platform == "linux/arm64" {
			infra.Spec.NodeSelector = map[string]string{"kubernetes.io/arch": "arm64"}
			infra.Spec.Os = string(harness.InfraOsLinux)
		}
		stage.Infrastructure = &infra
	}
	if d.githubConnector != "" {
		stage.Clone = true
	}
	return stage
}

func (d *Converter) getStages(p rio.Pipeline) *harness.Stages {
	if len(p.Machine.TargetPlatforms) == 0 {
		// assume linux/amd64 - one stage
		return &harness.Stages{
			Stage: &harness.Stage{
				Name: p.Name,
				ID:   d.convertNameToID(p.Name),
				Spec: d.convertCIStage(p, ""),
				Type: harness.StageTypeCI,
			},
		}
	}
	// parallel stages with different target platforms
	stages := make([]*harness.Stages, 0)
	for _, platform := range p.Machine.TargetPlatforms {
		stageName := fmt.Sprintf("%s_%s", p.Name, d.convertNameToID(platform))
		stage := &harness.Stages{
			Stage: &harness.Stage{
				Name: stageName,
				ID:   d.convertNameToID(stageName),
				Spec: d.convertCIStage(p, platform),
				Type: harness.StageTypeCI,
			},
		}
		stages = append(stages, stage)
	}
	return &harness.Stages{
		Parallel: stages,
	}
}

func (d *Converter) convertNameToID(name string) (ID string) {
	ID = strings.ReplaceAll(name, " ", "_")
	ID = strings.ReplaceAll(ID, "-", "_")
	ID = strings.ReplaceAll(ID, "/", "_")
	return ID
}

func (d *Converter) getNotifications() []harness.NotificationRules {
	notificationRules := []harness.NotificationRules{
		{
			Name: "user_group_notification",
			Id:   "user_group_notification",
			PipelineEvents: []harness.NotificationPipelineEvent{
				{Type: "PipelineSuccess"},
				{Type: "PipelineFailed"},
			},
			NotificationMethod: harness.NotificationMethod{
				Type: "Email",
				Spec: harness.NotificationSpec{
					UserGroups: []string{d.notifyUserGroup},
				},
			},
			Enabled: true,
		},
	}
	return notificationRules
}

// converts converts a Rio pipeline to a Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {
	// create the harness pipeline spec
	pipeline := &harness.Pipeline{}
	for _, p := range ctx.config.Pipelines {
		pipeline.Stages = append(pipeline.Stages, d.getStages(p))
	}
	pipeline.Name = d.pipelineName
	pipeline.ID = d.pipelineId
	pipeline.Org = d.pipelineOrg
	pipeline.Project = d.pipelineProj
	if ctx.config.Timeout != 0 {
		pipeline.Timeout = fmt.Sprintf("%dm", ctx.config.Timeout)
	}
	if ctx.config.Notify.Email.Enabled {
		pipeline.NotificationRules = d.getNotifications()
	}
	if d.githubConnector != "" {
		pipeline.Props.CI.Codebase = harness.Codebase{
			Conn:  d.githubConnector,
			Build: "<+input>",
		}
	}

	// create the harness pipeline resource
	config := &harness.Config{Pipeline: *pipeline}

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

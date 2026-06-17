package pipelineconverter

import (
	convertexpressions "github.com/drone/go-convert/convert/convertexpressions"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

type PipelineConverter struct {
	stageCtx      *convert_helpers.StageConversionContext
	stepInfoByFQN map[string]*convertexpressions.StepInfoFQN // v1 FQN -> step info (sole step registry)
}

func NewPipelineConverter() *PipelineConverter {
	return &PipelineConverter{
		stageCtx:      convert_helpers.NewStageConversionContext(),
		stepInfoByFQN: make(map[string]*convertexpressions.StepInfoFQN),
	}
}

// GetStepInfoByFQN returns the accumulated FQN-keyed step (and step group)
// info map, keyed by each step's full v1 FQN.
func (c *PipelineConverter) GetStepInfoByFQN() map[string]*convertexpressions.StepInfoFQN {
	return c.stepInfoByFQN
}

// ConvertPipeline converts a v0 Pipeline to a v1 Pipeline.
func (c *PipelineConverter) ConvertPipeline(src *v0.Pipeline) *v1.Pipeline {
	if src == nil {
		return nil
	}

	dst := &v1.Pipeline{
		Id:   src.ID,
		Name: src.Name,
	}

	// Check for Template - if template exists, set it and continue converting rest
	if src.Template != nil {
		dst.Template = c.convertPipelineTemplate(src)
	}

	var barriers []string
	if src.FlowControl != nil {
		barriers = convertBarriers(src.FlowControl.Barriers)
	}

	inputs := c.convertVariables(src.Variables)
	stages := c.convertStages(src.Stages, "pipeline")

	clone := c.convertCodebase(src.Props.CI.Codebase)
	dst.Inputs = inputs
	dst.Stages = stages
	dst.Barriers = barriers
	dst.Clone = clone
	dst.Notifications = convert_helpers.ConvertNotifications(src.NotificationRules)
	dst.Delegate = convert_helpers.ConvertDelegate(src.DelegateSelectors, nil)
	dst.Timeout = src.Timeout
	dst.AllowStageExecutions = src.AllowStageExecutions
	dst.Tags = src.Tags
	dst.Desc = src.Desc
	dst.FixedInputsOnRerun = src.FixedInputsOnRerun

	return dst
}

func (c *PipelineConverter) convertCodebase(src *v0.Codebase) *v1.Clone {
	if src == nil {
		return &v1.Clone{
			Enabled: false,
		}
	}

	clone := &v1.Clone{
		Enabled:   true,
		Repo:      src.Name,
		Connector: src.Conn,
	}

	// Handle Build field - can be either a string expression or a Build struct
	if !src.Build.IsNil() {
		if build, ok := src.Build.AsStruct(); ok {
			// Build is a struct with Type and Spec
			cloneRef := &v1.CloneRef{}

			// Extract name from Spec based on type
			if build.Type == "branch" && build.Spec.Branch != "" {
				cloneRef.Name = build.Spec.Branch
				cloneRef.Type = "branch"
			} else if build.Type == "tag" && build.Spec.Tag != "" {
				cloneRef.Name = build.Spec.Tag
				cloneRef.Type = "tag"
			} else if build.Type == "PR" && build.Spec.Number != nil {
				cloneRef.Number = build.Spec.Number
				cloneRef.Type = "pull-request"
			} else if build.Type == "commitSha" && build.Spec.CommitSha != "" {
				cloneRef.Sha = build.Spec.CommitSha
				cloneRef.Type = "commit"
			}

			clone.Ref = cloneRef
		}
	}
	clone.Depth = src.Depth
	clone.Lfs = src.Lfs

	clone.Tags = src.FetchTags
	clone.Trace = src.Debug
	clone.CloneDir = src.CloneDirectory

	clone.Submodules = src.SubmoduleStrategy
	clone.PersistCredentials = src.PersistCredentials
	clone.SparseCheckout = src.SparseCheckout
	clone.PreFetchCommand = src.PreFetchCommand
	clone.User = src.RunAsUser

	switch src.PrCloneStrategy {
	case "MergeCommit":
		clone.Strategy = "merge"
	case "SourceBranch":
		clone.Strategy = "source-branch"
	}

	clone.Insecure = src.SslVerify

	clone.Resources = convert_helpers.ConvertContainerResources(src.Resources)

	return clone
}

// convertVariables converts a list of v0 Variables to v1 Inputs.
// Delegates to convert_helpers.ConvertVariables.
func (c *PipelineConverter) convertVariables(src []*v0.Variable) map[string]*v1.Input {
	return convert_helpers.ConvertVariables(src)
}

// convertBarriers converts a list of v0 Barriers to v1 Barriers.
func convertBarriers(src []*v0.Barrier) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, 0, len(src))
	for _, barrier := range src {
		if barrier == nil {
			continue
		}
		dst = append(dst, barrier.Name)
	}
	return dst
}

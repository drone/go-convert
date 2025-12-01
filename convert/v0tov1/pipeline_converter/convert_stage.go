package pipelineconverter

import (
	"log"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// convertStages converts a list of v0 Stages to v1 Stages.
func (c *PipelineConverter) convertStages(src []*v0.Stages) []*v1.Stage {
	dst := make([]*v1.Stage, 0)

	for _, stages := range src {
		if stages.Stage != nil {
			dst = append(dst, c.convertStage(stages.Stage))
			continue
		}
		if stages.Parallel != nil {
			stage := v1.Stage{
				Parallel: &v1.StageGroup{
					Stages: c.convertStages(stages.Parallel),
				},
			}
			dst = append(dst, &stage)
		}
	}
	return dst
}

// convertStage converts a v0 Stage to a v1 Stage.
func (c *PipelineConverter) convertStage(src *v0.Stage) *v1.Stage {
	if src == nil {
		return nil
	}

	var steps []*v1.Step
	var rollback []*v1.Step
	var service *v1.ServiceRef
	var environment *v1.EnvironmentRef

	switch spec := src.Spec.(type) {

	case *v0.StageApproval:
		steps = c.ConvertSteps(spec.Execution.Steps)

	case *v0.StageCustom:
		steps = c.ConvertSteps(spec.Execution.Steps)
		environment = convert_helpers.ConvertEnvironment(spec.Environment)

	case *v0.StageCI:
		steps = c.ConvertSteps(spec.Execution.Steps)
		clone := convert_helpers.ConvertCloneCodebase(spec.Clone)
		cache := convert_helpers.ConvertCaching(spec.Cache)
		buildIntelligence := convert_helpers.ConvertBuildIntelligence(spec.BuildIntelligence)
		runtime := convert_helpers.ConvertInfrastructure(spec.Infrastructure)
		platform := convert_helpers.ConvertPlatform(spec.Platform)

		// Create stage with CI-specific fields
		stage := &v1.Stage{
			Id:                src.ID,
			Name:              src.Name,
			Steps:             steps,
			Clone:             clone,
			Cache:             cache,
			BuildIntelligence: buildIntelligence,
			Runtime:           runtime,
			Platform:          platform,
			OnFailure:         convert_helpers.ConvertFailureStrategies(src.FailureStrategies),
			Inputs:            c.convertVariables(src.Vars),
			Delegate:          convert_helpers.ConvertDelegate(src.DelegateSelectors),
			Strategy:          convert_helpers.ConvertStrategy(src.Strategy),
		}
		return stage

	case *v0.StageDeployment:
		// Convert deployment steps
		if spec.Execution != nil {
			steps = c.ConvertSteps(spec.Execution.Steps)
			// Convert rollback steps with different path
			if spec.Execution.RollbackSteps != nil {
				rollback = c.ConvertSteps(spec.Execution.RollbackSteps)
			}
		}
		deprecatedInfraDefinition := false
		// Convert environment configuration - check all possible sources
		if spec.Environment != nil {
			environment = convert_helpers.ConvertEnvironment(spec.Environment)
		} else if spec.Environments != nil {
			environment = convert_helpers.ConvertEnvironments(spec.Environments)
		} else if spec.EnvironmentGroup != nil {
			environment = convert_helpers.ConvertEnvironmentGroup(spec.EnvironmentGroup)
		} else if spec.Infrastructure != nil {
			// Convert infrastructure to environment configuration
			log.Printf("Warning!!! Deprecated infrastructure definition found in Deployment stage: %s, infrastructure and service definition will be skipped\n", src.ID)
			deprecatedInfraDefinition = true
		}

		// Convert service configuration
		if spec.Service != nil {
			service = convert_helpers.ConvertDeploymentService(spec.Service)
		} else if spec.Services != nil {
			service = convert_helpers.ConvertDeploymentServices(spec.Services)
		} else if spec.ServiceConfig != nil && !deprecatedInfraDefinition {
			service = convert_helpers.ConvertDeploymentServiceConfig(spec.ServiceConfig)
		}

	default:
		log.Printf("Warning!!! stage type: %s (stage: %s) is not yet supported!\n", src.Type, src.ID)
	}

	onFailure := convert_helpers.ConvertFailureStrategies(src.FailureStrategies)
	strategy := convert_helpers.ConvertStrategy(src.Strategy)
	inputs := c.convertVariables(src.Vars)
	return &v1.Stage{
		Id:          src.ID,
		Name:        src.Name,
		Steps:       steps,
		Rollback:    rollback,
		Service:     service,
		Environment: environment,
		OnFailure:   onFailure,
		Inputs:      inputs,
		Delegate:    convert_helpers.ConvertDelegate(src.DelegateSelectors),
		Strategy:    strategy,
	}
}

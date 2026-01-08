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

func (c *PipelineConverter) convertStage(src *v0.Stage) *v1.Stage {
	if src == nil {
		return nil
	}

	stage := &v1.Stage{
		Id:        src.ID,
		Name:      src.Name,
		OnFailure: convert_helpers.ConvertFailureStrategies(src.FailureStrategies),
		Inputs:    c.convertVariables(src.Vars),
		Delegate:  convert_helpers.ConvertDelegate(src.DelegateSelectors),
		Strategy:  convert_helpers.ConvertStrategy(src.Strategy),
		If:        convert_helpers.ConvertStageWhen(src.When),
	}
	stage.Env = convertStageInputsToEnv(stage.Inputs)

	switch spec := src.Spec.(type) {
	case *v0.StageApproval:
		stage.Steps = c.ConvertSteps(spec.Execution.Steps)
		stage.Timeout = spec.Timeout

	case *v0.StageCustom:
		stage.Steps = c.ConvertSteps(spec.Execution.Steps)
		stage.Environment = convert_helpers.ConvertEnvironment(spec.Environment)
		stage.Timeout = spec.Timeout

	case *v0.StageCI:
		stage.Steps = c.ConvertSteps(spec.Execution.Steps)

		// Convert service dependencies to background steps and prepend them
		if len(spec.Services) > 0 {
			backgroundSteps := convert_helpers.ConvertServiceDependenciesToBackgroundSteps(spec.Services)
			if len(backgroundSteps) > 0 {
				stage.Steps = append(backgroundSteps, stage.Steps...)
			}
		}

		stage.Clone = convert_helpers.ConvertCloneCodebase(spec.Clone)
		stage.Cache = convert_helpers.ConvertCaching(spec.Cache)
		stage.BuildIntelligence = convert_helpers.ConvertBuildIntelligence(spec.BuildIntelligence)
		stage.Timeout = spec.Timeout

		// Handle volumes
		volumes := convert_helpers.ConvertSharedPaths(spec.SharedPaths)
		if spec.Infrastructure != nil {
			stage.Runtime = convert_helpers.ConvertInfrastructureToRuntime(spec.Infrastructure)
			volumes = append(volumes, convert_helpers.ConvertInfrastructureToVolumes(spec.Infrastructure)...)
		} else if spec.Runtime != nil {
			stage.Runtime = convert_helpers.ConvertRuntime(spec.Runtime)
		} else {
			log.Printf("Warning!!! No runtime or infrastructure found in CI stage: %s\n", src.ID)
		}
		stage.Volumes = volumes
		stage.Platform = convert_helpers.ConvertPlatform(spec.Platform)

	case *v0.StageDeployment:
		// Convert deployment steps
		if spec.Execution != nil {
			stage.Steps = c.ConvertSteps(spec.Execution.Steps)
			if spec.Execution.RollbackSteps != nil {
				stage.Rollback = c.ConvertSteps(spec.Execution.RollbackSteps)
			}
		}

		// Convert environment configuration
		deprecatedInfraDefinition := false
		if spec.Environment != nil {
			stage.Environment = convert_helpers.ConvertEnvironment(spec.Environment)
		} else if spec.Environments != nil {
			stage.Environment = convert_helpers.ConvertEnvironments(spec.Environments)
		} else if spec.EnvironmentGroup != nil {
			stage.Environment = convert_helpers.ConvertEnvironmentGroup(spec.EnvironmentGroup)
		} else if spec.Infrastructure != nil {
			log.Printf("Warning!!! Deprecated infrastructure definition found in Deployment stage: %s, infrastructure and service definition will be skipped\n", src.ID)
			deprecatedInfraDefinition = true
		}

		// Convert service configuration
		if spec.Service != nil {
			stage.Service = convert_helpers.ConvertDeploymentService(spec.Service)
		} else if spec.Services != nil {
			stage.Service = convert_helpers.ConvertDeploymentServices(spec.Services)
		} else if spec.ServiceConfig != nil && !deprecatedInfraDefinition {
			stage.Service = convert_helpers.ConvertDeploymentServiceConfig(spec.ServiceConfig)
		}
		stage.Timeout = spec.Timeout

	default:
		log.Printf("Warning!!! stage type: %s (stage: %s) is not yet supported!\n", src.Type, src.ID)
	}

	return stage
}

func convertStageInputsToEnv(inputs map[string]*v1.Input) map[string]interface{} {
	env := map[string]interface{}{}
	for input_name, input := range inputs {
		if input.Value != nil {
			env[input_name] = input.Value
		} else if input.Default != nil {
			env[input_name] = input.Default
		}
	}
	return env
}

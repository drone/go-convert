package pipelineconverter

import (
	"log"
	"reflect"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// convertSteps converts a list of v0.Steps to list of v1.Step.
func (c *PipelineConverter) ConvertSteps(src []*v0.Steps) []*v1.Step {
	if len(src) == 0 {
		return nil
	}
	dst := make([]*v1.Step, 0, len(src))
	for _, s := range src {
		if s == nil {
			continue
		}
		if s.Step != nil {
			if step := c.ConvertSingleStep(s.Step); step != nil {
				dst = append(dst, step)
			}
			continue
		}
		if s.Parallel != nil {
			parallel_group := &v1.StepGroup{
				Steps: c.ConvertSteps(s.Parallel),
			}
			parallel := &v1.Step{
				Parallel: parallel_group,
			}
			dst = append(dst, parallel)
			continue
		}
		if s.StepGroup != nil {
			// TODO: handle conditionals
			group := &v1.Step{
				Name: s.StepGroup.Name,
				Id:   s.StepGroup.ID,
				Env:  s.StepGroup.Env,
				Group: &v1.StepGroup{
					Steps: c.ConvertSteps(s.StepGroup.Steps),
				},
				OnFailure: convert_helpers.ConvertFailureStrategies(s.StepGroup.FailureStrategies),
				Strategy:  convert_helpers.ConvertStrategy(s.StepGroup.Strategy),
				Timeout:   s.StepGroup.Timeout,
				Delegate:  convert_helpers.ConvertDelegate(s.StepGroup.DelegateSelectors),
			}
			dst = append(dst, group)
		}
	}
	return dst
}

func (c *PipelineConverter) ConvertSingleStep(src *v0.Step) *v1.Step {
	if src == nil {
		return nil
	}

	step := &v1.Step{
		Id:   src.ID,
		Name: src.Name,
	}

	// Convert step-specific settings
	switch src.Type {
	case v0.StepTypeAction:
		step.Action = convert_helpers.ConvertStepAction(src)
	case v0.StepTypeJiraCreate:
		step.Template = convert_helpers.ConvertStepJiraCreate(src)
	case v0.StepTypeJiraUpdate:
		step.Template = convert_helpers.ConvertStepJiraUpdate(src)
	case v0.StepTypeRun:
		step.Run = convert_helpers.ConvertStepRun(src)
	case v0.StepTypeHarnessApproval:
		step.Approval = convert_helpers.ConvertStepHarnessApproval(src)
	case v0.StepTypeK8sRollingDeploy:
		step.Template = convert_helpers.ConvertStepK8sRollingDeploy(src)
	case v0.StepTypeK8sRollingRollback:
		step.Template = convert_helpers.ConvertStepK8sRollingRollback(src)
	case v0.StepTypeK8sApply:
		step.Template = convert_helpers.ConvertStepK8sApply(src)
	case v0.StepTypeK8sBGSwapServices:
		step.Template = convert_helpers.ConvertStepK8sBGSwapServices(src)
	case v0.StepTypeK8sBlueGreenStageScaleDown:
		step.Template = convert_helpers.ConvertStepK8sBlueGreenStageScaleDown(src)
	case v0.StepTypeK8sCanaryDelete:
		step.Template = convert_helpers.ConvertStepK8sCanaryDelete(src)
	case v0.StepTypeK8sDiff:
		step.Template = convert_helpers.ConvertStepK8sDiff(src)
	case v0.StepTypeK8sDelete:
		step.Template = convert_helpers.ConvertStepK8sDelete(src)
	case v0.StepTypeK8sRollout:
		step.Template = convert_helpers.ConvertStepK8sRollout(src)
	case v0.StepTypeK8sScale:
		step.Template = convert_helpers.ConvertStepK8sScale(src)
	case v0.StepTypeK8sDryRun:
		step.Template = convert_helpers.ConvertStepK8sDryRun(src)
	case v0.StepTypeK8sTrafficRouting:
		step.Template = convert_helpers.ConvertStepK8sTrafficRouting(src)
	case v0.StepTypeK8sCanaryDeploy:
		step.Template = convert_helpers.ConvertStepK8sCanaryDeploy(src)
	case v0.StepTypeK8sBlueGreenDeploy:
		step.Template = convert_helpers.ConvertStepK8sBlueGreenDeploy(src)
	case v0.StepTypeHelmBGDeploy:
		step.Template = convert_helpers.ConvertStepHelmBGDeploy(src)
	case v0.StepTypeHelmBlueGreenSwapStep:
		step.Template = convert_helpers.ConvertStepHelmBlueGreenSwapStep(src)
	case v0.StepTypeHelmCanaryDeploy:
		step.Template = convert_helpers.ConvertStepHelmCanaryDeploy(src)
	case v0.StepTypeHelmDelete:
		step.Template = convert_helpers.ConvertStepHelmDelete(src)
	case v0.StepTypeHelmDeploy:
		step.Template = convert_helpers.ConvertStepHelmDeploy(src)
	case v0.StepTypeHelmRollback:
		step.Template = convert_helpers.ConvertStepHelmRollback(src)
	case v0.StepTypeWait:
		step.Wait = convert_helpers.ConvertStepWait(src)
	case v0.StepTypeHTTP:
		step.Run = convert_helpers.ConvertStepHTTP(src)
	case v0.StepTypeShellScript:
		step.Run = convert_helpers.ConvertStepShellScript(src)
	case v0.StepTypeBarrier:
		step.Barrier = convert_helpers.ConvertStepBarrier(src)
	case v0.StepTypeQueue:
		step.Queue = convert_helpers.ConvertStepQueue(src)
	case v0.StepTypeCustomApproval:
		step.Approval = convert_helpers.ConvertStepCustomApproval(src)
	case v0.StepTypeJiraApproval:
		step.Approval = convert_helpers.ConvertStepJiraApproval(src)
	case v0.StepTypeServiceNowApproval:
		step.Approval = convert_helpers.ConvertStepServiceNowApproval(src)
	case v0.StepTypeEmail:
		step.Template = convert_helpers.ConvertStepEmail(src)
	case v0.StepTypeArtifactoryUpload:
		step.Template = convert_helpers.ConvertStepArtifactoryUpload(src)
	case v0.StepTypeSaveCacheS3:
		step.Template = convert_helpers.ConvertStepSaveCacheS3(src)
	case v0.StepTypeSaveCacheGCS:
		step.Template = convert_helpers.ConvertStepSaveCacheGCS(src)
	case v0.StepTypeRestoreCacheGCS:
		step.Template = convert_helpers.ConvertStepRestoreCacheGCS(src)
	case v0.StepTypeRestoreCacheS3:
		step.Template = convert_helpers.ConvertStepRestoreCacheS3(src)
	case v0.StepTypeBuildAndPushECR:
		step.Template = convert_helpers.ConvertStepBuildAndPushECR(src)
	case v0.StepTypeGCSUpload:
		step.Template = convert_helpers.ConvertStepGCSUpload(src)
	case v0.StepTypeS3Upload:
		step.Template = convert_helpers.ConvertStepS3Upload(src)
	case v0.StepTypeBuildAndPushGAR:
		step.Template = convert_helpers.ConvertStepBuildAndPushGAR(src)
	case v0.StepTypeBuildAndPushDockerRegistry:
		step.Template = convert_helpers.ConvertStepBuildAndPushDockerRegistry(src)
	case v0.StepTypePlugin:
		step.Run = convert_helpers.ConvertStepPlugin(src)
	case v0.StepTypeTest:
		step.RunTest = convert_helpers.ConvertStepTestIntelligence(src)
	default:
		// Unknown step type, return nil
		log.Println("Warning!!! step type: " + src.Type + " is not yet supported!")
		step.Template = &v1.StepTemplate{
			Uses: src.Type,
		}
	}

	// Convert common step settings
	convertCommonStepSettings(src, step)
	return step
}

// convertCommonStepSettings converts common step settings like timeout, failure strategies, etc.
func convertCommonStepSettings(src *v0.Step, dst *v1.Step) {
	// Convert timeout
	if src.Timeout != "" {
		dst.Timeout = src.Timeout
	}

	// Convert failure strategies
	if src.FailureStrategies != nil {
		dst.OnFailure = convert_helpers.ConvertFailureStrategies(src.FailureStrategies)
	}

	// Convert environment variables
	if len(src.Env) > 0 {
		dst.Env = src.Env
	}

	// Convert when conditions
	if src.When != nil {
		dst.If = convertStepWhen(src.When)
	}

	// Convert strategies
	if src.Strategy != nil {
		dst.Strategy = convert_helpers.ConvertStrategy(src.Strategy)
	}

	// Convert delegate selectors

	// extract delegate selectors and includeInfraSelectors from src using reflection
	var delegate_selectors *flexible.Field[[]string]
	var include_infra_selectors bool

	if src.Spec != nil {
		// Use reflection to find embedded CommonStepSpec
		specValue := reflect.ValueOf(src.Spec)
		if specValue.Kind() == reflect.Ptr {
			specValue = specValue.Elem()
		}

		if specValue.Kind() == reflect.Struct {
			// Look for embedded CommonStepSpec fields
			specType := specValue.Type()
			for i := 0; i < specValue.NumField(); i++ {
				field := specValue.Field(i)
				fieldType := specType.Field(i)

				// Check if this field is an embedded CommonStepSpec
				if fieldType.Anonymous && fieldType.Type.Name() == "CommonStepSpec" {
					// Extract DelegateSelectors (FlexibleField[[]string])
					if delegateField := field.FieldByName("DelegateSelectors"); delegateField.IsValid() {
						// Check if it's a pointer to FlexibleField[[]string]
						if delegateField.Kind() == reflect.Ptr && !delegateField.IsNil() {
							elemType := delegateField.Type().Elem()
							if elemType.String() == "flexible.Field[[]string]" {
								delegate_selectors = delegateField.Interface().(*flexible.Field[[]string])
							}
						}

						// Extract IncludeInfraSelectors
						if infraField := field.FieldByName("IncludeInfraSelectors"); infraField.IsValid() {
							if infraField.Kind() == reflect.Bool {
								include_infra_selectors = infraField.Bool()
							}
						}
						break
					}
				}
			}
		}
	}
	// Convert delegate using the extracted values
	delegate := convert_helpers.ConvertDelegate(delegate_selectors)

	// Handle includeInfraSelectors for struct-based delegates
	if include_infra_selectors && delegate != nil {
		if str, ok := delegate.AsString(); !ok || str == "" {
			if delegateStruct, ok := delegate.AsStruct(); ok {
				delegateStruct.Inherit = true
				delegate.Set(delegateStruct)
			}
		}
	}

	dst.Delegate = delegate
}

// convertStepWhen converts v0 step when conditions to v1 format
func convertStepWhen(when *v0.StepWhen) string {
	if when == nil {
		return ""
	}

	if when.Condition != "" {
		return when.Condition
	}

	return ""
}

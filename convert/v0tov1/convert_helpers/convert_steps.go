package converthelpers

import (
	"fmt"
	"reflect"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertSteps converts a list of v0.Steps to list of v1.Step.
func ConvertSteps(src []*v0.Steps) []*v1.Step {
	if len(src) == 0 {
		return nil
	}
	dst := make([]*v1.Step, 0, len(src))
	for _, s := range src {
		if s == nil {
			continue
		}
		if s.Step != nil {
			if step := ConvertSingleStep(s.Step); step != nil {
				dst = append(dst, step)
			}
			continue
		}
		if s.Parallel != nil {
			// TODO: handle s.Parallel
			parallel_group := &v1.StepGroup{
				Steps: ConvertSteps(s.Parallel),
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
					Steps: ConvertSteps(s.StepGroup.Steps),
				},
			}
			dst = append(dst, group)
		}
	}
	return dst
}

// convertSingleStep is a factory that dispatches to the appropriate
// per-type converter based on the v0 step Type.
func ConvertSingleStep(src *v0.Step) *v1.Step {
	if src == nil {
		return nil
	}

	// Create base step with common fields
	step := &v1.Step{
		Id:   src.ID,
		Name: src.Name,
	}

	// Convert step-specific settings
	switch src.Type {
	case v0.StepTypeAction:
		step.Action = ConvertStepAction(src)
	case v0.StepTypeJiraCreate:
		step.Action = ConvertStepJiraCreate(src)
	case v0.StepTypeJiraUpdate:
		step.Action = ConvertStepJiraUpdate(src)
	case v0.StepTypeRun:
		step.Run = ConvertStepRun(src)
	case v0.StepTypeHarnessApproval:
		step.Approval = ConvertStepHarnessApproval(src)
	case v0.StepTypeK8sRollingDeploy:
		step.Template = ConvertStepK8sRollingDeploy(src)
	case v0.StepTypeK8sRollingRollback:
		step.Template = ConvertStepK8sRollingRollback(src)
	case v0.StepTypeK8sApply:
		step.Template = ConvertStepK8sApply(src)
	case v0.StepTypeK8sBGSwapServices:
		step.Template = ConvertStepK8sBGSwapServices(src)
	case v0.StepTypeK8sBlueGreenStageScaleDown:
		step.Template = ConvertStepK8sBlueGreenStageScaleDown(src)
	case v0.StepTypeK8sCanaryDelete:
		step.Template = ConvertStepK8sCanaryDelete(src)
	case v0.StepTypeK8sDiff:
		step.Template = ConvertStepK8sDiff(src)
	case v0.StepTypeK8sDelete:
		step.Template = ConvertStepK8sDelete(src)
	case v0.StepTypeK8sRollout:
		step.Template = ConvertStepK8sRollout(src)
	case v0.StepTypeK8sScale:
		step.Template = ConvertStepK8sScale(src)
	case v0.StepTypeK8sDryRun:
		step.Template = ConvertStepK8sDryRun(src)
	case v0.StepTypeK8sTrafficRouting:
		step.Template = ConvertStepK8sTrafficRouting(src)
	case v0.StepTypeK8sCanaryDeploy:
		step.Template = ConvertStepK8sCanaryDeploy(src)
	case v0.StepTypeK8sBlueGreenDeploy:
		step.Template = ConvertStepK8sBlueGreenDeploy(src)
	case v0.StepTypeHelmBGDeploy:
		step.Template = ConvertStepHelmBGDeploy(src)
	case v0.StepTypeHelmBlueGreenSwapStep:
		step.Template = ConvertStepHelmBlueGreenSwapStep(src)
	case v0.StepTypeHelmCanaryDeploy:
		step.Template = ConvertStepHelmCanaryDeploy(src)
	case v0.StepTypeHelmDelete:
		step.Template = ConvertStepHelmDelete(src)
	case v0.StepTypeHelmDeploy:
		step.Template = ConvertStepHelmDeploy(src)
	case v0.StepTypeHelmRollback:
		step.Template = ConvertStepHelmRollback(src)
	case v0.StepTypeWait:
		step.Action = ConvertStepWait(src)
	case v0.StepTypeShellScript:
		step.Run = ConvertStepShellScript(src)
	case v0.StepTypeBarrier:
		step.Barrier = ConvertStepBarrier(src)
	case v0.StepTypeQueue:
		step.Queue = ConvertStepQueue(src)
	default:
		// Unknown step type, return nil
		fmt.Println("step type: " + src.Type + " is not yet supported!")
		step.Template = &v1.StepTemplate{
			Uses: src.Type,
			With: "to be implemented",
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
		dst.OnFailure = ConvertFailureStrategies(src.FailureStrategies)
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
		dst.Strategy = ConvertStrategy(src.Strategy)
	}

	// Convert delegate selectors

	// extract delegate selectors and includeInfraSelectors from src using reflection
	var delegate_selectors v0.FlexibleField[[]string]
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
						if delegateField.Type().Name() == "FlexibleField[[]string]" {
							// Copy the entire FlexibleField
							delegate_selectors = delegateField.Interface().(v0.FlexibleField[[]string])
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
	// Convert delegate using the extracted values
	delegate := ConvertDelegate(delegate_selectors)

	// Handle includeInfraSelectors for struct-based delegates
	if include_infra_selectors && delegate != nil && !delegate.IsExpression() {
		if delegateStruct, ok := delegate.AsStruct(); ok {
			delegateStruct.Inherit = true
			delegate.Set(delegateStruct)
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

	// Convert stage status to condition
	switch when.StageStatus {
	case "Success":
		return "success()"
	case "Failure":
		return "failure()"
	default:
		return ""
	}
}

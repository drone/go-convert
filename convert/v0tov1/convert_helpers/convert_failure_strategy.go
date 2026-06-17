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

package converthelpers

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertFailureStrategies converts a flexible.Field containing failure strategies from v0 to v1 format.
// This handles both expression strings and arrays of failure strategies.
func ConvertFailureStrategies(src *flexible.Field[[]*v0.FailureStrategy]) *flexible.Field[[]*v1.FailureStrategy] {
	if src == nil || src.IsNil() {
		return nil
	}

	result := &flexible.Field[[]*v1.FailureStrategy]{}

	// Handle expression strings
	if str, ok := src.AsString(); ok {
		result.SetString(str)
		return result
	}

	// Handle arrays of failure strategies
	if strategies, ok := src.AsStruct(); ok {
		convertedStrategies := make([]*v1.FailureStrategy, len(strategies))
		for i, strategy := range strategies {
			convertedStrategies[i] = ConvertFailureStrategy(strategy)
		}
		result.Set(convertedStrategies)
		return result
	}

	return nil
}

// ConvertFailureStrategy converts a single v0 failure strategy to v1 format.
func ConvertFailureStrategy(src *v0.FailureStrategy) *v1.FailureStrategy {
	if src == nil || src.OnFailure == nil {
		return nil
	}

	var errors interface{}
	converted := ConvertErrorNames(src.OnFailure.Errors)
	if len(converted) == 0 {
		errors = []string{} // default to empty array
	} else {
		errors = converted
	}

	return &v1.FailureStrategy{
		Errors: errors,
		Action: ConvertFailureAction(src.OnFailure.Action),
	}
}

// ConvertFailureStrategiesArray converts a list of v0 FailureStrategies to v1 FailureStrategies.
func ConvertFailureStrategiesArray(src []*v0.FailureStrategy) []*v1.FailureStrategy {
	if len(src) == 0 {
		return nil
	}

	dst := make([]*v1.FailureStrategy, 0, len(src))
	for _, strategy := range src {
		if converted := ConvertFailureStrategy(strategy); converted != nil {
			dst = append(dst, converted)
		}
	}

	return dst
}

// ConvertErrorNames converts v0 error names to v1 format using type enums.
// Filters out empty/sentinel values produced by unsupported error types.
func ConvertErrorNames(errors []v0.FailureType) []v1.FailureType {
	if len(errors) == 0 {
		return nil
	}

	converted := make([]v1.FailureType, 0, len(errors))
	for _, err := range errors {
		v1Type := ConvertErrorName(err)
		if v1Type != v1.FailureTypeNone {
			converted = append(converted, v1Type)
		}
	}
	if len(converted) == 0 {
		return nil
	}
	return converted
}

// ConvertErrorName converts a single error name from v0 to v1 format.
func ConvertErrorName(errorType v0.FailureType) v1.FailureType {
	switch errorType {
	case v0.FailureTypeDelegateRestart:
		return v1.FailureTypeDelegateRestart
	case v0.FailureTypeAuthentication:
		return v1.FailureTypeAuthentication
	case v0.FailureTypeConnectivity:
		return v1.FailureTypeConnectivity
	case v0.FailureTypeInputTimeoutError:
		return v1.FailureTypeInputTimeout
	case v0.FailureTypeTaskFailure:
		return v1.FailureTypeUnknown // Map TaskFailure to Unknown in v1
	case v0.FailureTypePolicyEvaluationFailure:
		return v1.FailureTypePolicyEvaluation
	case v0.FailureTypeUnknown:
		return v1.FailureTypeUnknown
	case v0.FailureTypeVerification:
		return v1.FailureTypeVerification
	case v0.FailureTypeDelegateProvisioning:
		return v1.FailureTypeDelegateProvisioning
	case v0.FailureTypeAuthorization:
		return v1.FailureTypeAuthorization
	case v0.FailureTypeApprovalRejection:
		return v1.FailureTypeApprovalRejection
	case v0.FailureTypeTimeout:
		return v1.FailureTypeTimeout
	case v0.FailureTypeUserMarkFail:
		return v1.FailureTypeUserMarkFail
	case v0.FailureTypeAll:
		return v1.FailureTypeAll
	case v0.FailureTypeInfrastructureFailure:
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_ERROR_TYPE",
			fmt.Sprintf("v0 error type %q has no v1 equivalent; dropped", string(errorType)),
			messagelog.WithContext(map[string]string{"error_type": string(errorType)}),
		)
		return v1.FailureTypeNone
	case v0.FailureTypePluginImageFailure:
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_ERROR_TYPE",
			fmt.Sprintf("v0 error type %q has no v1 equivalent; dropped", string(errorType)),
			messagelog.WithContext(map[string]string{"error_type": string(errorType)}),
		)
		return v1.FailureTypeNone
	case v0.FailureTypeResourceLimitsFailure:
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_ERROR_TYPE",
			fmt.Sprintf("v0 error type %q has no v1 equivalent; dropped", string(errorType)),
			messagelog.WithContext(map[string]string{"error_type": string(errorType)}),
		)
		return v1.FailureTypeNone
	case v0.FailureTypeConfigurationFailure:
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_ERROR_TYPE",
			fmt.Sprintf("v0 error type %q has no v1 equivalent; dropped", string(errorType)),
			messagelog.WithContext(map[string]string{"error_type": string(errorType)}),
		)
		return v1.FailureTypeNone
	case v0.FailureTypeRetryableTransientFailure:
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_ERROR_TYPE",
			fmt.Sprintf("v0 error type %q has no v1 equivalent; dropped", string(errorType)),
			messagelog.WithContext(map[string]string{"error_type": string(errorType)}),
		)
		return v1.FailureTypeNone
	default:
		messagelog.GetMessageLogger().LogError(
			"UNKNOWN_ERROR_TYPE",
			fmt.Sprintf("unknown failure error type %q; dropped", string(errorType)),
			messagelog.WithContext(map[string]string{"error_type": string(errorType)}),
		)
		return v1.FailureTypeNone
	}
}

// ConvertFailureAction converts v0 failure action to v1 format.
// Returns either a string for simple actions or an object for complex actions.
// Actions with no v1 equivalent are logged via messagelog and dropped (returns nil).
func ConvertFailureAction(action *v0.Action) interface{} {
	if action == nil {
		return nil
	}

	switch action.Type {
	case v0.ActionTypeMarkAsSuccess:
		return v1.ActionTypeSuccess
	case v0.ActionTypeIgnore:
		return v1.ActionTypeIgnore
	case v0.ActionTypeFail:
		return v1.ActionTypeFail
	case v0.ActionTypeStageRollback:
		return v1.ActionTypeStageRollback
	case v0.ActionTypePipelineRollback:
		return v1.ActionTypePipelineRollback

	case v0.ActionTypeProceedWithDefaultValues:
		return convertSimpleActionWithFailAll(action, string(v1.ActionTypeProceedWithDefaultValues))
	case v0.ActionTypeAbort:
		return convertSimpleActionWithFailAll(action, string(v1.ActionTypeAbort))
	case v0.ActionTypeMarkAsFailure:
		return convertSimpleActionWithFailAll(action, string(v1.ActionTypeFail))

	case v0.ActionTypeRetryStepGroup:
		return map[string]interface{}{
			"retry-step-group": ConvertRetryStepGroupSpec(action),
		}
	case v0.ActionTypeRetry:
		return convertRetryActionWithFailAll(action)
	case v0.ActionTypeManualIntervention:
		return convertManualInterventionActionWithFailAll(action)

	case v0.ActionTypeStepGroupRollback:
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_ACTION_TYPE",
			fmt.Sprintf("v0 action type %q has no v1 equivalent; dropped", string(action.Type)),
			messagelog.WithContext(map[string]string{"action_type": string(action.Type)}),
		)
		return nil

	default:
		messagelog.GetMessageLogger().LogError(
			"UNKNOWN_ACTION_TYPE",
			fmt.Sprintf("unknown failure action type %q; dropped", string(action.Type)),
			messagelog.WithContext(map[string]string{"action_type": string(action.Type)}),
		)
		return nil
	}
}

// convertSimpleActionWithFailAll returns the v1 action type string, or a map
// with fail-all: true if the v0 spec has failAll set.
func convertSimpleActionWithFailAll(action *v0.Action, v1ActionName string) interface{} {
	if action.Spec != nil {
		if spec, ok := action.Spec.(*v0.FailureSpecConfig); ok && spec.FailAll {
			return map[string]interface{}{
				v1ActionName: map[string]interface{}{},
				"fail-all":   true,
			}
		}
	}
	return v1.ActionType(v1ActionName)
}

// convertRetryActionWithFailAll wraps ConvertRetrySpec and adds fail-all if needed.
func convertRetryActionWithFailAll(action *v0.Action) interface{} {
	result := map[string]interface{}{
		"retry": ConvertRetrySpec(action),
	}
	if action.Spec != nil {
		if spec, ok := action.Spec.(*v0.RetrySpec); ok && spec.FailAll {
			result["fail-all"] = true
		}
	}
	return result
}

// convertManualInterventionActionWithFailAll wraps ConvertManualInterventionSpec and adds fail-all if needed.
func convertManualInterventionActionWithFailAll(action *v0.Action) interface{} {
	result := map[string]interface{}{
		"manual-intervention": ConvertManualInterventionSpec(action),
	}
	if action.Spec != nil {
		if spec, ok := action.Spec.(*v0.ManualInterventionSpec); ok && spec.FailAll {
			result["fail-all"] = true
		}
	}
	return result
}

// ConvertRetrySpec converts v0 retry action to v1 retry format.
func ConvertRetrySpec(action *v0.Action) *v1.ActionRetry {
	if action == nil || action.Spec == nil {
		return &v1.ActionRetry{
			Attempts:      1,
			Interval:      []string{"10s"},
			FailureAction: v1.ActionTypeIgnore,
		}
	}

	spec, ok := action.Spec.(*v0.RetrySpec)
	if !ok {
		return &v1.ActionRetry{
			Attempts:      1,
			Interval:      []string{"10s"},
			FailureAction: v1.ActionTypeIgnore,
		}
	}

	retryAction := &v1.ActionRetry{
		Attempts:      int64(spec.RetryCount),
		FailureAction: v1.ActionTypeIgnore, // Default failure action
	}

	// Set interval from retry intervals if available
	if len(spec.RetryIntervals) > 0 {
		retryAction.Interval = spec.RetryIntervals
	} else {
		retryAction.Interval = []string{"10s"} // Default interval
	}

	// Convert onRetryFailure action
	if spec.OnRetryFailure != nil && spec.OnRetryFailure.Action != nil {
		retryAction.FailureAction = ConvertFailureAction(spec.OnRetryFailure.Action)
	}

	// Log and drop condition field (no v1 equivalent)
	if spec.Condition != "" {
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_FIELD",
			fmt.Sprintf("retry condition %q has no v1 equivalent; dropped", spec.Condition),
			messagelog.WithContext(map[string]string{"field": "condition", "value": spec.Condition}),
		)
	}

	return retryAction
}

// ConvertRetryStepGroupSpec converts v0 RetryStepGroup action to v1 format.
// Unlike Retry, RetryStepGroup in v1 only has attempts and interval (no failure-action).
func ConvertRetryStepGroupSpec(action *v0.Action) *v1.ActionRetryStepGroup {
	if action == nil || action.Spec == nil {
		return &v1.ActionRetryStepGroup{
			Attempts: 1,
			Interval: []string{"10s"},
		}
	}

	spec, ok := action.Spec.(*v0.RetrySpec)
	if !ok {
		return &v1.ActionRetryStepGroup{
			Attempts: 1,
			Interval: []string{"10s"},
		}
	}

	result := &v1.ActionRetryStepGroup{
		Attempts: int64(spec.RetryCount),
	}

	if len(spec.RetryIntervals) > 0 {
		result.Interval = spec.RetryIntervals
	} else {
		result.Interval = []string{"10s"}
	}

	// Log and drop condition field (no v1 equivalent)
	if spec.Condition != "" {
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_FIELD",
			fmt.Sprintf("retry-step-group condition %q has no v1 equivalent; dropped", spec.Condition),
			messagelog.WithContext(map[string]string{"field": "condition", "value": spec.Condition}),
		)
	}

	return result
}

// ConvertManualInterventionSpec converts v0 manual intervention action to v1 format.
func ConvertManualInterventionSpec(action *v0.Action) *v1.ActionManual {
	if action == nil || action.Spec == nil {
		return &v1.ActionManual{
			Timeout:       "1h",
			TimeoutAction: "abort",
		}
	}

	spec, ok := action.Spec.(*v0.ManualInterventionSpec)
	if !ok {
		return &v1.ActionManual{
			Timeout:       "1h",
			TimeoutAction: "abort",
		}
	}

	manualAction := &v1.ActionManual{
		Timeout:       spec.Timeout,
		TimeoutAction: "abort", // Default timeout action
	}

	if spec.OnTimeout != nil {
		manualAction.TimeoutAction = ConvertFailureAction(spec.OnTimeout.Action)
	}

	// Convert availableActions: v0 enum strings → v1 action config objects
	if len(spec.AvailableActions) > 0 {
		manualAction.AvailableActions = ConvertAvailableActions(spec.AvailableActions)
	}

	return manualAction
}

// ConvertAvailableActions converts v0 available action enum strings to v1 action config objects.
// v0: ["Ignore", "Retry", "MarkAsSuccess"] → v1: [{"ignore": {}}, {"retry": {}}, {"success": {}}]
// Actions with no v1 equivalent are logged and skipped.
func ConvertAvailableActions(actions []string) []interface{} {
	if len(actions) == 0 {
		return nil
	}

	v0ToV1ActionName := map[string]string{
		"Ignore":                   "ignore",
		"Abort":                    "abort",
		"MarkAsSuccess":            "success",
		"MarkAsFailure":            "fail",
		"StageRollback":            "stage-rollback",
		"PipelineRollback":         "pipeline-rollback",
		"ProceedWithDefaultValues": "proceed-with-default",
		"Retry":                    "retry",
		"RetryStepGroup":           "retry-step-group",
		"ManualIntervention":       "manual-intervention",
	}

	result := make([]interface{}, 0, len(actions))
	for _, action := range actions {
		v1Name, ok := v0ToV1ActionName[action]
		if !ok {
			messagelog.GetMessageLogger().LogError(
				"UNSUPPORTED_ACTION_TYPE",
				fmt.Sprintf("available action %q has no v1 equivalent; dropped", action),
				messagelog.WithContext(map[string]string{"action": action}),
			)
			continue
		}
		result = append(result, map[string]interface{}{
			v1Name: map[string]interface{}{},
		})
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

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
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertFailureStrategyFlexible converts a FlexibleField containing failure strategies from v0 to v1 format.
// This handles both expression strings and arrays of failure strategies.
func ConvertFailureStrategies(src *v0.FlexibleField[[]*v0.FailureStrategy]) *v1.FlexibleField[[]*v1.FailureStrategy] {
	if src == nil || src.IsNil() {
		return nil
	}

	result := &v1.FlexibleField[[]*v1.FailureStrategy]{}

	// Handle expression strings
	if src.IsExpression() {
		result.SetExpression(src.AsString())
		return result
	}

	// Handle struct arrays
	if strategies, ok := src.AsStruct(); ok {
		converted := ConvertFailureStrategiesArray(strategies)
		if converted != nil {
			result.Set(converted)
			return result
		}
	}

	return nil
}

// ConvertFailureStrategy converts a single v0 failure strategy to v1 format.
func ConvertFailureStrategy(src *v0.FailureStrategy) *v1.FailureStrategy {
	if src == nil || src.OnFailure == nil {
		return nil
	}

	return &v1.FailureStrategy{
		Errors: ConvertErrorNames(src.OnFailure.Errors),
		Action: ConvertFailureAction(src.OnFailure.Action),
	}
}

// ConvertFailureStrategies converts a list of v0 FailureStrategies to v1 FailureStrategies.
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
func ConvertErrorNames(errors []v0.FailureType) []v1.FailureType {
	if len(errors) == 0 {
		return nil
	}

	converted := make([]v1.FailureType, len(errors))
	for i, err := range errors {
		converted[i] = ConvertErrorName(err)
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
	default:
		fmt.Println("Unknown error type: " + string(errorType))
		return v1.FailureTypeNone
	}
}

// ConvertFailureAction converts v0 failure action to v1 format.
// Returns either a string for simple actions or an object for complex actions.
func ConvertFailureAction(action *v0.Action) interface{} {
	if action == nil {
		return nil
	}

	switch action.Type {
	case v0.ActionTypeMarkAsSuccess:
		return v1.ActionTypeSuccess
	case v0.ActionTypeIgnore:
		return v1.ActionTypeIgnore
	case v0.ActionTypeAbort:
		return v1.ActionTypeAbort
	case v0.ActionTypeFail:
		return v1.ActionTypeFail
	case v0.ActionTypeStageRollback:
		return v1.ActionTypeStageRollback
	case v0.ActionTypePipelineRollback:
		return v1.ActionTypePipelineRollback
	case v0.ActionTypeRetryStepGroup:
		return v1.ActionTypeRetryStepGroup
	case v0.ActionTypeRetry:
		return map[string]interface{}{
			"retry": ConvertRetrySpec(action),
		}
	case v0.ActionTypeManualIntervention:
		return map[string]interface{}{
			"manual-intervention": ConvertManualInterventionSpec(action),
		}
	default:
		// Default to ignore for unknown action types
		return v1.ActionTypeIgnore
	}
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

	return retryAction
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

	return manualAction
}

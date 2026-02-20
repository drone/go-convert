package yaml

import (
	"encoding/json"
	"fmt"
)

type FailureType string

// FailureType enumeration for v0.
const (
	FailureTypeNone                    FailureType = ""
	FailureTypeAll                     FailureType = "AllErrors"
	FailureTypeApprovalRejection       FailureType = "ApprovalRejection"
	FailureTypeAuthentication          FailureType = "Authentication"
	FailureTypeAuthorization           FailureType = "Authorization"
	FailureTypeConnectivity            FailureType = "Connectivity"
	FailureTypeDelegateProvisioning    FailureType = "DelegateProvisioning"
	FailureTypeDelegateRestart         FailureType = "DelegateRestart"
	FailureTypeInputTimeoutError       FailureType = "InputTimeoutError"
	FailureTypePolicyEvaluationFailure FailureType = "PolicyEvaluationFailure"
	FailureTypeTaskFailure             FailureType = "TaskFailure"
	FailureTypeTimeout                 FailureType = "Timeout"
	FailureTypeUnknown                 FailureType = "Unknown"
	FailureTypeVerification            FailureType = "Verification"
	FailureTypeUserMarkFail            FailureType = "UserMarkedFailure"
)

type ActionType string

// ActionType enumeration for v0.
const (
	ActionTypeNone               ActionType = ""
	ActionTypeAbort              ActionType = "Abort"
	ActionTypeFail               ActionType = "Fail"
	ActionTypeIgnore             ActionType = "Ignore"
	ActionTypeManualIntervention ActionType = "ManualIntervention"
	ActionTypeMarkAsSuccess      ActionType = "MarkAsSuccess"
	ActionTypePipelineRollback   ActionType = "PipelineRollback"
	ActionTypeRetry              ActionType = "Retry"
	ActionTypeRetryStepGroup     ActionType = "RetryStepGroup"
	ActionTypeStageRollback      ActionType = "StageRollback"
	ActionTypeMarkAsFailure      ActionType = "MarkAsFailure"
	ActionTypeProceedWithDefaultValues ActionType = "ProceedWithDefaultValues"
)

type FailureStrategy struct {
	OnFailure *OnFailure `json:"onFailure,omitempty" yaml:"onFailure,omitempty"`
}

type OnFailure struct {
	Errors []FailureType `json:"errors,omitempty" yaml:"errors,omitempty"`
	Action *Action       `json:"action,omitempty" yaml:"action,omitempty"`
}

type Action struct {
	Type ActionType  `json:"type,omitempty" yaml:"type,omitempty"`
	Spec interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// UnmarshalJSON implement the json.Unmarshaler interface.
func (a *Action) UnmarshalJSON(data []byte) error {
	type A Action
	type T struct {
		*A
		Spec json.RawMessage `json:"spec"`
	}

	obj := &T{A: (*A)(a)}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}

	switch a.Type {
	case ActionTypeRetry,ActionTypeRetryStepGroup:
		a.Spec = new(RetrySpec)
	case ActionTypeManualIntervention:
		a.Spec = new(ManualInterventionSpec)
	case ActionTypeMarkAsFailure,ActionTypeMarkAsSuccess, ActionTypeIgnore, ActionTypeAbort, ActionTypeStageRollback, ActionTypePipelineRollback, ActionTypeProceedWithDefaultValues:
		// These actions don't have specs
		return nil
	default:
		return fmt.Errorf("unknown action type %s", a.Type)
	}

	if obj.Spec != nil {
		return json.Unmarshal(obj.Spec, a.Spec)
	}
	return nil
}

type RetrySpec struct {
	RetryCount     int             `json:"retryCount,omitempty" yaml:"retryCount,omitempty"`
	RetryIntervals []string        `json:"retryIntervals,omitempty" yaml:"retryIntervals,omitempty"`
	OnRetryFailure *OnRetryFailure `json:"onRetryFailure,omitempty" yaml:"onRetryFailure,omitempty"`
}

type OnRetryFailure struct {
	Action *Action `json:"action,omitempty" yaml:"action,omitempty"`
}

type ManualInterventionSpec struct {
	Timeout       string  `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	OnTimeout *OnTimeout `json:"onTimeout,omitempty" yaml:"onTimeout,omitempty"`
}

type OnTimeout struct {
	Action *Action `json:"action,omitempty" yaml:"action,omitempty"`
}

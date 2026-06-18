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
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertStepHarnessApproval converts a v0 HarnessApproval step to v1 approval format
func ConvertStepHarnessApproval(src *v0.Step) *v1.StepApproval {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Type assert the spec to StepHarnessApproval
	spec, ok := src.Spec.(*v0.StepHarnessApproval)
	if !ok {
		return nil
	}

	dst := &v1.StepApproval{
		Uses: "harness",
		With: make(map[string]interface{}),
	}

	// Map approvalMessage to message
	if spec.ApprovalMessage != "" {
		dst.With["message"] = spec.ApprovalMessage
	}

	// Map includePipelineExecutionHistory to execution-details
	if spec.IncludePipelineExecutionHistory != nil {
		dst.With["execution-details"] = spec.IncludePipelineExecutionHistory
	}

	// Map isAutoRejectEnabled to auto-reject
	if spec.IsAutoRejectEnabled != nil {
		dst.With["auto-reject"] = spec.IsAutoRejectEnabled
	}
	
	// Map callbackId to callback
	if spec.CallbackId != "" {
		dst.With["callback"] = spec.CallbackId
	}

	// Map approvers fields
	if spec.Approvers != nil {
		// Map minimumCount to approvers-min-count
		if spec.Approvers.MinimumCount != nil {
			dst.With["approvers-min-count"] = spec.Approvers.MinimumCount
		}

		// Map disallowPipelineExecutor to block-executor
		if spec.Approvers.DisallowPipelineExecutor != nil {
			dst.With["block-executor"] = spec.Approvers.DisallowPipelineExecutor
		}

		// Map userGroups to user-groups
		if spec.Approvers.UserGroups != nil {
			dst.With["user-groups"] = spec.Approvers.UserGroups
		}

		// Map serviceAccounts to service-accounts
		if spec.Approvers.ServiceAccounts != nil {
			dst.With["service-accounts"] = spec.Approvers.ServiceAccounts
		}
	}

	// Map approverInputs to inputs
	if len(spec.ApproverInputs) > 0 {
		inputs := convertApproverInputs(spec.ApproverInputs)
		if len(inputs) > 0 {
			dst.With["inputs"] = inputs
		}
	}

	dst.With["auto-approve"] = false //step does not work without this default value
	// Map autoApproval fields
	if spec.AutoApproval != nil {
		// Map action to auto-approve (true if action is APPROVE)
		if spec.AutoApproval.Action == "APPROVE" {
			dst.With["auto-approve"] = true
		}

		// Map comments
		if spec.AutoApproval.Comments != "" {
			dst.With["comments"] = spec.AutoApproval.Comments
		}

		// Map scheduledDeadline fields
		if spec.AutoApproval.ScheduledDeadline != nil {
			if spec.AutoApproval.ScheduledDeadline.Time != "" {
				dst.With["deadline"] = spec.AutoApproval.ScheduledDeadline.Time
			}
			if spec.AutoApproval.ScheduledDeadline.TimeZone != "" {
				dst.With["timezone"] = spec.AutoApproval.ScheduledDeadline.TimeZone
			}
		}

		// Map callbackId to callback
		if spec.AutoApproval.CallbackId != "" {
			dst.With["callback"] = spec.AutoApproval.CallbackId
		}
	}

	return dst
}

// convertApproverInputs converts v0 approverInputs array to v1 inputs map.
// v0: [{name, defaultValue, regex, allowedValues, selectOneFrom, required, description}]
// v1: map[name]ApproverInputDef{required, description, default, enum, multi-select, pattern}
func convertApproverInputs(inputs []*v0.ApproverInput) map[string]*v1.ApproverInputDef {
	result := make(map[string]*v1.ApproverInputDef)
	for _, input := range inputs {
		if input == nil || input.Name == "" {
			continue
		}

		def := &v1.ApproverInputDef{
			Required:    input.Required,
			Description: input.Description,
			Default:     input.DefaultValue,
			Pattern:     input.Regex,
		}

		// selectOneFrom → enum (single-select dropdown)
		if input.SelectOneFrom != nil {
			def.Enum = input.SelectOneFrom
		}

		// allowedValues → enum with multi-select
		if input.AllowedValues != nil {
			def.Enum = input.AllowedValues
			ms := &flexible.Field[bool]{}
			ms.Set(true)
			def.MultiSelect = ms
		}

		result[input.Name] = def
	}
	return result
}

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
	dst.With["execution-details"] = spec.IncludePipelineExecutionHistory

	// Map isAutoRejectEnabled to auto-reject
	dst.With["auto-reject"] = spec.IsAutoRejectEnabled

	// Map approvers fields
	if spec.Approvers != nil {
		// Map minimumCount to approvers-min-count
		if spec.Approvers.MinimumCount != nil {
			dst.With["approvers-min-count"] = spec.Approvers.MinimumCount
		}

		// Map disallowPipelineExecutor to block-executor
		dst.With["block-executor"] = spec.Approvers.DisallowPipelineExecutor

		// Map userGroups to user-groups
		if spec.Approvers.UserGroups != nil {
			dst.With["user-groups"] = spec.Approvers.UserGroups
		}
	}

	// Map approverInputs to params
	dst.With["params"] = [][]map[string]string{} //step does not work without this default value
	if len(spec.ApproverInputs) > 0 {
		params := make([]map[string]string, 0, len(spec.ApproverInputs))
		for _, input := range spec.ApproverInputs {
			if input != nil && input.Name != "" {
				param := map[string]string{
					input.Name: input.DefaultValue,
				}
				params = append(params, param)
			}
		}
		if len(params) > 0 {
			dst.With["params"] = params
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
	}

	return dst
}

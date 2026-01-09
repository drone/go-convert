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
	"reflect"
	"strings"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertNotifications converts v0 notification rules to v1 notifications
func ConvertNotifications(v0Rules []*v0.NotificationRule) []*v1.Notification {
	if len(v0Rules) == 0 {
		return nil
	}

	var notifications []*v1.Notification
	for _, rule := range v0Rules {
		if rule == nil {
			continue
		}

		notification := &v1.Notification{
			ID:       rule.Identifier,
			Name:     rule.Name,
			Disabled: !rule.Enabled,
		}

		// Convert notification method to uses and with
		if rule.NotificationMethod != nil {
			notification.Uses = convertNotificationMethodType(rule.NotificationMethod.Type)
			notification.With = convertNotificationMethodSpec(rule.NotificationMethod)
		}

		// Convert pipeline events to "on" field
		notification.On = convertPipelineEvents(rule.PipelineEvents)

		notification.NotificationTemplate = convertNotificationTemplate(rule.NotificationTemplate)

		notifications = append(notifications, notification)
	}

	return notifications
}

func convertNotificationTemplate(template *v0.NotificationTemplate) *v1.NotificationTemplate {
	if template == nil {
		return nil
	}

	with := make(map[string]interface{})
	with["version"] = template.VersionLabel

	return &v1.NotificationTemplate{
		Uses: template.TemplateRef,
		With: with,
	}
}

// convertNotificationMethodType converts v0 notification method type to v1 uses field
func convertNotificationMethodType(methodType string) string {
	switch methodType {
	case "Webhook":
		return "webhook"
	case "Slack":
		return "slack"
	case "PagerDuty":
		return "pagerduty"
	case "MsTeams":
		return "ms-teams"
	case "Email":
		return "email"
	case "Datadog":
		return "datadog"
	default:
		return strings.ToLower(methodType)
	}
}

// convertNotificationMethodSpec converts v0 notification method spec to v1 with field
func convertNotificationMethodSpec(method *v0.NotificationMethod) map[string]interface{} {
	if method == nil || method.Spec == nil {
		return nil
	}

	with := make(map[string]interface{})

    // Use reflection to extract CommonNotificationSpec fields
    specValue := reflect.ValueOf(method.Spec)
    if specValue.Kind() == reflect.Struct {
        // Try to get ExecuteOnDelegate field
        if executeOnDelegateField := specValue.FieldByName("ExecuteOnDelegate"); executeOnDelegateField.IsValid() {
            if executeOnDelegateField.Kind() == reflect.Bool {
                if executeOnDelegate := executeOnDelegateField.Bool(); executeOnDelegate {
                    with["execute-on-delegate"] = executeOnDelegate
                }
            }
        }

		// Try to get DelegateSelectors field
		if delegateSelectorsField := specValue.FieldByName("DelegateSelectors"); delegateSelectorsField.IsValid() {
			if delegateSelectors, ok := delegateSelectorsField.Interface().([]string); ok && len(delegateSelectors) > 0 {
				with["delegate-selectors"] = delegateSelectors
			}
		}
	}

	switch method.Type {
	case "Webhook":
		if spec, ok := method.Spec.(v0.WebhookNotificationSpec); ok {
			with["url"] = spec.WebhookUrl
			if len(spec.Headers) > 0 {
				with["headers"] = spec.Headers
			}
		}
	case "Slack":
		if spec, ok := method.Spec.(v0.SlackNotificationSpec); ok {
			with["webhook"] = spec.WebhookUrl
			if len(spec.UserGroups) > 0 {
				with["user-groups"] = spec.UserGroups
			}
		}
	case "PagerDuty":
		if spec, ok := method.Spec.(v0.PagerDutyNotificationSpec); ok {
			with["key"] = spec.IntegrationKey
			if len(spec.UserGroups) > 0 {
				with["user-groups"] = spec.UserGroups
			}
		}
	case "MsTeams":
		if spec, ok := method.Spec.(v0.MsTeamsNotificationSpec); ok {
			if len(spec.MsTeamKeys) > 0 {
				with["keys"] = spec.MsTeamKeys
			}
			if len(spec.UserGroups) > 0 {
				with["user-groups"] = spec.UserGroups
			}
		}
	case "Email":
		if spec, ok := method.Spec.(v0.EmailNotificationSpec); ok {
			if len(spec.Recipients) > 0 {
				with["recipients"] = spec.Recipients
			}
			if len(spec.UserGroups) > 0 {
				with["user-groups"] = spec.UserGroups
			}
		}
	case "Datadog":
		if spec, ok := method.Spec.(v0.DatadogNotificationSpec); ok {
			with["api-key"] = spec.ApiKey
			with["url"] = spec.Url
			if len(spec.Headers) > 0 {
				with["headers"] = spec.Headers
			}
		}
	}

	return with
}

// convertPipelineEvents converts v0 pipeline events to v1 "on" field
func convertPipelineEvents(events []*v0.PipelineEvent) []*v1.NotificationOn {
	if len(events) == 0 {
		return nil
	}

	// Check for AllEvents first - if present, use only the AllEvents format
	for _, event := range events {
		if event != nil && event.Type == "AllEvents" {
			return []*v1.NotificationOn{
				{Pipeline: "all"},
				{Stage: "all"},
				{Step: "all"},
			}
		}
	}

	// Group events by type (only if AllEvents is not present)
	pipelineEvents := []string{}
	stageEvents := make(map[string]interface{})
	stepEvents := []string{}

	for _, event := range events {
		if event == nil {
			continue
		}

		switch event.Type {
		case "PipelineStart":
			pipelineEvents = append(pipelineEvents, "start")
		case "PipelineEnd":
			pipelineEvents = append(pipelineEvents, "end")
		case "PipelineSuccess":
			pipelineEvents = append(pipelineEvents, "success")
		case "PipelineFailed":
			pipelineEvents = append(pipelineEvents, "failed")
		case "StageStart":
			addStageEvent(stageEvents, "start", event.ForStages)
		case "StageSuccess":
			addStageEvent(stageEvents, "success", event.ForStages)
		case "StageFailed":
			addStageEvent(stageEvents, "failed", event.ForStages)
		case "StepFailed":
			stepEvents = append(stepEvents, "failed")
		}
	}

	// Build the "on" array
	var onEvents []*v1.NotificationOn

	// Add pipeline events
	if len(pipelineEvents) > 0 {
		on := &v1.NotificationOn{}
		if len(pipelineEvents) == 1 {
			on.Pipeline = pipelineEvents[0]
		} else {
			on.Pipeline = pipelineEvents
		}
		onEvents = append(onEvents, on)
	}

	// Add stage events
	if len(stageEvents) > 0 {
		on := &v1.NotificationOn{
			Stage: stageEvents,
		}
		onEvents = append(onEvents, on)
	}

	// Add step events
	if len(stepEvents) > 0 {
		on := &v1.NotificationOn{}
		if len(stepEvents) == 1 {
			on.Step = stepEvents[0]
		} else {
			on.Step = stepEvents
		}
		onEvents = append(onEvents, on)
	}

	return onEvents
}

// addStageEvent adds a stage event to the stage events map
func addStageEvent(stageEvents map[string]interface{}, eventType string, forStages []string) {
	if len(forStages) == 0 {
		return
	}

	// Check if forStages contains "AllStages"
	for _, stage := range forStages {
		if stage == "AllStages" {
			stageEvents[eventType] = "all"
			return
		}
	}

	// If not AllStages, use the specific stage list
	if len(forStages) == 1 {
		stageEvents[eventType] = forStages[0]
	} else {
		stageEvents[eventType] = forStages
	}
}
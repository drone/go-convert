package yaml

import (
	"encoding/json"
	"fmt"
)

type (
	NotificationRule struct {
		Identifier         string              `json:"identifier,omitempty" yaml:"identifier,omitempty"`
		Name               string              `json:"name,omitempty" yaml:"name,omitempty"`
		Enabled            bool                `json:"enabled,omitempty" yaml:"enabled,omitempty"`
		PipelineEvents     []*PipelineEvent    `json:"pipelineEvents,omitempty" yaml:"pipelineEvents,omitempty"`
		NotificationMethod *NotificationMethod `json:"notificationMethod,omitempty" yaml:"notificationMethod,omitempty"`
	}

	PipelineEvent struct {
		Type      string   `json:"type,omitempty" yaml:"type,omitempty"`
		ForStages []string `json:"forStages,omitempty" yaml:"forStages,omitempty"`
	}

	NotificationMethod struct {
		Type string      `json:"type,omitempty" yaml:"type,omitempty"`
		Spec interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// NotificationMethodSpec defines the base interface for notification method specifications
	NotificationMethodSpec interface{}

	// WebhookNotificationSpec defines webhook notification configuration
	WebhookNotificationSpec struct {
		WebhookUrl string            `json:"webhookUrl,omitempty" yaml:"webhookUrl,omitempty"`
		Headers    map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	}

	// PagerDutyNotificationSpec defines PagerDuty notification configuration
	PagerDutyNotificationSpec struct {
		UserGroups     []string `json:"userGroups,omitempty" yaml:"userGroups,omitempty"`
		IntegrationKey string   `json:"integrationKey,omitempty" yaml:"integrationKey,omitempty"`
	}

	// SlackNotificationSpec defines Slack notification configuration
	SlackNotificationSpec struct {
		UserGroups []string `json:"userGroups,omitempty" yaml:"userGroups,omitempty"`
		WebhookUrl string   `json:"webhookUrl,omitempty" yaml:"webhookUrl,omitempty"`
	}

	// MsTeamsNotificationSpec defines Microsoft Teams notification configuration
	MsTeamsNotificationSpec struct {
		UserGroups []string `json:"userGroups,omitempty" yaml:"userGroups,omitempty"`
		MsTeamKeys []string `json:"msTeamKeys,omitempty" yaml:"msTeamKeys,omitempty"`
	}

	// EmailNotificationSpec defines Email notification configuration
	EmailNotificationSpec struct {
		UserGroups []string `json:"userGroups,omitempty" yaml:"userGroups,omitempty"`
		Recipients []string `json:"recipients,omitempty" yaml:"recipients,omitempty"`
	}

	// DatadogNotificationSpec defines Datadog notification configuration
	DatadogNotificationSpec struct {
		ApiKey  string            `json:"apiKey,omitempty" yaml:"apiKey,omitempty"`
		Url     string            `json:"url,omitempty" yaml:"url,omitempty"`
		Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	}
)

func (nm *NotificationMethod) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type string          `json:"type"`
		Spec json.RawMessage `json:"spec"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	nm.Type = aux.Type

	if len(aux.Spec) == 0 {
		nm.Spec = nil
		return nil
	}

	switch nm.Type {
	case "Webhook":
		var spec WebhookNotificationSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	case "PagerDuty":
		var spec PagerDutyNotificationSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	case "Slack":
		var spec SlackNotificationSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	case "MsTeams":
		var spec MsTeamsNotificationSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	case "Email":
		var spec EmailNotificationSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	case "Datadog":
		var spec DatadogNotificationSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	default:
		// For unknown types, keep as raw JSON
		fmt.Println("unknown notification method type: " + nm.Type)
		var spec interface{}
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		nm.Spec = spec
	}

	return nil
}

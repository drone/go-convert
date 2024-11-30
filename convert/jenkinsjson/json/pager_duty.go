package json

import (
	"encoding/json"
	"log"

	harness "github.com/drone/spec/dist/go"
)

// ConvertPagerduty creates a Harness step for nunit plugin.
func ConvertPagerDuty(node Node, arguments map[string]interface{}) *harness.Step {
	incidentSource, _ := arguments["incidentSource"].(string)
	resolve, _ := arguments["resolve"].(bool)
	dedupKey, _ := arguments["dedupKey"].(string)
	incidentSummary, _ := arguments["incidentSummary"].(string)
	incidentSeverity, _ := arguments["incidentSeverity"].(string)
	routingKey, _ := arguments["routingKey"].(string)

	convertPagerduty := &harness.Step{
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Name: "Pagerduty",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/pagerduty",
			With: map[string]interface{}{
				"log_level":         "info",
				"routing_key":       routingKey,
				"incident_summary":  incidentSummary,
				"incident_source":   incidentSource,
				"incident_severity": incidentSeverity,
				"resolve":           resolve,
				"job_status":        "<+pipeline.status>",
				"dedup_key":         dedupKey,
			},
		},
	}

	return convertPagerduty
}

// ConvertPagerDutyChangeEvent creates a Harness step for nunit plugin.
func ConvertPagerDutyChangeEvent(node Node, arguments map[string]interface{}) *harness.Step {
	incidentSource, _ := arguments["incidentSource"].(string)
	incidentSummary, _ := arguments["summaryText"].(string)
	integrationkey, _ := arguments["integrationKey"].(string)
	customDetails, _ := arguments["customDetails"].(map[string]interface{})

	customDetailsStr := ""
	if customDetails != nil {
		jsonBytes, err := json.Marshal(customDetails)
		if err != nil {
			log.Printf("Failed to marshal customDetails: %v", err)
		} else {
			customDetailsStr = string(jsonBytes)
		}
	}

	convertPagerdutyChangeEvent := &harness.Step{
		Name: "Pagerduty_Change_Event",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/pagerduty",
			With: map[string]interface{}{
				"log_level":           "info",
				"integration_key":     integrationkey,
				"incident_summary":    incidentSummary,
				"incident_source":     incidentSource,
				"create_change_event": true,
				"custom_details":      customDetailsStr,
			},
		},
	}

	return convertPagerdutyChangeEvent
}

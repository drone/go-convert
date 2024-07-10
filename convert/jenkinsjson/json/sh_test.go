package json

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertSh(t *testing.T) {
	tests := []struct {
		json string
		want Node
	}{
		{
			json: `
{
  "spanId": "63af111b367be749",
  "traceId": "5452e5febc747713f30e334cb298db4f",
  "parent": "first",
  "all-info": "span(name: sh, spanId: 63af111b367be749, parentSpanId: c89b0adcaaf5f60c, traceId: 5452e5febc747713f30e334cb298db4f, attr: ci.pipeline.run.user:SYSTEM;harness-attribute:{\n  \"script\" : \"exit 1\"\n};harness-others:-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long;jenkins.pipeline.step.id:7;jenkins.pipeline.step.name:Shell Script;jenkins.pipeline.step.plugin.name:workflow-durable-task-step;jenkins.pipeline.step.plugin.version:1360.v82d13453da_a_f;jenkins.pipeline.step.type:sh;)",
  "name": "first #5",
  "attributesMap": {
    "harness-others": "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
    "jenkins.pipeline.step.name": "Shell Script",
    "ci.pipeline.run.user": "SYSTEM",
    "jenkins.pipeline.step.id": "7",
    "jenkins.pipeline.step.type": "sh",
    "harness-attribute": "{\n  \"script\" : \"exit 1\"\n}",
    "jenkins.pipeline.step.plugin.name": "workflow-durable-task-step",
    "jenkins.pipeline.step.plugin.version": "1360.v82d13453da_a_f"
  },
  "type": "Run Phase Span",
  "parentSpanId": "c89b0adcaaf5f60c",
  "parameterMap": {"script": "exit 1"},
  "spanName": "sh"
}
			`,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "7",
					"jenkins.pipeline.step.name":           "Shell Script",
					"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
					"jenkins.pipeline.step.plugin.version": "1360.v82d13453da_a_f",
					"jenkins.pipeline.step.type":           "sh",
					"harness-attribute":                    "{\n  \"script\" : \"exit 1\"\n}",
					"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
				},
				Name:         "first #5",
				Parent:       "first",
				ParentSpanId: "c89b0adcaaf5f60c",
				SpanId:       "63af111b367be749",
				SpanName:     "sh",
				TraceId:      "5452e5febc747713f30e334cb298db4f",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"script": string("exit 1")},
			},
		},
	}
	for i, test := range tests {
		got := new(Node)
		if err := json.Unmarshal([]byte(test.json), got); err != nil {
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

package json

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertReadFile(t *testing.T) {
	tests := []struct {
		json string
		want Node
	}{
		{
			json: `
{
  "spanId": "484461c194e15443",
  "traceId": "900da4b8783a6df372c284b2fbdcf8a8",
  "parent": "readfile",
  "all-info": "span(name: readFile, spanId: 484461c194e15443, parentSpanId: 5af3d9cada8b7ebe, traceId: 900da4b8783a6df372c284b2fbdcf8a8, attr: ci.pipeline.run.user:SYSTEM;harness-attribute:{\n  \"file\" : \"output1.txt\"\n};harness-others:;jenkins.pipeline.step.id:14;jenkins.pipeline.step.name:Read file from workspace;jenkins.pipeline.step.plugin.name:workflow-basic-steps;jenkins.pipeline.step.plugin.version:1058.vcb_fc1e3a_21a_9;jenkins.pipeline.step.type:readFile;)",
  "name": "readfile #1",
  "attributesMap": {
    "harness-others": "",
    "jenkins.pipeline.step.name": "Read file from workspace",
    "ci.pipeline.run.user": "SYSTEM",
    "jenkins.pipeline.step.id": "14",
    "jenkins.pipeline.step.type": "readFile",
    "harness-attribute": "{\n  \"file\" : \"output1.txt\"\n}",
    "jenkins.pipeline.step.plugin.name": "workflow-basic-steps",
    "jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9"
  },
  "type": "Run Phase Span",
  "parentSpanId": "5af3d9cada8b7ebe",
  "parameterMap": {"file": "output1.txt"},
  "spanName": "readFile"
}
			`,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "14",
					"jenkins.pipeline.step.name":           "Read file from workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "readFile",
					"harness-attribute":                    "{\n  \"file\" : \"output1.txt\"\n}",
					"harness-others":                       "",
				},
				Name:         "readfile #1",
				Parent:       "readfile",
				ParentSpanId: "5af3d9cada8b7ebe",
				SpanId:       "484461c194e15443",
				SpanName:     "readFile",
				TraceId:      "900da4b8783a6df372c284b2fbdcf8a8",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": string("output1.txt")},
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

package json

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertWriteFile(t *testing.T) {
	tests := []struct {
		json string
		want Node
	}{
		{
			json: `
{
  "spanId": "705399e59b21cf47",
  "traceId": "900da4b8783a6df372c284b2fbdcf8a8",
  "parent": "writefile",
  "all-info": "span(name: writeFile, spanId: 705399e59b21cf47, parentSpanId: ffaea7a0704eaf64, traceId: 900da4b8783a6df372c284b2fbdcf8a8, attr: ci.pipeline.run.user:SYSTEM;harness-attribute:{\n  \"file\" : \"output1.txt\",\n  \"text\" : \"line1 \\n\\nline2\\n\\nline3\"\n};harness-others:;jenkins.pipeline.step.id:7;jenkins.pipeline.step.name:Write file to workspace;jenkins.pipeline.step.plugin.name:workflow-basic-steps;jenkins.pipeline.step.plugin.version:1058.vcb_fc1e3a_21a_9;jenkins.pipeline.step.type:writeFile;)",
  "name": "writefile #1",
  "attributesMap": {
    "harness-others": "",
    "jenkins.pipeline.step.name": "Write file to workspace",
    "ci.pipeline.run.user": "SYSTEM",
    "jenkins.pipeline.step.id": "7",
    "jenkins.pipeline.step.type": "writeFile",
    "harness-attribute": "{\n  \"file\" : \"output1.txt\",\n  \"text\" : \"line1 \\n\\nline2\\n\\nline3\"\n}",
    "jenkins.pipeline.step.plugin.name": "workflow-basic-steps",
    "jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9"
  },
  "type": "Run Phase Span",
  "parentSpanId": "ffaea7a0704eaf64",
  "parameterMap": {
    "file": "output1.txt",
    "text": "line1 \n\nline2\n\nline3"
  },
  "spanName": "writeFile"
}
			`,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "7",
					"jenkins.pipeline.step.name":           "Write file to workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "writeFile",
					"harness-attribute":                    "{\n  \"file\" : \"output1.txt\",\n  \"text\" : \"line1 \\n\\nline2\\n\\nline3\"\n}",
					"harness-others":                       "",
				},
				Name:         "writefile #1",
				Parent:       "writefile",
				ParentSpanId: "ffaea7a0704eaf64",
				SpanId:       "705399e59b21cf47",
				SpanName:     "writeFile",
				TraceId:      "900da4b8783a6df372c284b2fbdcf8a8",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": string("output1.txt"), "text": string("line1 \n\nline2\n\nline3")},
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

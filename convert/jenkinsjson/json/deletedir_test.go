package json

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertDir(t *testing.T) {
	tests := []struct {
		json string
		want Node
	}{
		{
			// children are not considered here for unit tests. Refer TestConvertDeleteDir as it is dir's child node
			json: `
{
  "spanId": "fd7dfe5971c096aa",
  "traceId": "473b5dc91e544902871080a25554e963",
  "parent": "Dir",
  "all-info": "span(name: Stage: null, spanId: fd7dfe5971c096aa, parentSpanId: 9c5d0bbf29308002, traceId: 473b5dc91e544902871080a25554e963, attr: harness-attribute:{\n  \"path\" : \"target\"\n};harness-others:;jenkins.pipeline.step.type:dir;)",
  "name": "Dir #1",
  "attributesMap": {
    "harness-others": "",
    "jenkins.pipeline.step.type": "dir",
    "harness-attribute": "{\n  \"path\" : \"target\"\n}"
  },
  "type": "Run Phase Span",
  "parentSpanId": "9c5d0bbf29308002",
  "parameterMap": {"path": "target"},
  "spanName": "Stage: null"
}
			`,
			want: Node{
				AttributesMap: map[string]string{
					"jenkins.pipeline.step.type": "dir",
					"harness-attribute":          "{\n  \"path\" : \"target\"\n}",
					"harness-others":             "",
				},
				Name:         "Dir #1",
				Parent:       "Dir",
				ParentSpanId: "9c5d0bbf29308002",
				SpanId:       "fd7dfe5971c096aa",
				SpanName:     "Stage: null",
				TraceId:      "473b5dc91e544902871080a25554e963",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"path": string("target")},
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

func TestConvertDeleteDir(t *testing.T) {
	tests := []struct {
		json string
		want Node
	}{
		{
			json: `
{
  "spanId": "0f283dfd620daa10",
  "traceId": "473b5dc91e544902871080a25554e963",
  "parent": "deleteDir",
  "all-info": "span(name: deleteDir, spanId: 0f283dfd620daa10, parentSpanId: 58a6a7b62dc5eb76, traceId: 473b5dc91e544902871080a25554e963, attr: ci.pipeline.run.user:SYSTEM;harness-others:;jenkins.pipeline.step.id:92;jenkins.pipeline.step.name:Recursively delete the current directory from the workspace;jenkins.pipeline.step.plugin.name:workflow-basic-steps;jenkins.pipeline.step.plugin.version:1058.vcb_fc1e3a_21a_9;jenkins.pipeline.step.type:deleteDir;)",
  "name": "deleteDir #1",
  "attributesMap": {
    "harness-others": "",
    "jenkins.pipeline.step.name": "Recursively delete the current directory from the workspace",
    "ci.pipeline.run.user": "SYSTEM",
    "jenkins.pipeline.step.id": "1",
    "jenkins.pipeline.step.type": "deleteDir",
    "jenkins.pipeline.step.plugin.name": "workflow-basic-steps",
    "jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9"
  },
  "type": "Run Phase Span",
  "parentSpanId": "58a6a7b62dc5eb76",
  "spanName": "deleteDir"
}
			`,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "1",
					"jenkins.pipeline.step.name":           "Recursively delete the current directory from the workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "deleteDir",
					"harness-others":                       "",
				},
				Name:         "deleteDir #1",
				Parent:       "deleteDir",
				ParentSpanId: "58a6a7b62dc5eb76",
				SpanId:       "0f283dfd620daa10",
				SpanName:     "deleteDir",
				TraceId:      "473b5dc91e544902871080a25554e963",
				Type:         "Run Phase Span",
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

package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertJiraBuild(t *testing.T) {
	// Get the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Update file path according to the location of fileOps_test.json
	filePath := filepath.Join(workingDir, "../convertTestFiles/jira/jira_build/jirabuild_snippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON into a Node struct
	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json Node
		want Node
	}{
		{
			json: node1,
			want: Node{
				SpanId:  "dba869e5a94c6aa2",
				TraceId: "9df67eabac7871cc66120d0402b63706",
				Parent:  "jira-plugin-pipeline",
				Name:    "jira-plugin-pipeline #3",
				AttributesMap: map[string]string{
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@45b0b6d0": "org.jenkinsci.plugins.workflow.actions.TimingAction@45b0b6d0",
					"harness-others":             "",
					"jenkins.pipeline.step.name": "Atlassian Jira Software Cloud Jenkins Integration (Build)",
					"ci.pipeline.run.user":       "SYSTEM",
					"jenkins.pipeline.step.id":   "8",
					"jenkins.pipeline.step.type": "jiraSendBuildInfo",
					"harness-attribute":          "{\r\n  \"branch\" : \"main\",\r\n  \"site\" : \"fossops.atlassian.net\"\r\n}",
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@22b5e87c": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@22b5e87c",
					"jenkins.pipeline.step.plugin.name":    "atlassian-jira-software-cloud",
					"jenkins.pipeline.step.plugin.version": "2.0.15",
				},
				ParentSpanId: "0191118758698812",
				SpanName:     "jiraSendBuildInfo",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"site":   "fossops.atlassian.net",
					"branch": "main",
				},
			},
		},
	}

	for i, test := range tests {
		got := test.json
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestConvertJiraDeployment(t *testing.T) {
	// Get the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Update file path according to the location of fileOps_test.json
	filePath := filepath.Join(workingDir, "../convertTestFiles/jira/jira_deploy/jiradeploy_snippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON into a Node struct
	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json Node
		want Node
	}{
		{
			json: node1,
			want: Node{
				SpanId:  "d3f4788be1429db2",
				TraceId: "9df67eabac7871cc66120d0402b63706",
				Parent:  "jira-plugin-pipeline",
				Name:    "jira-plugin-pipeline #3",
				AttributesMap: map[string]string{
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@43706be5": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@43706be5",
					"harness-others": "",
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@7e77f4ff": "org.jenkinsci.plugins.workflow.actions.TimingAction@7e77f4ff",
					"jenkins.pipeline.step.name":           "Atlassian Jira Software Cloud Jenkins Integration (Deployment)",
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "20",
					"jenkins.pipeline.step.type":           "jiraSendDeploymentInfo",
					"harness-attribute":                    "{\r\n  \"issueKeys\" : [ \"SCRUM-1\", \"SCRUM-2\", \"SCRUM-3\" ],\r\n  \"environmentId\" : \"prod-env-1\",\r\n  \"state\" : \"successful\",\r\n  \"environmentName\" : \"production\",\r\n  \"environmentType\" : \"production\"\r\n}",
					"jenkins.pipeline.step.plugin.name":    "atlassian-jira-software-cloud",
					"jenkins.pipeline.step.plugin.version": "2.0.15",
				},
				ParentSpanId: "ff029e3e25d645ef",
				SpanName:     "jiraSendDeploymentInfo",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"issueKeys": []any{
						"SCRUM-1",
						"SCRUM-2",
						"SCRUM-3",
					},
					"environmentId":   "prod-env-1",
					"environmentName": "production",
					"environmentType": "production",
					"state":           "successful",
				},
			},
		},
	}

	for i, test := range tests {
		got := test.json
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

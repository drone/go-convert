package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

type jirarunner struct {
	name      string
	inputNode Node
	inputVars map[string]string
	wantStep  *harness.Step
}

// Helper function to prepare test cases from JSON files
func jiraprepare(t *testing.T, filename string, step *harness.Step) jirarunner {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/jira", filename+".json"))
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var inputNode Node

	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	return jirarunner{
		name:      filename,
		inputNode: inputNode,
		wantStep:  step,
	}
}

func TestConvertJiraBuildInfo(t *testing.T) {
	// Define test cases for ConvertJiraBuildInfo
	var tests []jirarunner

	tests = append(tests, jiraprepare(t, "jira_build/jirabuild_snippet", &harness.Step{
		Id:   "jiraSendBuildInfodba869",
		Name: "jiraSendBuildInfo",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/jira",
			With: map[string]interface{}{
				"connect_key": "<+secrets.getValue(\"JIRA_CONNECT_KEY\")>",
				"project":     "$JIRA_PROJECT",
				"branch":      "main",
				"instance":    "fossops.atlassian.net",
			},
		},
	}))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertJiraBuildInfo(tt.inputNode, tt.inputVars)
			if diff := cmp.Diff(got, tt.wantStep); diff != "" {
				t.Errorf("ConvertJiraBuildInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertJiraDeploymentInfo(t *testing.T) {
	// Define test cases for ConvertJiraDeploymentInfo
	var tests []jirarunner

	tests = append(tests, jiraprepare(t, "jira_deploy/jiradeploy_snippet", &harness.Step{
		Id:   "jiraSendDeploymentInfod3f478",
		Name: "jiraSendDeploymentInfo",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/jira",
			With: map[string]interface{}{
				"connect_key":      "<+secrets.getValue(\"JIRA_CONNECT_KEY\")>",
				"project":          "$JIRA_PROJECT",
				"instance":         "$JIRA_SITE_ID",
				"environment_id":   "prod-env-1",
				"environment_type": "production",
				"environment_name": "production",
				"state":            "successful",
				"issuekeys":        []string{"SCRUM-1", "SCRUM-2", "SCRUM-3"},
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertJiraDeploymentInfo(tt.inputNode, tt.inputVars)
			if diff := cmp.Diff(got, tt.wantStep); diff != "" {
				t.Errorf("ConvertJiraDeploymentInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

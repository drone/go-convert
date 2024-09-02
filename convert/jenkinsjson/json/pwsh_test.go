package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertPwsh(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Path to the test JSON file for pwsh
	filePath := filepath.Join(workingDir, "../convertTestFiles/pwsh/pwshSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON data into a Node object
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
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "7",
					"jenkins.pipeline.step.name":           "PowerShell Core Script",
					"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
					"jenkins.pipeline.step.plugin.version": "1353.v1891a_b_01da_18",
					"jenkins.pipeline.step.type":           "pwsh",
					"harness-attribute":                    "{\n  \"script\" : \"Write-Output \\\"Hello from powershell\\\"\"\n}",
					"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
				},
				Name:         "pwsh-test #4",
				Parent:       "pwsh-test",
				ParentSpanId: "ae814e400d403a68",
				SpanId:       "48add07c4c323c37",
				SpanName:     "pwsh",
				TraceId:      "c7ae9b130161e9586c66041e8022c888",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"script": "Write-Output \"Hello from powershell\""},
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
package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertAnchore(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/anchore/anchoreSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

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
					"jenkins.pipeline.step.id":             "9",
					"jenkins.pipeline.step.name":           "Anchore Container Image Scanner",
					"jenkins.pipeline.step.plugin.name":    "anchore-container-scanner",
					"jenkins.pipeline.step.plugin.version": "3.1.0",
					"jenkins.pipeline.step.type":           "anchore",
					"harness-attribute": `{
  "delegate" : {
    "symbol" : "anchore",
    "klass" : null,
    "arguments" : {
      "bailOnFail" : "true",
      "forceAnalyze" : "false",
      "name" : "anchore_images",
      "policyBundleId" : "77fff4af-3bfb-421a-87ab-3ee4dd520b76"
    },
    "model" : null
  }
}`,
					"harness-others": "-delegate-field org.jenkinsci.plugins.workflow.steps.CoreStep delegate-org.jenkinsci.plugins.workflow.steps.CoreStep.delegate-interface jenkins.tasks.SimpleBuildStep",
				},
				Name:         "anchoretest #1",
				Parent:       "anchoretest",
				ParentSpanId: "b55ad86a016ca66c",
				SpanId:       "31ac6cc6db919edc",
				SpanName:     "anchore",
				TraceId:      "9ce8e6e532bea908b06be25432621f33",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"delegate": map[string]any{
						"symbol": "anchore",
						"klass":  nil,
						"arguments": map[string]any{
							"forceAnalyze":   "false",
							"name":           "anchore_images",
							"bailOnFail":     "false",
							"policyBundleId": "77fff4af-3bfb-421a-87ab-3ee4dd520b76",
						},
						"model": nil,
					},
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

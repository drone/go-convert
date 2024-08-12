package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertSynopsysDetect(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/synopsysDetect/synopsysDetectSnippet.json")
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
					"jenkins.pipeline.step.id":             "7",
					"jenkins.pipeline.step.name":           "Synopsys Detect",
					"jenkins.pipeline.step.plugin.name":    "blackduck-detect",
					"jenkins.pipeline.step.plugin.version": "9.0.0",
					"jenkins.pipeline.step.type":           "synopsys_detect",
					"harness-attribute":                    "{\n  \"detectProperties\" : \"--detect.project.name=jenkinstest --detect.project.version.name=v1.0\"\n}",
					"harness-others":                       "-DISPLAY_NAME-staticField com.synopsys.integration.jenkins.detect.extensions.pipeline.DetectPipelineStep DISPLAY_NAME-com.synopsys.integration.jenkins.detect.extensions.pipeline.DetectPipelineStep.DISPLAY_NAME-class java.lang.String-PIPELINE_NAME-staticField com.synopsys.integration.jenkins.detect.extensions.pipeline.DetectPipelineStep PIPELINE_NAME-com.synopsys.integration.jenkins.detect.extensions.pipeline.DetectPipelineStep.PIPELINE_NAME-class java.lang.String",
				},
				Name:         "synopsysDetect #1",
				Parent:       "synopsysDetect",
				ParentSpanId: "65d2aa46fdde017d",
				SpanId:       "de0464382a8230c3",
				SpanName:     "synopsys_detect",
				TraceId:      "bc02fa7897ebd5d4f79323768403caa2",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"detectProperties": string("--detect.project.name=jenkinstest --detect.project.version.name=v1.0")},
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

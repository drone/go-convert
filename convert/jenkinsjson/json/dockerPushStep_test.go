package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertDockerPushStep(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/dockerPushStep/dockerPushStepSnippet.json")
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
				SpanId:       "8cfe7ca204c15164",
				TraceId:      "bd776003e36aba7758e42f247d45df49",
				Parent:       "InternalTesting/w-state",
				Name:         "InternalTesting » w-state #3",
				Type:         "Run Phase Span",
				ParentSpanId: "b6b093c1cc0a406a",
				SpanName:     "#3",
				AttributesMap: map[string]string{
					"harness-others":                       "",
					"jenkins.pipeline.step.name":           "Artifactory docker push",
					"jenkins.pipeline.step.id":             "53",
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.type":           "dockerPushStep",
					"jenkins.pipeline.step.plugin.name":    "artifactory",
					"jenkins.pipeline.step.plugin.version": "4.0.7",
					"harness-attribute": `{
  "image" : "docker-snapshot.cbp.dhs.gov/cbp/cloud/warren-testing-konvoy/w-state:DEV.3.20240724",
  "server" : "UNSERIALIZABLE",
  "targetRepo" : "docker-snapshot",
  "javaArgs" : null,
  "buildInfo" : "UNSERIALIZABLE",
  "host" : null,
  "properties" : "UNSERIALIZABLE"
}`,
					"harness-attribute-extra-pip: com.cloudbees.workflow.rest.endpoints.FlowNodeAPI@5ff7bbc5":             "com.cloudbees.workflow.rest.endpoints.FlowNodeAPI@5ff7bbc5",
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@41b67bf0": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@41b67bf0",
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@448496ae":           "org.jenkinsci.plugins.workflow.actions.TimingAction@448496ae",
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

package jenkinsjson

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	jenkinsjson "github.com/jamie-harness/go-convert/convert/jenkinsjson/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertMaven(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../jenkinsjson/convertTestFiles/maven/mavenSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node1 jenkinsjson.Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json jenkinsjson.Node
		want jenkinsjson.Node
	}{
		{
			json: node1,
			want: jenkinsjson.Node{
				SpanId:  "615aa9215d575030",
				TraceId: "473b5dc91e544902871080a25554e963",
				Parent:  "CombinedPipeline",
				Children: []jenkinsjson.Node{
					{
						SpanId:  "f73b459b535f48e5",
						TraceId: "473b5dc91e544902871080a25554e963",
						Parent:  "CombinedPipeline",
						Name:    "CombinedPipeline #9",
						AttributesMap: map[string]string{
							"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
							"jenkins.pipeline.step.name":           "Shell Script",
							"ci.pipeline.run.user":                 "SYSTEM",
							"jenkins.pipeline.step.id":             "32",
							"jenkins.pipeline.step.type":           "sh",
							"harness-attribute":                    "{\n  \"script\" : \"mvn clean install                        package\"\n}",
							"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
							"jenkins.pipeline.step.plugin.version": "1336.v768003e07199",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "615aa9215d575030",
						ParameterMap: map[string]interface{}{"script": "mvn clean install                        package"},
						SpanName:     "sh",
					},
					{
						SpanId:  "42c1f9cc515bf008",
						TraceId: "473b5dc91e544902871080a25554e963",
						Parent:  "CombinedPipeline",
						Name:    "CombinedPipeline #9",
						AttributesMap: map[string]string{
							"harness-others":             "",
							"jenkins.pipeline.step.type": "withMaven",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "615aa9215d575030",
						SpanName:     "Stage: null",
					},
				},
				Name: "CombinedPipeline #9",
				AttributesMap: map[string]string{
					"harness-others":             "",
					"jenkins.pipeline.step.type": "withMaven",
				},
				Type:         "Run Phase Span",
				ParentSpanId: "ee2544e3979cb300",
				SpanName:     "Stage: null",
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

func TestConvertGradle(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../jenkinsjson/convertTestFiles/gradle/graddleSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node1 jenkinsjson.Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json jenkinsjson.Node
		want jenkinsjson.Node
	}{
		{
			json: node1,
			want: jenkinsjson.Node{
				SpanId:  "e4af455880389160",
				TraceId: "473b5dc91e544902871080a25554e963",
				Parent:  "CombinedPipeline",
				Children: []jenkinsjson.Node{
					{
						SpanId:  "0c6722682d8efbae",
						TraceId: "473b5dc91e544902871080a25554e963",
						Parent:  "CombinedPipeline",
						Name:    "CombinedPipeline #9",
						AttributesMap: map[string]string{
							"harness-others":             "",
							"jenkins.pipeline.step.type": "withGradle",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "e4af455880389160",
						SpanName:     "Stage: null",
					},
					{
						SpanId:  "7dc539f24df8acc2",
						TraceId: "473b5dc91e544902871080a25554e963",
						Parent:  "CombinedPipeline",
						Name:    "CombinedPipeline #9",
						AttributesMap: map[string]string{
							"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
							"jenkins.pipeline.step.name":           "Shell Script",
							"ci.pipeline.run.user":                 "SYSTEM",
							"jenkins.pipeline.step.id":             "42",
							"jenkins.pipeline.step.type":           "sh",
							"harness-attribute":                    "{\n  \"script\" : \"gradle build\"\n}",
							"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
							"jenkins.pipeline.step.plugin.version": "1336.v768003e07199",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "e4af455880389160",
						ParameterMap: map[string]interface{}{"script": "gradle build"},
						SpanName:     "sh",
					},
				},
				Name: "CombinedPipeline #9",
				AttributesMap: map[string]string{
					"harness-others":             "",
					"jenkins.pipeline.step.type": "withGradle",
				},
				Type:         "Run Phase Span",
				ParentSpanId: "e67f7aa85a834c2e",
				SpanName:     "Stage: null",
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

func TestConvertAnt(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../jenkinsjson/convertTestFiles/ant/antSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node1 jenkinsjson.Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json jenkinsjson.Node
		want jenkinsjson.Node
	}{
		{
			json: node1,
			want: jenkinsjson.Node{
				SpanId:  "70802fd5c6f6d623",
				TraceId: "473b5dc91e544902871080a25554e963",
				Parent:  "CombinedPipeline",
				Children: []jenkinsjson.Node{
					{
						SpanId:  "e3d404237356f65d",
						TraceId: "473b5dc91e544902871080a25554e963",
						Parent:  "CombinedPipeline",
						Name:    "CombinedPipeline #9",
						AttributesMap: map[string]string{
							"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
							"jenkins.pipeline.step.name":           "Shell Script",
							"ci.pipeline.run.user":                 "SYSTEM",
							"jenkins.pipeline.step.id":             "37",
							"jenkins.pipeline.step.type":           "sh",
							"harness-attribute":                    "{\n  \"script\" : \"ant build\"\n}",
							"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
							"jenkins.pipeline.step.plugin.version": "1336.v768003e07199",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "70802fd5c6f6d623",
						ParameterMap: map[string]interface{}{"script": "ant build"},
						SpanName:     "sh",
					},
					{
						SpanId:  "e696a9520bd0250e",
						TraceId: "473b5dc91e544902871080a25554e963",
						Parent:  "CombinedPipeline",
						Name:    "CombinedPipeline #9",
						AttributesMap: map[string]string{
							"harness-others":             "",
							"jenkins.pipeline.step.type": "wrap",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "70802fd5c6f6d623",
						SpanName:     "Stage: null",
					},
				},
				Name: "CombinedPipeline #9",
				AttributesMap: map[string]string{
					"harness-others":             "",
					"jenkins.pipeline.step.type": "wrap",
				},
				Type:         "Run Phase Span",
				ParentSpanId: "e818845986d7bccd",
				SpanName:     "Stage: null",
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

func TestConvertSonarqube(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../jenkinsjson/convertTestFiles/sonarqube/sonarqubeSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node1 jenkinsjson.Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json jenkinsjson.Node
		want jenkinsjson.Node
	}{
		{
			json: node1,
			want: jenkinsjson.Node{
				SpanId:  "c9943fda45b8312f",
				TraceId: "f59bb3660fe13cc634fb357b59c0f1d7",
				Parent:  "sonarcube",
				Children: []jenkinsjson.Node{
					{
						SpanId:  "952bc3bfafb5bdcf",
						TraceId: "f59bb3660fe13cc634fb357b59c0f1d7",
						Parent:  "sonarcube",
						Name:    "sonarcube #9",
						AttributesMap: map[string]string{
							"harness-others":             "",
							"jenkins.pipeline.step.type": "wrap",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "c9943fda45b8312f",
						SpanName:     "Stage: null",
					},
					{
						SpanId:  "65bb92aeb82fcdc5",
						TraceId: "f59bb3660fe13cc634fb357b59c0f1d7",
						Parent:  "sonarcube",
						Name:    "sonarcube #9",
						AttributesMap: map[string]string{
							"harness-others":                       "-WATCHING_RECURRENCE_PERIOD-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep WATCHING_RECURRENCE_PERIOD-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.WATCHING_RECURRENCE_PERIOD-long-USE_WATCHING-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep USE_WATCHING-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.USE_WATCHING-boolean-REMOTE_TIMEOUT-staticField org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep REMOTE_TIMEOUT-org.jenkinsci.plugins.workflow.steps.durable_task.DurableTaskStep.REMOTE_TIMEOUT-long",
							"jenkins.pipeline.step.name":           "Shell Script",
							"ci.pipeline.run.user":                 "SYSTEM",
							"jenkins.pipeline.step.id":             "19",
							"jenkins.pipeline.step.type":           "sh",
							"harness-attribute":                    "{\n  \"script\" : \"/Users/rakshitagarwal/.jenkins/tools/hudson.plugins.sonar.SonarRunnerInstallation/sonar_cube/bin/sonar-scanner                             -Dsonar.projectKey=jenkins-SonarCube                             -Dsonar.projectName=jenkins-SonarCube                             -Dsonar.sources=.\"\n}",
							"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
							"jenkins.pipeline.step.plugin.version": "1336.v768003e07199",
						},
						Type:         "Run Phase Span",
						ParentSpanId: "c9943fda45b8312f",
						ParameterMap: map[string]interface{}{"script": "/Users/rakshitagarwal/.jenkins/tools/hudson.plugins.sonar.SonarRunnerInstallation/sonar_cube/bin/sonar-scanner                             -Dsonar.projectKey=jenkins-SonarCube                             -Dsonar.projectName=jenkins-SonarCube                             -Dsonar.sources=."},
						SpanName:     "sh",
					},
				},
				Name: "sonarcube #9",
				AttributesMap: map[string]string{
					"harness-others":             "",
					"jenkins.pipeline.step.type": "wrap",
				},
				Type:         "Run Phase Span",
				ParentSpanId: "32035723074a4821",
				SpanName:     "Stage: null",
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

func TestNewConverter_Defaults(t *testing.T) {
	converter := New()
	if converter.kubeNamespace != "default" {
		t.Errorf("Expected default namespace to be 'default', got '%s'", converter.kubeNamespace)
	}
	if converter.kubeEnabled {
		t.Error("Expected kubeEnabled to be false by default")
	}
}

//func TestConvert_ValidJenkinsJSON(t *testing.T) {
//	converter := New()
//	inputJSON := `{"some": "valid", "jenkins": "json"}`
//	expectedYAML := `version: 1
//kind: pipeline
//spec:
//  stages: [...]` // Replace with expected output
//
//	output, err := converter.ConvertString(inputJSON)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//	if string(output) != expectedYAML {
//		t.Errorf("Expected output to be '%s', got '%s'", expectedYAML, string(output))
//	}
//}

func TestConvert_EmptyJenkinsJSON(t *testing.T) {
	converter := New()
	inputJSON := `{}`
	expectedYAML := `version: 1
kind: pipeline
type: ""
name: ""
spec:
  stages:
  - desc: ""
    id: build
	name: build
	strategy: null
	delegate: []
	status: null
	type: ci
	when: null
	failure: null
	inputs: {}
	spec:
      cache: null
      clone: null
      platform: null
      runtime: null
      steps: []
      envs: {}
      volumes: []
  inputs: {}
  options: null
`

	output, err := converter.ConvertString(inputJSON)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	normalize := func(s string) string {
		return strings.Join(strings.Fields(s), " ")
	}

	normalizedExpectedYAML := normalize(expectedYAML)
	normalizedOutput := normalize(string(output))

	if normalizedOutput != normalizedExpectedYAML {
		t.Errorf("Expected output to be '%s', got '%s'", normalizedExpectedYAML, normalizedOutput)
	}
}

func TestRecursiveParseJsonToStages(t *testing.T) {
	mockNode := jenkinsjson.Node{
		Name:   string("MockName"),
		Parent: string("MockParent"),
	}
	dst := &harness.Pipeline{}
	processedTools := &ProcessedTools{}
	variables := map[string]string{}

	recursiveParseJsonToStages(&mockNode, dst, processedTools, variables)

	if len(dst.Stages) == 0 {
		t.Error("Expected stages to be populated, got none")
	}
}

func TestHandleTool_ValidToolNodeMaven(t *testing.T) {
	toolNode := jenkinsjson.Node{
		AttributesMap: map[string]string{
			"jenkins.pipeline.step.type": "tool",
			"harness-attribute":          `{"type": "$MavenInstallation"}`,
		},
	}
	processedTools := &ProcessedTools{}
	handleTool(toolNode, processedTools)

	if !processedTools.MavenPresent {
		t.Error("Expected Maven to be marked as present")
	}
}

func TestHandleTool_ValidToolNodeGradle(t *testing.T) {
	toolNode := jenkinsjson.Node{
		AttributesMap: map[string]string{
			"jenkins.pipeline.step.type": "tool",
			"harness-attribute":          `{"type": "$GradleInstallation"}`,
		},
	}
	processedTools := &ProcessedTools{}
	handleTool(toolNode, processedTools)

	if !processedTools.GradlePresent {
		t.Error("Expected Gradle to be marked as present")
	}
}

func TestHandleTool_ValidToolNodeAnt(t *testing.T) {
	toolNode := jenkinsjson.Node{
		AttributesMap: map[string]string{
			"jenkins.pipeline.step.type": "tool",
			"harness-attribute":          `{"type": "$AntInstallation"}`,
		},
	}
	processedTools := &ProcessedTools{}
	handleTool(toolNode, processedTools)

	if !processedTools.AntPresent {
		t.Error("Expected Ant to be marked as present")
	}
}

func TestExtractEnvironmentVariables(t *testing.T) {
	node := jenkinsjson.Node{
		ParameterMap: map[string]interface{}{
			"overrides": []interface{}{"KEY=value"},
		},
	}

	envVars := ExtractEnvironmentVariables(node)
	if len(envVars) == 0 || envVars["KEY"] != "value" {
		t.Error("Expected environment variable KEY to have value 'value'")
	}
}

func TestMergeMaps(t *testing.T) {
	a := map[string]string{"key1": "value1"}
	b := map[string]string{"key2": "value2", "key1": "value2"}

	merged := mergeMaps(a, b)
	if len(merged) != 2 || merged["key1"] != "value2" || merged["key2"] != "value2" {
		t.Error("Expected merged map to contain key1=value2 and key2=value2")
	}
}

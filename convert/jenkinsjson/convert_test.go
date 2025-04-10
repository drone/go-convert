package jenkinsjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	jenkinsjson "github.com/drone/go-convert/convert/jenkinsjson/json"
	harness "github.com/drone/spec/dist/go"

	"github.com/google/go-cmp/cmp"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "test-build-scan-push",
			input: "./convertTestFiles/convert/test-build-scan-push.json",
			want:  "./convertTestFiles/convert/test-build-scan-push.yaml",
		},
		{
			name:  "build-and-multiple-deploy",
			input: "./convertTestFiles/convert/build-and-multiple-deploy.json",
			want:  "./convertTestFiles/convert/build-and-multiple-deploy.yaml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			converter := Converter{}
			got, err := converter.ConvertFile(tc.input)
			if err != nil {
				t.Error("Failed to convert file", tc.input, err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Error("Failed to read the expected output file", tc.want, err)
			}

			fmt.Println(got)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("TestConvert mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

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
name: default
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

	expected := normalize(expectedYAML)
	got := normalize(string(output))

	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("TestConvert_EmptyJenkinsJSON mismatch (-want +got):\n%s", diff)
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

	// create a converter object and call recursiveParseJsonToStages method
	d := &Converter{}
	d.recursiveParseJsonToStages(&mockNode, dst, processedTools, variables)

	if len(dst.Stages) == 0 {
		t.Error("Expected stages to be populated, got none")
	}
}

func TestHandleTool_ValidToolNodeMaven(t *testing.T) {
	toolNode := jenkinsjson.Node{
		AttributesMap: map[string]string{
			"jenkins.pipeline.step.type": "tool",
		},
		ParameterMap: map[string]interface{}{"type": "$MavenInstallation"},
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
		},
		ParameterMap: map[string]interface{}{"type": "$GradleInstallation"},
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
		},
		ParameterMap: map[string]interface{}{"type": "$AntInstallation"},
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

func TestMergeRunSteps(t *testing.T) {
	createStepExec := func(run, image, shell string, envs map[string]string) *harness.StepExec {
		return &harness.StepExec{
			Run:        run,
			Image:      image,
			Shell:      shell,
			Envs:       envs,
			Connector:  "",
			Args:       []string{},
			Privileged: false,
			Network:    "",
		}
	}

	createStep := func(name, stepType string, spec *harness.StepExec) *harness.Step {
		return &harness.Step{
			Name: name,
			Type: stepType,
			Spec: spec,
		}
	}

	t.Run("Empty steps slice", func(t *testing.T) {
		steps := []StepWithID{}
		mergeRunSteps(&steps)
		if len(steps) != 0 {
			t.Errorf("expected empty slice, got length %d", len(steps))
		}
	})

	t.Run("Single step", func(t *testing.T) {
		exec := createStepExec("echo hello", "", "sh", nil)
		step := createStep("step1", "script", exec)
		steps := []StepWithID{{Step: step}}

		mergeRunSteps(&steps)
		if len(steps) != 1 {
			t.Errorf("expected length 1, got %d", len(steps))
		}
		if steps[0].Step.Name != "step1" {
			t.Errorf("expected name 'step1', got %s", steps[0].Step.Name)
		}
	})

	t.Run("Mergeable steps", func(t *testing.T) {
		steps := []StepWithID{
			{
				Step: createStep("step1", "script", createStepExec("echo hello", "", "sh", nil)),
			},
			{
				Step: createStep("step2", "script", createStepExec("echo world", "", "sh", nil)),
			},
		}

		mergeRunSteps(&steps)
		if len(steps) != 1 {
			t.Errorf("expected length 1, got %d", len(steps))
		}
		if steps[0].Step.Name != "step1_step2" {
			t.Errorf("expected name 'step1_step2', got %s", steps[0].Step.Name)
		}
		execSpec := steps[0].Step.Spec.(*harness.StepExec)
		expectedRun := "echo hello\necho world"
		if execSpec.Run != expectedRun {
			t.Errorf("expected run command '%s', got '%s'", expectedRun, execSpec.Run)
		}
	})

	t.Run("Non-mergeable steps - different types", func(t *testing.T) {
		steps := []StepWithID{
			{
				Step: createStep("step1", "script", createStepExec("echo hello", "", "sh", nil)),
			},
			{
				Step: createStep("step2", "non-script", createStepExec("echo world", "", "sh", nil)),
			},
		}

		mergeRunSteps(&steps)
		if len(steps) != 2 {
			t.Errorf("expected length 2, got %d", len(steps))
		}
		if steps[0].Step.Name != "step1" {
			t.Errorf("expected first step name 'step1', got %s", steps[0].Step.Name)
		}
		if steps[1].Step.Name != "step2" {
			t.Errorf("expected second step name 'step2', got %s", steps[1].Step.Name)
		}
	})
}

func TestCanMergeSteps(t *testing.T) {
	t.Run("Different step types", func(t *testing.T) {
		step1 := &harness.Step{Type: "script"}
		step2 := &harness.Step{Type: "non-script"}
		if canMergeSteps(step1, step2) {
			t.Error("expected steps with different types to not be mergeable")
		}
	})

	t.Run("Invalid spec type", func(t *testing.T) {
		step1 := &harness.Step{Type: "script", Spec: "invalid"}
		step2 := &harness.Step{Type: "script", Spec: &harness.StepExec{}}
		if canMergeSteps(step1, step2) {
			t.Error("expected steps with invalid spec to not be mergeable")
		}
	})

	t.Run("Different images", func(t *testing.T) {
		step1 := &harness.Step{
			Type: "script",
			Spec: &harness.StepExec{Image: "image1"},
		}
		step2 := &harness.Step{
			Type: "script",
			Spec: &harness.StepExec{Image: "image2"},
		}
		if canMergeSteps(step1, step2) {
			t.Error("expected steps with different images to not be mergeable")
		}
	})

	t.Run("Matching steps", func(t *testing.T) {
		exec := &harness.StepExec{
			Image:      "",
			Shell:      "sh",
			Connector:  "",
			Envs:       map[string]string{"KEY": "VALUE"},
			Args:       []string{"arg1"},
			Privileged: false,
			Network:    "",
		}
		step1 := &harness.Step{Type: "script", Spec: exec}
		step2 := &harness.Step{Type: "script", Spec: exec}
		if !canMergeSteps(step1, step2) {
			t.Error("expected identical steps to be mergeable")
		}
	})
}

func TestENVmapsEqual(t *testing.T) {
	t.Run("Empty maps", func(t *testing.T) {
		if !ENVmapsEqual(nil, nil) {
			t.Error("expected nil maps to be equal")
		}
		if !ENVmapsEqual(map[string]string{}, map[string]string{}) {
			t.Error("expected empty maps to be equal")
		}
	})

	t.Run("Different sizes", func(t *testing.T) {
		m1 := map[string]string{"a": "1"}
		m2 := map[string]string{"a": "1", "b": "2"}
		if ENVmapsEqual(m1, m2) {
			t.Error("expected maps of different sizes to not be equal")
		}
	})

	t.Run("Same content", func(t *testing.T) {
		m1 := map[string]string{"a": "1", "b": "2"}
		m2 := map[string]string{"a": "1", "b": "2"}
		if !ENVmapsEqual(m1, m2) {
			t.Error("expected maps with same content to be equal")
		}
	})

	t.Run("Different values", func(t *testing.T) {
		m1 := map[string]string{"a": "1", "b": "2"}
		m2 := map[string]string{"a": "1", "b": "3"}
		if ENVmapsEqual(m1, m2) {
			t.Error("expected maps with different values to not be equal")
		}
	})
}

func TestARGSslicesEqual(t *testing.T) {
	t.Run("Empty slices", func(t *testing.T) {
		if !ARGSslicesEqual(nil, nil) {
			t.Error("expected nil slices to be equal")
		}
		if !ARGSslicesEqual([]string{}, []string{}) {
			t.Error("expected empty slices to be equal")
		}
	})

	t.Run("Different lengths", func(t *testing.T) {
		s1 := []string{"a"}
		s2 := []string{"a", "b"}
		if ARGSslicesEqual(s1, s2) {
			t.Error("expected slices of different lengths to not be equal")
		}
	})

	t.Run("Same content", func(t *testing.T) {
		s1 := []string{"a", "b", "c"}
		s2 := []string{"a", "b", "c"}
		if !ARGSslicesEqual(s1, s2) {
			t.Error("expected slices with same content to be equal")
		}
	})

	t.Run("Different content", func(t *testing.T) {
		s1 := []string{"a", "b", "c"}
		s2 := []string{"a", "b", "d"}
		if ARGSslicesEqual(s1, s2) {
			t.Error("expected slices with different content to not be equal")
		}
	})
}

// TestS3UploadWithGzipDisabled tests the s3Upload functionality with a single entry where "gzipFiles" is set to false.
//
// Expected Behavior:
// - The "gzipFiles" flag is set to false, so only one step is expected:
// - Step 1: The S3 upload step should be created without any archive (gzip) process.
//
// Validations Performed:
// - Ensures only one step is created when gzip is false.
// - Compares individual attributes like bucket, region, and source for correctness.
//
// Example File Structure (s3upload_snippet.json):
// [
//   {
//     "gzipFiles": false,
//     "bucket": "bucket-1",
//     "source": "*.txt",
//     "region": "us-west-1",
//     ...
//   }
// ]
//
// This test ensures that the appropriate S3 upload step is created for a single S3 entry with gzip disabled.

func TestS3UploadWithGzipDisabled(t *testing.T) {
	// Step 1: Read the JSON file
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../jenkinsjson/convertTestFiles/s3publisher/s3upload/s3upload_snippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Step 2: Unmarshal the JSON into the appropriate structure
	var node jenkinsjson.Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// Initialize the list to collect steps
	stepGroupWithId := []StepGroupWithID{}
	stepWithIDList := []StepWithID{}
	processedTools := &ProcessedTools{}
	variables := map[string]string{}
	timeout := "10m"
	dockerImage := "plugin/s3upload"

	// Call the function collectStepsWithID to collect the steps
	collectStepsWithID(node, &stepGroupWithId, &stepWithIDList, processedTools, variables, timeout, dockerImage)

	// Expecting single step for s3uplaod without gzip
	if len(stepWithIDList) != 1 {
		t.Fatalf("Expected 1 steps, but got %d", len(stepWithIDList))
	}
	// Validate the first step (s3Upload without gzip for the single entry file)
	firstStep := stepWithIDList[0].Step
	if firstStep.Spec.(*harness.StepPlugin).With["bucket"] != "bucket-1" {
		t.Errorf("Expected bucket 'bucket-1', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["bucket"])
	}
	if firstStep.Spec.(*harness.StepPlugin).With["region"] != "us-west-1" {
		t.Errorf("Expected region 'us-west-1', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["region"])
	}
	if firstStep.Spec.(*harness.StepPlugin).With["target"] != "<+input>" {
		t.Errorf("Expected target '<+input>', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["target"])
	}
	if firstStep.Spec.(*harness.StepPlugin).With["exclude"] != "2.txt" {
		t.Errorf("Expected target '2.txt', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["exclude"])
	}
}

// TestCollectStepsWithIDS3UploadMultipleEntry tests the s3Upload functionality for multiple entries.
// It simulates two entries in the provided JSON file for S3 uploads.
//
// Expected Behavior: Total 3 steps:
// - Entry 1: "gzipFiles" is set to false, so only a single upload step(0) is expected.
// - Entry 2: "gzipFiles" is set to true, meaning two steps are expected:
//   - Step(1): First step should be the gzip (archive) process.
//   - Step(2): Second step should be the actual S3 upload.
//
// Validations Performed:
// - Ensures correct number of steps are created based on the gzip flag.
// - Compares individual attributes like target, bucket, source, glob, and exclude for correctness.
//
// Example File Structure (s3upload_multiple-entries_snippet.json):
// [
//
//	{
//	  "gzipFiles": false,
//	  "bucket": "bucket-1",
//	  "source": "*.txt",
//	  "region": "us-west-1",
//	  ...
//	},
//	{
//	  "gzipFiles": true,
//	  "bucket": "bucket-2",
//	  "source": "*.txt",
//	  "exclude": "2.txt",
//	  "region": "us-east-1",
//	  ...
//	}
//
// ]
func TestCollectStepsWithIDS3UploadMultipleEntry(t *testing.T) {
	// Step 1: Read the JSON file
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../jenkinsjson/convertTestFiles/s3publisher/s3upload/s3upload_multiple-entries_snippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON into the appropriate structure
	var node jenkinsjson.Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// Initialize the list to collect steps
	stepGroupWithId := []StepGroupWithID{}
	stepWithIDList1 := []StepWithID{}
	processedTools := &ProcessedTools{}
	variables := make(map[string]string)
	timeout := "10m"
	dockerImage := "plugin/s3upload"

	// Call the function collectStepsWithID to collect the steps
	collectStepsWithID(node, &stepGroupWithId, &stepWithIDList1, processedTools, variables, timeout, dockerImage)

	// Expecting 3 steps: step1: for the gzip false case (only s3upload)
	// step2:  gzip true(Converts3Archive) and step3: s3upload .
	if len(stepWithIDList1) != 3 {
		t.Fatalf("Expected 3 steps, but got %d", len(stepWithIDList1))
	}
	// Validate the first step (s3Upload without gzip for the first step)
	firstStep := stepWithIDList1[0].Step
	if firstStep.Spec.(*harness.StepPlugin).With["bucket"] != "bucket-1" {
		t.Errorf("Expected bucket 'bucket-1', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["bucket"])
	}
	if firstStep.Spec.(*harness.StepPlugin).With["region"] != "us-west-1" {
		t.Errorf("Expected region 'us-west-1', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["region"])
	}
	if firstStep.Spec.(*harness.StepPlugin).With["target"] != "<+input>" {
		t.Errorf("Expected target '<+input>', but got '%s'", firstStep.Spec.(*harness.StepPlugin).With["target"])
	}

	// Validate the second  step (Converts3Archive for gzipFiles true)
	secondStep := stepWithIDList1[1].Step
	if secondStep.Spec.(*harness.StepPlugin).With["target"] != "s3Upload.gzip" {
		t.Errorf("Expected target 's3Upload.gzip', but got '%s'", secondStep.Spec.(*harness.StepPlugin).With["target"])
	}
	if secondStep.Spec.(*harness.StepPlugin).With["source"] != "." {
		t.Errorf("Expected source '.', but got '%s'", secondStep.Spec.(*harness.StepPlugin).With["source"])
	}
	if secondStep.Spec.(*harness.StepPlugin).With["glob"] != "*.txt" {
		t.Errorf("Expected glob '*.txt', but got '%s'", secondStep.Spec.(*harness.StepPlugin).With["glob"])
	}
	if secondStep.Spec.(*harness.StepPlugin).With["exclude"] != "2.txt" {
		t.Errorf("Expected exclude '2.txt', but got '%s'", secondStep.Spec.(*harness.StepPlugin).With["exclude"])
	}

	// Validate the third step (s3Upload with gzip)
	thirdStep := stepWithIDList1[2].Step
	if thirdStep.Spec.(*harness.StepPlugin).With["bucket"] != "bucket-2" {
		t.Errorf("Expected bucket 'bucket-2', but got '%s'", thirdStep.Spec.(*harness.StepPlugin).With["bucket"])
	}
	if thirdStep.Spec.(*harness.StepPlugin).With["region"] != "us-west-1" {
		t.Errorf("Expected region 'us-west-1', but got '%s'", thirdStep.Spec.(*harness.StepPlugin).With["region"])
	}
	if thirdStep.Spec.(*harness.StepPlugin).With["target"] != "<+input>" {
		t.Errorf("Expected target '<+input>', but got '%s'", thirdStep.Spec.(*harness.StepPlugin).With["target"])
	}
	if thirdStep.Spec.(*harness.StepPlugin).With["exclude"] != "2.txt" {
		t.Errorf("Expected target '2.txt', but got '%s'", thirdStep.Spec.(*harness.StepPlugin).With["exclude"])
	}
}

// Test for infrastructure configuration specified in the CLI
func TestInfrastructureOptions(t *testing.T) {
	tests := []struct {
		name           string
		infrastructure string
		arch           string
		os             string
		expectedErr    bool
	}{
		{
			name:           "Valid infrastructure",
			infrastructure: "k8s",
			arch:           "amd64",
			os:             "linux",
			expectedErr:    false,
		},
		{
			name:           "Invalid infrastructure",
			infrastructure: "invalid",
			arch:           "amd64",
			os:             "linux",
			expectedErr:    true,
		},
		{
			name:           "Valid arch",
			infrastructure: "k8s",
			arch:           "arm64",
			os:             "linux",
			expectedErr:    false,
		},
		{
			name:           "Invalid arch",
			infrastructure: "k8s",
			arch:           "invalid",
			os:             "linux",
			expectedErr:    true,
		},
		{
			name:           "Valid os",
			infrastructure: "k8s",
			arch:           "amd64",
			os:             "darwin",
			expectedErr:    false,
		},
		{
			name:           "Invalid os",
			infrastructure: "k8s",
			arch:           "amd64",
			os:             "invalid",
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := New(
				WithInfrastructure(tt.infrastructure),
				WithOS(tt.os),
				WithArch(tt.arch),
			)

			err := converter.ValidateInfrastructureOptions()

			if tt.expectedErr && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

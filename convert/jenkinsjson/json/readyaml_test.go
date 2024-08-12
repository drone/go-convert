package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertReadYaml(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/readyaml/readyamlSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read Yaml file: %v", err)
	}

	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode Yaml: %v", err)
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
					"jenkins.pipeline.step.name":           "Read yaml from files in the workspace or text.",
					"jenkins.pipeline.step.plugin.name":    "pipeline-utility-steps",
					"jenkins.pipeline.step.plugin.version": "2.17.0",
					"jenkins.pipeline.step.type":           "readYaml",
					"harness-attribute":                    "{\n  \"file\" : \"abc.yml\"\n}",
					"harness-others":                       "-LIBRARY_DEFAULT_CODE_POINT_LIMIT-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep LIBRARY_DEFAULT_CODE_POINT_LIMIT-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.LIBRARY_DEFAULT_CODE_POINT_LIMIT-int-MAX_CODE_POINT_LIMIT_PROPERTY-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep MAX_CODE_POINT_LIMIT_PROPERTY-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.MAX_CODE_POINT_LIMIT_PROPERTY-class java.lang.String-DEFAULT_CODE_POINT_LIMIT_PROPERTY-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep DEFAULT_CODE_POINT_LIMIT_PROPERTY-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.DEFAULT_CODE_POINT_LIMIT_PROPERTY-class java.lang.String-HARDCODED_CEILING_MAX_ALIASES_FOR_COLLECTIONS-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep HARDCODED_CEILING_MAX_ALIASES_FOR_COLLECTIONS-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.HARDCODED_CEILING_MAX_ALIASES_FOR_COLLECTIONS-int-LIBRARY_DEFAULT_MAX_ALIASES_FOR_COLLECTIONS-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep LIBRARY_DEFAULT_MAX_ALIASES_FOR_COLLECTIONS-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.LIBRARY_DEFAULT_MAX_ALIASES_FOR_COLLECTIONS-int-MAX_MAX_ALIASES_PROPERTY-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep MAX_MAX_ALIASES_PROPERTY-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.MAX_MAX_ALIASES_PROPERTY-class java.lang.String-DEFAULT_MAX_ALIASES_PROPERTY-staticField org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep DEFAULT_MAX_ALIASES_PROPERTY-org.jenkinsci.plugins.pipeline.utility.steps.conf.ReadYamlStep.DEFAULT_MAX_ALIASES_PROPERTY-class java.lang.String",
				},
				Name:         "ag-readJSON #21",
				Parent:       "ag-readJSON",
				ParentSpanId: "77f5f4e014f800c7",
				SpanId:       "89f1da1fb2dbe5fd",
				SpanName:     "readYaml",
				TraceId:      "596ff9663df204a94745f24c97828730",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": "abc.yml"},
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

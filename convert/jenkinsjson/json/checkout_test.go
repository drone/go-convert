package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertCheckout(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/checkout/checkoutSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var rawNode map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	node1 := Node{}
	if err := mapToNode(rawNode, &node1); err != nil {
		t.Fatalf("failed to convert raw node to Node: %v", err)
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
					"jenkins.pipeline.step.id":             "29",
					"jenkins.pipeline.step.name":           "Check out from version control",
					"rpc.service":                          "git",
					"http.url":                             "https://github.com/Anshika2203/CombinedAllBuilds.git",
					"git.branch":                           "*/main",
					"rpc.system":                           "https",
					"git.repository":                       "Anshika2203/CombinedAllBuilds",
					"net.peer.name":                        "github.com",
					"peer.service":                         "github.com",
					"jenkins.pipeline.step.plugin.name":    "workflow-scm-step",
					"jenkins.pipeline.step.plugin.version": "427.v4ca_6512e7df1",
					"http.method":                          "POST",
					"git.clone.depth":                      "0",
					"jenkins.pipeline.step.type":           "checkout",
					"rpc.method":                           "checkout",
					"git.clone.shallow":                    "false",
					"harness-attribute":                    "{\n  \"scm\" : {\n    \"$class\" : \"GitSCM\",\n    \"branches\" : [ {\n      \"name\" : \"*/main\"\n    } ],\n    \"doGenerateSubmoduleConfigurations\" : false,\n    \"extensions\" : [ ],\n    \"submoduleCfg\" : [ ],\n    \"userRemoteConfigs\" : [ {\n      \"url\" : \"https://github.com/Anshika2203/CombinedAllBuilds.git\"\n    } ]\n  }\n}",
					"harness-others":                       "",
				},
				Name:         "CombinedPipeline #9",
				Parent:       "CombinedPipeline",
				ParentSpanId: "da753c3c600edc4d",
				SpanId:       "dbdf4b330ba1e0d7",
				SpanName:     "checkout: github.com/Anshika2203/CombinedAllBuilds",
				TraceId:      "473b5dc91e544902871080a25554e963",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"scm": map[string]any{
					"$class":                            "GitSCM",
					"extensions":                        []any{},
					"submoduleCfg":                      []any{},
					"userRemoteConfigs":                 []any{map[string]any{"url": string("https://github.com/Anshika2203/CombinedAllBuilds.git")}},
					"doGenerateSubmoduleConfigurations": false,
					"branches":                          []any{map[string]any{"name": string("*/main")}},
				}},
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

func mapToNode(raw map[string]interface{}, node *Node) error {
	node.AttributesMap = make(map[string]string)
	for k, v := range raw["attributesMap"].(map[string]interface{}) {
		node.AttributesMap[k] = fmt.Sprintf("%v", v)
	}

	node.Name = raw["name"].(string)
	node.Parent = raw["parent"].(string)
	node.ParentSpanId = raw["parentSpanId"].(string)
	node.SpanId = raw["spanId"].(string)
	node.SpanName = raw["spanName"].(string)
	node.TraceId = raw["traceId"].(string)
	node.Type = raw["type"].(string)
	node.ParameterMap = raw["parameterMap"].(map[string]interface{})

	return nil
}

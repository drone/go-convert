package json

import (
	"encoding/json"
	"fmt"
	harness "github.com/drone/spec/dist/go"
	"log"
)

const (
	ArtifactUploaderImage = "plugins/nexus-publish:latest"
)

var JenkinsToNexusArtifactUploaderParamMapperList = []JenkinsToDroneParamMapper{
	{"nexusVersion", "nexus_version", StringType, nil},
	{"nexusUrl", "server_url", StringType, nil},
	{"protocol", "protocol", StringType, nil},
	{"groupId", "group_id", StringType, nil},
	{"repository", "repository", StringType, nil},
}

func ConvertNexusArtifactUploader(node Node, variables map[string]string) *harness.Step {

	step := ConvertToStepWithProperties(&node, variables, JenkinsToNexusArtifactUploaderParamMapperList,
		ArtifactUploaderImage)

	tmpStepPlugin := step.Spec.(*harness.StepPlugin)
	tmpStepPlugin.With["username"] = "<+input>"
	tmpStepPlugin.With["password"] = "<+input>"
	artifactsListStr, err := GetArtifactsListString(node)
	if err != nil {
		log.Println("Error getting artifacts list string:", err)
		return nil
	}

	tmpStepPlugin.With["artifacts"] = artifactsListStr
	return step
}

func GetArtifactsListString(node Node) (string, error) {

	attr, ok := node.AttributesMap[HarnessAttribute]
	if !ok {
		log.Printf("harness-attribute missing for spanName %s", node.SpanName)
		return "", fmt.Errorf("harness-attribute missing for spanName %s", node.SpanName)
	}

	attrMap, err := ToMapFromJsonString[map[string]interface{}](attr)
	if err != nil {
		log.Printf("Failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		return "", fmt.Errorf("Failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
	}

	artifactsInfoMapList, ok := attrMap["artifacts"].([]interface{})
	if !ok {
		log.Printf("Error artifacts attribute missing in node %s", node.SpanName)
		return "", fmt.Errorf("Error artifacts attribute missing in node %s", node.SpanName)
	}

	var combinedArtifacts []map[string]interface{}

	for _, artifactInfoMap := range artifactsInfoMapList {
		v, ok := artifactInfoMap.(map[string]interface{})
		if !ok {
			log.Printf("Error invalid artifact info for %s ", node.SpanName)
			return "", fmt.Errorf("Error invalid artifacts list info for %s ", node.SpanName)
		}
		combinedArtifacts = append(combinedArtifacts, v)
	}

	jsonData, err := json.Marshal(combinedArtifacts)
	if err != nil {
		log.Printf("Error converting artifacts list to JSON string: %s ", err)
		return "", fmt.Errorf("Error converting artifacts list to JSON string: %s ", err)
	}

	return string(jsonData), nil
}

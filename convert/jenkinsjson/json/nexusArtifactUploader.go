package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"log"
)

const (
	ArtifactUploaderImage = "harnesscommunity/drone-nexus-publish"
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
	tmpStepPlugin.With["username"] = <+input>
	tmpStepPlugin.With["password"] = <+input>
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
		return "", nil
	}

	attrMap, err := ToMapFromJsonString[map[string]interface{}](attr)
	if err != nil {
		log.Printf("Failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		return "", nil
	}

	artifactsInfoMapList, ok := attrMap["artifacts"].([]interface{})
	if !ok {
		log.Println("Error converting artifacts to map")
		return "", nil
	}

	var combinedArtifacts []map[string]interface{}

	for _, artifactInfoMap := range artifactsInfoMapList {
		v, ok := artifactInfoMap.(map[string]interface{})
		if !ok {
			log.Println("Error converting artifact info to map")
			return "", nil
		}
		combinedArtifacts = append(combinedArtifacts, v)
	}

	jsonData, err := json.Marshal(combinedArtifacts)
	if err != nil {
		log.Println("Error converting artifacts list to JSON string:", err)
		return "", nil
	}

	return string(jsonData), nil
}

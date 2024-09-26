package json

import (
	"encoding/json"
	"fmt"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertJiraBuildInfo(node Node, variables map[string]string) *harness.Step {
	var branch, instance string

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if br, ok := attrMap["branch"].(string); ok {
				branch = br
			}
			if site, ok := attrMap["site"].(string); ok {
				instance = site
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	withProperties := make(map[string]interface{})
	withProperties["connect_key"] = "<+secrets.getValue(\"JIRA_CONNECT_KEY\")>"
	withProperties["project"] = "$JIRA_PROJECT"
	if branch != "" {
		withProperties["branch"] = branch
	}
	if instance != "" {
		withProperties["instance"] = instance
	}
	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/jira",
			With:  withProperties,
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}
	/*yamlData, err := yaml.Marshal(step)
	if err != nil {
		log.Fatalf("error marshaling to yaml: %v", err)
	}
	log.Printf("YAML Build Output: %s", string(yamlData))*/

	return step
}

func ConvertJiraDeploymentInfo(node Node, variables map[string]string) *harness.Step {
	var envId, envType, envName, site, state string
	var issueKeys []string

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if keys, ok := attrMap["issueKeys"].([]interface{}); ok {
				// Convert each interface{} element into a string and append to the string slice
				for _, key := range keys {
					if keyStr, ok := key.(string); ok {
						issueKeys = append(issueKeys, keyStr)
					} else {
						log.Println("Invalid type found in issueKeys, expected string.")
					}
				}
			}
			if br, ok := attrMap["environmentId"].(string); ok {
				fmt.Println("AM I here")
				envId = br
			}
			if eType, ok := attrMap["environmentType"].(string); ok {
				envType = eType
			}
			if ename, ok := attrMap["environmentName"].(string); ok {
				envName = ename
			}
			if s, ok := attrMap["site"].(string); ok {
				site = s
			}
			if st, ok := attrMap["state"].(string); ok {
				state = st
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	withProperties := make(map[string]interface{})

	withProperties["connect_key"] = "<+secrets.getValue(\"JIRA_CONNECT_KEY\")>"
	withProperties["project"] = "$JIRA_PROJECT"
	withProperties["instance"] = "$JIRA_SITE_ID"
	if envId != "" {
		withProperties["environment_id"] = envId
	}
	if envType != "" {
		withProperties["environment_type"] = envType
	}
	if envName != "" {
		withProperties["environment_name"] = envName
	}
	if site != "" {
		withProperties["instance"] = site
	}
	if state != "" {
		withProperties["state"] = state
	}
	if len(issueKeys) > 0 {
		fmt.Println("Inside Issue Keys")
		fmt.Println(issueKeys)
		withProperties["issuekeys"] = issueKeys
	}
	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/jira",
			With:  withProperties,
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}
	return step
}

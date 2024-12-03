package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ConvertAnsiblePlaybook creates a Harness step for Ansible Playbook plugin.
func ConvertAnsiblePlaybook(node Node, arguments map[string]interface{}) *harness.Step {
	playbook, _ := arguments["playbook"].(string)
	become, _ := arguments["become"].(bool)
	becomeUser, _ := arguments["becomeUser"].(string)
	checkMode, _ := arguments["checkMode"].(bool)
	forks, _ := arguments["forks"].(int)
	hostKeyChecking, _ := arguments["hostKeyChecking"].(string)
	inventory, _ := arguments["inventory"].(string)
	limit, _ := arguments["limit"].(string)
	skippedTags, _ := arguments["skippedTags"].(string)
	startAtTask, _ := arguments["startAtTask"].(string)
	tags, _ := arguments["tags"].(string)
	vaultCredentialsId, _ := arguments["vaultCredentialsId"].(string)

	// Extract extraVars from the arguments map
	extraVarsMap, ok := arguments["extraVars"].(map[string]interface{})
	if !ok {
		fmt.Println("Failed to cast extraVars")
		return nil
	}

	// Convert map to JSON string
	extraVarsBytes, err := json.Marshal(extraVarsMap)
	if err != nil {
		fmt.Println("Failed to convert extraVars to JSON")
		return nil
	}

	// Assign JSON string to extraVars
	extraVars := string(extraVarsBytes)

	convertAnsiblePlaybook := &harness.Step{
		Name: "Ansible_Playbook",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With: map[string]interface{}{
				"playbook":       playbook,
				"become":         become,
				"become_user":    becomeUser,
				"check":          checkMode,
				"extra_vars":     extraVars,
				"forks":          forks,
				"list_hosts":     hostKeyChecking,
				"inventory":      inventory,
				"limit":          limit,
				"skip_tags":      skippedTags,
				"start_at_task":  startAtTask,
				"tags":           tags,
				"vault_id":       vaultCredentialsId,
				"vault_password": "<secrets.getValue('vault-password')>",
			},
		},
	}

	return convertAnsiblePlaybook
}

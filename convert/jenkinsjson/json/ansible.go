package json

import (
	"fmt"
	"strconv"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertAnsiblePlaybook creates a Harness step for Ansible Playbook plugin.
func ConvertAnsiblePlaybook(node Node, arguments map[string]interface{}) *harness.Step {
	playbook, _ := arguments["playbook"].(string)
	become, _ := arguments["become"].(bool)
	becomeUser, _ := arguments["becomeUser"].(string)
	checkMode, _ := arguments["checkMode"].(bool)
	hostKeyChecking, _ := arguments["hostKeyChecking"].(string)
	inventory, _ := arguments["inventory"].(string)
	limit, _ := arguments["limit"].(string)
	skippedTags, _ := arguments["skippedTags"].(string)
	startAtTask, _ := arguments["startAtTask"].(string)
	tags, _ := arguments["tags"].(string)
	vaultCredentialsId, _ := arguments["vaultCredentialsId"].(string)

	var forks int
	if value, ok := arguments["forks"]; ok {
		switch v := value.(type) {
		case float64:
			forks = int(v) // JSON numbers are often float64
			fmt.Println("forks is float64, converted to int:", forks)
		case int:
			forks = v
			fmt.Println("forks is int:", forks)
		case string:
			// Convert string to int
			if intValue, err := strconv.Atoi(v); err == nil {
				forks = intValue
				fmt.Println("forks is string, converted to int:", forks)
			} else {
				fmt.Println("Failed to convert forks from string to int:", err)
			}
		default:
			fmt.Println("forks has an unexpected type:", v)
			forks = 0 // Default value if type assertion fails
		}
	} else {
		fmt.Println("forks key not found in arguments")
	}

	// Handle extraVars
	var extraVars string
	if value, ok := arguments["extraVars"]; ok {
		switch v := value.(type) {
		case []string:
			// Already in correct format
			extraVars = strings.Join(v, " ")
		case map[string]interface{}:
			// Convert map to "key=value" strings
			var extraVarsList []string
			for key, val := range v {
				extraVarsList = append(extraVarsList, fmt.Sprintf("%s=%v", key, val))
			}
			extraVars = strings.Join(extraVarsList, ",")
		default:
			fmt.Println("Unexpected type for extraVars:", v)
		}
	} else {
		fmt.Println("extraVars key not found in arguments")
	}

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

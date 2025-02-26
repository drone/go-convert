package json

import (
	"fmt"
	"strconv"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertAnsiblePlaybook creates a Harness step for Ansible Playbook plugin with conditional parameters.
func ConvertAnsiblePlaybook(node Node, arguments map[string]interface{}) *harness.Step {
	playbook, _ := arguments["playbook"].(string)
	become, _ := arguments["become"].(bool)
	becomeUser, _ := arguments["becomeUser"].(string)
	checkMode, _ := arguments["checkMode"].(bool)
	hostKeyChecking, _ := arguments["hostKeyChecking"].(bool)
	disableHostKeyChecking, _ := arguments["disableHostKeyChecking"].(bool)
	dynamicInventory, _ := arguments["dynamicInventory"].(bool)
	startAtTask, _ := arguments["startAtTask"].(string)
	extras, _ := arguments["extras"].(string)
	credentialsId, _ := arguments["credentialsId"].(string)
	inventory, _ := arguments["inventory"].(string)
	tags, _ := arguments["tags"].(string)
	sudoUser, _ := arguments["sudoUser"].(string)
	sudo, _ := arguments["sudo"].(bool)
	limit, _ := arguments["limit"].(string)
	skippedTags, _ := arguments["skippedTags"].(string)
	installation, _ := arguments["installation"].(string)
	inventoryContent, _ := arguments["inventoryContent"].(string)

	// Use helper functions for forks and extraVars
	forks := HandleForks(arguments)
	extraVars := HandleExtraVars(arguments)

	// Create the "with" map dynamically
	withMap := map[string]interface{}{
		"mode":                      "playbook",
		"playbook":                  playbook,
		"become":                    become,
		"become_user":               becomeUser,
		"check":                     checkMode,
		"forks":                     forks,
		"host_key_checking":         hostKeyChecking,
		"disable_host_key_checking": disableHostKeyChecking,
		"dynamic_inventory":         dynamicInventory,
		"start_at_task":             startAtTask,
		"extras":                    extras,
		"private_key":               credentialsId,
		"inventory":                 inventory,
		"inventory_content":         inventoryContent,
		"tags":                      tags,
		"sudo_user":                 sudoUser,
		"sudo":                      sudo,
		"limit":                     limit,
		"skip_tags":                 skippedTags,
		"installation":              installation,
		"extra_vars":                extraVars,
	}

	// Conditionally add vault-related parameters
	if _, exists := arguments["vaultCredentialsId"]; exists {
		vaultCredentialsId, _ := arguments["vaultCredentialsId"].(string)
		withMap["vault_id"] = vaultCredentialsId
	}
	if _, exists := arguments["vaultTmpPath"]; exists {
		vaultTmpPath, _ := arguments["vaultTmpPath"].(string)
		withMap["vault_tmp_path"] = vaultTmpPath
	}

	return &harness.Step{
		Name: "Ansible_Playbook",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With:  withMap,
		},
	}
}

// ConvertAnsibleAdhoc creates a Harness step for Ansible Ad-Hoc plugin with dynamic parameters.
func ConvertAnsibleAdhoc(node Node, arguments map[string]interface{}) *harness.Step {
	hosts, _ := arguments["hosts"].(string)
	module, _ := arguments["module"].(string)
	moduleArguments, _ := arguments["moduleArguments"].(string)
	dynamicInventory, _ := arguments["dynamicInventory"].(bool)
	inventory, _ := arguments["inventory"].(string)
	inventoryContent, _ := arguments["inventoryContent"].(string)
	become, _ := arguments["become"].(bool)
	becomeUser, _ := arguments["becomeUser"].(string)
	extras, _ := arguments["extras"].(string)
	hostKeyChecking, _ := arguments["hostKeyChecking"].(bool)
	credentialsId, _ := arguments["credentialsId"].(string)
	installation, _ := arguments["installation"].(string)

	// Use helper functions for forks and extraVars
	forks := HandleForks(arguments)
	extraVars := HandleExtraVars(arguments)

	// Create the "with" map dynamically
	withMap := map[string]interface{}{
		"mode":              "adhoc",
		"hosts":             hosts,
		"module":            module,
		"module_args":       moduleArguments,
		"dynamic_inventory": dynamicInventory,
		"inventory":         inventory,
		"inventory_content": inventoryContent,
		"become":            become,
		"become_user":       becomeUser,
		"forks":             forks,
		"extras":            extras,
		"extra_vars":        extraVars,
		"private_key":       credentialsId,
		"host_key_checking": hostKeyChecking,
		"installation":      installation,
	}

	// Conditionally add vault-related parameters
	if _, exists := arguments["vaultCredentialsId"]; exists {
		vaultCredentialsId, _ := arguments["vaultCredentialsId"].(string)
		withMap["vault_credentials_key"] = vaultCredentialsId
	}
	if _, exists := arguments["vaultTmpPath"]; exists {
		vaultTmpPath, _ := arguments["vaultTmpPath"].(string)
		withMap["vault_tmp_path"] = vaultTmpPath
	}

	return &harness.Step{
		Name: "Ansible_Adhoc",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With:  withMap,
		},
	}
}

// ConvertAnsibleVault creates a Harness step for Ansible Vault plugin with dynamic parameters.
func ConvertAnsibleVault(node Node, arguments map[string]interface{}) *harness.Step {
	action, _ := arguments["action"].(string)
	input, _ := arguments["input"].(string)
	output, _ := arguments["output"].(string)
	vaultCredentialsId, _ := arguments["vaultCredentialsId"].(string)
	newVaultCredentialsId, _ := arguments["newVaultCredentialsId"].(string)
	vaultTmpPath, _ := arguments["vaultTmpPath"].(string)
	content, _ := arguments["content"].(string)
	installation, _ := arguments["installation"].(string)

	// Create the "with" map dynamically
	withMap := map[string]interface{}{
		"mode":         "vault",
		"action":       action,
		"input":        input,
		"output":       output,
		"content":      content,
		"installation": installation,
	}

	// Conditionally add vault-related parameters
	if _, exists := arguments["vaultCredentialsId"]; exists {
		withMap["vault_credentials_key"] = vaultCredentialsId
	}
	if _, exists := arguments["newVaultCredentialsId"]; exists {
		withMap["new_vault_credentials_key"] = newVaultCredentialsId
	}
	if _, exists := arguments["vaultTmpPath"]; exists {
		withMap["vault_tmp_path"] = vaultTmpPath
	}

	return &harness.Step{
		Name: "Ansible_Vault",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With:  withMap,
		},
	}
}

// HandleForks converts the forks parameter into an integer.
func HandleForks(arguments map[string]interface{}) int {
	var forks int
	if value, ok := arguments["forks"]; ok {
		switch v := value.(type) {
		case float64:
			forks = int(v)
		case int:
			forks = v
		case string:
			intValue, err := strconv.Atoi(v)
			if err == nil {
				forks = intValue
			} else {
				fmt.Println("Error converting forks from string to int:", err)
				forks = 0
			}
		default:
			fmt.Println("Unexpected type for forks:", v)
			forks = 0
		}
	}
	return forks
}

// HandleExtraVars converts extraVars into a comma-separated string of key-value pairs.
func HandleExtraVars(arguments map[string]interface{}) string {
	var extraVars string
	if value, ok := arguments["extraVars"]; ok {
		switch v := value.(type) {
		case []map[string]string:
			var extraVarsList []string
			for _, pair := range v {
				extraVarsList = append(extraVarsList, fmt.Sprintf("%s=%s", pair["key"], pair["value"]))
			}
			extraVars = strings.Join(extraVarsList, ",")
		default:
			fmt.Println("Unexpected type for extraVars:", v)
		}
	}
	return extraVars
}

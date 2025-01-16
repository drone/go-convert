package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertAnsiblePlaybook(t *testing.T) {
	var tests []runner
	tests = append(tests, prepare(t, "/ansible/ansible-playbook/ansible_playbook_snippet", &harness.Step{
		Id:   "ansiblePlaybookb8e0a2",
		Name: "Ansible_Playbook",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With: map[string]interface{}{
				"become":                    false,
				"become_user":               "root",
				"check":                     false,
				"disable_host_key_checking": true,
				"dynamic_inventory":         false,
				"extra_vars":                "",
				"extras":                    "--timeout=30",
				"forks":                     5,
				"host_key_checking":         false,
				"installation":              "",
				"inventory":                 "inventory",
				"inventory_content":         "",
				"limit":                     "",
				"mode":                      "playbook",
				"playbook":                  "playbook.yml",
				"private_key":               "",
				"skip_tags":                 "",
				"start_at_task":             "",
				"sudo":                      false,
				"sudo_user":                 "",
				"tags":                      "",
				"vault_id":                  "",
				"vault_tmp_path":            "/tmp",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertAnsiblePlaybook(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertAnsiblePlaybook() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertAnsibleAdhoc(t *testing.T) {
	var tests []runner
	tests = append(tests, prepare(t, "/ansible/ansible-adhoc/ansible_adhoc_snippet", &harness.Step{
		Id:   "ansibleAdhoc4eb052",
		Name: "Ansible_Adhoc",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With: map[string]interface{}{
				"become":                false,
				"become_user":           "root",
				"dynamic_inventory":     false,
				"extra_vars":            "",
				"extras":                "--timeout=30",
				"forks":                 5,
				"host_key_checking":     false,
				"hosts":                 "all",
				"installation":          "",
				"inventory":             "inventory",
				"inventory_content":     "[all]\nlocalhost ansible_connection=local",
				"mode":                  "adhoc",
				"module":                "ping",
				"module_args":           "",
				"private_key":           "",
				"vault_credentials_key": "",
				"vault_tmp_path":        "/tmp",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertAnsibleAdhoc(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertAnsibleAdhoc() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertAnsibleVault(t *testing.T) {
	var tests []runner
	tests = append(tests, prepare(t, "/ansible/ansible-vault/ansible_vault_snippet", &harness.Step{
		Id:   "ansibleVault60a60b",
		Name: "Ansible_Vault",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With: map[string]interface{}{
				"action":                "encrypt",
				"content":               "",
				"input":                 "content_to_encrypt.txt",
				"installation":          "",
				"mode":                  "vault",
				"output":                "encrypted_content.txt",
				"vault_credentials_key": "test-vault-id",
				"vault_tmp_path":        "/tmp",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertAnsibleVault(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertAnsibleVault() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertAnsiblePlaybook(t *testing.T) {
	var tests []runner
	tests = append(tests, prepare(t, "/ansible/ansible_snippet", &harness.Step{
		Id:   "ansiblePlaybookb2fcf1",
		Name: "Ansible_Playbook",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/ansible",
			With: map[string]interface{}{
				"become":         true,
				"become_user":    "root",
				"check":          false,
				"extra_vars":     "{\"var1\":\"value1\",\"var2\":\"value2\"}",
				"forks":          0,
				"inventory":      "inventory.ini",
				"limit":          "localhost",
				"list_hosts":     "",
				"playbook":       "example-playbook.yml",
				"skip_tags":      "debug",
				"start_at_task":  "",
				"tags":           "setup,deploy",
				"vault_id":       "",
				"vault_password": "<secrets.getValue('vault-password')>",
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

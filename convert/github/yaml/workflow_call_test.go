package yaml

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestOnConditionsWorkflowCall(t *testing.T) {
	tests := []struct {
		yaml string
		want OnConditions
	}{
		{
			yaml: `
workflow_call:
  inputs:
    username:
      description: A username passed from the caller workflow
      default: john-doe
      required: false
      type: string
`,
			want: OnConditions{
				WorkflowCall: &WorkflowCallCondition{
					Inputs: map[string]interface{}{
						"username": map[string]interface{}{
							"description": "A username passed from the caller workflow",
							"default":     "john-doe",
							"required":    false,
							"type":        "string",
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_call:
  outputs:
    workflow_output1:
      description: The first job output
      value: ${{ jobs.my_job.outputs.job_output1 }}
    workflow_output2:
      description: The second job output
      value: ${{ jobs.my_job.outputs.job_output2 }}

`,
			want: OnConditions{
				WorkflowCall: &WorkflowCallCondition{
					Outputs: map[string]interface{}{
						"workflow_output1": map[string]interface{}{
							"description": "The first job output",
							"value":       "${{ jobs.my_job.outputs.job_output1 }}",
						},
						"workflow_output2": map[string]interface{}{
							"description": "The second job output",
							"value":       "${{ jobs.my_job.outputs.job_output2 }}",
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_call:
  secrets:
    access-token:
      description: A token passed from the caller workflow
      required: false
`,
			want: OnConditions{
				WorkflowCall: &WorkflowCallCondition{
					Secrets: map[string]WorkflowSecrets{
						"access-token": {
							Description: "A token passed from the caller workflow",
							Required:    false,
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(OnConditions)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Log(test.yaml)
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Log(test.yaml)
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestOnConditionsWorkflowCall_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Concurrency))
	if err == nil || err.Error() != "failed to unmarshal concurrency" {
		t.Errorf("Expect error, got %s", err)
	}
}

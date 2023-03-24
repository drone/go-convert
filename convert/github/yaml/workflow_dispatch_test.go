package yaml

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestOnConditionsWorkflowDispatch(t *testing.T) {
	tests := []struct {
		yaml string
		want OnConditions
	}{
		{
			yaml: `
workflow_dispatch:
  inputs:
    logLevel:
      description: Log level
      required: true
      default: warning
      type: choice
      options:
        - info
        - warning
        - debug
    print_tags:
      description: True to print to STDOUT
      required: true
      type: boolean
    tags:
      description: Test scenario tags
      required: true
      type: string
    environment:
      description: Environment to run tests against
      type: environment
`,
			want: OnConditions{
				WorkflowDispatch: &WorkflowDispatchCondition{
					Inputs: map[string]InputDefinition{
						"logLevel": {
							Description: "Log level",
							Required:    true,
							Default:     "warning",
							Type:        "choice",
							Options: []interface{}{
								"info",
								"warning",
								"debug",
							},
						},
						"print_tags": {
							Description: "True to print to STDOUT",
							Required:    true,
							Type:        "boolean",
						},
						"tags": {
							Description: "Test scenario tags",
							Required:    true,
							Type:        "string",
						},
						"environment": {
							Description: "Environment to run tests against",
							Type:        "environment",
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_dispatch:
  inputs:
    logLevel:
      description: Log level
      required: true
      default: warning
      type: choice
      options:
        - info
        - warning
        - debug
    print_tags:
      description: True to print to STDOUT
      required: true
      type: boolean
    tags:
      description: Test scenario tags
      required: true
      type: string
    environment:
      description: Environment to run tests against
      type: environment
      required: true

`,
			want: OnConditions{
				WorkflowDispatch: &WorkflowDispatchCondition{
					Inputs: map[string]InputDefinition{
						"logLevel": {
							Description: "Log level",
							Required:    true,
							Default:     "warning",
							Type:        "choice",
							Options: []interface{}{
								"info",
								"warning",
								"debug",
							},
						},
						"print_tags": {
							Description: "True to print to STDOUT",
							Required:    true,
							Type:        "boolean",
						},
						"tags": {
							Description: "Test scenario tags",
							Required:    true,
							Type:        "string",
						},
						"environment": {
							Description: "Environment to run tests against",
							Type:        "environment",
							Required:    true,
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

func TestOnConditionsWorkflowDispatch_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Concurrency))
	if err == nil || err.Error() != "failed to unmarshal concurrency" {
		t.Errorf("Expect error, got %s", err)
	}
}

package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepAwsSamBuild(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected *v1.StepTemplate
	}{
		{
			name: "full AWS SAM build",
			step: &v0.Step{
				Spec: &v0.StepAwsSamBuild{
					ConnectorRef:        "account.harnessImage",
					Image:               "harnessdev/aws-sam-dev:0.0.12",
					PreExecution:        "npm install",
					BuildCommandOptions: &flexible.Field[[]string]{Value: []string{"--use-container", "--parallel"}},
					RunAsUser:           &flexible.Field[int]{Value: 1000},
					Privileged:          &flexible.Field[bool]{Value: true},
					ImagePullPolicy:     "IfNotPresent",
				},
			},
			expected: &v1.StepTemplate{
				Uses: "awsSamBuildStep",
				With: map[string]interface{}{
					"container_registry": "account.harnessImage",
					"image":              "harnessdev/aws-sam-dev:0.0.12",
					"cmd_opts":           "--use-container --parallel",
					"pre_exec_cmd":       "npm install",
				},
				Container: &v1.Container{
					User:       &flexible.Field[int]{Value: 1000},
					Privileged: &flexible.Field[bool]{Value: true},
					Pull:       "if-not-exists",
				},
			},
		},
		{
			name: "minimal AWS SAM build",
			step: &v0.Step{
				Spec: &v0.StepAwsSamBuild{
					ConnectorRef: "account.harnessImage",
				},
			},
			expected: &v1.StepTemplate{
				Uses: "awsSamBuildStep",
				With: map[string]interface{}{
					"container_registry": "account.harnessImage",
				},
			},
		},
		{
			name: "build command options as expression",
			step: &v0.Step{
				Spec: &v0.StepAwsSamBuild{
					BuildCommandOptions: &flexible.Field[[]string]{Value: "<+input>"},
				},
			},
			expected: &v1.StepTemplate{
				Uses: "awsSamBuildStep",
				With: map[string]interface{}{
					"cmd_opts": "<+input>",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertStepAwsSamBuild(tt.step)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("ConvertStepAwsSamBuild() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepAwsSamDeploy(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected *v1.StepTemplate
	}{
		{
			name: "full AWS SAM deploy",
			step: &v0.Step{
				Spec: &v0.StepAwsSamDeploy{
					ConnectorRef:         "account.harnessImage",
					Image:                "harnessdev/aws-sam-dev:0.0.12",
					StackName:            "my-sam-stack",
					PreExecution:         "npm install",
					DeployCommandOptions: &flexible.Field[[]string]{Value: []string{"--no-confirm-changeset"}},
				},
			},
			expected: &v1.StepTemplate{
				Uses: "awsSamDeployStep",
				With: map[string]interface{}{
					"container_registry": "account.harnessImage",
					"image":              "harnessdev/aws-sam-dev:0.0.12",
					"stack":              "my-sam-stack",
					"cmd_opts":           "--no-confirm-changeset",
					"pre_exec_cmd":       "npm install",
				},
			},
		},
		{
			name: "minimal AWS SAM deploy with stack only",
			step: &v0.Step{
				Spec: &v0.StepAwsSamDeploy{
					StackName: "my-sam-stack",
				},
			},
			expected: &v1.StepTemplate{
				Uses: "awsSamDeployStep",
				With: map[string]interface{}{
					"stack": "my-sam-stack",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertStepAwsSamDeploy(tt.step)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("ConvertStepAwsSamDeploy() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepServerlessAwsLambdaDeployV2(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected *v1.StepTemplate
	}{
		{
			name: "full serverless deploy",
			step: &v0.Step{
				Spec: &v0.StepServerlessAwsLambdaDeployV2{
					ConnectorRef:         "account.harnessImage",
					Image:                "harnessdev/serverless-plugin:0.0.2",
					PreExecution:         "npm install",
					DeployCommandOptions: &flexible.Field[[]string]{Value: []string{"--verbose"}},
					RunAsUser:            &flexible.Field[int]{Value: 0},
					Privileged:           &flexible.Field[bool]{Value: false},
					ImagePullPolicy:      "Always",
				},
			},
			expected: &v1.StepTemplate{
				Uses: "serverlessDeployStep",
				With: map[string]interface{}{
					"container_registry": "account.harnessImage",
					"image":              "harnessdev/serverless-plugin:0.0.2",
					"cmd_opts":           "--verbose",
					"pre_exec_cmd":       "npm install",
				},
				Container: &v1.Container{
					User:       &flexible.Field[int]{Value: 0},
					Privileged: &flexible.Field[bool]{Value: false},
					Pull:       "always",
				},
			},
		},
		{
			name: "minimal serverless deploy",
			step: &v0.Step{
				Spec: &v0.StepServerlessAwsLambdaDeployV2{
					Image: "harnessdev/serverless-plugin:0.0.2",
				},
			},
			expected: &v1.StepTemplate{
				Uses: "serverlessDeployStep",
				With: map[string]interface{}{
					"image": "harnessdev/serverless-plugin:0.0.2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertStepServerlessAwsLambdaDeployV2(tt.step)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("ConvertStepServerlessAwsLambdaDeployV2() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepServerlessAwsLambdaPackageV2(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected *v1.StepTemplate
	}{
		{
			name: "full serverless package",
			step: &v0.Step{
				Spec: &v0.StepServerlessAwsLambdaPackageV2{
					ConnectorRef:          "account.harnessImage",
					Image:                 "harnessdev/serverless-plugin:0.0.2",
					PreExecution:          "pip install",
					PackageCommandOptions: &flexible.Field[[]string]{Value: []string{"--verbose", "--stage", "dev"}},
				},
			},
			expected: &v1.StepTemplate{
				Uses: "serverlessPackageStep",
				With: map[string]interface{}{
					"container_registry": "account.harnessImage",
					"image":              "harnessdev/serverless-plugin:0.0.2",
					"cmd_opts":           "--verbose --stage dev",
					"pre_exec_cmd":       "pip install",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertStepServerlessAwsLambdaPackageV2(tt.step)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("ConvertStepServerlessAwsLambdaPackageV2() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepServerlessAwsLambdaRollbackV2(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected *v1.StepTemplate
	}{
		{
			name: "full serverless rollback",
			step: &v0.Step{
				Spec: &v0.StepServerlessAwsLambdaRollbackV2{
					ConnectorRef: "account.harnessImage",
					Image:        "harnessdev/serverless-plugin:0.0.2",
					PreExecution: "npm install",
				},
			},
			expected: &v1.StepTemplate{
				Uses: "serverlessRollbackStep",
				With: map[string]interface{}{
					"container_registry": "account.harnessImage",
					"image":              "harnessdev/serverless-plugin:0.0.2",
					"pre_exec_cmd":       "npm install",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertStepServerlessAwsLambdaRollbackV2(tt.step)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("ConvertStepServerlessAwsLambdaRollbackV2() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

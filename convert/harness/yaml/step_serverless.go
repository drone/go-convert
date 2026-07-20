package yaml

import "github.com/drone/go-convert/internal/flexible"

// v0 struct ports of the harness-core AWS SAM and Serverless (AWS Lambda V2)
// container-based CD step specs. Field yaml tags match the authoritative
// *StepInfo.java JSON keys.
type (
	// CD: AWS SAM Build (harness-core AwsSamBuildStepInfo, extends AwsSamBaseStepInfo)
	StepAwsSamBuild struct {
		CommonStepSpec
		Image                              string                    `json:"image,omitempty"                              yaml:"image,omitempty"`
		ConnectorRef                       string                    `json:"connectorRef,omitempty"                       yaml:"connectorRef,omitempty"`
		Resources                          *Resources                `json:"resources,omitempty"                          yaml:"resources,omitempty"`
		EnvVariables                       *flexible.Field[map[string]string] `json:"envVariables,omitempty"              yaml:"envVariables,omitempty"`
		Privileged                         *flexible.Field[bool]     `json:"privileged,omitempty"                         yaml:"privileged,omitempty"`
		RunAsUser                          *flexible.Field[int]      `json:"runAsUser,omitempty"                          yaml:"runAsUser,omitempty"`
		ImagePullPolicy                    string                    `json:"imagePullPolicy,omitempty"                    yaml:"imagePullPolicy,omitempty"`
		SamVersion                         string                    `json:"samVersion,omitempty"                         yaml:"samVersion,omitempty"`
		PreExecution                       string                    `json:"preExecution,omitempty"                       yaml:"preExecution,omitempty"`
		BuildCommandOptions                *flexible.Field[[]string] `json:"buildCommandOptions,omitempty"                yaml:"buildCommandOptions,omitempty"`
		SamBuildDockerRegistryConnectorRef string                    `json:"samBuildDockerRegistryConnectorRef,omitempty" yaml:"samBuildDockerRegistryConnectorRef,omitempty"`
	}

	// CD: AWS SAM Deploy (harness-core AwsSamDeployStepInfo, extends AwsSamBaseStepInfo)
	StepAwsSamDeploy struct {
		CommonStepSpec
		Image                string                             `json:"image,omitempty"                yaml:"image,omitempty"`
		ConnectorRef         string                             `json:"connectorRef,omitempty"         yaml:"connectorRef,omitempty"`
		Resources            *Resources                         `json:"resources,omitempty"            yaml:"resources,omitempty"`
		EnvVariables         *flexible.Field[map[string]string] `json:"envVariables,omitempty"         yaml:"envVariables,omitempty"`
		Privileged           *flexible.Field[bool]              `json:"privileged,omitempty"           yaml:"privileged,omitempty"`
		RunAsUser            *flexible.Field[int]               `json:"runAsUser,omitempty"            yaml:"runAsUser,omitempty"`
		ImagePullPolicy      string                             `json:"imagePullPolicy,omitempty"      yaml:"imagePullPolicy,omitempty"`
		SamVersion           string                             `json:"samVersion,omitempty"           yaml:"samVersion,omitempty"`
		PreExecution         string                             `json:"preExecution,omitempty"         yaml:"preExecution,omitempty"`
		DeployCommandOptions *flexible.Field[[]string]          `json:"deployCommandOptions,omitempty" yaml:"deployCommandOptions,omitempty"`
		StackName            string                             `json:"stackName,omitempty"            yaml:"stackName,omitempty"`
	}

	// CD: AWS SAM Rollback (harness-core AwsSamRollbackStepInfo, extends AwsSamRollbackBaseStepInfo).
	// Only delegateSelectors (handled via CommonStepSpec) is a YAML field.
	StepAwsSamRollback struct {
		CommonStepSpec
	}

	// CD: Serverless AWS Lambda Deploy V2 (harness-core ServerlessAwsLambdaDeployV2StepInfo,
	// extends ServerlessAwsLambdaV2BaseStepInfo)
	StepServerlessAwsLambdaDeployV2 struct {
		CommonStepSpec
		Image                string                             `json:"image,omitempty"                yaml:"image,omitempty"`
		ConnectorRef         string                             `json:"connectorRef,omitempty"         yaml:"connectorRef,omitempty"`
		Resources            *Resources                         `json:"resources,omitempty"            yaml:"resources,omitempty"`
		EnvVariables         *flexible.Field[map[string]string] `json:"envVariables,omitempty"         yaml:"envVariables,omitempty"`
		Privileged           *flexible.Field[bool]              `json:"privileged,omitempty"           yaml:"privileged,omitempty"`
		RunAsUser            *flexible.Field[int]               `json:"runAsUser,omitempty"            yaml:"runAsUser,omitempty"`
		ImagePullPolicy      string                             `json:"imagePullPolicy,omitempty"      yaml:"imagePullPolicy,omitempty"`
		ServerlessVersion    string                             `json:"serverlessVersion,omitempty"    yaml:"serverlessVersion,omitempty"`
		PreExecution         string                             `json:"preExecution,omitempty"         yaml:"preExecution,omitempty"`
		DeployCommandOptions *flexible.Field[[]string]          `json:"deployCommandOptions,omitempty" yaml:"deployCommandOptions,omitempty"`
	}

	// CD: Serverless AWS Lambda Package V2 (harness-core ServerlessAwsLambdaPackageV2StepInfo,
	// extends ServerlessAwsLambdaV2BaseStepInfo)
	StepServerlessAwsLambdaPackageV2 struct {
		CommonStepSpec
		Image                 string                             `json:"image,omitempty"                 yaml:"image,omitempty"`
		ConnectorRef          string                             `json:"connectorRef,omitempty"          yaml:"connectorRef,omitempty"`
		Resources             *Resources                         `json:"resources,omitempty"             yaml:"resources,omitempty"`
		EnvVariables          *flexible.Field[map[string]string] `json:"envVariables,omitempty"          yaml:"envVariables,omitempty"`
		Privileged            *flexible.Field[bool]              `json:"privileged,omitempty"            yaml:"privileged,omitempty"`
		RunAsUser             *flexible.Field[int]               `json:"runAsUser,omitempty"             yaml:"runAsUser,omitempty"`
		ImagePullPolicy       string                             `json:"imagePullPolicy,omitempty"       yaml:"imagePullPolicy,omitempty"`
		ServerlessVersion     string                             `json:"serverlessVersion,omitempty"     yaml:"serverlessVersion,omitempty"`
		PreExecution          string                             `json:"preExecution,omitempty"          yaml:"preExecution,omitempty"`
		PackageCommandOptions *flexible.Field[[]string]          `json:"packageCommandOptions,omitempty" yaml:"packageCommandOptions,omitempty"`
	}

	// CD: Serverless AWS Lambda Prepare Rollback V2 (harness-core
	// ServerlessAwsLambdaPrepareRollbackV2StepInfo, extends ServerlessAwsLambdaV2BaseStepInfo).
	// Same v0 YAML fields as the Rollback V2 step; differs only by type + internal FQN.
	StepServerlessAwsLambdaPrepareRollbackV2 struct {
		CommonStepSpec
		Image             string                             `json:"image,omitempty"             yaml:"image,omitempty"`
		ConnectorRef      string                             `json:"connectorRef,omitempty"      yaml:"connectorRef,omitempty"`
		Resources         *Resources                         `json:"resources,omitempty"         yaml:"resources,omitempty"`
		EnvVariables      *flexible.Field[map[string]string] `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Privileged        *flexible.Field[bool]              `json:"privileged,omitempty"        yaml:"privileged,omitempty"`
		RunAsUser         *flexible.Field[int]               `json:"runAsUser,omitempty"         yaml:"runAsUser,omitempty"`
		ImagePullPolicy   string                             `json:"imagePullPolicy,omitempty"   yaml:"imagePullPolicy,omitempty"`
		ServerlessVersion string                             `json:"serverlessVersion,omitempty" yaml:"serverlessVersion,omitempty"`
		PreExecution      string                             `json:"preExecution,omitempty"      yaml:"preExecution,omitempty"`
	}

	// CD: Serverless AWS Lambda Rollback V2 (harness-core ServerlessAwsLambdaRollbackV2StepInfo,
	// extends ServerlessAwsLambdaV2BaseStepInfo)
	StepServerlessAwsLambdaRollbackV2 struct {
		CommonStepSpec
		Image             string                             `json:"image,omitempty"             yaml:"image,omitempty"`
		ConnectorRef      string                             `json:"connectorRef,omitempty"      yaml:"connectorRef,omitempty"`
		Resources         *Resources                         `json:"resources,omitempty"         yaml:"resources,omitempty"`
		EnvVariables      *flexible.Field[map[string]string] `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Privileged        *flexible.Field[bool]              `json:"privileged,omitempty"        yaml:"privileged,omitempty"`
		RunAsUser         *flexible.Field[int]               `json:"runAsUser,omitempty"         yaml:"runAsUser,omitempty"`
		ImagePullPolicy   string                             `json:"imagePullPolicy,omitempty"   yaml:"imagePullPolicy,omitempty"`
		ServerlessVersion string                             `json:"serverlessVersion,omitempty" yaml:"serverlessVersion,omitempty"`
		PreExecution      string                             `json:"preExecution,omitempty"      yaml:"preExecution,omitempty"`
	}
)

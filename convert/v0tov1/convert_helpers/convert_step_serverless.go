package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// joinCommandOptions renders a v0 command-options list (*flexible.Field[[]string])
// into the single-string form expected by the template `cmd_opts` input. A list
// value is space-joined; an expression is passed through unchanged.
func joinCommandOptions(opts *flexible.Field[[]string]) string {
	if opts == nil {
		return ""
	}
	if list, ok := opts.AsStruct(); ok {
		return strings.Join(list, " ")
	}
	if expr, ok := opts.AsString(); ok {
		return expr
	}
	return ""
}

// ConvertStepAwsSamBuild converts a v0 AwsSamBuild step to the v1 awsSamBuildStep template.
func ConvertStepAwsSamBuild(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepAwsSamBuild)
	if !ok || spec == nil {
		return nil
	}

	with := make(map[string]interface{})

	// v0 connectorRef/image identify the plugin container image + its registry connector,
	// which map to the template container_registry/image inputs.
	if spec.ConnectorRef != "" {
		with["container_registry"] = spec.ConnectorRef
	}
	if spec.Image != "" {
		with["image"] = spec.Image
	}
	if opts := joinCommandOptions(spec.BuildCommandOptions); opts != "" {
		with["cmd_opts"] = opts
	}
	if spec.PreExecution != "" {
		with["pre_exec_cmd"] = spec.PreExecution
	}

	
	// FEATURE GAP: v0 samVersion, and samBuildDockerRegistryConnectorRef
	// have no awsSamBuildStep template input
	// (samBuildDockerRegistryConnectorRef is a connector, but the template exposes only
	// registry_url/username/pwd strings). Template inputs connector (AWS), region, path,
	// work_dir, timeout, docker_retry_count, registry_url/username/pwd, and log_level are
	// derived from infra/service or defaulted by the template, so they have no v0
	// step-level source.

	return &v1.StepTemplate{
		Uses:      v1.StepTypeAwsSamBuild,
		With:      with,
		Container: ConvertTemplateContainer(
			spec.RunAsUser,
			spec.Resources,
			WithPrivileged(spec.Privileged),
			WithImagePullPolicy(spec.ImagePullPolicy),
		),
	}
}

// ConvertStepAwsSamDeploy converts a v0 AwsSamDeploy step to the v1 awsSamDeployStep template.
func ConvertStepAwsSamDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepAwsSamDeploy)
	if !ok || spec == nil {
		return nil
	}

	with := make(map[string]interface{})

	// v0 connectorRef/image identify the plugin container image + its registry connector,
	// which map to the template container_registry/image inputs.
	if spec.ConnectorRef != "" {
		with["container_registry"] = spec.ConnectorRef
	}
	if spec.Image != "" {
		with["image"] = spec.Image
	}
	// stack is a required template input (maps from v0 stackName).
	if spec.StackName != "" {
		with["stack"] = spec.StackName
	}
	if opts := joinCommandOptions(spec.DeployCommandOptions); opts != "" {
		with["cmd_opts"] = opts
	}
	if spec.PreExecution != "" {
		with["pre_exec_cmd"] = spec.PreExecution
	}

	
	// FEATURE GAP: v0 samVersion have no
	// awsSamDeployStep template input. Template inputs connector (AWS), region, path,
	// work_dir, timeout, registry_url/username/pwd, and log_level are derived from
	// infra/service or defaulted by the template.

	return &v1.StepTemplate{
		Uses:      v1.StepTypeAwsSamDeploy,
		With:      with,
		Container: ConvertTemplateContainer(
			spec.RunAsUser,
			spec.Resources,
			WithPrivileged(spec.Privileged),
			WithImagePullPolicy(spec.ImagePullPolicy),
		),
	}
}

// ConvertStepServerlessAwsLambdaDeployV2 converts a v0 ServerlessAwsLambdaDeployV2 step
// to the v1 serverlessDeployStep template.
func ConvertStepServerlessAwsLambdaDeployV2(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepServerlessAwsLambdaDeployV2)
	if !ok || spec == nil {
		return nil
	}

	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["container_registry"] = spec.ConnectorRef
	}
	if spec.Image != "" {
		with["image"] = spec.Image
	}
	if opts := joinCommandOptions(spec.DeployCommandOptions); opts != "" {
		with["cmd_opts"] = opts
	}
	if spec.PreExecution != "" {
		with["pre_exec_cmd"] = spec.PreExecution
	}

	
	// FEATURE GAP: v0 serverlessVersion have
	// no serverlessDeployStep template input. Template inputs connector, region, stage,
	// path, log_level, timeout, client_path, and work_dir are derived from infra/service
	// or defaulted by the template.

	return &v1.StepTemplate{
		Uses:      v1.StepTypeServerlessDeploy,
		With:      with,
		Container: ConvertTemplateContainer(
			spec.RunAsUser,
			spec.Resources,
			WithPrivileged(spec.Privileged),
			WithImagePullPolicy(spec.ImagePullPolicy),
		),
	}
}

// ConvertStepServerlessAwsLambdaPackageV2 converts a v0 ServerlessAwsLambdaPackageV2 step
// to the v1 serverlessPackageStep template.
func ConvertStepServerlessAwsLambdaPackageV2(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepServerlessAwsLambdaPackageV2)
	if !ok || spec == nil {
		return nil
	}

	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["container_registry"] = spec.ConnectorRef
	}
	if spec.Image != "" {
		with["image"] = spec.Image
	}
	if opts := joinCommandOptions(spec.PackageCommandOptions); opts != "" {
		with["cmd_opts"] = opts
	}
	if spec.PreExecution != "" {
		with["pre_exec_cmd"] = spec.PreExecution
	}

	
	// FEATURE GAP: v0 serverlessVersion have
	// no serverlessPackageStep template input. Template inputs connector, region, stage,
	// path, log_level, timeout, client_path, and work_dir are derived from infra/service
	// or defaulted by the template.

	return &v1.StepTemplate{
		Uses:      v1.StepTypeServerlessPackage,
		With:      with,
		Container: ConvertTemplateContainer(
			spec.RunAsUser,
			spec.Resources,
			WithPrivileged(spec.Privileged),
			WithImagePullPolicy(spec.ImagePullPolicy),
		),
	}
}

// ConvertStepServerlessAwsLambdaRollbackV2 converts a v0 ServerlessAwsLambdaRollbackV2 step
// to the v1 serverlessRollbackStep template.
func ConvertStepServerlessAwsLambdaRollbackV2(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepServerlessAwsLambdaRollbackV2)
	if !ok || spec == nil {
		return nil
	}

	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["container_registry"] = spec.ConnectorRef
	}
	if spec.Image != "" {
		with["image"] = spec.Image
	}
	if spec.PreExecution != "" {
		with["pre_exec_cmd"] = spec.PreExecution
	}

	
	// FEATURE GAP: v0 serverlessVersion have
	// no serverlessRollbackStep template input. Template inputs connector, region,
	// log_level, timeout, and stack (from ${{rollback.data.PLUGIN_STACK_DETAILS}}) are
	// derived at runtime or defaulted by the template.

	return &v1.StepTemplate{
		Uses:      v1.StepTypeServerlessRollback,
		With:      with,
		Container: ConvertTemplateContainer(
			spec.RunAsUser,
			spec.Resources,
			WithPrivileged(spec.Privileged),
			WithImagePullPolicy(spec.ImagePullPolicy),
		),
	}
}

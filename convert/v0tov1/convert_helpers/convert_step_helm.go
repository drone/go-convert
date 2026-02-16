package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// Helm step with configuration structs with JSON tags
type HelmBGDeployWith struct {
	Flags                      []string            `json:"flags,omitempty"`
	IgnoreFailedReleaseHistory *flexible.Field[bool] `json:"ignore_failed_release,omitempty"`
	SkipSteadyStateCheck       *flexible.Field[bool] `json:"skip_deploy_steady_check,omitempty"`
	UseUpgradeWithInstall      *flexible.Field[bool] `json:"upgrade_with_install,omitempty"`
	Envvars                    []map[string]string   `json:"deploy_env_vars,omitempty"`
	ChartTest                  *flexible.Field[bool] `json:"deploy_test,omitempty"`
}

type HelmBlueGreenSwapStepWith struct {
	Flags []string `json:"flags,omitempty"`
}

type HelmCanaryDeployWith struct {
	Flags                      []string            `json:"flags,omitempty"`
	Envvars                    []map[string]string `json:"deploy_env_vars,omitempty"`
	IgnoreFailedReleaseHistory *flexible.Field[bool] `json:"ignore_failed_release,omitempty"`
	SkipSteadyStateCheck       *flexible.Field[bool] `json:"skip_deploy_steady_check,omitempty"`
	UseUpgradeWithInstall      *flexible.Field[bool] `json:"upgrade_with_install,omitempty"`
	InstanceUnitType           string              `json:"instanceunittype,omitempty"`
	Instances                  *flexible.Field[int]               `json:"instances,omitempty"`
    ChartTest             *flexible.Field[bool]                 `json:"deploy_test,omitempty"`
}

type HelmDeleteWith struct {
	ReleaseName string              `json:"release,omitempty"`
	DryRun      *flexible.Field[bool]                 `json:"dry_run,omitempty"`
	Flags       []string            `json:"deploy_flags,omitempty"`
	Envvars     []map[string]string `json:"deploy_env_vars,omitempty"`
}

type HelmDeployWith struct {
    Flags                 []string            `json:"flags,omitempty"`
    DeployEnvVars         []map[string]string `json:"deploy_env_vars,omitempty"`      
    IgnoreFailedRelease   *flexible.Field[bool]                 `json:"ignore_failed_release,omitempty"` 
    SkipDeploySteadyCheck *flexible.Field[bool]                 `json:"skip_deploy_steady_check,omitempty"` 
    UpgradeWithInstall    *flexible.Field[bool]                 `json:"upgrade_with_install,omitempty"`  
    ChartTest             *flexible.Field[bool]                 `json:"deploy_test,omitempty"`           
}

type HelmRollbackWith struct {
	Flags                []string            `json:"flags,omitempty"`
	Envvars              []map[string]string `json:"envvars,omitempty"`
	SkipSteadyStateCheck *flexible.Field[bool]                 `json:"skip_steady_check,omitempty"`
	ChartTest            *flexible.Field[bool]                 `json:"test,omitempty"`
}

func ConvertStepHelmBGDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepHelmBGDeploy)
	if !ok || spec == nil {
		return nil
	}

	// Convert environment variables to envvars format
	var envvars []map[string]string
	for key, value := range spec.EnvironmentVariables {
		envvars = append(envvars, map[string]string{
			"key":   key,
			"value": value,
		})
	}

	// Build the with parameters
	with := HelmBGDeployWith{
		Flags:   []string{},
		Envvars: envvars,
	}

	// Add boolean flags based on v0 spec
	with.IgnoreFailedReleaseHistory = spec.IgnoreReleaseHistFailStatus
	with.SkipSteadyStateCheck = spec.SkipSteadyStateCheck
	with.UseUpgradeWithInstall = spec.UseUpgradeInstall
	with.ChartTest = spec.RunChartTests

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmBGDeploy,
		With: with,
	}
}

// ConvertStepHelmBlueGreenSwapStep converts v0 HelmBlueGreenSwapStep to v1 template
func ConvertStepHelmBlueGreenSwapStep(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	_, ok := src.Spec.(*v0.StepHelmBlueGreenSwapStep)
	if !ok {
		return nil
	}

	with := HelmBlueGreenSwapStepWith{
		Flags: []string{},
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmBlueGreenSwapStep,
		With: with,
	}
}

// ConvertStepHelmCanaryDeploy converts v0 HelmCanaryDeploy to v1 helmDeployCanaryStep@1.0.0
func ConvertStepHelmCanaryDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepHelmCanaryDeploy)
	if !ok || spec == nil {
		return nil
	}

	// instance selection mapping
	unitType := ""
	var instances *flexible.Field[int]
	if sel := spec.InstanceSelection; sel != nil {
		switch sel.Type {
		case "Count":
			unitType = "count"
			if sel.Spec != nil {
				instances = sel.Spec.Count
			}
		case "Percentage":
			unitType = "percentage"
			if sel.Spec != nil {
				instances = sel.Spec.Percentage
			}
		}
	}

	// env vars
	envvars := make([]map[string]string, 0)
	for k, v := range spec.EnvironmentVariables {
		envvars = append(envvars, map[string]string{"key": k, "value": v})
	}

	with := HelmCanaryDeployWith{
		Flags:   []string{},
		Envvars: envvars, // include empty array if none
	}

	with.IgnoreFailedReleaseHistory = spec.IgnoreReleaseHistFailStatus
	with.SkipSteadyStateCheck = spec.SkipSteadyStateCheck
	with.UseUpgradeWithInstall = spec.UseUpgradeInstall 
	with.InstanceUnitType = unitType
	with.Instances = instances
	with.ChartTest = spec.RunChartTests

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmCanaryDeploy,
		With: with,
	}
}

// ConvertStepHelmDelete converts v0 HelmDelete to v1 helmDeleteStep@1.0.0
func ConvertStepHelmDelete(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepHelmDelete)
	if !ok || sp == nil {
		return nil
	}

	// map command flags
	flags := make([]string, 0, len(sp.CommandFlags))
	flags = append(flags, sp.CommandFlags...)

	// map env vars
	envvars := make([]map[string]string, 0, len(sp.EnvironmentVariables))
	for k, v := range sp.EnvironmentVariables {
		envvars = append(envvars, map[string]string{
			"key":   k,
			"value": v,
		})
	}

	with := HelmDeleteWith{
		ReleaseName: sp.ReleaseName,
		DryRun:      sp.DryRun,
		Flags:       flags,
		Envvars:     envvars,
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmDelete,
		With: with,
	}
}

// ConvertStepHelmDeploy converts v0 HelmDeploy (basic) to v1 helmDeployBasicStep@1.0.0
func ConvertStepHelmDeploy(src *v0.Step) *v1.StepTemplate {
    if src == nil || src.Spec == nil {
        return nil
    }
    sp, ok := src.Spec.(*v0.StepHelmDeploy)
    if !ok || sp == nil {
        return nil
    }

    // map env vars with new field name
    envvars := make([]map[string]string, 0, len(sp.EnvironmentVariables))
    for k, v := range sp.EnvironmentVariables {
        envvars = append(envvars, map[string]string{
            "key":   k,
            "value": v,
        })
    }

    with := HelmDeployWith{
        Flags:         []string{},
        DeployEnvVars: envvars, // Changed from Envvars
    }

	with.IgnoreFailedRelease = sp.IgnoreReleaseHistFailStatus 
	with.SkipDeploySteadyCheck = sp.SkipSteadyStateCheck 
	with.UpgradeWithInstall = sp.UseUpgradeInstall	 
	with.ChartTest = sp.RunChartTests	 

    // Note: skipDryRun and skipCleanup have NO v1 equivalents - intentionally omitted

    return &v1.StepTemplate{
        Uses: v1.StepTypeHelmDeploy,
        With: with,
    }
}

// ConvertStepHelmRollback converts v0 HelmRollback to v1 helmRollbackStep@1.0.0
func ConvertStepHelmRollback(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepHelmRollback)
	if !ok || sp == nil {
		return nil
	}

	// map env vars
	envvars := make([]map[string]string, 0, len(sp.EnvironmentVariables))
	for k, v := range sp.EnvironmentVariables {
		envvars = append(envvars, map[string]string{
			"key":   k,
			"value": v,
		})
	}

	with := HelmRollbackWith{
		Flags:   []string{},
		Envvars: envvars,
	}
	with.SkipSteadyStateCheck = sp.SkipSteadyStateCheck
	with.ChartTest = sp.RunChartTests

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmRollback,
		With: with,
	}
}

func ConvertStepHelmCanaryDelete(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepHelmCanaryDelete)
	if !ok || sp == nil {
		return nil
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmCanaryDelete,
	}
}
package converthelpers

import (
	"strconv"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// Helm step with configuration structs with JSON tags
type HelmBGDeployWith struct {
	Flags                      []string            `json:"flags,omitempty"`
	IgnoreFailedReleaseHistory bool                `json:"ignorefailedreleasehistory,omitempty"`
	SkipSteadyStateCheck       bool                `json:"skipsteadystatecheck,omitempty"`
	UseUpgradeWithInstall      bool                `json:"useupgradewithinstall,omitempty"`
	Envvars                    []map[string]string `json:"envvars,omitempty"`
}

type HelmBlueGreenSwapStepWith struct {
	Flags []string `json:"flags,omitempty"`
}

type HelmCanaryDeployWith struct {
	Flags                      []string            `json:"flags,omitempty"`
	Envvars                    []map[string]string `json:"envvars,omitempty"`
	IgnoreFailedReleaseHistory bool                `json:"ignorefailedreleasehistory,omitempty"`
	SkipSteadyStateCheck       bool                `json:"skipsteadystatecheck,omitempty"`
	UseUpgradeWithInstall      bool                `json:"useupgradewithinstall,omitempty"`
	InstanceUnitType           string              `json:"instanceunittype,omitempty"`
	Instances                  string              `json:"instances,omitempty"`
}

type HelmDeleteWith struct {
	ReleaseName string              `json:"releasename,omitempty"`
	DryRun      bool                `json:"dryrun,omitempty"`
	Flags       []string            `json:"flags,omitempty"`
	Envvars     []map[string]string `json:"envvars,omitempty"`
}

type HelmDeployWith struct {
	Flags                      []string            `json:"flags,omitempty"`
	Envvars                    []map[string]string `json:"envvars,omitempty"`
	IgnoreFailedReleaseHistory bool                `json:"ignorefailedreleasehistory,omitempty"`
	SkipSteadyStateCheck       bool                `json:"skipsteadystatecheck,omitempty"`
	UseUpgradeWithInstall      bool                `json:"useupgradewithinstall,omitempty"`
}

type HelmRollbackWith struct {
	Flags                []string            `json:"flags,omitempty"`
	Envvars              []map[string]string `json:"envvars,omitempty"`
	SkipSteadyStateCheck bool                `json:"skipsteadystatecheck,omitempty"`
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
	if spec.IgnoreReleaseHistFailStatus {
		with.IgnoreFailedReleaseHistory = true
	}

	if spec.SkipSteadyStateCheck {
		with.SkipSteadyStateCheck = true
	}

	if spec.UseUpgradeInstall {
		with.UseUpgradeWithInstall = true
	}

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
	instances := ""
	if sel := spec.InstanceSelection; sel != nil {
		switch sel.Type {
		case "Count":
			unitType = "count"
			if sel.Spec != nil {
				instances = strconv.Itoa(sel.Spec.Count)
			}
		case "Percentage":
			unitType = "percentage"
			if sel.Spec != nil {
				instances = strconv.Itoa(sel.Spec.Percentage)
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

	if spec.IgnoreReleaseHistFailStatus {
		with.IgnoreFailedReleaseHistory = true
	}
	if spec.SkipSteadyStateCheck {
		with.SkipSteadyStateCheck = true
	}
	if spec.UseUpgradeInstall {
		with.UseUpgradeWithInstall = true
	}
	if unitType != "" {
		with.InstanceUnitType = unitType
	}
	if instances != "" {
		with.Instances = instances
	}

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

	// map env vars
	envvars := make([]map[string]string, 0, len(sp.EnvironmentVariables))
	for k, v := range sp.EnvironmentVariables {
		envvars = append(envvars, map[string]string{
			"key":   k,
			"value": v,
		})
	}

	with := HelmDeployWith{
		Flags:   []string{},
		Envvars: envvars,
	}

	if sp.IgnoreReleaseHistFailStatus {
		with.IgnoreFailedReleaseHistory = true
	}
	if sp.SkipSteadyStateCheck {
		with.SkipSteadyStateCheck = true
	}
	if sp.UseUpgradeInstall {
		with.UseUpgradeWithInstall = true
	}
	// TODO
	// if sp.SkipDryRun {
	// 	with["skipdryrun"] = true
	// }
	// if sp.SkipCleanup {
	// 	with["skipcleanup"] = true
	// }

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
	if sp.SkipSteadyStateCheck {
		with.SkipSteadyStateCheck = true
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmRollback,
		With: with,
	}
}

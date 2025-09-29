package converthelpers

import (
	"strconv"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

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
	with := map[string]interface{}{
		"flags": []string{},
	}

	// Add boolean flags based on v0 spec
	if spec.IgnoreReleaseHistFailStatus {
		with["ignorefailedreleasehistory"] = true
	}

	if spec.SkipSteadyStateCheck {
		with["skipsteadystatecheck"] = true
	}

	if spec.UseUpgradeInstall {
		with["useupgradewithinstall"] = true
	}

	// Add environment variables. If none present, set to empty array to match expected YAML shape
	with["envvars"] = envvars

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

	with := map[string]interface{}{
		"flags": []string{},
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

	with := map[string]interface{}{
		"flags":   []string{},
		"envvars": envvars, // include empty array if none
	}

	if spec.IgnoreReleaseHistFailStatus {
		with["ignorefailedreleasehistory"] = true
	}
	if spec.SkipSteadyStateCheck {
		with["skipsteadystatecheck"] = true
	}
	if spec.UseUpgradeInstall {
		with["useupgradewithinstall"] = true
	}
	if unitType != "" {
		with["instanceunittype"] = unitType
	}
	if instances != "" {
		with["instances"] = instances
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

	with := map[string]interface{}{
		"releasename": sp.ReleaseName,
		"dryrun":      sp.DryRun,
		"flags":       flags,
		"envvars":     envvars,
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

	with := map[string]interface{}{
		"flags":   []string{},
		"envvars": envvars,
	}

	if sp.IgnoreReleaseHistFailStatus {
		with["ignorefailedreleasehistory"] = true
	}
	if sp.SkipSteadyStateCheck {
		with["skipsteadystatecheck"] = true
	}
	if sp.UseUpgradeInstall {
		with["useupgradewithinstall"] = true
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

	with := map[string]interface{}{
		"flags":   []string{},
		"envvars": envvars,
	}
	if sp.SkipSteadyStateCheck {
		with["skipsteadystatecheck"] = true
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmRollback,
		With: with,
	}
}

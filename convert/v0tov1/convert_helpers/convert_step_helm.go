package converthelpers

import (
	"fmt"
	"sort"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// HelmFlag is a single entry of the v1 template `flags` input ({command, flag}).
type HelmFlag struct {
	Command string `json:"command,omitempty"`
	Flag    string `json:"flag,omitempty"`
}

// helmCommandTypeToTemplate maps a v0 HelmCommandFlagType enum value (e.g. "Template",
// "Install", "Delete") to the v1 template `flags` command option value (e.g. "template",
// "install", "uninstall"). Unknown values fall back to a lowercased passthrough.
func helmCommandTypeToTemplate(commandType string) string {
	switch commandType {
	case "Template":
		return "template"
	case "Install":
		return "install"
	case "Upgrade":
		return "upgrade"
	case "Rollback":
		return "rollback"
	case "History":
		return "history"
	case "List":
		return "list"
	case "Test":
		return "test"
	// helm `delete` is the deprecated alias of `uninstall`.
	case "Uninstall", "Delete":
		return "uninstall"
	default:
		return strings.ToLower(commandType)
	}
}

// convertHelmCommandFlags maps v0 commandFlags ([]{commandType, flag}) to the v1
// template `flags` input ([]{command, flag}).
func convertHelmCommandFlags(flags []v0.HelmCommandFlag) []HelmFlag {
	out := make([]HelmFlag, 0, len(flags))
	for _, f := range flags {
		out = append(out, HelmFlag{
			Command: helmCommandTypeToTemplate(f.CommandType),
			Flag:    f.Flag,
		})
	}
	return out
}

// helmEnvVars converts a v0 environmentVariables map to the v1 `env_vars` list format.
// Entries are sorted by key for deterministic output.
func helmEnvVars(env map[string]string) []map[string]string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	envvars := make([]map[string]string, 0, len(env))
	for _, k := range keys {
		envvars = append(envvars, map[string]string{"key": k, "value": env[k]})
	}
	return envvars
}

// Helm step with configuration structs with JSON tags
type HelmBGDeployWith struct {
	Flags                      []HelmFlag            `json:"flags,omitempty"`
	IgnoreFailedReleaseHistory *flexible.Field[bool] `json:"ignore_history,omitempty"`
	SkipSteadyStateCheck       *flexible.Field[bool] `json:"skip_steady_check,omitempty"`
	Envvars                    []map[string]string   `json:"env_vars,omitempty"`
	ChartTest                  *flexible.Field[bool] `json:"test,omitempty"`
}

type HelmBlueGreenSwapStepWith struct {
	Flags []HelmFlag `json:"flags,omitempty"`
}

type HelmCanaryDeployWith struct {
	Flags                      []HelmFlag          `json:"flags,omitempty"`
	Envvars                    []map[string]string `json:"env_vars,omitempty"`
	IgnoreFailedReleaseHistory *flexible.Field[bool] `json:"ignore_history,omitempty"`
	SkipSteadyStateCheck       *flexible.Field[bool] `json:"skip_steady_check,omitempty"`
	InstanceUnitType           string              `json:"unit,omitempty"`
	Instances                  string              `json:"instances,omitempty"`
	ChartTest                  *flexible.Field[bool]                 `json:"test,omitempty"`
}

type HelmCanaryDeleteWith struct {
	Flags []HelmFlag `json:"flags,omitempty"`
}

type HelmDeleteWith struct {
	ReleaseName string              `json:"release,omitempty"`
	DryRun      *flexible.Field[bool]                 `json:"dry_run,omitempty"`
	Flags       []string            `json:"flags,omitempty"`
	Envvars     []map[string]string `json:"env_vars,omitempty"`
}

type HelmDeployWith struct {
	Flags                 []HelmFlag          `json:"flags,omitempty"`
	DeployEnvVars         []map[string]string `json:"env_vars,omitempty"`
	IgnoreFailedRelease   *flexible.Field[bool]                 `json:"ignore_history,omitempty"`
	SkipDeploySteadyCheck *flexible.Field[bool]                 `json:"skip_steady_check,omitempty"`
	SkipCleanup           *flexible.Field[bool]                 `json:"skip_cleanup,omitempty"`
	ChartTest             *flexible.Field[bool]                 `json:"test,omitempty"`
}

type HelmRollbackWith struct {
	Flags                []HelmFlag          `json:"flags,omitempty"`
	Envvars              []map[string]string `json:"env_vars,omitempty"`
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

	// Build the with parameters
	with := HelmBGDeployWith{
		Flags:   convertHelmCommandFlags(spec.CommandFlags),
		Envvars: helmEnvVars(spec.EnvironmentVariables),
	}

	// Add boolean flags based on v0 spec
	with.IgnoreFailedReleaseHistory = spec.IgnoreReleaseHistFailStatus
	with.SkipSteadyStateCheck = spec.SkipSteadyStateCheck
	with.ChartTest = spec.RunChartTests

	// FEATURE GAP: v0 useUpgradeInstall has no template input (dropped).

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
	spec, ok := src.Spec.(*v0.StepHelmBlueGreenSwapStep)
	if !ok || spec == nil {
		return nil
	}

	with := HelmBlueGreenSwapStepWith{
		Flags: convertHelmCommandFlags(spec.CommandFlags),
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmBlueGreenSwapStep,
		With: with,
	}
}

// ConvertStepHelmCanaryDeploy converts v0 HelmCanaryDeploy to v1 helmCanaryDeployStep
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

	with := HelmCanaryDeployWith{
		Flags:   convertHelmCommandFlags(spec.CommandFlags),
		Envvars: helmEnvVars(spec.EnvironmentVariables),
	}

	with.IgnoreFailedReleaseHistory = spec.IgnoreReleaseHistFailStatus
	with.SkipSteadyStateCheck = spec.SkipSteadyStateCheck
	with.InstanceUnitType = unitType
	// template `instances` input is a string
	if instances != nil {
		if s, ok := instances.AsString(); ok {
			with.Instances = s
		} else if v, ok := instances.AsStruct(); ok {
			with.Instances = fmt.Sprintf("%d", v)
		}
	}
	with.ChartTest = spec.RunChartTests

	// FEATURE GAP: v0 useUpgradeInstall has no template input (dropped).

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmCanaryDeploy,
		With: with,
	}
}

// ConvertStepHelmDelete converts v0 HelmDelete to v1 helmDeleteStep
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

	with := HelmDeleteWith{
		ReleaseName: sp.ReleaseName,
		DryRun:      sp.DryRun,
		Flags:       flags,
		Envvars:     helmEnvVars(sp.EnvironmentVariables),
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmDelete,
		With: with,
	}
}

// ConvertStepHelmDeploy converts v0 HelmDeploy (basic) to v1 helmBasicDeployStep
func ConvertStepHelmDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepHelmDeploy)
	if !ok || sp == nil {
		return nil
	}

	with := HelmDeployWith{
		Flags:         convertHelmCommandFlags(sp.CommandFlags),
		DeployEnvVars: helmEnvVars(sp.EnvironmentVariables),
	}

	with.IgnoreFailedRelease = sp.IgnoreReleaseHistFailStatus
	with.SkipDeploySteadyCheck = sp.SkipSteadyStateCheck
	with.SkipCleanup = sp.SkipCleanup
	with.ChartTest = sp.RunChartTests

	// FEATURE GAP: v0 useUpgradeInstall has no template input (dropped).

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmDeploy,
		With: with,
	}
}

// ConvertStepHelmRollback converts v0 HelmRollback to v1 helmRollbackStep
func ConvertStepHelmRollback(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepHelmRollback)
	if !ok || sp == nil {
		return nil
	}

	with := HelmRollbackWith{
		Flags:   convertHelmCommandFlags(sp.CommandFlags),
		Envvars: helmEnvVars(sp.EnvironmentVariables),
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

	with := HelmCanaryDeleteWith{
		Flags: convertHelmCommandFlags(sp.CommandFlags),
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeHelmCanaryDelete,
		With: with,
	}
}

package yaml

import "github.com/drone/go-convert/internal/flexible"

type (

	// CD: Helm Blue Green Deploy
	StepHelmBGDeploy struct {
		CommonStepSpec
		EnvironmentVariables        map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		IgnoreReleaseHistFailStatus *flexible.Field[bool]              `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		SkipSteadyStateCheck        *flexible.Field[bool]               `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           *flexible.Field[bool]               `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
		RunChartTests               *flexible.Field[bool]               `json:"runChartTests,omitempty" yaml:"runChartTests,omitempty"` 
	}

	// CD: Helm Blue Green Swap (no spec fields in provided example)
	StepHelmBlueGreenSwapStep struct {
		CommonStepSpec
	}

	// CD: Helm Canary Deploy
	StepHelmCanaryDeploy struct {
		CommonStepSpec
		EnvironmentVariables        map[string]string      `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		IgnoreReleaseHistFailStatus *flexible.Field[bool]                    `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		SkipSteadyStateCheck        *flexible.Field[bool]                    `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           *flexible.Field[bool]                    `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
		InstanceSelection           *HelmInstanceSelection `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
		RunChartTests               *flexible.Field[bool]               `json:"runChartTests,omitempty" yaml:"runChartTests,omitempty"` 
	}

	HelmInstanceSelection struct {
		Type string                     `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *HelmInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	HelmInstanceSelectionSpec struct {
		Count      *flexible.Field[int]  `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage *flexible.Field[int]  `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	StepHelmCanaryDelete struct {
		CommonStepSpec
	}

	// CD: Helm Delete
	StepHelmDelete struct {
		CommonStepSpec
		DryRun               *flexible.Field[bool]               `json:"dryRun,omitempty" yaml:"dryRun,omitempty"`
		CommandFlags         []string          `json:"commandFlags,omitempty" yaml:"commandFlags,omitempty"`
		EnvironmentVariables map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		ReleaseName          string            `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	}

	// CD: Helm Deploy (Basic)
	StepHelmDeploy struct {
		CommonStepSpec
		SkipDryRun                  *flexible.Field[bool]               `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		IgnoreReleaseHistFailStatus *flexible.Field[bool]               `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		EnvironmentVariables        map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		SkipCleanup                 *flexible.Field[bool]               `json:"skipCleanup,omitempty" yaml:"skipCleanup,omitempty"`
		SkipSteadyStateCheck        *flexible.Field[bool]               `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           *flexible.Field[bool]               `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
		RunChartTests               *flexible.Field[bool]               `json:"runChartTests,omitempty" yaml:"runChartTests,omitempty"` 
	}

	// CD: Helm Rollback
	StepHelmRollback struct {
		CommonStepSpec
		SkipDryRun           *flexible.Field[bool]               `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		EnvironmentVariables map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		SkipSteadyStateCheck *flexible.Field[bool]               `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		RunChartTests               *flexible.Field[bool]               `json:"runChartTests,omitempty" yaml:"runChartTests,omitempty"` 
	}
)
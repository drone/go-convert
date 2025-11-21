package yaml

type (

	// CD: Helm Blue Green Deploy
	StepHelmBGDeploy struct {
		CommonStepSpec
		EnvironmentVariables        map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		IgnoreReleaseHistFailStatus bool              `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		SkipSteadyStateCheck        bool              `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           bool              `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
	}

	// CD: Helm Blue Green Swap (no spec fields in provided example)
	StepHelmBlueGreenSwapStep struct {
		CommonStepSpec
	}

	// CD: Helm Canary Deploy
	StepHelmCanaryDeploy struct {
		CommonStepSpec
		EnvironmentVariables        map[string]string      `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		IgnoreReleaseHistFailStatus bool                   `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		SkipSteadyStateCheck        bool                   `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           bool                   `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
		InstanceSelection           *HelmInstanceSelection `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
	}

	HelmInstanceSelection struct {
		Type string                     `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *HelmInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	HelmInstanceSelectionSpec struct {
		Count      int `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage int `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	// CD: Helm Delete
	StepHelmDelete struct {
		CommonStepSpec
		DryRun               bool              `json:"dryRun,omitempty" yaml:"dryRun,omitempty"`
		CommandFlags         []string          `json:"commandFlags,omitempty" yaml:"commandFlags,omitempty"`
		EnvironmentVariables map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		ReleaseName          string            `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	}

	// CD: Helm Deploy (Basic)
	StepHelmDeploy struct {
		CommonStepSpec
		SkipDryRun                  bool              `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		IgnoreReleaseHistFailStatus bool              `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		EnvironmentVariables        map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		SkipCleanup                 bool              `json:"skipCleanup,omitempty" yaml:"skipCleanup,omitempty"`
		SkipSteadyStateCheck        bool              `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           bool              `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
	}

	// CD: Helm Rollback
	StepHelmRollback struct {
		CommonStepSpec
		SkipDryRun           bool              `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		EnvironmentVariables map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		SkipSteadyStateCheck bool              `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
	}
)
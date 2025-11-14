package yaml

import "github.com/drone/go-convert/internal/flexible"

type (

	StepK8sRollingDeploy struct {
		CommonStepSpec
		SkipDryRun     *flexible.Field[bool] `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		PruningEnabled *flexible.Field[bool] `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
	}

	StepK8sRollingRollback struct {
		CommonStepSpec
		PruningEnabled *flexible.Field[bool] `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
	}

	StepK8sApply struct {
		// TODO: Handle remote manifests and ovrrides
		CommonStepSpec
		SkipDryRun           *flexible.Field[bool]          `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		SkipSteadyStateCheck *flexible.Field[bool]          `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		SkipRendering        *flexible.Field[bool]          `json:"skipRendering,omitempty" yaml:"skipRendering,omitempty"`
		// Overrides            *flexible.Field[[]interface{}] `json:"overrides,omitempty" yaml:"overrides,omitempty"`
		FilePaths            []string     `json:"filePaths,omitempty" yaml:"filePaths,omitempty"`
	}

	StepK8sBGSwapServices struct {
		CommonStepSpec
		
	}

	StepK8sBlueGreenStageScaleDown struct {
		CommonStepSpec
		DeleteResources *flexible.Field[bool] `json:"deleteResources,omitempty" yaml:"deleteResources,omitempty"`
	}

	// CD: K8s Delete
	StepK8sDelete struct {
		CommonStepSpec
		DeleteResources *K8sDeleteResources `json:"deleteResources,omitempty" yaml:"deleteResources,omitempty"`
	}

	// K8sDeleteResources captures the delete selection and its spec.
	// Type is one of: ResourceName | ManifestPath | ReleaseName
	K8sDeleteResources struct {
		Type string                  `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sDeleteResourcesSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// K8sDeleteResourcesSpec holds the possible selectors. Only one list is expected
	// to be populated depending on the Type above.
	K8sDeleteResourcesSpec struct {
		ResourceNames []string `json:"resourceNames,omitempty" yaml:"resourceNames,omitempty"`
		ManifestPaths []string `json:"manifestPaths,omitempty" yaml:"manifestPaths,omitempty"`
		ReleaseNames  []string `json:"releaseNames,omitempty" yaml:"releaseNames,omitempty"`
	}

	// CD: K8s Canary Delete 
	StepK8sCanaryDelete struct {
		CommonStepSpec
	}

	// CD: K8s Diff
	StepK8sDiff struct {
		CommonStepSpec
	}

	// CD: K8s Rollout
	StepK8sRollout struct {
		CommonStepSpec
		Command   string               `json:"command,omitempty" yaml:"command,omitempty"`
		Resources *K8sRolloutResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	}

	K8sRolloutResources struct {
		Type string                   `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sRolloutResourcesSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sRolloutResourcesSpec struct {
		ResourceNames []string `json:"resourceNames,omitempty" yaml:"resourceNames,omitempty"`
		ManifestPaths []string `json:"manifestPaths,omitempty" yaml:"manifestPaths,omitempty"`
	}


	// CD: K8s Scale
	StepK8sScale struct {
		CommonStepSpec
		InstanceSelection    *K8sScaleInstanceSelection `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
		SkipSteadyStateCheck bool                       `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		Workload             string                     `json:"workload,omitempty" yaml:"workload,omitempty"`
	}

	K8sScaleInstanceSelection struct {
		Type string                         `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sScaleInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sScaleInstanceSelectionSpec struct {
		Count      *flexible.Field[int] `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage *flexible.Field[int] `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	// CD: K8s Dry Run
	StepK8sDryRun struct {
		CommonStepSpec
		EncryptYamlOutput bool `json:"encryptYamlOutput,omitempty" yaml:"encryptYamlOutput,omitempty"`
	}

	// CD: K8s Traffic Routing
	StepK8sTrafficRouting struct {
		CommonStepSpec
		Type           string                   `json:"type,omitempty" yaml:"type,omitempty"`
		TrafficRouting *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	K8sTrafficRoutingConfig struct {
		Provider string                 `json:"provider,omitempty" yaml:"provider,omitempty"`
		Spec     *K8sTrafficRoutingSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sTrafficRoutingSpec struct {
		Name        string                    `json:"name,omitempty" yaml:"name,omitempty"`
		RootService string                    `json:"rootService,omitempty" yaml:"rootService,omitempty"`
		Hosts       interface{}               `json:"hosts,omitempty" yaml:"hosts,omitempty"`
		Gateways    interface{}               `json:"gateways,omitempty" yaml:"gateways,omitempty"`
		Routes      []*K8sTrafficRoutingRoute `json:"routes,omitempty" yaml:"routes,omitempty"`
	}

	K8sTrafficRoutingRoute struct {
		Route *K8sTrafficRoutingRouteSpec `json:"route,omitempty" yaml:"route,omitempty"`
	}

	K8sTrafficRoutingRouteSpec struct {
		Type         string                          `json:"type,omitempty" yaml:"type,omitempty"`
		Name         string                          `json:"name,omitempty" yaml:"name,omitempty"`
		Destinations []*K8sTrafficRoutingDestination `json:"destinations,omitempty" yaml:"destinations,omitempty"`
	}

	K8sTrafficRoutingDestination struct {
		Destination *K8sTrafficRoutingDestinationSpec `json:"destination,omitempty" yaml:"destination,omitempty"`
	}

	K8sTrafficRoutingDestinationSpec struct {
		Host   string `json:"host,omitempty" yaml:"host,omitempty"`
		Weight int    `json:"weight,omitempty" yaml:"weight,omitempty"`
	}

	// CD: K8s Canary Deploy
	StepK8sCanaryDeploy struct {
		CommonStepSpec
		SkipDryRun        bool                     `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		InstanceSelection *K8sInstanceSelection    `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
		TrafficRouting    *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	K8sInstanceSelection struct {
		Type string                    `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sInstanceSelectionSpec struct {
		Count      int `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage int `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	// CD: K8s Blue Green Deploy
	StepK8sBlueGreenDeploy struct {
		CommonStepSpec
		SkipDryRun            bool                     `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		PruningEnabled        bool                     `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
		SkipUnchangedManifest bool                     `json:"skipUnchangedManifest,omitempty" yaml:"skipUnchangedManifest,omitempty"`
		TrafficRouting        *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}
)
package yaml

import (
	"github.com/drone/go-convert/internal/flexible"
	"encoding/json"
	"fmt"
)

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
		// ReleaseNames  []string `json:"releaseNames,omitempty" yaml:"releaseNames,omitempty"`
		DeleteNamespace *flexible.Field[bool] `json:"deleteNamespace" yaml:"deleteNamespace"`
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
		SkipSteadyStateCheck *flexible.Field[bool]                       `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
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
		EncryptYamlOutput *flexible.Field[bool] `json:"encryptYamlOutput,omitempty" yaml:"encryptYamlOutput,omitempty"`
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
		Port   int    `json:"port,omitempty" yaml:"port,omitempty"`
	}

	// CD: K8s Canary Deploy
	StepK8sCanaryDeploy struct {
		CommonStepSpec
		SkipDryRun        *flexible.Field[bool]                     `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		InstanceSelection *K8sScaleInstanceSelection    `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
		TrafficRouting    *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	// CD: K8s Blue Green Deploy
	StepK8sBlueGreenDeploy struct {
		CommonStepSpec
		SkipDryRun            *flexible.Field[bool]                      `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		PruningEnabled        *flexible.Field[bool]                      `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
		SkipUnchangedManifest *flexible.Field[bool]                      `json:"skipUnchangedManifest,omitempty" yaml:"skipUnchangedManifest,omitempty"`
		TrafficRouting        *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	// CD: K8s Patch
	StepK8sPatch struct {
		CommonStepSpec
		Workload             string                `json:"workload,omitempty" yaml:"workload,omitempty"`
		MergeStrategy        string                `json:"mergeStrategyType,omitempty" yaml:"mergeStrategyType,omitempty"` // Json | Strategic | Merge
		Source               *RemoteSource          `json:"source,omitempty" yaml:"source,omitempty"`
		RecordChangeCause    *flexible.Field[bool] `json:"recordChangeCause,omitempty" yaml:"recordChangeCause,omitempty"`
		SkipSteadyStateCheck *flexible.Field[bool] `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
	}

	RemoteSource struct {
		Type string      `json:"type,omitempty" yaml:"type,omitempty"`
		Spec interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	}
	// source type: Inline
	SourceSpecInline struct {
		Content string `json:"content,omitempty" yaml:"content,omitempty"`
	}
	// source type: Git | GitLab | Github | Bitbucket | HarnessCode
	RemoteSourceSpecGitRepo struct {
		ConnectorRef string `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		GitFetchType string `json:"gitFetchType,omitempty" yaml:"gitFetchType,omitempty"` // Commit | Branch
		Branch       string `json:"branch,omitempty" yaml:"branch,omitempty"`
		CommitId     string `json:"commitId,omitempty" yaml:"commitId,omitempty"`
		RepoName     string `json:"repoName,omitempty" yaml:"repoName,omitempty"`
		Paths       []string `json:"paths,omitempty" yaml:"paths,omitempty"`
	}
	// source type: Harness
	RemoteSourceSpecHarness struct {
		Files []string `json:"files,omitempty" yaml:"files,omitempty"`
	}
)

func (p *RemoteSource) UnmarshalJSON(data []byte) error {
	type Alias RemoteSource
	aux := &struct {
		Type string          `json:"type"`
		Spec json.RawMessage `json:"spec"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.Type = aux.Type

	if len(aux.Spec) == 0 {
		return nil
	}

	switch aux.Type {
	case "Inline":
		var spec SourceSpecInline
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Inline spec: %w", err)
		}
		p.Spec = &spec

	case "Git", "GitLab", "Github", "Bitbucket", "HarnessCode":
		var spec RemoteSourceSpecGitRepo
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal GitRepo spec: %w", err)
		}
		p.Spec = &spec

	case "Harness":
		var spec RemoteSourceSpecHarness
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Harness spec: %w", err)
		}
		p.Spec = &spec

	default:
		return fmt.Errorf("unknown RemoteSource type: %s", aux.Type)
	}

	return nil
}

func (p *RemoteSource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Alias RemoteSource
	aux := &struct {
		Type string      `yaml:"type"`
		Spec interface{} `yaml:"spec"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}

	p.Type = aux.Type

	if aux.Spec == nil {
		return nil
	}

	switch aux.Type {
	case "Inline":
		var spec SourceSpecInline
		specBytes, err := json.Marshal(aux.Spec)
		if err != nil {
			return fmt.Errorf("failed to marshal Inline spec: %w", err)
		}
		if err := json.Unmarshal(specBytes, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Inline spec: %w", err)
		}
		p.Spec = &spec

	case "Git", "GitLab", "Github", "Bitbucket", "HarnessCode":
		var spec RemoteSourceSpecGitRepo
		specBytes, err := json.Marshal(aux.Spec)
		if err != nil {
			return fmt.Errorf("failed to marshal GitRepo spec: %w", err)
		}
		if err := json.Unmarshal(specBytes, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal GitRepo spec: %w", err)
		}
		p.Spec = &spec

	case "Harness":
		var spec RemoteSourceSpecHarness
		specBytes, err := json.Marshal(aux.Spec)
		if err != nil {
			return fmt.Errorf("failed to marshal Harness spec: %w", err)
		}
		if err := json.Unmarshal(specBytes, &spec); err != nil {
			return fmt.Errorf("failed to unmarshal Harness spec: %w", err)
		}
		p.Spec = &spec

	default:
		return fmt.Errorf("unknown RemoteSource type: %s", aux.Type)
	}

	return nil
}
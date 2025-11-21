package converthelpers

import (
	"strconv"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// K8s step with configuration structs with JSON tags
type K8sRollingDeployWith struct {
	Flags      []interface{} `json:"flags,omitempty"`
	SkipDryRun *flexible.Field[bool]          `json:"skipdryrun,omitempty"`
	Pruning    *flexible.Field[bool]          `json:"pruning,omitempty"`
}

type K8sRollingRollbackWith struct {
	Pruning *flexible.Field[bool]          `json:"pruning,omitempty"`
	Flags   []interface{} `json:"flags,omitempty"`
}

type K8sApplyWith struct {
	Manifests            []interface{} `json:"manifests,omitempty"`
	SkipDryRun           *flexible.Field[bool]          `json:"skipdryrun,omitempty"`
	SkipSteadyStateCheck *flexible.Field[bool]          `json:"skipsteadystatecheck,omitempty"`
	Flags                []interface{} `json:"flags,omitempty"`
}

type K8sBGSwapServicesWith struct {
	// Empty struct for steps with no specific configuration
}

type K8sBlueGreenStageScaleDownWith struct {
	Pruning *flexible.Field[bool] `json:"pruning,omitempty"`
}

type K8sCanaryDeleteWith struct {
	// Empty struct for steps with no specific configuration
}

type K8sDiffWith struct {
	// Empty struct for steps with no specific configuration
}

type K8sRolloutWith struct {
	Command                string        `json:"command,omitempty"`
	SelectRolloutResources string        `json:"selectrolloutresources,omitempty"`
	Flags                  []interface{} `json:"flags,omitempty"`
	Resources              []interface{} `json:"resources,omitempty"`
	Manifests              []interface{} `json:"manifests,omitempty"`
}

type K8sScaleWith struct {
	UnitType             string `json:"unittype,omitempty"`
	Instances            *flexible.Field[int] `json:"instances,omitempty"`
	Workload             string `json:"workload,omitempty"`
	SkipSteadyStateCheck bool   `json:"skipsteadystatecheck,omitempty"`
}

type K8sDryRunWith struct {
	EncryptYamlOutput bool `json:"encryptyamloutput,omitempty"`
}

type K8sDeleteWith struct {
	SelectDeleteResources string        `json:"selectdeleteresources,omitempty"`
	Resources             []interface{} `json:"resources,omitempty"`
	Manifests             []interface{} `json:"manifests,omitempty"`
	Releasename           []interface{} `json:"releasename,omitempty"`
	Flags                 []interface{} `json:"flags,omitempty"`
}

type K8sTrafficRoutingWith struct {
	Config       string        `json:"config,omitempty"`
	Provider     string        `json:"provider,omitempty"`
	Hosts        []interface{} `json:"hosts,omitempty"`
	Gateways     []interface{} `json:"gateways,omitempty"`
	Routes       string        `json:"routes,omitempty"`
	ResourceName string        `json:"resourcename,omitempty"`
}

type K8sCanaryDeployWith struct {
	Provider     string        `json:"provider,omitempty"`
	UnitType     string        `json:"unittype,omitempty"`
	Instances    string        `json:"instances,omitempty"`
	ResourceName string        `json:"resourcename,omitempty"`
	Hosts        []interface{} `json:"hosts,omitempty"`
	Gateways     []interface{} `json:"gateways,omitempty"`
	Routes       string        `json:"routes,omitempty"`
	SkipDryRun   bool          `json:"skipdryrun,omitempty"`
	Flags        []interface{} `json:"flags,omitempty"`
}

type K8sBlueGreenDeployWith struct {
	Provider              string        `json:"provider,omitempty"`
	ResourceName          string        `json:"resourcename,omitempty"`
	Hosts                 []interface{} `json:"hosts,omitempty"`
	Gateways              []interface{} `json:"gateways,omitempty"`
	Routes                string        `json:"routes,omitempty"`
	SkipDryRun            bool          `json:"skipdryrun,omitempty"`
	Pruning               bool          `json:"pruning,omitempty"`
	SkipUnchangedManifest bool          `json:"skipunchangedmanifest,omitempty"`
	Flags                 []interface{} `json:"flags,omitempty"`
}

// ConvertStepK8sRollingDeploy converts a v0 K8sRollingDeploy step to v1 template spec only
func ConvertStepK8sRollingDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sRollingDeploy)
	if !ok {
		return nil
	}
	with := K8sRollingDeployWith{
		Flags:      []interface{}{},
		SkipDryRun: spec.SkipDryRun,
		Pruning:    spec.PruningEnabled,
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sRollingDeploy,
		With: with,
	}
}

// ConvertStepK8sRollingRollback converts a v0 K8sRollingRollback step to v1 template spec only
func ConvertStepK8sRollingRollback(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sRollingRollback)
	if !ok {
		return nil
	}

	with := K8sRollingRollbackWith{
		Pruning: spec.PruningEnabled,
		Flags:   []interface{}{},
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sRollingRollback,
		With: with,
	}
}

// ConvertStepK8sApply converts a v0 K8sApply step to v1 template spec only
func ConvertStepK8sApply(src *v0.Step) *v1.StepTemplate {
	// TODO: handle overrides and remote manifests
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sApply)
	if !ok {
		return nil
	}

	// Map filePaths to manifests (list)
	manifests := make([]interface{}, 0, len(spec.FilePaths))
	for _, p := range spec.FilePaths {
		manifests = append(manifests, p)
	}

	with := K8sApplyWith{
		Manifests:            manifests,
		SkipDryRun:           spec.SkipDryRun,
		SkipSteadyStateCheck: spec.SkipSteadyStateCheck,
		Flags:                []interface{}{},
	}
	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sApply,
		With: with,
	}
}

// ConvertStepK8sBGSwapServices converts a v0 K8sBGSwapServices step to v1 template spec only
func ConvertStepK8sBGSwapServices(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}
	// Spec is empty per v0 example; we still assert type for safety when present.
	if src.Spec != nil {
		if _, ok := src.Spec.(*v0.StepK8sBGSwapServices); !ok {
			return nil
		}
	}

	with := K8sBGSwapServicesWith{}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sBGSwapServices,
		With: with,
	}
}

// ConvertStepK8sBlueGreenStageScaleDown converts a v0 K8sBlueGreenStageScaleDown step to v1 template spec only
func ConvertStepK8sBlueGreenStageScaleDown(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}
	// Typed spec (contains deleteResources)
	if spec, ok := src.Spec.(*v0.StepK8sBlueGreenStageScaleDown); ok {
		return &v1.StepTemplate{
			Uses: v1.StepTypeK8sBlueGreenStageScaleDown,
			With: K8sBlueGreenStageScaleDownWith{
				Pruning: spec.DeleteResources,
			},
		}
	} else {
		return nil
	}
}

// ConvertStepK8sCanaryDelete converts a v0 K8sCanaryDelete step to v1 template spec only
func ConvertStepK8sCanaryDelete(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}
	// assert spec type when present (spec is empty per example)
	if src.Spec != nil {
		if _, ok := src.Spec.(*v0.StepK8sCanaryDelete); !ok {
			return nil
		}
	}

	with := K8sCanaryDeleteWith{}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sCanaryDelete,
		With: with,
	}
}

// ConvertStepK8sDiff converts a v0 K8sDiff step to v1 template spec only
func ConvertStepK8sDiff(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}
	// spec is empty per example; type-assert when present
	if src.Spec != nil {
		if _, ok := src.Spec.(*v0.StepK8sDiff); !ok {
			return nil
		}
	}

	with := K8sDiffWith{}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sDiff,
		With: with,
	}
}

// ConvertStepK8sRollout converts a v0 K8sRollout step to v1 template spec only
func ConvertStepK8sRollout(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepK8sRollout)
	if !ok {
		return nil
	}

	sel := ""
	// If ReleaseName is selected, v1 wants a separate releasename input.
	// releasenameNeeded := false
	// Hold optional list outputs for resources/manifests
	var resourcesList []interface{}
	var manifestsList []interface{}
	if sp.Resources != nil {
		switch sp.Resources.Type {
		case "ResourceName":
			sel = "resourcename"
			if sp.Resources.Spec != nil {
				for _, r := range sp.Resources.Spec.ResourceNames {
					resourcesList = append(resourcesList, r)
				}
			}
		case "ManifestPath":
			sel = "manifestpath"
			if sp.Resources.Spec != nil {
				for _, m := range sp.Resources.Spec.ManifestPaths {
					manifestsList = append(manifestsList, m)
				}
			}
		case "ReleaseName":
			sel = "releasename"
			// releasenameNeeded = true
		}
	}

	with := K8sRolloutWith{
		Command:                sp.Command,
		SelectRolloutResources: sel,
		Flags:                  []interface{}{},
	}
	// if releasenameNeeded {
	// 	with["releasename"] = "<+input>"
	// }
	if len(resourcesList) > 0 {
		with.Resources = resourcesList
	}
	if len(manifestsList) > 0 {
		with.Manifests = manifestsList
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sRollout,
		With: with,
	}
}

// ConvertStepK8sScale converts a v0 K8sScale step to v1 template spec only
func ConvertStepK8sScale(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepK8sScale)
	if !ok {
		return nil
	}

	unittype := ""
	instances := &flexible.Field[int]{}
	if sel := sp.InstanceSelection; sel != nil {
		switch sel.Type {
		case "Count":
			unittype = "count"
			if sel.Spec != nil {
				instances = sel.Spec.Count
			}
		case "Percentage":
			unittype = "percentage"
			if sel.Spec != nil {
				instances = sel.Spec.Percentage
			}
		}
	}

	with := K8sScaleWith{
		UnitType:             unittype,
		Instances:            instances,
		Workload:             sp.Workload,
		SkipSteadyStateCheck: sp.SkipSteadyStateCheck,
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sScale,
		With: with,
	}
}

// ConvertStepK8sDryRun converts a v0 K8sDryRun step to v1 template spec only
func ConvertStepK8sDryRun(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepK8sDryRun)
	if !ok {
		return nil
	}

	with := K8sDryRunWith{
		EncryptYamlOutput: sp.EncryptYamlOutput,
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sDryRun,
		With: with,
	}
}

// ConvertStepK8sDelete converts a v0 K8sDelete step to a v1 template-based step.
// It supports delete by ResourceName, ManifestPath, or ReleaseName.
func ConvertStepK8sDelete(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepK8sDelete)
	if !ok {
		return nil
	}

	sel := ""
	var items []string
	if sp.DeleteResources != nil {
		switch sp.DeleteResources.Type {
		case "ResourceName":
			sel = "resources"
			if sp.DeleteResources.Spec != nil {
				items = sp.DeleteResources.Spec.ResourceNames
			}
		case "ManifestPath":
			sel = "manifests"
			if sp.DeleteResources.Spec != nil {
				items = sp.DeleteResources.Spec.ManifestPaths
			}
		case "ReleaseName":
			sel = "releasename"
			if sp.DeleteResources.Spec != nil {
				items = sp.DeleteResources.Spec.ReleaseNames
			}
		}
	}

	// cast items to []interface{} for generic map
	resources := make([]interface{}, 0, len(items))
	for _, it := range items {
		resources = append(resources, it)
	}

	with := K8sDeleteWith{
		SelectDeleteResources: sel,
		Flags:                 []interface{}{},
	}

	// Set the appropriate field based on selection type
	switch sel {
	case "resources":
		with.Resources = resources
	case "manifests":
		with.Manifests = resources
	case "releasename":
		with.Releasename = resources
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sDelete,
		With: with,
	}
}

// ConvertStepK8sTrafficRouting converts a v0 K8sTrafficRouting step to v1 template spec only
func ConvertStepK8sTrafficRouting(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sTrafficRouting)
	if !ok {
		return nil
	}

	// Default values
	hosts := []interface{}{}
	gateways := []interface{}{}
	provider := ""
	resourceName := ""
	routes := "[]"

	// Extract traffic routing configuration
	if spec.TrafficRouting != nil {
		provider = spec.TrafficRouting.Provider

		if spec.TrafficRouting.Spec != nil {
			routingSpec := spec.TrafficRouting.Spec
			resourceName = routingSpec.Name

			// Handle hosts - can be string, array, or <+input>
			if routingSpec.Hosts != nil {
				hosts = []interface{}{}
			}

			// Handle gateways - can be string, array, or <+input>
			if routingSpec.Gateways != nil {
				gateways = []interface{}{}
			}

			// Convert routes using the reusable function
			if len(routingSpec.Routes) > 0 {
				routes = ConvertTrafficRoutingRoutes(routingSpec.Routes)
			}
		}
	}

	with := K8sTrafficRoutingWith{
		Config:       "new",
		Provider:     provider,
		Hosts:        hosts,
		Gateways:     gateways,
		Routes:       routes,
		ResourceName: resourceName,
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sTrafficRouting,
		With: with,
	}
}

// ConvertTrafficRoutingRoutes converts v0 traffic routing routes to v1 JSON string format
// This is a reusable function for converting route configurations
func ConvertTrafficRoutingRoutes(routes []*v0.K8sTrafficRoutingRoute) string {
	if len(routes) == 0 {
		return "[]"
	}

	// Build the routes array for JSON serialization
	var routesArray []map[string]interface{}
	for _, route := range routes {
		if route == nil || route.Route == nil {
			continue
		}

		routeSpec := route.Route
		routeMap := map[string]interface{}{
			"name": routeSpec.Name,
		}

		// Convert destinations
		if len(routeSpec.Destinations) > 0 {
			var destinations []map[string]interface{}
			for _, dest := range routeSpec.Destinations {
				if dest == nil || dest.Destination == nil {
					continue
				}
				destMap := map[string]interface{}{
					"host":   dest.Destination.Host,
					"weight": dest.Destination.Weight,
				}
				destinations = append(destinations, destMap)
			}
			routeMap["destinations"] = destinations
		}

		routesArray = append(routesArray, routeMap)
	}

	// Convert to JSON string format as expected by v1
	// For simplicity, we'll build the JSON string manually since the format is predictable
	if len(routesArray) == 0 {
		return "[]"
	}

	// Build JSON string manually for the expected format
	jsonStr := "["
	for i, route := range routesArray {
		if i > 0 {
			jsonStr += ","
		}
		jsonStr += `{"name":"` + route["name"].(string) + `"`
		if destinations, ok := route["destinations"].([]map[string]interface{}); ok && len(destinations) > 0 {
			jsonStr += `,"destinations":[`
			for j, dest := range destinations {
				if j > 0 {
					jsonStr += ","
				}
				jsonStr += `{"host":"` + dest["host"].(string) + `","weight":` + strconv.Itoa(dest["weight"].(int)) + `}`
			}
			jsonStr += `]`
		}
		jsonStr += `}`
	}
	jsonStr += "]"

	return jsonStr
}

// ConvertStepK8sCanaryDeploy converts a v0 K8sCanaryDeploy step to v1 template spec only
func ConvertStepK8sCanaryDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sCanaryDeploy)
	if !ok {
		return nil
	}

	// Default values
	hosts := []interface{}{}
	gateways := []interface{}{}
	provider := ""
	resourceName := ""
	routes := "[]"
	unitType := ""
	instances := ""

	// Extract instance selection
	if spec.InstanceSelection != nil {
		switch spec.InstanceSelection.Type {
		case "Count":
			unitType = "count"
			if spec.InstanceSelection.Spec != nil {
				instances = strconv.Itoa(spec.InstanceSelection.Spec.Count)
			}
		case "Percentage":
			unitType = "percentage"
			if spec.InstanceSelection.Spec != nil {
				instances = strconv.Itoa(spec.InstanceSelection.Spec.Percentage)
			}
		}
	}

	// Extract traffic routing configuration (reusing logic from K8sTrafficRouting)
	if spec.TrafficRouting != nil {
		provider = spec.TrafficRouting.Provider

		if spec.TrafficRouting.Spec != nil {
			routingSpec := spec.TrafficRouting.Spec
			resourceName = routingSpec.Name

			// Convert routes using the reusable function
			if len(routingSpec.Routes) > 0 {
				routes = ConvertTrafficRoutingRoutes(routingSpec.Routes)
			}
		}
	}

	with := K8sCanaryDeployWith{
		Provider:     provider,
		UnitType:     unitType,
		Instances:    instances,
		ResourceName: resourceName,
		Hosts:        hosts,
		Gateways:     gateways,
		Routes:       routes,
		SkipDryRun:   spec.SkipDryRun,
		Flags:        []interface{}{},
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sCanaryDeploy,
		With: with,
	}
}

// ConvertStepK8sBlueGreenDeploy converts a v0 K8sBlueGreenDeploy step to v1 template spec only
func ConvertStepK8sBlueGreenDeploy(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sBlueGreenDeploy)
	if !ok {
		return nil
	}

	// Default values for simple blue-green deploy
	hosts := []interface{}{}
	gateways := []interface{}{}
	provider := ""
	resourceName := ""
	routes := "[]"

	// Check if traffic routing is configured
	if spec.TrafficRouting != nil {
		provider = spec.TrafficRouting.Provider

		if spec.TrafficRouting.Spec != nil {
			routingSpec := spec.TrafficRouting.Spec
			resourceName = routingSpec.Name

			if routingSpec.Hosts != nil {
				if hostList, ok := routingSpec.Hosts.([]interface{}); ok {
					hosts = hostList
				}
			}

			if routingSpec.Gateways != nil {
				if gatewayList, ok := routingSpec.Gateways.([]interface{}); ok {
					gateways = gatewayList
				}
			}

			if len(routingSpec.Routes) > 0 {
				routes = ConvertTrafficRoutingRoutes(routingSpec.Routes)
			}
		}
	}

	with := K8sBlueGreenDeployWith{
		Provider:              provider,
		ResourceName:          resourceName,
		Hosts:                 hosts,
		Gateways:              gateways,
		Routes:                routes,
		SkipDryRun:            spec.SkipDryRun,
		Pruning:               spec.PruningEnabled,
		SkipUnchangedManifest: spec.SkipUnchangedManifest,
		Flags:                 []interface{}{},
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sBlueGreenDeploy,
		With: with,
	}
}

package converthelpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// K8s step with configuration structs with JSON tags
type K8sRollingDeployWith struct {
	Flags      []interface{}         `json:"flags,omitempty"`
	SkipDryRun *flexible.Field[bool] `json:"skip_dry_run,omitempty"`
	Pruning    *flexible.Field[bool] `json:"pruning,omitempty"`
}

type K8sRollingRollbackWith struct {
	Pruning *flexible.Field[bool] `json:"pruning,omitempty"`
	Flags   []interface{}         `json:"flags,omitempty"`
}

type K8sApplyWith struct {
	Manifests            []interface{}         `json:"manifests,omitempty"`
	SkipDryRun           *flexible.Field[bool] `json:"skip_dry_run,omitempty"`
	SkipSteadyStateCheck *flexible.Field[bool] `json:"skip_steady_state_check,omitempty"`
	Flags                []interface{}         `json:"flags,omitempty"`
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
        	Flags []interface{} `json:"flags,omitempty"`
}

type K8sRolloutWith struct {
	Command                string        `json:"command,omitempty"`
	SelectRolloutResources string        `json:"select_rollout_resources,omitempty"`
	Flags                  []interface{} `json:"flags,omitempty"`
	Resources              []interface{} `json:"resources,omitempty"`
	Manifests              []interface{} `json:"manifests,omitempty"`
}

type K8sScaleWith struct {
	UnitType             string                `json:"unit,omitempty"`
	Instances            string                `json:"instances,omitempty"`
	Workload             string                `json:"workload,omitempty"`
	SkipSteadyStateCheck *flexible.Field[bool] `json:"skip_steady_state_check,omitempty"`
}

type K8sDryRunWith struct {
	EncryptYamlOutput *flexible.Field[bool] `json:"encrypt_yaml_output,omitempty"`
	Flags             []interface{}         `json:"flags,omitempty"`
}

type K8sDeleteWith struct {
	SelectDeleteResources string                `json:"select_delete_resources,omitempty"`
	Resources             []string              `json:"resources,omitempty"`
	Manifests             []string              `json:"manifests,omitempty"`
	Releasename           string                `json:"release,omitempty"`
	Flags                 []interface{}         `json:"flags,omitempty"`
	DeleteNamespaces      *flexible.Field[bool] `json:"delete_namespaces,omitempty"`
}

type K8sTrafficRoutingWith struct {
	Config       string                    `json:"config_type,omitempty"`
	Provider     string                    `json:"provider,omitempty"`
	Hosts        *flexible.Field[[]string] `json:"hosts,omitempty"`
	Gateways     *flexible.Field[[]string] `json:"gateways,omitempty"`
	Routes       string                    `json:"routes,omitempty"`
	ResourceName string                    `json:"resource_name,omitempty"`
}

type K8sCanaryDeployWith struct {
	Provider     string                    `json:"provider,omitempty"`
	UnitType     string                    `json:"unit,omitempty"`
	Instances    string                    `json:"instances,omitempty"`
	ResourceName string                    `json:"resource_name,omitempty"`
	Hosts        *flexible.Field[[]string] `json:"hosts,omitempty"`
	Gateways     *flexible.Field[[]string] `json:"gateways,omitempty"`
	Routes       string                    `json:"routes,omitempty"`
	SkipDryRun   *flexible.Field[bool]     `json:"skip_dry_run,omitempty"`
	TrafficShift bool                      `json:"traffic_shift,omitempty"`
	Flags        []interface{}             `json:"flags,omitempty"`
}

type K8sBlueGreenDeployWith struct {
	Provider              string                    `json:"provider,omitempty"`
	ResourceName          string                    `json:"resource_name,omitempty"`
	Hosts                 *flexible.Field[[]string] `json:"hosts,omitempty"`
	Gateways              *flexible.Field[[]string] `json:"gateways,omitempty"`
	Routes                string                    `json:"routes,omitempty"`
	SkipDryRun            *flexible.Field[bool]     `json:"skip_dry_run,omitempty"`
	Pruning               *flexible.Field[bool]     `json:"pruning,omitempty"`
	SkipUnchangedManifest *flexible.Field[bool]     `json:"skip_unchanged_manifest,omitempty"`
	TrafficShift          bool                      `json:"traffic_shift,omitempty"`
	Flags                 []interface{}             `json:"flags,omitempty"`
}

type K8sPatchWith struct {
	Workload             string                `json:"workload,omitempty" yaml:"workload,omitempty"`
	SkipSteadyStateCheck *flexible.Field[bool] `json:"skip_steady_state_check,omitempty" yaml:"skip_steady_state_check,omitempty"`
	Content              string                `json:"content,omitempty" yaml:"content,omitempty"`
	MergeStrategy        string                `json:"strategy,omitempty" yaml:"strategy,omitempty"` // merge | strategic | json
	Flags                []interface{}         `json:"flags,omitempty" yaml:"flags,omitempty"`
}

// convertK8sCommandFlags maps v0 commandFlags (commandType + flag) to the
// template's `flags` input shape (a list of {command, flag} objects). It always
// returns a non-nil slice so the `flags` key is emitted consistently.
func convertK8sCommandFlags(flags []*v0.K8sStepCommandFlag) []interface{} {
	out := make([]interface{}, 0, len(flags))
	for _, cf := range flags {
		if cf == nil {
			continue
		}
		out = append(out, map[string]string{
			"command": cf.CommandType,
			"flag":    cf.Flag,
		})
	}
	return out
}

// instancesToString renders a v0 *flexible.Field[int] instance count/percentage
// as a string, matching the template's string-typed `instances` input. Harness
// expressions are passed through unchanged; integer values are formatted.
func instancesToString(f *flexible.Field[int]) string {
	if f == nil {
		return ""
	}
	if expr, ok := f.AsString(); ok {
		return expr
	}
	if v, ok := f.AsStruct(); ok {
		return strconv.Itoa(v)
	}
	return ""
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
		Flags:      convertK8sCommandFlags(spec.CommandFlags),
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
		Flags:   convertK8sCommandFlags(spec.CommandFlags),
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sRollingRollback,
		With: with,
	}
}

// ConvertStepK8sApply converts a v0 K8sApply step to v1 template spec only
func ConvertStepK8sApply(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sApply)
	if !ok {
		return nil
	}

	// Remote manifest sources and overrides are not supported yet; flag them so
	// the user knows the conversion is incomplete.
	if spec.ManifestSource != nil || spec.Overrides != nil {
		messagelog.GetMessageLogger().LogError(
			"UNSUPPORTED_K8S_APPLY_REMOTE_SOURCE",
			"remote manifest sources and overrides in K8s Apply step are not supported by the converter; only inline filePaths are converted",
			messagelog.WithStep(src.ID, src.Type),
		)
	}

	// Map filePaths to manifests (list)
	manifests := make([]interface{}, 0, len(spec.FilePaths))
	for _, p := range spec.FilePaths {
		p = "<+runtime.manifestPath>/" + p
		manifests = append(manifests, p)
	}

	// spec.SkipRendering is a feature gap: no template input exists for it, so it
	// is intentionally left unmapped.
	with := K8sApplyWith{
		Manifests:            manifests,
		SkipDryRun:           spec.SkipDryRun,
		SkipSteadyStateCheck: spec.SkipSteadyStateCheck,
		Flags:                convertK8sCommandFlags(spec.CommandFlags),
	}
	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sApply,
		With: with,
	}
}

// ConvertStepK8sBGSwapServices converts a v0 K8sBGSwapServices step to v1 template spec only
func ConvertStepK8sBGSwapServices(src *v0.Step, isRollback bool) *v1.StepTemplate {
	if src == nil {
		return nil
	}
	// Spec is empty per v0 example; we still assert type for safety when present.
	if src.Spec != nil {
		if _, ok := src.Spec.(*v0.StepK8sBGSwapServices); !ok {
			return nil
		}
	}

	with := map[string]interface{}{}

	if !isRollback {
		with["stable_service"] = "<+exportedVariables.getValue(\"stage.bluegreenprepareactionoutput.PLUGIN_STABLE_SERVICE\")>"
		with["stage_service"] = "<+exportedVariables.getValue(\"stage.bluegreenprepareactionoutput.PLUGIN_STAGE_SERVICE\")>"
		with["is_openshift"] = "<+exportedVariables.getValue(\"stage.bluegreenapplyactionoutput.HARNESS_IS_OPENSHIFT\")>"
	} else {
		with["stable_service"] = "${{rollback.data.PLUGIN_STABLE_SERVICE}}"
		with["stage_service"] = "${{rollback.data.PLUGIN_STAGE_SERVICE}}"
		with["is_openshift"] = "${{rollback.data.HARNESS_IS_OPENSHIFT}}"
	}

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
func ConvertStepK8sCanaryDelete(src *v0.Step, isRollback bool) *v1.StepTemplate {
	if src == nil {
		return nil
	}
	// assert spec type when present (spec is empty per example)
	if src.Spec != nil {
		if _, ok := src.Spec.(*v0.StepK8sCanaryDelete); !ok {
			return nil
		}
	}

	with := map[string]interface{}{}

	if !isRollback {
		with["resources"] = "<+exportedVariables.getValue(\"stage.canaryprepareactionoutput.PLUGIN_CANARY_WORKLOADS\")>"
		with["is_openshift"] = "<+exportedVariables.getValue(\"stage.canaryapplyactionoutput.HARNESS_IS_OPENSHIFT\")>"
	} else {
		with["resources"] = "${{rollback.data.PLUGIN_CANARY_WORKLOADS}}"
		with["is_openshift"] = "${{rollback.data.HARNESS_IS_OPENSHIFT}}"
	}

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

	with := K8sDiffWith{Flags: []interface{}{}}

	// type-assert when present so we can map command flags
	if src.Spec != nil {
		spec, ok := src.Spec.(*v0.StepK8sDiff)
		if !ok {
			return nil
		}
		with.Flags = convertK8sCommandFlags(spec.CommandFlags)
	}

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
	// Hold optional list outputs for resources/manifests
	var resourcesList []interface{}
	var manifestsList []interface{}
	if sp.Resources != nil {
		switch sp.Resources.Type {
		case "ResourceName":
			sel = "resources"
			if sp.Resources.Spec != nil {
				for _, r := range sp.Resources.Spec.ResourceNames {
					resourcesList = append(resourcesList, r)
				}
			}
		case "ManifestPath":
			sel = "manifests"
			if sp.Resources.Spec != nil {
				for _, m := range sp.Resources.Spec.ManifestPaths {
					manifestsList = append(manifestsList, "<+runtime.manifestPath>/"+m)
				}
			}
		case "ReleaseName":
			sel = "release name"
		}
	}

	with := K8sRolloutWith{
		Command:                sp.Command,
		SelectRolloutResources: sel,
		Flags:                  convertK8sCommandFlags(sp.CommandFlags),
	}

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
	var instances *flexible.Field[int]
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
		Instances:            instancesToString(instances),
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
		Flags:             convertK8sCommandFlags(sp.CommandFlags),
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
	var deleteNamespace *flexible.Field[bool]
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
				for _, manifest_path := range sp.DeleteResources.Spec.ManifestPaths {
					items = append(items, "<+runtime.manifestPath>/"+manifest_path)
				}
			}
		case "ReleaseName":
			sel = "release name"
			if sp.DeleteResources.Spec != nil {
				// items = sp.DeleteResources.Spec.ReleaseNames
				deleteNamespace = sp.DeleteResources.Spec.DeleteNamespace
			}
		}
	}

	with := K8sDeleteWith{
		SelectDeleteResources: sel,
		DeleteNamespaces:      deleteNamespace,
		Flags:                 convertK8sCommandFlags(sp.CommandFlags),
	}

	// Set the appropriate field based on selection type
	switch sel {
	case "resources":
		with.Resources = items
	case "manifests":
		with.Manifests = items
	case "release name":
		with.Releasename = "${{infra.releaseName}}"
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
	var hosts *flexible.Field[[]string]
	var gateways *flexible.Field[[]string]
	provider := ""
	resourceName := ""
	routes := ""

	// Map the v0 config type (inherit|config) to the plugin config_type
	// (update|new): a fresh "config" creates the resource, while "inherit"
	// patches weights on an existing one.
	config := mapTrafficConfigType(spec.Type)

	// Extract traffic routing configuration
	if spec.TrafficRouting != nil {
		provider = mapTrafficProvider(spec.TrafficRouting.Provider, src)

		if spec.TrafficRouting.Spec != nil {
			routingSpec := spec.TrafficRouting.Spec
			resourceName = routingSpec.Name
			hosts = routingSpec.Hosts
			gateways = routingSpec.Gateways
			routes = trafficRoutesToString(routingSpec.Routes)
		}

		// "inherit" config type carries routes directly under trafficRouting.
		if routes == "" {
			routes = trafficRoutesToString(spec.TrafficRouting.Routes)
		}
	}

	with := K8sTrafficRoutingWith{
		Config:       config,
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

// mapTrafficConfigType maps the v0 traffic routing config type (inherit|config)
// to the plugin's config_type (update|new).
func mapTrafficConfigType(configType string) string {
	if strings.EqualFold(strings.TrimSpace(configType), "inherit") {
		return "update"
	}
	return "new"
}

// mapTrafficProvider maps the v0 provider (istio|smi) to the value accepted by
// the k8s-traffic-shift plugin (istio|k8s-native). SMI is not supported: it is
// reported as an error and not emitted. Unknown values (including expressions)
// pass through.
//
// Feature gaps: the SMI-only `rootService` and the Istio-only `delegateService`
// fields on K8sTrafficRoutingSpec have no template input and are not converted.
func mapTrafficProvider(provider string, src *v0.Step) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "smi":
		if src != nil {
			messagelog.GetMessageLogger().LogError(
				"UNSUPPORTED_K8S_TRAFFIC_ROUTING_PROVIDER",
				"traffic routing provider \"smi\" is not supported by the k8s-traffic-shift plugin",
				messagelog.WithStep(src.ID, src.Type),
			)
		}
		return ""
	case "istio":
		return "istio"
	default:
		return provider
	}
}

// trafficRoutesToString converts a flexible routes field (struct list or
// expression string) into the plugin's PLUGIN_ROUTES JSON string.
func trafficRoutesToString(routes *flexible.Field[[]*v0.K8sTrafficRoutingRoute]) string {
	if routes == nil {
		return ""
	}
	if routesList, ok := routes.AsStruct(); ok && len(routesList) > 0 {
		return ConvertTrafficRoutingRoutes(routesList)
	}
	if routeExpr, ok := routes.AsString(); ok {
		return routeExpr
	}
	return ""
}

// The following structs mirror the JSON contract consumed by the
// k8s-traffic-shift plugin (PLUGIN_ROUTES -> []RouteImpl). They are the
// serialization target of ConvertTrafficRoutingRoutes.
type trafficRouteJSON struct {
	Type         string                   `json:"type,omitempty"`
	Name         string                   `json:"name,omitempty"`
	Matches      []trafficMatchJSON       `json:"matches,omitempty"`
	Filters      []trafficFilterJSON      `json:"filters,omitempty"`
	Destinations []trafficDestinationJSON `json:"destinations,omitempty"`
}

type trafficDestinationJSON struct {
	Host   string `json:"host,omitempty"`
	Port   *int   `json:"port,omitempty"`
	Weight *int   `json:"weight,omitempty"`
}

type trafficMatchJSON struct {
	Path      *trafficMatchValueJSON   `json:"path,omitempty"`
	Scheme    *trafficMatchValueJSON   `json:"scheme,omitempty"`
	Method    *trafficMatchValueJSON   `json:"method,omitempty"`
	Authority *trafficMatchValueJSON   `json:"authority,omitempty"`
	Headers   []trafficMatchHeaderJSON `json:"headers,omitempty"`
	Port      *int                     `json:"port,omitempty"`
}

type trafficMatchValueJSON struct {
	MatchType string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
}

type trafficMatchHeaderJSON struct {
	MatchType string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	Value     string `json:"value,omitempty"`
}

type trafficFilterJSON struct {
	URLRewrite *trafficURLRewriteJSON `json:"url-rewrite,omitempty"`
}

type trafficURLRewriteJSON struct {
	Hostname *string `json:"hostname,omitempty"`
	Path     *string `json:"path,omitempty"`
}

// ConvertTrafficRoutingRoutes converts v0 traffic routing routes to the JSON
// string format consumed by the k8s-traffic-shift plugin (PLUGIN_ROUTES). It
// emits each route's destinations, match conditions (from v0 rules) and URL
// rewrite filter (from rewriteRule). This is a reusable function shared by the
// standalone, canary and blue-green converters.
func ConvertTrafficRoutingRoutes(routes []*v0.K8sTrafficRoutingRoute) string {
	if len(routes) == 0 {
		return ""
	}

	var out []trafficRouteJSON
	for _, route := range routes {
		if route == nil || route.Route == nil {
			continue
		}
		routeSpec := route.Route

		rj := trafficRouteJSON{
			Type: "http",
			Name: routeSpec.Name,
		}

		for _, dest := range routeSpec.Destinations {
			if dest == nil || dest.Destination == nil {
				continue
			}
			rj.Destinations = append(rj.Destinations, trafficDestinationJSON{
				Host:   dest.Destination.Host,
				Port:   dest.Destination.Port,
				Weight: dest.Destination.Weight,
			})
		}

		rj.Matches = convertTrafficRoutingMatches(routeSpec.Rules, routeSpec.MatchAllRules)

		if routeSpec.RewriteRule != "" {
			path := routeSpec.RewriteRule
			rj.Filters = append(rj.Filters, trafficFilterJSON{
				URLRewrite: &trafficURLRewriteJSON{Path: &path},
			})
		}

		out = append(out, rj)
	}

	if len(out) == 0 {
		return ""
	}

	encoded, err := json.Marshal(out)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// convertTrafficRoutingMatches translates v0 route rules into the plugin's
// match conditions. When matchAllRules is true all rules are ANDed into a
// single match; otherwise each rule becomes its own (ORed) match entry.
func convertTrafficRoutingMatches(rules []*v0.K8sTrafficRoutingRule, matchAllRules bool) []trafficMatchJSON {
	if len(rules) == 0 {
		return nil
	}

	apply := func(m *trafficMatchJSON, rule *v0.K8sTrafficRoutingRule) {
		if rule == nil || rule.Rule == nil || rule.Rule.Spec == nil {
			return
		}
		spec := rule.Rule.Spec
		matchType := normalizeMatchType(spec.MatchType)
		switch strings.ToLower(rule.Rule.Type) {
		case "uri":
			m.Path = &trafficMatchValueJSON{MatchType: matchType, Value: trafficValueToString(spec.Value)}
		case "scheme":
			m.Scheme = &trafficMatchValueJSON{MatchType: matchType, Value: trafficValueToString(spec.Value)}
		case "method":
			m.Method = &trafficMatchValueJSON{MatchType: matchType, Value: trafficValueToString(spec.Value)}
		case "authority":
			m.Authority = &trafficMatchValueJSON{MatchType: matchType, Value: trafficValueToString(spec.Value)}
		case "port":
			if port, ok := trafficValueToInt(spec.Value); ok {
				m.Port = &port
			}
		case "headers":
			for _, header := range spec.Values {
				if header == nil {
					continue
				}
				m.Headers = append(m.Headers, trafficMatchHeaderJSON{
					MatchType: normalizeMatchType(header.MatchType),
					Name:      header.Key,
					Value:     header.Value,
				})
			}
		}
	}

	if matchAllRules {
		var match trafficMatchJSON
		for _, rule := range rules {
			apply(&match, rule)
		}
		return []trafficMatchJSON{match}
	}

	var matches []trafficMatchJSON
	for _, rule := range rules {
		var match trafficMatchJSON
		apply(&match, rule)
		matches = append(matches, match)
	}
	return matches
}

// normalizeMatchType lowercases the v0 match type to the plugin's expected
// values (exact|prefix|regex), defaulting to "exact" to match harness-core.
func normalizeMatchType(matchType string) string {
	if strings.TrimSpace(matchType) == "" {
		return "exact"
	}
	return strings.ToLower(matchType)
}

// trafficValueToString renders a rule value (string or numeric) as a string.
func trafficValueToString(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// trafficValueToInt extracts an int from a rule value (used for the port rule).
func trafficValueToInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n, true
		}
	}
	return 0, false
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
	var hosts *flexible.Field[[]string]
	var gateways *flexible.Field[[]string]
	provider := ""
	resourceName := ""
	routes := ""
	unitType := ""
	var instances *flexible.Field[int]

	// Extract instance selection
	if spec.InstanceSelection != nil {
		switch spec.InstanceSelection.Type {
		case "Count":
			unitType = "count"
			if spec.InstanceSelection.Spec != nil {
				instances = spec.InstanceSelection.Spec.Count
			}
		case "Percentage":
			unitType = "percentage"
			if spec.InstanceSelection.Spec != nil {
				instances = spec.InstanceSelection.Spec.Percentage
			}
		}
	}

	// Extract traffic routing configuration (reusing logic from K8sTrafficRouting)
	trafficShift := false
	if spec.TrafficRouting != nil {
		trafficShift = true
		provider = mapTrafficProvider(spec.TrafficRouting.Provider, src)

		if spec.TrafficRouting.Spec != nil {
			routingSpec := spec.TrafficRouting.Spec
			resourceName = routingSpec.Name
			hosts = routingSpec.Hosts
			gateways = routingSpec.Gateways
			routes = trafficRoutesToString(routingSpec.Routes)
		}

		// "inherit" config type carries routes directly under trafficRouting.
		if routes == "" {
			routes = trafficRoutesToString(spec.TrafficRouting.Routes)
		}
	}

	with := K8sCanaryDeployWith{
		Provider:     provider,
		UnitType:     unitType,
		Instances:    instancesToString(instances),
		ResourceName: resourceName,
		Hosts:        hosts,
		Gateways:     gateways,
		Routes:       routes,
		SkipDryRun:   spec.SkipDryRun,
		TrafficShift: trafficShift,
		Flags:        convertK8sCommandFlags(spec.CommandFlags),
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
	var hosts *flexible.Field[[]string]
	var gateways *flexible.Field[[]string]
	provider := ""
	resourceName := ""
	routes := ""

	// Check if traffic routing is configured
	trafficShift := false
	if spec.TrafficRouting != nil {
		trafficShift = true
		provider = mapTrafficProvider(spec.TrafficRouting.Provider, src)

		if spec.TrafficRouting.Spec != nil {
			routingSpec := spec.TrafficRouting.Spec
			resourceName = routingSpec.Name
			hosts = routingSpec.Hosts
			gateways = routingSpec.Gateways
			routes = trafficRoutesToString(routingSpec.Routes)
		}

		// "inherit" config type carries routes directly under trafficRouting.
		if routes == "" {
			routes = trafficRoutesToString(spec.TrafficRouting.Routes)
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
		TrafficShift:          trafficShift,
		Flags:                 convertK8sCommandFlags(spec.CommandFlags),
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sBlueGreenDeploy,
		With: with,
	}
}

// ConvertStepK8sPatch converts a v0 K8sPatch step to v1 template spec only
func ConvertStepK8sPatch(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}
	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepK8sPatch)
	if !ok {
		return nil
	}

	with := K8sPatchWith{
		Workload:             spec.Workload,
		SkipSteadyStateCheck: spec.SkipSteadyStateCheck,
		MergeStrategy:        strings.ToLower(spec.MergeStrategy),
		Flags:                convertK8sCommandFlags(spec.CommandFlags),
	}
	if spec.Source != nil {
		switch spec.Source.Type {
		case "Inline":
			if source, ok := spec.Source.Spec.(*v0.SourceSpecInline); ok {
				with.Content = source.Content
			}
		default:
			// Non-inline (remote/git/Harness) patch sources are not supported; the
			// other source types are ignored.
			messagelog.GetMessageLogger().LogError(
				"UNSUPPORTED_K8S_PATCH_SOURCE",
				fmt.Sprintf("conversion for Source type %q in K8s Patch step is not supported", spec.Source.Type),
				messagelog.WithStep(src.ID, src.Type),
			)
		}
	}

	return &v1.StepTemplate{
		Uses: v1.StepTypeK8sPatch,
		With: with,
	}

}

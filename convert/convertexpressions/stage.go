package convertexpressions

var StageConversionRules = []ConversionRule{
	{"identifier", "id"},
}

var DeploymentStageSpecConversionRules = []ConversionRule{}

// Deployment stage fields service/manifests/configFiles/infra/env/artifacts move
// under stage.steps in v1. The field rule sets below are node-relative and are
// attached to BOTH the spec-child node (v1Name "steps.<field>") and the bare
// alias node (v1Name "<field>"), so full/spec-prefixed paths gain the "steps"
// prefix while bare relative forms stay at root. See pipeline_trie.go.
var EnvFieldRules = []ConversionRule{
	{"identifier", "id"},
	{"envGroupName", "group.name"},
	{"envGroupRef", "group.id"},
}

var ServiceFieldRules = []ConversionRule{
	{"identifier", "id"},
	{"serviceInputs", "(with.overlay)"},
}

var InfraFieldRules = []ConversionRule{
	{"connectorRef", "connector"},
	{"infraInputs", "(with.overlay)"},
}


// manifests / configFiles / artifacts have no field renames; the spec-child node
// alone injects the "steps" prefix and the bare alias passes through at root.

// TemplateFieldRules rename a template reference's inputs to the v1 overlay form:
// template.templateInputs -> template.with.overlay. Attached to a "template" child
// at pipeline/stage/stepGroup/step level plus a standalone "template" alias, so
// both FQN (pipeline.stages.S...steps.X.template.templateInputs) and bare
// relative (template.templateInputs) forms convert. See pipeline_trie.go.
var TemplateFieldRules = []ConversionRule{
	{"templateInputs", "(with.overlay)"},
}

var CIStageSpecConversionRules = []ConversionRule{
	// Infrastructure
	// KubernetesDirect -> runtime.kubernetes
	{"infrastructure.spec.connectorRef", "runtime.kubernetes.connector"},
	{"infrastructure.spec.namespace", "runtime.kubernetes.namespace"},
	{"infrastructure.spec.annotations", "runtime.kubernetes.annotations"},
	{"infrastructure.spec.labels", "runtime.kubernetes.labels"},
	{"infrastructure.spec.serviceAccountName", "runtime.kubernetes.service-account"},
	{"infrastructure.spec.initTimeout", "runtime.kubernetes.timeout"},
	{"infrastructure.spec.nodeSelector", "runtime.kubernetes.node"},
	{"infrastructure.spec.hostNames", "runtime.kubernetes.host"},
	{"infrastructure.spec.tolerations", "runtime.kubernetes.tolerations"},
	{"infrastructure.spec.automountServiceAccountToken", "runtime.kubernetes.automount-service-token"},
	{"infrastructure.spec.containerSecurityContext", "runtime.kubernetes.security-context"},
	{"infrastructure.spec.containerSecurityContext.allowPrivilegeEscalation", "runtime.kubernetes.security-context.allow-privilege-escalation"},
	{"infrastructure.spec.containerSecurityContext.procMount", "runtime.kubernetes.security-context.proc-mount"},
	{"infrastructure.spec.containerSecurityContext.privileged", "runtime.kubernetes.security-context.privileged"},
	{"infrastructure.spec.containerSecurityContext.readOnlyRootFilesystem", "runtime.kubernetes.security-context.read-only-root-file-system"},
	{"infrastructure.spec.containerSecurityContext.runAsNonRoot", "runtime.kubernetes.security-context.run-as-non-root"},
	{"infrastructure.spec.containerSecurityContext.runAsGroup", "runtime.kubernetes.security-context.run-as-group"},
	{"infrastructure.spec.containerSecurityContext.runAsUser", "runtime.kubernetes.security-context.user"},
	{"infrastructure.spec.containerSecurityContext.capabilities", "runtime.kubernetes.security-context.capabilities"},
	{"infrastructure.spec.containerSecurityContext.capabilities.add", "runtime.kubernetes.security-context.capabilities.add"},
	{"infrastructure.spec.containerSecurityContext.capabilities.drop", "runtime.kubernetes.security-context.capabilities.drop"},
	{"infrastructure.spec.priorityClassName", "runtime.kubernetes.priority-class"},
	{"infrastructure.spec.os", "runtime.kubernetes.os"},
	{"infrastructure.spec.harnessImageConnectorRef", "runtime.kubernetes.harness-image-connector"},
	{"infrastructure.spec.imagePullPolicy", "runtime.kubernetes.pull"},
	{"infrastructure.spec.podSpecOverlay", "runtime.kubernetes.pod-spec-overlay"},
	{"infrastructure.spec.runAsUser", "runtime.kubernetes.user"},
	{"infrastructure.spec.volumes", "runtime.kubernetes.volumes"},

	// VM -> runtime.vm
	{"infrastructure.spec.pool.poolName", "runtime.vm.pool"},
	{"infrastructure.spec.pool.identifier", "runtime.vm.pool"}, // fallback if poolName not set
	{"infrastructure.spec.pool.os", "runtime.vm.os"},
	{"infrastructure.spec.pool.harnessImageConnectorRef", "runtime.vm.harness-image-connector"},
	{"infrastructure.spec.pool.timeout", "runtime.vm.timeout"},

	// runtime - spec is skipped in v1, cloud is the v1 prefix
	{"runtime.spec.size", "runtime.-.(cloud.size)"},
	{"runtime.spec.(imageSpec.imageName)", "runtime.-.(cloud.image)"},
}

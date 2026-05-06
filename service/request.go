package service

// ConvertRequest is the request body shared by all single-entity conversion endpoints.
type ConvertRequest struct {
	YAML string `json:"yaml"`

	// TemplateRefMapping rewrites template references in the converted
	// output. Applied to `template.uses` (the ref portion of
	// "ref@version"), and to legacy `templateRef` / `template_ref` keys.
	TemplateRefMapping map[string]string `json:"template_ref_mapping,omitempty"`

	// PipelineRefMapping rewrites pipeline identifiers in the converted
	// output. Applied to `pipeline.id`, to the pipeline segment of a
	// `chain.uses` value ("org/project/pipeline"), and to a trigger's
	// `pipelineIdentifier`. For triggers, the map is also applied
	// recursively to the embedded `inputYaml`.
	PipelineRefMapping map[string]string `json:"pipeline_ref_mapping,omitempty"`

	// ContextPipelineYAML is an optional raw v0 pipeline YAML used purely as
	// expression-postprocess context for template / input-set / trigger
	// conversions. When provided the server parses + structurally converts
	// this pipeline (suppressing its diagnostic messages), harvests the
	// resulting step-type map, and uses it with FQN=true when walking the
	// requested entity for expression conversion. The pipeline endpoint
	// ignores this field. If empty (or omitted), postprocess runs without
	// FQN context — equivalent to the previous behaviour.
	ContextPipelineYAML string `json:"context_pipeline_yaml,omitempty"`
}

// ConvertResponse is the response body for a successful single-entity
// conversion. Report bundles converter messages, the unrecognised-fields
// list (input keys that don't match the v0 schema), and per-expression
// conversions for this entity.
type ConvertResponse struct {
	YAML     string            `json:"yaml"`
	Checksum string            `json:"checksum"`
	Report   *ConversionReport `json:"report,omitempty"`
}

// BatchConvertRequest is the request body for POST /api/v1/convert/batch.
type BatchConvertRequest struct {
	Items []BatchItem `json:"items"`
}

// BatchItem is one entity to convert inside a BatchConvertRequest.
type BatchItem struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"` // "pipeline" | "template" | "input-set" | "trigger"
	YAML       string `json:"yaml"`

	// TemplateRefMapping — same semantics as ConvertRequest.TemplateRefMapping.
	TemplateRefMapping map[string]string `json:"template_ref_mapping,omitempty"`

	// PipelineRefMapping — same semantics as ConvertRequest.PipelineRefMapping.
	PipelineRefMapping map[string]string `json:"pipeline_ref_mapping,omitempty"`

	// ContextPipelineYAML — same semantics as ConvertRequest.ContextPipelineYAML.
	// Ignored when EntityType == "pipeline".
	ContextPipelineYAML string `json:"context_pipeline_yaml,omitempty"`
}

// BatchConvertResponse is the response body for POST /api/v1/convert/batch.
type BatchConvertResponse struct {
	Results []BatchResult `json:"results"`
}

// BatchResult is the outcome for a single item in a batch conversion.
// On success YAML and Checksum are set and Error is nil.
// On failure YAML and Checksum are nil and Error is set.
// Report mirrors the per-entity ConversionReport from the single-entity
// endpoints and is nil when conversion failed.
type BatchResult struct {
	ID         string            `json:"id"`
	EntityType string            `json:"entity_type"`
	YAML       *string           `json:"yaml"`
	Checksum   *string           `json:"checksum"`
	Error      *string           `json:"error"`
	Report     *ConversionReport `json:"report,omitempty"`
}

// ChecksumRequest is the request body for POST /api/v1/checksum.
type ChecksumRequest struct {
	YAML string `json:"yaml"`
}

// ChecksumResponse is the response body for POST /api/v1/checksum.
type ChecksumResponse struct {
	Checksum string `json:"checksum"`
}

// ErrorResponse is the standard error body returned on non-2xx responses.
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ExpressionConvertRequest is the request body for POST /api/v1/convert/expression.
type ExpressionConvertRequest struct {
	// Expression is a single expression to convert (use this OR Expressions, not both)
	Expression string `json:"expression,omitempty"`

	// Expressions is a list of expressions to convert (use this OR Expression, not both)
	Expressions []string `json:"expressions,omitempty"`

	// Context provides optional metadata for context-aware conversion
	Context *ExpressionContextRequest `json:"context,omitempty"`
}

// ExpressionContextRequest holds the context needed for expression conversion.
type ExpressionContextRequest struct {
	// CurrentStepID is the ID of the current step we're inside (optional)
	CurrentStepID string `json:"current_step_id,omitempty"`

	// CurrentStepType is the type of the step we're currently inside (e.g., "Run", "Action", "Plugin")
	CurrentStepType string `json:"current_step_type,omitempty"`

	// CurrentStepV1Path is the v1 FQN base path to the current step
	// Example: "pipeline.stages.build.steps.restoreCache"
	CurrentStepV1Path string `json:"current_step_v1_path,omitempty"`

	// StepTypeMap maps step ID to step type for all steps in the pipeline
	StepTypeMap map[string]string `json:"step_type_map,omitempty"`

	// StepV1PathMap maps step ID to its v1 FQN base path
	StepV1PathMap map[string]string `json:"step_v1_path_map,omitempty"`

	// UseFQN enables fully qualified name mode for step expressions
	UseFQN bool `json:"use_fqn,omitempty"`
}

// ExpressionConvertResponse is the response body for POST /api/v1/convert/expression.
type ExpressionConvertResponse struct {
	// Expression is the converted expression (when single expression was provided)
	Expression string `json:"expression,omitempty"`

	// Expressions is a map of original expression to converted expression (when multiple were provided)
	Expressions map[string]string `json:"expressions,omitempty"`

	// Checksum is the SHA-256 checksum of the converted expression(s)
	Checksum string `json:"checksum"`
}

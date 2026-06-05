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

	// ContextPipelineYAML is an optional v1 pipeline YAML used as expression
	// context for template / input-set / trigger conversions: it builds the
	// FQN-keyed step lookup (StepInfoByFQN) and enables FQN mode. Ignored by
	// the pipeline endpoint; if empty, postprocess runs without FQN context.
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

	// RemoteFile is raw file contents (manifest, values.yaml, etc.) with
	// embedded Harness v0 expressions. The server converts every <+...> /
	// ${{...}} occurrence and returns the file with all expressions replaced.
	RemoteFile string `json:"remote_file,omitempty"`

	// ContextPipelineYAML is an optional v1 pipeline YAML. When set, it builds
	// the FQN-keyed step lookup (StepInfoByFQN) and enables FQN mode.
	ContextPipelineYAML string `json:"context_pipeline_yaml,omitempty"`

	// CurrentFQN is the v1 FQN of callsite of the expression.
	CurrentFQN string `json:"current_fqn,omitempty"`
}

// ExpressionConvertResponse is the response body for POST /api/v1/convert/expression.
type ExpressionConvertResponse struct {
	// Expression is the converted expression (when single expression was provided)
	Expression string `json:"expression,omitempty"`

	// Expressions is a map of original expression to converted expression (when multiple were provided)
	Expressions map[string]string `json:"expressions,omitempty"`

	// RemoteFile is the file contents with all embedded expressions converted
	RemoteFile string `json:"remote_file,omitempty"`

	// Warnings holds non-fatal diagnostics, e.g. ambiguous step types resolved
	// via best-match fallback or unmapped template/approval uses.
	Warnings []string `json:"warnings,omitempty"`

	// Checksum is the SHA-256 checksum of the converted expression(s)
	Checksum string `json:"checksum"`
}

package service

// ConvertRequest is the request body shared by all single-entity conversion endpoints.
type ConvertRequest struct {
	YAML             string            `json:"yaml"`
	EntityRefMapping map[string]string `json:"entity_ref_mapping,omitempty"`
}

// ConvertResponse is the response body for a successful single-entity conversion.
type ConvertResponse struct {
	YAML     string `json:"yaml"`
	Checksum string `json:"checksum"`
}

// BatchConvertRequest is the request body for POST /api/v1/convert/batch.
type BatchConvertRequest struct {
	Items []BatchItem `json:"items"`
}

// BatchItem is one entity to convert inside a BatchConvertRequest.
type BatchItem struct {
	ID               string            `json:"id"`
	EntityType       string            `json:"entity_type"` // "pipeline" | "template" | "input-set"
	YAML             string            `json:"yaml"`
	EntityRefMapping map[string]string `json:"entity_ref_mapping,omitempty"`
}

// BatchConvertResponse is the response body for POST /api/v1/convert/batch.
type BatchConvertResponse struct {
	Results []BatchResult `json:"results"`
}

// BatchResult is the outcome for a single item in a batch conversion.
// On success YAML and Checksum are set and Error is nil.
// On failure YAML and Checksum are nil and Error is set.
type BatchResult struct {
	ID         string  `json:"id"`
	EntityType string  `json:"entity_type"`
	YAML       *string `json:"yaml"`
	Checksum   *string `json:"checksum"`
	Error      *string `json:"error"`
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
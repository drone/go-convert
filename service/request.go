package service

// ConvertRequest is the request body shared by all single-entity conversion endpoints.
type ConvertRequest struct {
	YAML string `json:"yaml"`
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
	ID         string `json:"id"`
	EntityType string `json:"entity_type"` // "pipeline" | "template" | "input-set"
	YAML       string `json:"yaml"`
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

// ErrorResponse is the standard error body returned on non-2xx responses.
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

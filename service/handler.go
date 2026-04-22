package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/drone/go-convert/service/converter"
)

const (
	entityPipeline = "pipeline"
	entityTemplate = "template"
	entityInputSet = "input-set"
)

// Handler holds the HTTP handler methods for the conversion service.
type Handler struct {
	maxBatchSize int
	maxYAMLBytes int64
}

// NewHandler creates a Handler with the given limits.
func NewHandler(maxBatchSize int, maxYAMLBytes int64) *Handler {
	return &Handler{
		maxBatchSize: maxBatchSize,
		maxYAMLBytes: maxYAMLBytes,
	}
}

// Healthz handles GET /healthz.
func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ConvertPipeline handles POST /api/v1/convert/pipeline.
func (h *Handler) ConvertPipeline(w http.ResponseWriter, r *http.Request) {
	h.convertSingle(w, r, entityPipeline)
}

// ConvertTemplate handles POST /api/v1/convert/template.
func (h *Handler) ConvertTemplate(w http.ResponseWriter, r *http.Request) {
	h.convertSingle(w, r, entityTemplate)
}

// ConvertInputSet handles POST /api/v1/convert/input-set.
func (h *Handler) ConvertInputSet(w http.ResponseWriter, r *http.Request) {
	h.convertSingle(w, r, entityInputSet)
}

// ConvertBatch handles POST /api/v1/convert/batch.
// The response is always HTTP 200; per-item errors are reported inline.
func (h *Handler) ConvertBatch(w http.ResponseWriter, r *http.Request) {
	var req BatchConvertRequest
	if err := h.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error(), nil)
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "MISSING_FIELD", "'items' must not be empty", nil)
		return
	}
	if len(req.Items) > h.maxBatchSize {
		writeError(w, http.StatusBadRequest, "BATCH_TOO_LARGE",
			fmt.Sprintf("batch size %d exceeds maximum of %d", len(req.Items), h.maxBatchSize), nil)
		return
	}

	results := make([]BatchResult, 0, len(req.Items))
	for _, item := range req.Items {
		result := BatchResult{ID: item.ID, EntityType: item.EntityType}
		outBytes, err := dispatch(item.EntityType, item.YAML, item.EntityRefMapping)
		if err != nil {
			e := err.Error()
			result.Error = &e
		} else {
			s := string(outBytes)
			cs := Checksum(outBytes)
			result.YAML = &s
			result.Checksum = &cs
		}
		results = append(results, result)
	}
	writeJSON(w, http.StatusOK, BatchConvertResponse{Results: results})
}

// ComputeChecksum handles POST /api/v1/checksum.
func (h *Handler) ComputeChecksum(w http.ResponseWriter, r *http.Request) {
	var req ChecksumRequest
	if err := h.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error(), nil)
		return
	}
	if strings.TrimSpace(req.YAML) == "" {
		writeError(w, http.StatusBadRequest, "MISSING_FIELD", "'yaml' field is required and must not be empty", nil)
		return
	}

	writeJSON(w, http.StatusOK, ChecksumResponse{
		Checksum: Checksum([]byte(req.YAML)),
	})
}

// convertSingle is the shared implementation for single-entity endpoints.
func (h *Handler) convertSingle(w http.ResponseWriter, r *http.Request, entityType string) {
	var req ConvertRequest
	if err := h.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error(), nil)
		return
	}
	if strings.TrimSpace(req.YAML) == "" {
		writeError(w, http.StatusBadRequest, "MISSING_FIELD", "'yaml' field is required and must not be empty", nil)
		return
	}

	outBytes, err := dispatch(entityType, req.YAML, req.EntityRefMapping)
	if err != nil {
		code, status := classifyError(err)
		writeError(w, status, code, err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, ConvertResponse{
		YAML:     string(outBytes),
		Checksum: Checksum(outBytes),
	})
}

// dispatch routes a conversion request to the appropriate converter function.
func dispatch(entityType, yamlStr string, refMapping map[string]string) ([]byte, error) {
	switch entityType {
	case entityPipeline:
		return converter.Pipeline(yamlStr, refMapping)
	case entityTemplate:
		return converter.Template(yamlStr, refMapping)
	case entityInputSet:
		return converter.InputSet(yamlStr, refMapping)
	default:
		return nil, fmt.Errorf("unknown entity_type %q (must be pipeline, template, or input-set)", entityType)
	}
}

// classifyError maps a converter error to an HTTP status code and error code string.
func classifyError(err error) (code string, httpStatus int) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "expected top-level"):
		return "WRONG_ENTITY_TYPE", http.StatusBadRequest
	case strings.Contains(msg, "failed to parse"):
		return "INVALID_YAML", http.StatusBadRequest
	case strings.Contains(msg, "unsupported template type"),
		strings.Contains(msg, "conversion returned nil"),
		strings.Contains(msg, "produced no stages"),
		strings.Contains(msg, "produced no steps"):
		return "CONVERSION_FAILED", http.StatusUnprocessableEntity
	case strings.Contains(msg, "unknown entity_type"):
		return "MISSING_FIELD", http.StatusBadRequest
	default:
		return "INTERNAL_ERROR", http.StatusInternalServerError
	}
}

// decodeJSON decodes the JSON request body into v, enforcing the configured size limit.
func (h *Handler) decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, h.maxYAMLBytes)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(v); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return fmt.Errorf("request body exceeds maximum allowed size of %d bytes", h.maxYAMLBytes)
		}
		return err
	}
	return nil
}

// writeJSON serialises v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a standard ErrorResponse JSON body.
func writeError(w http.ResponseWriter, status int, code, message string, details interface{}) {
	writeJSON(w, status, ErrorResponse{Code: code, Message: message, Details: details})
}

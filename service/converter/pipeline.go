package converter

import (
	"fmt"
	"strings"
	"sync"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// buildContextFromPipelineYAML parses contextPipelineYAML as a v0 pipeline,
// runs structural conversion (without expression post-processing) on a fresh
// PipelineConverter, and returns the harvested step-type map. The boolean
// useFQN reports whether the caller should pass useFQN=true to
// PostProcessExpressions (true iff a usable map was built).
//
// Diagnostic messages emitted by the structural conversion of the context
// pipeline are intentionally suppressed via the global MessageLogger so that
// the caller's ConversionReport reflects only the requested entity. A
// CONTEXT_PIPELINE_PARSE_FAILED warning is added when the YAML is non-empty
// but unparseable; in that case the function returns (nil, false) so the
// caller falls back to postprocess without FQN context.
func buildContextFromPipelineYAML(contextPipelineYAML string) (map[string]*pipelineconverter.StepInfo, bool) {
	if strings.TrimSpace(contextPipelineYAML) == "" {
		return nil, false
	}
	cfg, _, err := v0.ParseStringWithUnknownFields(contextPipelineYAML)
	if err != nil || cfg == nil {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		pipelineconverter.GetMessageLogger().LogWarning(
			"CONTEXT_PIPELINE_PARSE_FAILED",
			"failed to parse 'pipeline_yaml' context; postprocess will run without FQN context",
			pipelineconverter.WithContext(map[string]string{"error": errMsg}),
		)
		return nil, false
	}

	// Suppress messages emitted while structurally converting the context
	// pipeline — they describe a pipeline the caller did not ask to convert
	// and would otherwise pollute the ConversionReport.
	msg := pipelineconverter.GetMessageLogger()
	msg.Disable()
	defer msg.Enable("")

	ctxC := pipelineconverter.NewPipelineConverter()
	_ = ctxC.ConvertPipeline(&cfg.Pipeline)
	stepTypeMap := ctxC.GetStepTypeMap()
	if len(stepTypeMap) == 0 {
		// No steps means nothing useful for FQN resolution; skip FQN mode.
		return nil, false
	}
	return stepTypeMap, true
}

// Result is the outcome of a successful conversion. YAML carries the
// converted v1 document; UnknownFields lists JSON paths for keys in the
// input that do not match any field in the v0 schema (parsing does not fail
// on unknown fields — they are surfaced for observability only). Summary
// bundles converter messages, unknown fields, and expression conversions.
type Result struct {
	YAML          []byte
	UnknownFields []string
	Summary       *pipelineconverter.ConversionSummary
}

// The three global loggers (expressions, unknown fields, messages) are
// process-wide singletons, so the service path serialises conversions with
// this mutex. Conversions are typically sub-second; any throughput regression
// can be revisited by making each logger support per-goroutine scoping.
var apiConversionMu sync.Mutex

// apiFileMarker is the per-conversion key used with the file-scoped loggers
// when running on the API path. The BuildSummary step uses the same marker.
const apiFileMarker = "api"

// beginAPIConversion locks the service-side logger and scopes the three
// loggers to apiFileMarker. Callers must defer the returned closer to
// release the lock and reset logger state.
func beginAPIConversion() func() {
	apiConversionMu.Lock()
	pipelineconverter.GetMessageLogger().Enable("")
	pipelineconverter.GetMessageLogger().SetBatchMode(false)
	pipelineconverter.GetMessageLogger().SetCurrentFile(apiFileMarker)
	pipelineconverter.GetExpressionLogger().Enable("")
	pipelineconverter.GetExpressionLogger().SetBatchMode(false)
	pipelineconverter.GetExpressionLogger().SetCurrentFile(apiFileMarker)
	return func() {
		pipelineconverter.GetMessageLogger().Clear()
		pipelineconverter.GetMessageLogger().Disable()
		pipelineconverter.GetExpressionLogger().Clear()
		pipelineconverter.GetExpressionLogger().Disable()
		pipelineconverter.GetUnknownFieldsLogger().Clear()
		apiConversionMu.Unlock()
	}
}

// buildAPISummary composes the per-conversion summary for the API path. The
// unknown-fields list is the parse-time authoritative one rather than the
// logger's entry, which is empty on the API path (we don't Record() since
// the list is already in hand).
//
// For API responses, expressions are deduplicated by (original, converted)
// pair and context is stripped to keep the response concise.
func buildAPISummary(unknownFields []string) *pipelineconverter.ConversionSummary {
	s := pipelineconverter.BuildSummary(apiFileMarker)
	if s == nil {
		s = &pipelineconverter.ConversionSummary{}
	}
	s.FilePath = ""
	if len(unknownFields) > 0 {
		s.UnknownFields = unknownFields
	}
	// Deduplicate expressions and strip context for API response
	s.Expressions = dedupeExpressions(s.Expressions)
	return s
}

// dedupeExpressions returns a deduplicated list of expressions with context
// stripped. Only unique (original, converted) pairs are kept.
func dedupeExpressions(exprs []pipelineconverter.ExpressionLogEntry) []pipelineconverter.ExpressionLogEntry {
	if len(exprs) == 0 {
		return nil
	}
	type key struct{ orig, conv string }
	seen := make(map[key]struct{})
	var out []pipelineconverter.ExpressionLogEntry
	for _, e := range exprs {
		k := key{e.Original, e.Converted}
		if _, exists := seen[k]; exists {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, pipelineconverter.ExpressionLogEntry{
			Original:  e.Original,
			Converted: e.Converted,
			// Context intentionally omitted for API response
		})
	}
	return out
}

// Pipeline converts a Harness v0 pipeline YAML string into v1 YAML bytes.
// The input must have a top-level "pipeline:" key.
// If refMapping is provided, template references in the output will be replaced.
func Pipeline(yamlStr string, refMapping map[string]string) (*Result, error) {
	if err := validateTopLevelKey(yamlStr, "pipeline"); err != nil {
		return nil, err
	}

	done := beginAPIConversion()
	defer done()

	v0Config, unknownFields, err := v0.ParseStringWithUnknownFields(yamlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse v0 pipeline: %w", err)
	}
	if v0Config == nil {
		return nil, fmt.Errorf("failed to parse v0 pipeline: result is nil")
	}

	c := pipelineconverter.NewPipelineConverter()
	v1Pipeline := c.ConvertPipeline(&v0Config.Pipeline)
	if v1Pipeline == nil {
		return nil, fmt.Errorf("conversion returned nil: check that the pipeline has at least one supported stage")
	}

	// Single-pass expression post-process with full pipeline context.
	pipelineconverter.PostProcessExpressions(v1Pipeline, c.GetStepTypeMap(), true)

	out, err := v1.MarshalPipeline(v1Pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal v1 pipeline: %w", err)
	}

	out, err = ReplaceTemplateRefs(out, refMapping)
	if err != nil {
		return nil, err
	}
	return &Result{
		YAML:          out,
		UnknownFields: unknownFields,
		Summary:       buildAPISummary(unknownFields),
	}, nil
}

// validateTopLevelKey returns an error when yamlStr does not start with "key:".
// It tolerates a leading YAML document separator (---).
func validateTopLevelKey(yamlStr, key string) error {
	s := strings.TrimSpace(yamlStr)
	if strings.HasPrefix(s, "---") {
		s = strings.TrimSpace(s[3:])
	}
	if !strings.HasPrefix(s, key+":") {
		return fmt.Errorf("expected top-level '%s:' key", key)
	}
	return nil
}

package service

import (
	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
)

// ConverterMessageDTO is the public-facing form of a converter notice.
// Fields mirror the proto ConverterMessage message (severity is stringified
// for JSON readability).
type ConverterMessageDTO struct {
	Severity string            `json:"severity"`
	Code     string            `json:"code"`
	Message  string            `json:"message"`
	Context  map[string]string `json:"context,omitempty"`
}

// ExpressionEntryDTO is the public-facing form of an expression conversion.
// Status is "SUCCESS" when the converter produced a different string and
// "NOT_CONVERTED" when the original was returned unchanged.
type ExpressionEntryDTO struct {
	Original  string `json:"original"`
	Converted string `json:"converted"`
	Status    string `json:"status"`
}

// ConversionReport is the unified per-conversion report returned on every
// successful convert call: structured messages, the list of input fields
// the converter did not recognise, and per-expression conversions.
type ConversionReport struct {
	Messages           []ConverterMessageDTO `json:"messages,omitempty"`
	UnrecognizedFields []string              `json:"unrecognized_fields,omitempty"`
	Expressions        []ExpressionEntryDTO  `json:"expressions,omitempty"`
}

// expressionStatus picks the ConversionStatus enum string for an entry.
func expressionStatus(original, converted string) string {
	if original == converted {
		return "NOT_CONVERTED"
	}
	return "SUCCESS"
}

// buildReport flattens an internal ConversionSummary plus the parse-time
// unrecognised-fields list into a public ConversionReport. Returns nil when
// there is nothing to report.
func buildReport(s *pipelineconverter.ConversionSummary, unrecognized []string) *ConversionReport {
	if s == nil && len(unrecognized) == 0 {
		return nil
	}
	r := &ConversionReport{UnrecognizedFields: unrecognized}
	if s != nil {
		if len(unrecognized) == 0 && len(s.UnknownFields) > 0 {
			r.UnrecognizedFields = s.UnknownFields
		}
		if len(s.Messages) > 0 {
			r.Messages = make([]ConverterMessageDTO, 0, len(s.Messages))
			for _, m := range s.Messages {
				r.Messages = append(r.Messages, ConverterMessageDTO{
					Severity: string(m.Severity),
					Code:     m.Code,
					Message:  m.Message,
					Context:  m.Context,
				})
			}
		}
		if len(s.Expressions) > 0 {
			type key struct{ orig, conv string }
			seen := make(map[key]struct{}, len(s.Expressions))
			for _, e := range s.Expressions {
				k := key{e.Original, e.Converted}
				if _, dup := seen[k]; dup {
					continue
				}
				seen[k] = struct{}{}
				r.Expressions = append(r.Expressions, ExpressionEntryDTO{
					Original:  e.Original,
					Converted: e.Converted,
					Status:    expressionStatus(e.Original, e.Converted),
				})
			}
		}
	}
	if r.Messages == nil && len(r.UnrecognizedFields) == 0 && r.Expressions == nil {
		return nil
	}
	return r
}

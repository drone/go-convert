// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pipelineconverter

// Severity, Message, and the associated constants are now defined in the
// leaf package convert/v0tov1/messagelog and re-exported via aliases in
// message_logger.go so existing call sites keep working.

// Counts summarises message totals per severity.
type Counts struct {
	Info    int `json:"info"`
	Warning int `json:"warning"`
	Error   int `json:"error"`
}

// ConversionSummary bundles everything emitted by one pipeline conversion:
// structured messages, the unknown-fields list, and expression conversions.
// It is returned in the API response and written to the CLI _summary.json
// sidecar.
type ConversionSummary struct {
	FilePath      string               `json:"file_path,omitempty"`
	Counts        Counts               `json:"counts"`
	Messages      []Message            `json:"messages,omitempty"`
	UnknownFields []string             `json:"unknown_fields,omitempty"`
	Expressions   []ExpressionLogEntry `json:"expressions,omitempty"`
}

// BuildSummary merges the per-file views of the three loggers into a single
// ConversionSummary. Callers invoke this after ConvertPipeline and before
// Clear().
func BuildSummary(filePath string) *ConversionSummary {
	s := &ConversionSummary{FilePath: filePath}

	if fl := GetMessageLogger().GetFileLog(filePath); fl != nil {
		s.Messages = fl.Messages
	}
	if uf := GetUnknownFieldsLogger().GetEntry(filePath); len(uf) > 0 {
		s.UnknownFields = uf
	}
	if el := GetExpressionLogger().GetFileLog(filePath); el != nil {
		s.Expressions = el.Expressions
	}
	s.Counts = computeCounts(s.Messages)
	return s
}

func computeCounts(msgs []Message) Counts {
	var c Counts
	for _, m := range msgs {
		switch m.Severity {
		case SeverityInfo:
			c.Info++
		case SeverityWarning:
			c.Warning++
		case SeverityError:
			c.Error++
		}
	}
	return c
}

// HasErrors reports whether the summary contains any ERROR-severity messages.
func (s *ConversionSummary) HasErrors() bool {
	return s != nil && s.Counts.Error > 0
}

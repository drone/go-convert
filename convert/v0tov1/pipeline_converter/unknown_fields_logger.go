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

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FileUnknownFieldsLog is the sidecar entry for a single converted file.
type FileUnknownFieldsLog struct {
	FilePath      string   `json:"file_path"`
	UnknownFields []string `json:"unknown_fields"`
}

// UnknownFieldsLogger accumulates unknown-field paths discovered during
// parsing and flushes them to a sidecar JSON file. Single-file mode emits
// one object; batch mode emits an array of per-file objects.
//
// The logger mirrors ExpressionLogger so CLI wiring is symmetric.
type UnknownFieldsLogger struct {
	mu          sync.Mutex
	enabled     bool
	logFilePath string
	batchMode   bool
	entries     []FileUnknownFieldsLog
}

var (
	globalUnknownFieldsLogger *UnknownFieldsLogger
	unknownFieldsLoggerOnce   sync.Once
)

// GetUnknownFieldsLogger returns the global singleton.
func GetUnknownFieldsLogger() *UnknownFieldsLogger {
	unknownFieldsLoggerOnce.Do(func() {
		globalUnknownFieldsLogger = &UnknownFieldsLogger{}
	})
	return globalUnknownFieldsLogger
}

// Enable turns on logging and sets the output path.
func (l *UnknownFieldsLogger) Enable(logFilePath string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = true
	l.logFilePath = logFilePath
}

// Disable turns off logging.
func (l *UnknownFieldsLogger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = false
}

// SetBatchMode configures whether Flush emits a single object or an array.
func (l *UnknownFieldsLogger) SetBatchMode(batch bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.batchMode = batch
}

// Record appends the unknown-fields list for a file. A nil/empty list is
// still recorded in batch mode for a complete per-file manifest; in
// single-file mode it is recorded too so the sidecar reflects a run that
// found nothing (the CLI callers decide whether to skip the sidecar).
func (l *UnknownFieldsLogger) Record(filePath string, unknownFields []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.enabled {
		return
	}
	l.entries = append(l.entries, FileUnknownFieldsLog{
		FilePath:      filePath,
		UnknownFields: unknownFields,
	})
}

// Flush writes accumulated entries to the configured path. Entries with no
// unknown fields are filtered out to keep the sidecar focused on drift.
// Returns nil (no error, no file written) when there is nothing to report.
func (l *UnknownFieldsLogger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.enabled || l.logFilePath == "" {
		return nil
	}

	var nonEmpty []FileUnknownFieldsLog
	for _, e := range l.entries {
		if len(e.UnknownFields) > 0 {
			nonEmpty = append(nonEmpty, e)
		}
	}
	if len(nonEmpty) == 0 {
		return nil
	}

	dir := filepath.Dir(l.logFilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	file, err := os.Create(l.logFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	if l.batchMode || len(nonEmpty) > 1 {
		return enc.Encode(nonEmpty)
	}
	return enc.Encode(nonEmpty[0])
}

// GetEntry returns a copy of the unknown-fields list recorded for filePath,
// or nil when nothing has been recorded. Used by BuildSummary to embed
// unknown fields in the per-pipeline ConversionSummary.
func (l *UnknownFieldsLogger) GetEntry(filePath string) []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range l.entries {
		if e.FilePath == filePath {
			if len(e.UnknownFields) == 0 {
				return nil
			}
			out := make([]string, len(e.UnknownFields))
			copy(out, e.UnknownFields)
			return out
		}
	}
	return nil
}

// Clear drops accumulated entries.
func (l *UnknownFieldsLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = nil
}

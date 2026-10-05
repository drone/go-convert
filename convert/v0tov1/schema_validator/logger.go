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

package schema_validator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// SchemaValidationLogger accumulates ValidationResult entries and writes
// them as a JSON array to a file. It is safe for concurrent use.
//
// In batch mode all results are collected in memory and written at once
// by Flush(). In non-batch (single-file) mode each Record() call also
// triggers a write.
type SchemaValidationLogger struct {
	mu        sync.Mutex
	enabled   bool
	batchMode bool
	filePath  string
	results   []*ValidationResult
}

var (
	schemaLoggerOnce     sync.Once
	schemaLoggerInstance *SchemaValidationLogger
)

// GetSchemaValidationLogger returns the package-level singleton logger.
func GetSchemaValidationLogger() *SchemaValidationLogger {
	schemaLoggerOnce.Do(func() {
		schemaLoggerInstance = &SchemaValidationLogger{}
	})
	return schemaLoggerInstance
}

// Enable turns on logging and sets the output file path.
func (l *SchemaValidationLogger) Enable(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = true
	l.filePath = path
}

// Disable turns off logging.
func (l *SchemaValidationLogger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = false
}

// SetBatchMode controls whether results are accumulated (true) or
// written immediately on each Record call (false).
func (l *SchemaValidationLogger) SetBatchMode(batch bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.batchMode = batch
}

// Clear removes all accumulated results.
func (l *SchemaValidationLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.results = nil
}

// Record adds a validation result. In non-batch mode it also writes
// the accumulated results to disk immediately.
func (l *SchemaValidationLogger) Record(r *ValidationResult) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.enabled || r == nil {
		return
	}
	l.results = append(l.results, r)
	if !l.batchMode {
		_ = l.writeFile()
	}
}

// Flush writes all accumulated results to the output file.
func (l *SchemaValidationLogger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.enabled || l.filePath == "" {
		return nil
	}
	return l.writeFile()
}

// Stats returns total, valid, and invalid counts.
func (l *SchemaValidationLogger) Stats() (total, valid, invalid int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	total = len(l.results)
	for _, r := range l.results {
		if r.Valid {
			valid++
		} else {
			invalid++
		}
	}
	return total, valid, total - valid
}

// writeFile serializes l.results to JSON and writes to l.filePath.
// Caller must hold l.mu.
func (l *SchemaValidationLogger) writeFile() error {
	if l.filePath == "" {
		return nil
	}
	dir := filepath.Dir(l.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(l.results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.filePath, data, 0o644)
}

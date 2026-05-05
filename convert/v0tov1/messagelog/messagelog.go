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

// Package messagelog provides the structured converter-message logger used
// throughout the v0->v1 pipeline converter. It is a leaf package so it can
// be imported from the v0 unmarshallers (convert/harness/yaml), the
// per-step converters (convert/v0tov1/convert_helpers), and the
// orchestrator (convert/v0tov1/pipeline_converter) without creating import
// cycles.
package messagelog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Severity classifies a converter message.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityError   Severity = "ERROR"
)

// Message is a single structured converter notice.
type Message struct {
	Severity Severity          `json:"severity"`
	Code     string            `json:"code"`
	Message  string            `json:"message"`
	Location string            `json:"location,omitempty"`
	Context  map[string]string `json:"context,omitempty"`
}

// FileMessageLog is the sidecar entry for one converted file.
type FileMessageLog struct {
	FilePath string    `json:"file_path"`
	Messages []Message `json:"messages"`
}

// MessageLogger is a process-global singleton that accumulates structured
// converter messages. It mirrors ExpressionLogger and UnknownFieldsLogger so
// CLI and service wiring stays symmetric. Per-pipeline scoping is achieved
// via SetCurrentFile + Clear.
//
// The logger is not concurrency-safe for cross-pipeline writes — callers
// running parallel conversions must serialize access (see service/converter
// for the HTTP/gRPC path).
type MessageLogger struct {
	mu          sync.Mutex
	enabled     bool
	logFilePath string
	currentFile string
	batchMode   bool
	fileLogs    map[string]*FileMessageLog
}

var (
	globalMessageLogger *MessageLogger
	messageLoggerOnce   sync.Once
)

// GetMessageLogger returns the global singleton.
func GetMessageLogger() *MessageLogger {
	messageLoggerOnce.Do(func() {
		globalMessageLogger = &MessageLogger{
			fileLogs: make(map[string]*FileMessageLog),
		}
	})
	return globalMessageLogger
}

// ResetMessageLogger recreates the singleton. Intended for tests.
func ResetMessageLogger() {
	globalMessageLogger = &MessageLogger{
		fileLogs: make(map[string]*FileMessageLog),
	}
}

// Enable turns on logging and sets the optional sidecar output path.
// An empty path is valid — accumulated messages are still readable via
// GetFileLog but Flush becomes a no-op.
func (l *MessageLogger) Enable(logFilePath string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = true
	l.logFilePath = logFilePath
}

// Disable turns off logging.
func (l *MessageLogger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = false
}

// IsEnabled returns whether logging is enabled.
func (l *MessageLogger) IsEnabled() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.enabled
}

// SetBatchMode configures whether Flush emits a single object or an array.
func (l *MessageLogger) SetBatchMode(batch bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.batchMode = batch
}

// SetCurrentFile scopes subsequent Log* calls to filePath.
func (l *MessageLogger) SetCurrentFile(filePath string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.currentFile = filePath
	if _, ok := l.fileLogs[filePath]; !ok {
		l.fileLogs[filePath] = &FileMessageLog{
			FilePath: filePath,
			Messages: []Message{},
		}
	}
}

// MessageOption mutates a Message before it is recorded.
type MessageOption func(*Message)

// WithLocation attaches a structural path (e.g. "stages[0].spec.execution.steps[2]").
func WithLocation(loc string) MessageOption {
	return func(m *Message) { m.Location = loc }
}

// WithContext merges key/value context pairs into the message.
func WithContext(kv map[string]string) MessageOption {
	return func(m *Message) {
		if m.Context == nil {
			m.Context = make(map[string]string, len(kv))
		}
		for k, v := range kv {
			m.Context[k] = v
		}
	}
}

// WithStep adds step_id and type to the message context.
func WithStep(stepID, stepType string) MessageOption {
	return WithContext(map[string]string{
		"step_id": stepID,
		"type":    stepType,
	})
}

// WithStage adds stage_id and type to the message context.
func WithStage(stageID, stageType string) MessageOption {
	return WithContext(map[string]string{
		"stage_id": stageID,
		"type":     stageType,
	})
}

// Log records a fully-formed message for the current file.
func (l *MessageLogger) Log(m Message) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.enabled || l.currentFile == "" {
		return
	}
	fl, ok := l.fileLogs[l.currentFile]
	if !ok {
		fl = &FileMessageLog{FilePath: l.currentFile}
		l.fileLogs[l.currentFile] = fl
	}
	fl.Messages = append(fl.Messages, m)
}

// LogInfo records an INFO-severity message.
func (l *MessageLogger) LogInfo(code, msg string, opts ...MessageOption) {
	l.emit(SeverityInfo, code, msg, opts...)
}

// LogWarning records a WARNING-severity message.
func (l *MessageLogger) LogWarning(code, msg string, opts ...MessageOption) {
	l.emit(SeverityWarning, code, msg, opts...)
}

// LogError records an ERROR-severity message.
func (l *MessageLogger) LogError(code, msg string, opts ...MessageOption) {
	l.emit(SeverityError, code, msg, opts...)
}

func (l *MessageLogger) emit(sev Severity, code, msg string, opts ...MessageOption) {
	m := Message{Severity: sev, Code: code, Message: msg}
	for _, opt := range opts {
		opt(&m)
	}
	l.Log(m)
}

// GetFileLog returns a copy of the messages recorded for filePath, or nil
// when nothing has been recorded.
func (l *MessageLogger) GetFileLog(filePath string) *FileMessageLog {
	l.mu.Lock()
	defer l.mu.Unlock()
	fl, ok := l.fileLogs[filePath]
	if !ok || len(fl.Messages) == 0 {
		return nil
	}
	msgs := make([]Message, len(fl.Messages))
	copy(msgs, fl.Messages)
	return &FileMessageLog{FilePath: fl.FilePath, Messages: msgs}
}

// Flush writes accumulated entries to the configured path. Entries with no
// messages are filtered out. Returns nil (no error, no file written) when
// there is nothing to report or logging was not configured with a path.
func (l *MessageLogger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.enabled || l.logFilePath == "" {
		return nil
	}

	var nonEmpty []FileMessageLog
	for _, fl := range l.fileLogs {
		if len(fl.Messages) > 0 {
			nonEmpty = append(nonEmpty, *fl)
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

// Clear drops accumulated entries and the currentFile marker.
func (l *MessageLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileLogs = make(map[string]*FileMessageLog)
	l.currentFile = ""
}

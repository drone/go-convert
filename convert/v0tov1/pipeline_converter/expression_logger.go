package pipelineconverter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/drone/go-convert/convert/convertexpressions"
)

// ConversionContextLog is a JSON-serializable version of ConversionContext
type ConversionContextLog struct {
	StepID            string            `json:"step_id,omitempty"`
	StepType          string            `json:"step_type,omitempty"`
	CurrentStepType   string            `json:"current_step_type,omitempty"`
	CurrentStepV1Path string            `json:"current_step_v1_path,omitempty"`
	UseFQN            bool              `json:"use_fqn"`
	StepTypeMap       map[string]string `json:"step_type_map,omitempty"`
	StepV1PathMap     map[string]string `json:"step_v1_path_map,omitempty"`
}

// ExpressionLogEntry represents a single expression conversion log entry
type ExpressionLogEntry struct {
	Original  string                `json:"original"`
	Converted string                `json:"converted"`
	Context   *ConversionContextLog `json:"context,omitempty"`
}

// FileExpressionLog represents all expression conversions for a single file
type FileExpressionLog struct {
	FilePath    string               `json:"file_path"`
	Expressions []ExpressionLogEntry `json:"expressions"`
}

// conversionKey is used for deduplication of original-converted pairs per file
type conversionKey struct {
	original  string
	converted string
}

// ExpressionLogger handles logging of expression conversions to a separate JSON log file
type ExpressionLogger struct {
	mu              sync.Mutex
	enabled         bool
	logFilePath     string
	currentFile     string
	fileLogs        map[string]*FileExpressionLog     // map of file path to its expression logs
	seenConversions map[string]map[conversionKey]bool // map of file path to set of seen conversions
	batchMode       bool                              // true when processing a directory
	includeContext  bool                              // whether to include context in logs (default true)
}

// Global expression logger instance
var globalExprLogger *ExpressionLogger
var loggerOnce sync.Once

// GetExpressionLogger returns the global expression logger instance
func GetExpressionLogger() *ExpressionLogger {
	loggerOnce.Do(func() {
		globalExprLogger = &ExpressionLogger{
			fileLogs:        make(map[string]*FileExpressionLog),
			seenConversions: make(map[string]map[conversionKey]bool),
			includeContext:  true, // default to including context
		}
	})
	return globalExprLogger
}

// ResetExpressionLogger resets the global logger (useful for testing)
func ResetExpressionLogger() {
	globalExprLogger = &ExpressionLogger{
		fileLogs:        make(map[string]*FileExpressionLog),
		seenConversions: make(map[string]map[conversionKey]bool),
		includeContext:  true,
	}
}

// Enable enables expression logging with the specified log file path
func (l *ExpressionLogger) Enable(logFilePath string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = true
	l.logFilePath = logFilePath
}

// Disable disables expression logging
func (l *ExpressionLogger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = false
}

// IsEnabled returns whether logging is enabled
func (l *ExpressionLogger) IsEnabled() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.enabled
}

// SetIncludeContext sets whether to include context in logs
func (l *ExpressionLogger) SetIncludeContext(include bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.includeContext = include
}

// SetBatchMode sets whether we're in batch mode (processing multiple files)
func (l *ExpressionLogger) SetBatchMode(batch bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.batchMode = batch
}

// SetCurrentFile sets the current file being processed
func (l *ExpressionLogger) SetCurrentFile(filePath string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.currentFile = filePath
	if _, exists := l.fileLogs[filePath]; !exists {
		l.fileLogs[filePath] = &FileExpressionLog{
			FilePath:    filePath,
			Expressions: []ExpressionLogEntry{},
		}
		l.seenConversions[filePath] = make(map[conversionKey]bool)
	}
}

// LogConversion logs a single expression conversion with full context
// All expressions are logged (including unchanged ones) for complete tracking
func (l *ExpressionLogger) LogConversion(original, converted string, ctx *convertexpressions.ConversionContext) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.enabled {
		return
	}

	// Check for duplicate original-converted pair in this file
	key := conversionKey{original: original, converted: converted}
	if l.currentFile != "" {
		if seen, exists := l.seenConversions[l.currentFile]; exists {
			if seen[key] {
				return // Skip duplicate
			}
			seen[key] = true
		}
	}

	entry := ExpressionLogEntry{
		Original:  original,
		Converted: converted,
	}

	// Include context if enabled and context is provided
	if l.includeContext && ctx != nil {
		entry.Context = &ConversionContextLog{
			StepID:            ctx.StepID,
			StepType:          ctx.StepType,
			CurrentStepType:   ctx.CurrentStepType,
			CurrentStepV1Path: ctx.CurrentStepV1Path,
			UseFQN:            ctx.UseFQN,
			StepTypeMap:       ctx.StepTypeMap,
			StepV1PathMap:     ctx.StepV1PathMap,
		}
	}

	if l.currentFile != "" {
		if fileLog, exists := l.fileLogs[l.currentFile]; exists {
			fileLog.Expressions = append(fileLog.Expressions, entry)
		}
	}
}

// Flush writes all accumulated logs to the log file (overwrites existing file)
func (l *ExpressionLogger) Flush() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.enabled || l.logFilePath == "" {
		return nil
	}

	// Collect all file logs
	var allLogs []FileExpressionLog
	for _, fileLog := range l.fileLogs {
		if len(fileLog.Expressions) > 0 {
			allLogs = append(allLogs, *fileLog)
		}
	}

	if len(allLogs) == 0 {
		return nil
	}

	// Create directory if needed
	dir := filepath.Dir(l.logFilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Write JSON log file (os.Create truncates/overwrites existing file)
	file, err := os.Create(l.logFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	// SetEscapeHTML(false) ensures <> are not escaped to \u003c\u003e
	encoder.SetEscapeHTML(false)

	// For batch mode, write array of file logs
	// For single file mode, write just the expressions
	if l.batchMode || len(allLogs) > 1 {
		return encoder.Encode(allLogs)
	} else if len(allLogs) == 1 {
		return encoder.Encode(allLogs[0])
	}

	return nil
}

// Clear clears all accumulated logs
func (l *ExpressionLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileLogs = make(map[string]*FileExpressionLog)
	l.seenConversions = make(map[string]map[conversionKey]bool)
	l.currentFile = ""
}

// GetLogFilePath returns the current log file path
func (l *ExpressionLogger) GetLogFilePath() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logFilePath
}

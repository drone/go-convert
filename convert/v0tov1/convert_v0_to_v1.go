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

package v0tov1

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipeline_converter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	"github.com/drone/go-convert/convert/v0tov1/schema_validator"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// globalSchemaValidator is initialized once in Main() when --schema_dir
// or HARNESS_SCHEMA_DIR is provided. Nil means schema validation is disabled.
var globalSchemaValidator *schema_validator.SchemaValidator

// loadDotEnv reads a .env file (KEY=VALUE per line) and sets any variables
// that are not already present in the environment. This lets users define
// settings in a .env file at the repo root without modifying their shell
// profile. Lines starting with # are comments.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // missing .env is fine — silently skip
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Strip optional surrounding quotes
		val = strings.Trim(val, `"'`)
		// Only set if not already in environment (explicit env takes precedence)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}

func Main() {
	// Load .env file from working directory (if present) before parsing flags,
	// so that HARNESS_SCHEMA_DIR and future variables are available.
	loadDotEnv(".env")

	baseDir := flag.String("base_dir", "", "Base directory containing v0 and v1 subdirectories")
	filePath := flag.String("file_path", "", "Single pipeline file path to convert")
	inputDir := flag.String("input_dir", "", "Input directory to recursively convert")
	outputDir := flag.String("output_dir", "", "Output directory for converted files")
	accountDir := flag.String("account_dir", "", "Account directory containing {org}/{project}/v0/ subdirectories")
	schemaDir := flag.String("schema_dir", "", "Path to harness-schema/v1/ directory for schema validation (or set HARNESS_SCHEMA_DIR env var)")
	flag.Parse()

	// Resolve schema directory: flag takes precedence over env var.
	schemaDirValue := *schemaDir
	if schemaDirValue == "" {
		schemaDirValue = os.Getenv("HARNESS_SCHEMA_DIR")
	}
	if schemaDirValue != "" {
		sv, err := schema_validator.NewSchemaValidator(schemaDirValue)
		if err != nil {
			log.Printf("Warning: schema validation disabled: %v", err)
		} else {
			globalSchemaValidator = sv
			log.Printf("Schema validation enabled from %s", schemaDirValue)
		}
	}

	// Validate flag combinations
	flagsSet := 0
	if *baseDir != "" {
		flagsSet++
	}
	if *filePath != "" {
		flagsSet++
	}
	if *inputDir != "" || *outputDir != "" {
		flagsSet++
	}
	if *accountDir != "" {
		flagsSet++
	}

	if flagsSet != 1 {
		log.Fatalf("Usage: %s --base_dir <directory> OR --file_path <file> OR --input_dir <dir> --output_dir <dir> OR --account_dir <dir>\n", os.Args[0])
	}

	// Validate input_dir and output_dir are used together
	if (*inputDir != "" && *outputDir == "") || (*inputDir == "" && *outputDir != "") {
		log.Fatalf("Both --input_dir and --output_dir must be specified together\n")
	}

	if *filePath != "" {
		convertSingleFile(*filePath)
	} else if *baseDir != "" {
		convertBaseDirectory(*baseDir)
	} else if *accountDir != "" {
		convertAccountDirectory(*accountDir)
	} else {
		convertRecursiveDirectory(*inputDir, *outputDir)
	}
}

func convertSingleFile(inputPath string) {
	// Validate file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", inputPath)
	}

	// Generate output path (same directory, with _v1 suffix)
	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimSuffix(inputPath, ext) + "_v1" + ext

	// Setup expression logging for single file
	exprLogPath := strings.TrimSuffix(inputPath, ext) + "_expressions.json"
	exprLogger := pipeline_converter.GetExpressionLogger()
	exprLogger.Enable(exprLogPath)
	exprLogger.SetBatchMode(false)
	exprLogger.SetCurrentFile(inputPath)
	defer func() {
		if err := exprLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write expression log: %v", err)
		}
		exprLogger.Clear()
		exprLogger.Disable()
	}()

	// Setup unknown-fields logging (sidecar JSON next to the converted file)
	unknownLogPath := strings.TrimSuffix(inputPath, ext) + "_unknown_fields.json"
	unknownLogger := pipeline_converter.GetUnknownFieldsLogger()
	unknownLogger.Enable(unknownLogPath)
	unknownLogger.SetBatchMode(false)
	defer func() {
		if err := unknownLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write unknown-fields log: %v", err)
		}
		unknownLogger.Clear()
		unknownLogger.Disable()
	}()

	// Setup structured message logging; scoped to this file.
	msgLogger := pipeline_converter.GetMessageLogger()
	msgLogger.Enable("")
	msgLogger.SetBatchMode(false)
	msgLogger.SetCurrentFile(inputPath)
	defer func() {
		msgLogger.Clear()
		msgLogger.Disable()
	}()

	// Setup schema validation logging for single file (sidecar JSON).
	schemaLogger := schema_validator.GetSchemaValidationLogger()
	if globalSchemaValidator != nil {
		schemaLogPath := strings.TrimSuffix(inputPath, ext) + "_schema_validation.json"
		schemaLogger.Enable(schemaLogPath)
		schemaLogger.SetBatchMode(false)
		defer func() {
			if err := schemaLogger.Flush(); err != nil {
				log.Printf("Warning: failed to write schema validation log: %v", err)
			}
			schemaLogger.Clear()
			schemaLogger.Disable()
		}()
	}
	summaryPath := strings.TrimSuffix(inputPath, ext) + "_summary.json"
	defer func() {
		summary := pipeline_converter.BuildSummary(inputPath)
		if err := writeSummaryFile(summaryPath, summary); err != nil {
			log.Printf("Warning: failed to write summary: %v", err)
		}
		printSummary(os.Stdout, summary)
	}()

	// Benchmark: Read v0
	readStart := time.Now()
	v0Config, unknownFields, err := v0.ParseFileWithUnknownFields(inputPath)
	readDur := time.Since(readStart)
	if err != nil {
		log.Fatalf("Failed to parse v0 file: %v", err)
	}
	unknownLogger.Record(inputPath, unknownFields)

	// Benchmark: Convert to v1
	convStart := time.Now()
	converter := pipeline_converter.NewPipelineConverter()

	// Auto-detect root node type and convert accordingly
	writeStart := time.Now()

	var entityType string

	if v0Config.Trigger != nil {
		entityType = schema_validator.EntityTrigger
		// Trigger conversion
		v1Trigger := converter.ConvertTrigger(v0Config.Trigger, nil, false)
		convDur := time.Since(convStart)
		if v1Trigger == nil {
			log.Fatalf("Failed to convert trigger to v1 format")
		}
		pipeline_converter.PostProcessExpressions(v1Trigger, nil, false)
		if err := v1.WriteTriggerFile(outputPath, v1Trigger); err != nil {
			log.Fatalf("Failed to write v1 trigger YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted trigger %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else if v0Config.InputSet != nil {
		entityType = schema_validator.EntityInputSet
		// InputSet conversion
		v1InputSet := converter.ConvertInputSet(v0Config.InputSet)
		convDur := time.Since(convStart)
		if v1InputSet == nil {
			log.Fatalf("Failed to convert inputset to v1 format")
		}
		pipeline_converter.PostProcessExpressions(v1InputSet, nil, false)
		if err := v1.WriteInputSetFile(outputPath, v1InputSet); err != nil {
			log.Fatalf("Failed to write v1 inputset YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted inputset %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else if v0Config.Template != nil {
		entityType = schema_validator.EntityTemplate
		// Template conversion
		v1Template := converter.ConvertTemplate(v0Config.Template)
		convDur := time.Since(convStart)
		if v1Template == nil {
			log.Fatalf("Failed to convert template to v1 format")
		}
		pipeline_converter.PostProcessExpressions(v1Template, nil, false)
		if err := v1.WriteTemplateFile(outputPath, v1Template); err != nil {
			log.Fatalf("Failed to write v1 template YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted template %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else {
		entityType = schema_validator.EntityPipeline
		// Pipeline conversion (default)
		v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
		convDur := time.Since(convStart)
		if v1Pipeline == nil {
			log.Fatalf("Failed to convert pipeline to v1 format")
		}
		pipeline_converter.PostProcessExpressions(v1Pipeline, converter.GetStepInfoByFQN(), true)
		if err := v1.WritePipelineFile(outputPath, v1Pipeline); err != nil {
			log.Fatalf("Failed to write v1 pipeline YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted pipeline %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	}

	// Schema validation (single-file mode: individual report).
	if result := validateV1Output(outputPath, entityType); result != nil {
		schemaLogger.Record(result)
		printSchemaValidationResult(result)
	}

}

func convertBaseDirectory(baseDir string) {
	inputDir := filepath.Join(baseDir, "v0")
	outputDir := filepath.Join(baseDir, "v1")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
	}

	// Setup expression logging for batch mode - single log file for the directory
	exprLogPath := filepath.Join(outputDir, "expressions.json")
	exprLogger := pipeline_converter.GetExpressionLogger()
	exprLogger.Enable(exprLogPath)
	exprLogger.SetBatchMode(true)
	defer func() {
		if err := exprLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write expression log: %v", err)
		}
		exprLogger.Clear()
		exprLogger.Disable()
	}()

	unknownLogPath := filepath.Join(outputDir, "unknown_fields.json")
	unknownLogger := pipeline_converter.GetUnknownFieldsLogger()
	unknownLogger.Enable(unknownLogPath)
	unknownLogger.SetBatchMode(true)
	defer func() {
		if err := unknownLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write unknown-fields log: %v", err)
		}
		unknownLogger.Clear()
		unknownLogger.Disable()
	}()

	// Setup structured message logging for batch mode.
	msgLogger := pipeline_converter.GetMessageLogger()
	msgLogger.Enable("")
	msgLogger.SetBatchMode(true)
	var batchSummaries []*pipeline_converter.ConversionSummary
	defer func() {
		if err := writeAggregateSummary(filepath.Join(outputDir, "summary.json"), batchSummaries); err != nil {
			log.Printf("Warning: failed to write aggregate summary: %v", err)
		}
		msgLogger.Clear()
		msgLogger.Disable()
	}()

	// Setup schema validation logging for batch mode (single combined file).
	schemaLogger := schema_validator.GetSchemaValidationLogger()
	if globalSchemaValidator != nil {
		schemaLogger.Enable(filepath.Join(outputDir, "schema_validation.json"))
		schemaLogger.SetBatchMode(true)
		defer func() {
			if err := schemaLogger.Flush(); err != nil {
				log.Printf("Warning: failed to write schema validation log: %v", err)
			}
			printBatchSchemaValidationStats()
			schemaLogger.Clear()
			schemaLogger.Disable()
		}()
	}

	// Log to stdout only (Python script captures this)
	log.SetOutput(os.Stdout)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory %s: %v", inputDir, err)
	}

	converted := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !(strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml")) {
			continue
		}

		inputPath := filepath.Join(inputDir, name)
		outputPath := filepath.Join(outputDir, name)

		// Emit sentinel before parsing - we'll determine the type after parsing
		log.Printf("CONVERTING %s", inputPath)

		// Scope per-pipeline loggers to this file.
		exprLogger.SetCurrentFile(inputPath)
		msgLogger.SetCurrentFile(inputPath)

		// Benchmark: Read v0
		readStart := time.Now()
		v0Config, unknownFields, err := v0.ParseFileWithUnknownFields(inputPath)
		readDur := time.Since(readStart)
		if err != nil {
			log.Printf("ERROR_PARSING %s: failed to parse v0 file: %v", inputPath, err)
			// Log parse error as structured message so it appears in summary.json
			msgLogger.LogError(
				"PARSE_ERROR",
				fmt.Sprintf("Failed to parse v0 YAML file: %v", err),
				pipeline_converter.WithContext(map[string]string{
					"error": err.Error(),
				}),
			)
			// Build summary even for parse failures
			summary := pipeline_converter.BuildSummary(inputPath)
			batchSummaries = append(batchSummaries, summary)
			printSummary(os.Stdout, summary)
			continue
		}
		unknownLogger.Record(inputPath, unknownFields)

		// Benchmark: Convert to v1
		convStart := time.Now()
		converter := pipeline_converter.NewPipelineConverter()

		// Auto-detect root node type and convert accordingly
		writeStart := time.Now()
		var convDur, writeDur time.Duration

		if v0Config.Trigger != nil {
			v1Trigger := converter.ConvertTrigger(v0Config.Trigger, nil, false)
			convDur = time.Since(convStart)
			if v1Trigger == nil {
				log.Printf("ERROR_TRIGGER %s: failed to convert trigger to v1 format", inputPath)
				continue
			}
			pipeline_converter.PostProcessExpressions(v1Trigger, nil, false)
			if err := v1.WriteTriggerFile(outputPath, v1Trigger); err != nil {
				log.Printf("ERROR_TRIGGER %s: failed to write v1 trigger YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_TRIGGER %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		} else if v0Config.InputSet != nil {
			v1InputSet := converter.ConvertInputSet(v0Config.InputSet)
			convDur = time.Since(convStart)
			if v1InputSet == nil {
				log.Printf("ERROR_INPUTSET %s: failed to convert inputset to v1 format", inputPath)
				continue
			}
			pipeline_converter.PostProcessExpressions(v1InputSet, nil, false)
			if err := v1.WriteInputSetFile(outputPath, v1InputSet); err != nil {
				log.Printf("ERROR_INPUTSET %s: failed to write v1 inputset YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_INPUTSET %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		} else if v0Config.Template != nil {
			v1Template := converter.ConvertTemplate(v0Config.Template)
			convDur = time.Since(convStart)
			if v1Template == nil {
				log.Printf("ERROR_TEMPLATE %s: failed to convert template to v1 format", inputPath)
				continue
			}
			pipeline_converter.PostProcessExpressions(v1Template, nil, false)
			if err := v1.WriteTemplateFile(outputPath, v1Template); err != nil {
				log.Printf("ERROR_TEMPLATE %s: failed to write v1 template YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_TEMPLATE %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		} else {
			v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
			convDur = time.Since(convStart)
			if v1Pipeline == nil {
				log.Printf("ERROR_PIPELINE %s: failed to convert pipeline to v1 format", inputPath)
				continue
			}
			pipeline_converter.PostProcessExpressions(v1Pipeline, converter.GetStepInfoByFQN(), true)
			if err := v1.WritePipelineFile(outputPath, v1Pipeline); err != nil {
				log.Printf("ERROR_PIPELINE %s: failed to write v1 pipeline YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_PIPELINE %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		}

		// Schema validation (batch mode: accumulated into single file).
		if result := validateV1Output(outputPath, ""); result != nil {
			schemaLogger.Record(result)
		}

		// Build per-file summary: write sidecar + accumulate for aggregate, print console line.
		summary := pipeline_converter.BuildSummary(inputPath)
		ext := filepath.Ext(outputPath)
		summaryPath := strings.TrimSuffix(outputPath, ext) + "_summary.json"
		if err := writeSummaryFile(summaryPath, summary); err != nil {
			log.Printf("Warning: failed to write summary for %s: %v", inputPath, err)
		}
		batchSummaries = append(batchSummaries, summary)
		printSummary(os.Stdout, summary)

		converted++
	}

	log.Printf("Converted %d file(s) into %s", converted, outputDir)
}

func convertRecursiveDirectory(inputDir, outputDir string) {
	// Validate input directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		log.Fatalf("Input directory does not exist: %s", inputDir)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
	}

	// Setup log file
	logFile, err := setupLogFile(outputDir)
	if err != nil {
		log.Fatalf("Failed to setup log file: %v", err)
	}
	defer logFile.Close()

	// Setup expression logging for batch mode - single log file for the directory
	exprLogPath := filepath.Join(outputDir, "expressions.json")
	exprLogger := pipeline_converter.GetExpressionLogger()
	exprLogger.Enable(exprLogPath)
	exprLogger.SetBatchMode(true)
	defer func() {
		if err := exprLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write expression log: %v", err)
		}
		exprLogger.Clear()
		exprLogger.Disable()
	}()

	unknownLogPath := filepath.Join(outputDir, "unknown_fields.json")
	unknownLogger := pipeline_converter.GetUnknownFieldsLogger()
	unknownLogger.Enable(unknownLogPath)
	unknownLogger.SetBatchMode(true)
	defer func() {
		if err := unknownLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write unknown-fields log: %v", err)
		}
		unknownLogger.Clear()
		unknownLogger.Disable()
	}()

	msgLogger := pipeline_converter.GetMessageLogger()
	msgLogger.Enable("")
	msgLogger.SetBatchMode(true)
	var batchSummaries []*pipeline_converter.ConversionSummary
	defer func() {
		if err := writeAggregateSummary(filepath.Join(outputDir, "summary.json"), batchSummaries); err != nil {
			log.Printf("Warning: failed to write aggregate summary: %v", err)
		}
		msgLogger.Clear()
		msgLogger.Disable()
	}()

	// Setup schema validation logging for batch mode.
	schemaLogger := schema_validator.GetSchemaValidationLogger()
	if globalSchemaValidator != nil {
		schemaLogger.Enable(filepath.Join(outputDir, "schema_validation.json"))
		schemaLogger.SetBatchMode(true)
		defer func() {
			if err := schemaLogger.Flush(); err != nil {
				log.Printf("Warning: failed to write schema validation log: %v", err)
			}
			printBatchSchemaValidationStats()
			schemaLogger.Clear()
			schemaLogger.Disable()
		}()
	}

	converted := 0
	skipped := 0

	// Walk through input directory recursively
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process YAML files
		if !(strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			return nil
		}

		// Calculate relative path from input directory
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			log.Printf("Failed to get relative path for %s: %v", path, err)
			skipped++
			return nil
		}

		// Create corresponding output path
		outputPath := filepath.Join(outputDir, relPath)

		// Create output subdirectories if needed
		outputSubDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputSubDir, 0o755); err != nil {
			log.Printf("Failed to create output subdirectory %s: %v", outputSubDir, err)
			skipped++
			return nil
		}

		// Convert the file
		success := convertFile(path, outputPath)

		// Always build summary (even for failures) so parse errors appear in reports
		summary := pipeline_converter.BuildSummary(path)
		ext := filepath.Ext(outputPath)
		summaryPath := strings.TrimSuffix(outputPath, ext) + "_summary.json"
		if err := writeSummaryFile(summaryPath, summary); err != nil {
			log.Printf("Warning: failed to write summary for %s: %v", path, err)
		}
		batchSummaries = append(batchSummaries, summary)
		printSummary(os.Stdout, summary)

		// Schema validation (batch mode).
		if success {
			if result := validateV1Output(outputPath, ""); result != nil {
				schemaLogger.Record(result)
			}
			converted++
		} else {
			skipped++
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk input directory: %v", err)
	}

	log.Printf("\nConversion complete:\n")
	log.Printf("  Input directory:  %s\n", inputDir)
	log.Printf("  Output directory: %s\n", outputDir)
	log.Printf("  Converted: %d file(s)\n", converted)
	log.Printf("  Skipped:   %d file(s)\n", skipped)
}

func convertAccountDirectory(accountDir string) {
	// Validate account directory exists
	if _, err := os.Stat(accountDir); os.IsNotExist(err) {
		log.Fatalf("Account directory does not exist: %s", accountDir)
	}

	// Setup expression logging for batch mode — account-level aggregate
	exprLogPath := filepath.Join(accountDir, "expressions.json")
	exprLogger := pipeline_converter.GetExpressionLogger()
	exprLogger.Enable(exprLogPath)
	exprLogger.SetBatchMode(true)
	defer func() {
		if err := exprLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write expression log: %v", err)
		}
		exprLogger.Clear()
		exprLogger.Disable()
	}()

	unknownLogPath := filepath.Join(accountDir, "unknown_fields.json")
	unknownLogger := pipeline_converter.GetUnknownFieldsLogger()
	unknownLogger.Enable(unknownLogPath)
	unknownLogger.SetBatchMode(true)
	defer func() {
		if err := unknownLogger.Flush(); err != nil {
			log.Printf("Warning: failed to write unknown-fields log: %v", err)
		}
		unknownLogger.Clear()
		unknownLogger.Disable()
	}()

	msgLogger := pipeline_converter.GetMessageLogger()
	msgLogger.Enable("")
	msgLogger.SetBatchMode(true)
	var allSummaries []*pipeline_converter.ConversionSummary
	defer func() {
		// Account-level aggregate summary
		if err := writeAggregateSummary(filepath.Join(accountDir, "summary.json"), allSummaries); err != nil {
			log.Printf("Warning: failed to write account-level aggregate summary: %v", err)
		}
		msgLogger.Clear()
		msgLogger.Disable()
	}()

	// Setup schema validation logging for batch mode (account-level aggregate).
	schemaLogger := schema_validator.GetSchemaValidationLogger()
	if globalSchemaValidator != nil {
		schemaLogger.Enable(filepath.Join(accountDir, "schema_validation.json"))
		schemaLogger.SetBatchMode(true)
		defer func() {
			if err := schemaLogger.Flush(); err != nil {
				log.Printf("Warning: failed to write schema validation log: %v", err)
			}
			printBatchSchemaValidationStats()
			schemaLogger.Clear()
			schemaLogger.Disable()
		}()
	}

	// Log to stdout (Python captures this)
	log.SetOutput(os.Stdout)

	totalConverted := 0
	totalSkipped := 0
	projectCount := 0

	// Walk the account directory tree looking for v0/ directories.
	// Expected layout: accountDir/{org}/{project}/v0/*.yaml
	err := filepath.Walk(accountDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return nil
		}

		// We only care about directories named "v0"
		if !info.IsDir() || info.Name() != "v0" {
			return nil
		}

		v0Dir := path
		projectDir := filepath.Dir(v0Dir)
		v1Dir := filepath.Join(projectDir, "v1")

		// Read YAML files from this v0 directory
		entries, err := os.ReadDir(v0Dir)
		if err != nil {
			log.Printf("Failed to read v0 directory %s: %v", v0Dir, err)
			return filepath.SkipDir
		}

		var yamlFiles []os.DirEntry
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
				yamlFiles = append(yamlFiles, e)
			}
		}

		if len(yamlFiles) == 0 {
			return filepath.SkipDir
		}

		// Create v1 output directory
		if err := os.MkdirAll(v1Dir, 0o755); err != nil {
			log.Printf("Failed to create v1 directory %s: %v", v1Dir, err)
			return filepath.SkipDir
		}

		// Derive org/project from path for logging
		relPath, _ := filepath.Rel(accountDir, projectDir)
		log.Printf("Converting %d file(s) in %s", len(yamlFiles), relPath)
		projectCount++

		var projectSummaries []*pipeline_converter.ConversionSummary
		converted := 0

		for _, entry := range yamlFiles {
			inputPath := filepath.Join(v0Dir, entry.Name())
			outputPath := filepath.Join(v1Dir, entry.Name())

			log.Printf("CONVERTING %s", inputPath)

			exprLogger.SetCurrentFile(inputPath)
			msgLogger.SetCurrentFile(inputPath)

			success := convertFile(inputPath, outputPath)

			// Always build summary (even for failures) so parse errors appear in reports
			summary := pipeline_converter.BuildSummary(inputPath)
			projectSummaries = append(projectSummaries, summary)
			allSummaries = append(allSummaries, summary)
			printSummary(os.Stdout, summary)

			if success {
				converted++
				// Schema validation (batch mode: accumulated into single file).
				if result := validateV1Output(outputPath, ""); result != nil {
					schemaLogger.Record(result)
				}
			} else {
				totalSkipped++
			}
		}

		totalConverted += converted

		// Write per-project aggregate summary.json in v1/
		if err := writeAggregateSummary(filepath.Join(v1Dir, "summary.json"), projectSummaries); err != nil {
			log.Printf("Warning: failed to write project summary for %s: %v", relPath, err)
		}

		// Write per-project v1/expressions.json so the Python aggregator
		// always has fresh data (mirrors what --base_dir mode produces).
		var projectInputPaths []string
		for _, entry := range yamlFiles {
			projectInputPaths = append(projectInputPaths, filepath.Join(v0Dir, entry.Name()))
		}
		if err := exprLogger.FlushForFiles(projectInputPaths, filepath.Join(v1Dir, "expressions.json")); err != nil {
			log.Printf("Warning: failed to write project expression log for %s: %v", relPath, err)
		}

		log.Printf("Converted %d/%d file(s) in %s", converted, len(yamlFiles), relPath)

		// Don't recurse into v0 subdirectories
		return filepath.SkipDir
	})

	if err != nil {
		log.Fatalf("Failed to walk account directory: %v", err)
	}

	log.Printf("\nAccount conversion complete:\n")
	log.Printf("  Account directory: %s\n", accountDir)
	log.Printf("  Projects processed: %d\n", projectCount)
	log.Printf("  Total converted: %d file(s)\n", totalConverted)
	log.Printf("  Total skipped:   %d file(s)\n", totalSkipped)
}

func convertFile(inputPath, outputPath string) bool {
	// Scope per-pipeline loggers to this file.
	pipeline_converter.GetExpressionLogger().SetCurrentFile(inputPath)
	pipeline_converter.GetMessageLogger().SetCurrentFile(inputPath)

	// Benchmark: Read v0
	log.Printf("Converting %s to %s", inputPath, outputPath)
	readStart := time.Now()
	v0Config, unknownFields, err := v0.ParseFileWithUnknownFields(inputPath)
	readDur := time.Since(readStart)

	if err != nil {
		log.Printf("Skipping %s: failed to parse v0 file: %v", inputPath, err)
		// Log parse error as structured message so it appears in summary.json and HTML
		pipeline_converter.GetMessageLogger().LogError(
			"PARSE_ERROR",
			fmt.Sprintf("Failed to parse v0 YAML file: %v", err),
			pipeline_converter.WithContext(map[string]string{
				"error": err.Error(),
			}),
		)
		return false
	}
	pipeline_converter.GetUnknownFieldsLogger().Record(inputPath, unknownFields)

	// Benchmark: Convert to v1
	convStart := time.Now()
	converter := pipeline_converter.NewPipelineConverter()

	// Auto-detect root node type and convert accordingly
	var writeErr error
	writeStart := time.Now()

	if v0Config.Trigger != nil {
		// Trigger conversion
		v1Trigger := converter.ConvertTrigger(v0Config.Trigger, nil, false)
		convDur := time.Since(convStart)
		if v1Trigger == nil {
			log.Printf("Skipping %s: failed to convert trigger to v1 format", inputPath)
			return false
		}
		pipeline_converter.PostProcessExpressions(v1Trigger, nil, false)
		writeErr = v1.WriteTriggerFile(outputPath, v1Trigger)
		writeDur := time.Since(writeStart)
		if writeErr != nil {
			log.Printf("Failed to write v1 trigger YAML for %s: %v", inputPath, writeErr)
			return false
		}
		fmt.Printf("Converted trigger %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else if v0Config.InputSet != nil {
		// InputSet conversion
		v1InputSet := converter.ConvertInputSet(v0Config.InputSet)
		convDur := time.Since(convStart)
		if v1InputSet == nil {
			log.Printf("Skipping %s: failed to convert inputset to v1 format", inputPath)
			return false
		}
		pipeline_converter.PostProcessExpressions(v1InputSet, nil, false)
		writeErr = v1.WriteInputSetFile(outputPath, v1InputSet)
		writeDur := time.Since(writeStart)
		if writeErr != nil {
			log.Printf("Failed to write v1 inputset YAML for %s: %v", inputPath, writeErr)
			return false
		}
		fmt.Printf("Converted inputset %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else if v0Config.Template != nil {
		// Template conversion
		v1Template := converter.ConvertTemplate(v0Config.Template)
		convDur := time.Since(convStart)
		if v1Template == nil {
			log.Printf("Skipping %s: failed to convert template to v1 format", inputPath)
			return false
		}
		pipeline_converter.PostProcessExpressions(v1Template, nil, false)
		writeErr = v1.WriteTemplateFile(outputPath, v1Template)
		writeDur := time.Since(writeStart)
		if writeErr != nil {
			log.Printf("Failed to write v1 template YAML for %s: %v", inputPath, writeErr)
			return false
		}
		fmt.Printf("Converted template %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else {
		// Pipeline conversion (default)
		v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
		convDur := time.Since(convStart)
		if v1Pipeline == nil {
			log.Printf("Skipping %s: failed to convert pipeline to v1 format", inputPath)
			return false
		}
		pipeline_converter.PostProcessExpressions(v1Pipeline, converter.GetStepInfoByFQN(), true)
		writeErr = v1.WritePipelineFile(outputPath, v1Pipeline)
		writeDur := time.Since(writeStart)
		if writeErr != nil {
			log.Printf("Failed to write v1 pipeline YAML for %s: %v", inputPath, writeErr)
			return false
		}
		fmt.Printf("Converted pipeline %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	}

	return true
}

func setupLogFile(outputDir string) (*os.File, error) {
	logFileName := fmt.Sprintf("conversion.log")
	logFilePath := filepath.Join(outputDir, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %v", err)
	}

	// Setup multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Printf("Starting conversion - logging to %s\n", logFilePath)
	return logFile, nil
}

// writeSummaryFile writes a single ConversionSummary to path as indented JSON.
// A nil or empty summary produces no file.
func writeSummaryFile(path string, s *pipeline_converter.ConversionSummary) error {
	if s == nil || isEmptySummary(s) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(s)
}

// writeAggregateSummary writes the batch-mode combined summary.json. Only
// non-empty per-file summaries are included.
func writeAggregateSummary(path string, summaries []*pipeline_converter.ConversionSummary) error {
	var kept []*pipeline_converter.ConversionSummary
	for _, s := range summaries {
		if s != nil && !isEmptySummary(s) {
			kept = append(kept, s)
		}
	}
	if len(kept) == 0 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(kept)
}

func isEmptySummary(s *pipeline_converter.ConversionSummary) bool {
	return len(s.Messages) == 0 && len(s.UnknownFields) == 0 && len(s.Expressions) == 0
}

// printSummary writes a one-line human-readable summary plus a bullet list
// of unknown step/stage types encountered, if any.
func printSummary(w io.Writer, s *pipeline_converter.ConversionSummary) {
	if s == nil {
		return
	}
	unknownSteps := collectCodeContext(s.Messages, "UNKNOWN_STEP_TYPE", "type")
	unknownStages := collectCodeContext(s.Messages, "UNKNOWN_STAGE_TYPE", "type")

	fmt.Fprintf(w, "Summary: %d info, %d warnings, %d errors; %d unknown fields\n",
		s.Counts.Info, s.Counts.Warning, s.Counts.Error, len(s.UnknownFields))
	if len(unknownSteps) > 0 {
		fmt.Fprintf(w, "  Unknown step types: %s\n", strings.Join(unknownSteps, ", "))
	}
	if len(unknownStages) > 0 {
		fmt.Fprintf(w, "  Unknown stage types: %s\n", strings.Join(unknownStages, ", "))
	}
}

// collectCodeContext pulls sorted, deduplicated values of context[key] for
// messages whose Code matches code.
func collectCodeContext(msgs []pipeline_converter.Message, code, key string) []string {
	seen := make(map[string]struct{})
	for _, m := range msgs {
		if m.Code != code {
			continue
		}
		if v, ok := m.Context[key]; ok && v != "" {
			seen[v] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for v := range seen {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// ---------------------------------------------------------------------------
// Schema validation helpers
// ---------------------------------------------------------------------------

// validateV1Output validates a written v1 YAML file against the JSON Schema.
// Returns nil when schema validation is disabled or the output file cannot
// be read.
func validateV1Output(outputPath, entityType string) *schema_validator.ValidationResult {
	if globalSchemaValidator == nil {
		return nil
	}

	yamlBytes, err := os.ReadFile(outputPath)
	if err != nil {
		log.Printf("Warning: cannot read v1 output for schema validation: %v", err)
		return nil
	}

	// Auto-detect entity type if not provided.
	if entityType == "" {
		entityType = schema_validator.DetectEntityType(yamlBytes)
	}

	result := globalSchemaValidator.Validate(yamlBytes, entityType)
	result.FilePath = outputPath
	return result
}

// printSchemaValidationResult prints a one-line console summary of
// the schema validation result.
func printSchemaValidationResult(result *schema_validator.ValidationResult) {
	if result == nil {
		return
	}
	if result.Valid {
		fmt.Printf("Schema validation: VALID (%s)\n", result.FilePath)
	} else {
		fmt.Printf("Schema validation: %d error(s) (%s)\n", len(result.SchemaErrors), result.FilePath)
	}
}

// printBatchSchemaValidationStats prints an aggregate one-line summary
// for batch conversion.
func printBatchSchemaValidationStats() {
	logger := schema_validator.GetSchemaValidationLogger()
	total, valid, invalid := logger.Stats()
	if total == 0 {
		return
	}
	fmt.Printf("Schema validation: %d/%d valid, %d with errors\n", valid, total, invalid)
}

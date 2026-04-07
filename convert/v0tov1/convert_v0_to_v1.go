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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	pipeline_converter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func Main() {
	baseDir := flag.String("base_dir", "", "Base directory containing v0 and v1 subdirectories")
	filePath := flag.String("file_path", "", "Single pipeline file path to convert")
	inputDir := flag.String("input_dir", "", "Input directory to recursively convert")
	outputDir := flag.String("output_dir", "", "Output directory for converted files")
	flag.Parse()

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

	if flagsSet != 1 {
		log.Fatalf("Usage: %s --base_dir <directory> OR --file_path <file> OR --input_dir <dir> --output_dir <dir>\n", os.Args[0])
	}

	// Validate input_dir and output_dir are used together
	if (*inputDir != "" && *outputDir == "") || (*inputDir == "" && *outputDir != "") {
		log.Fatalf("Both --input_dir and --output_dir must be specified together\n")
	}

	if *filePath != "" {
		convertSingleFile(*filePath)
	} else if *baseDir != "" {
		convertBaseDirectory(*baseDir)
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

	// Benchmark: Read v0
	readStart := time.Now()
	v0Config, err := v0.ParseFile(inputPath)
	readDur := time.Since(readStart)
	if err != nil {
		log.Fatalf("Failed to parse v0 file: %v", err)
	}

	// Benchmark: Convert to v1
	convStart := time.Now()
	converter := pipeline_converter.NewPipelineConverter()

	// Auto-detect root node type and convert accordingly
	writeStart := time.Now()

	if v0Config.InputSet != nil {
		// InputSet conversion
		v1InputSet := converter.ConvertInputSet(v0Config.InputSet)
		convDur := time.Since(convStart)
		if v1InputSet == nil {
			log.Fatalf("Failed to convert inputset to v1 format")
		}
		if err := v1.WriteInputSetFile(outputPath, v1InputSet); err != nil {
			log.Fatalf("Failed to write v1 inputset YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted inputset %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else if v0Config.Template != nil {
		// Template conversion
		v1Template := converter.ConvertTemplate(v0Config.Template)
		convDur := time.Since(convStart)
		if v1Template == nil {
			log.Fatalf("Failed to convert template to v1 format")
		}
		if err := v1.WriteTemplateFile(outputPath, v1Template); err != nil {
			log.Fatalf("Failed to write v1 template YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted template %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	} else {
		// Pipeline conversion (default)
		v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
		convDur := time.Since(convStart)
		if v1Pipeline == nil {
			log.Fatalf("Failed to convert pipeline to v1 format")
		}
		if err := v1.WritePipelineFile(outputPath, v1Pipeline); err != nil {
			log.Fatalf("Failed to write v1 pipeline YAML: %v", err)
		}
		writeDur := time.Since(writeStart)
		fmt.Printf("Converted pipeline %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
	}
}

func convertBaseDirectory(baseDir string) {
	inputDir := filepath.Join(baseDir, "v0")
	outputDir := filepath.Join(baseDir, "v1")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
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

		// Benchmark: Read v0
		readStart := time.Now()
		v0Config, err := v0.ParseFile(inputPath)
		readDur := time.Since(readStart)
		if err != nil {
			log.Printf("ERROR_PIPELINE %s: failed to parse v0 file: %v", inputPath, err)
			continue
		}

		// Benchmark: Convert to v1
		convStart := time.Now()
		converter := pipeline_converter.NewPipelineConverter()

		// Auto-detect root node type and convert accordingly
		writeStart := time.Now()
		var convDur, writeDur time.Duration

		if v0Config.InputSet != nil {
			// Sentinel: start of inputset conversion
			log.Printf("CONVERTING_INPUTSET %s", inputPath)

			v1InputSet := converter.ConvertInputSet(v0Config.InputSet)
			convDur = time.Since(convStart)
			if v1InputSet == nil {
				log.Printf("ERROR_INPUTSET %s: failed to convert inputset to v1 format", inputPath)
				continue
			}
			if err := v1.WriteInputSetFile(outputPath, v1InputSet); err != nil {
				log.Printf("ERROR_INPUTSET %s: failed to write v1 inputset YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_INPUTSET %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		} else if v0Config.Template != nil {
			// Sentinel: start of template conversion
			log.Printf("CONVERTING_TEMPLATE %s", inputPath)

			v1Template := converter.ConvertTemplate(v0Config.Template)
			convDur = time.Since(convStart)
			if v1Template == nil {
				log.Printf("ERROR_TEMPLATE %s: failed to convert template to v1 format", inputPath)
				continue
			}
			if err := v1.WriteTemplateFile(outputPath, v1Template); err != nil {
				log.Printf("ERROR_TEMPLATE %s: failed to write v1 template YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_TEMPLATE %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		} else {
			// Sentinel: start of pipeline conversion
			log.Printf("CONVERTING_PIPELINE %s", inputPath)

			v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
			convDur = time.Since(convStart)
			if v1Pipeline == nil {
				log.Printf("ERROR_PIPELINE %s: failed to convert pipeline to v1 format", inputPath)
				continue
			}
			if err := v1.WritePipelineFile(outputPath, v1Pipeline); err != nil {
				log.Printf("ERROR_PIPELINE %s: failed to write v1 pipeline YAML: %v", inputPath, err)
				continue
			}
			writeDur = time.Since(writeStart)
			log.Printf("CONVERTED_PIPELINE %s -> %s (read=%v, convert=%v, write=%v)", inputPath, outputPath, readDur, convDur, writeDur)
		}

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
		if convertFile(path, outputPath) {
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

func convertFile(inputPath, outputPath string) bool {
	// Benchmark: Read v0
	log.Printf("Converting %s to %s", inputPath, outputPath)
	readStart := time.Now()
	v0Config, err := v0.ParseFile(inputPath)
	readDur := time.Since(readStart)

	if err != nil {
		log.Printf("Skipping %s: failed to parse v0 file: %v", inputPath, err)
		return false
	}

	// Benchmark: Convert to v1
	convStart := time.Now()
	converter := pipeline_converter.NewPipelineConverter()

	// Auto-detect root node type and convert accordingly
	var writeErr error
	writeStart := time.Now()

	if v0Config.InputSet != nil {
		// InputSet conversion
		v1InputSet := converter.ConvertInputSet(v0Config.InputSet)
		convDur := time.Since(convStart)
		if v1InputSet == nil {
			log.Printf("Skipping %s: failed to convert inputset to v1 format", inputPath)
			return false
		}
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

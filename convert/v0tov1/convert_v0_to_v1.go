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
	flag.Parse()

	// Validate that exactly one flag is provided
	if (*baseDir == "" && *filePath == "") || (*baseDir != "" && *filePath != "") {
		log.Fatalf("Usage: %s --base_dir <directory> OR --file_path <file>\n", os.Args[0])
	}

	if *filePath != "" {
		// Single file mode
		convertSingleFile(*filePath)
	} else {
		// Base directory mode
		convertBaseDirectory(*baseDir)
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
		log.Fatalf("Failed to parse v0 pipeline file after %v: %v", readDur, err)
	}

	// Benchmark: Convert to v1
	convStart := time.Now()
	converter := pipeline_converter.NewPipelineConverter()
	v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
	convDur := time.Since(convStart)
	if v1Pipeline == nil {
		log.Fatalf("Failed to convert pipeline to v1 format (convert took %v)", convDur)
	}

	// Benchmark: Write v1 YAML
	writeStart := time.Now()
	if err := v1.WritePipelineFile(outputPath, v1Pipeline); err != nil {
		writeDur := time.Since(writeStart)
		log.Fatalf("Failed to write v1 pipeline YAML after %v: %v", writeDur, err)
	}
	writeDur := time.Since(writeStart)

	fmt.Printf("Converted %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
}

func convertBaseDirectory(baseDir string) {
	inputDir := filepath.Join(baseDir, "v0")
	outputDir := filepath.Join(baseDir, "v1")

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
	}

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
			log.Printf("Skipping %s: failed to parse v0 pipeline file after %v: %v", inputPath, readDur, err)
			continue
		}

		// Benchmark: Convert to v1
		convStart := time.Now()
		converter := pipeline_converter.NewPipelineConverter()
		v1Pipeline := converter.ConvertPipeline(&v0Config.Pipeline)
		convDur := time.Since(convStart)
		if v1Pipeline == nil {
			log.Printf("Skipping %s: failed to convert pipeline to v1 format (convert took %v)", inputPath, convDur)
			continue
		}

		// Benchmark: Write v1 YAML
		writeStart := time.Now()
		if err := v1.WritePipelineFile(outputPath, v1Pipeline); err != nil {
			writeDur := time.Since(writeStart)
			log.Printf("Failed to write v1 pipeline YAML for %s after %v: %v", inputPath, writeDur, err)
			continue
		}
		writeDur := time.Since(writeStart)

		fmt.Printf("Converted %s -> %s (read=%v, convert=%v, write=%v)\n", inputPath, outputPath, readDur, convDur, writeDur)
		converted++
	}

	fmt.Printf("Converted %d pipeline(s) into %s\n", converted, outputDir)
}

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

package command

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/drone/go-convert/convert/harness"

	"github.com/drone/go-convert/convert/harness/downgrader"
	"github.com/drone/go-convert/convert/jenkinsjson"

	"github.com/google/subcommands"
)

type JenkinsJson struct {
	name         string
	proj         string
	org          string
	repoName     string
	repoConn     string
	kubeName     string
	kubeConn     string
	dockerConn   string
	defaultImage string

	// Infrastructure configuration
	infrastructure string
	os             string
	arch           string

	downgrade       bool
	useIntelligence bool
	randomId        bool
	beforeAfter     bool
	outputDir       string
}

func (*JenkinsJson) Name() string     { return "jenkinsjson" }
func (*JenkinsJson) Synopsis() string { return "converts a jenkinsjson pipeline" }
func (*JenkinsJson) Usage() string {
	return `jenkinsjson [-downgrade] [-intelligence] [-random-id] [-infrastructure cloud|kubernetes|local] [-os linux|mac|windows] [-arch amd64|arm64] [jenkinsjson.json]
`
}

func (c *JenkinsJson) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.downgrade, "downgrade", false, "downgrade to the legacy yaml format")
	f.BoolVar(&c.useIntelligence, "intelligence", false, "Use Harness intelligence features")
	f.BoolVar(&c.randomId, "random-id", false, "Generate random ID for pipeline")
	f.BoolVar(&c.beforeAfter, "before-after", false, "print the befor and after")
	f.StringVar(&c.outputDir, "output-dir", "", "directory where the output should be saved")

	f.StringVar(&c.org, "org", "default", "harness organization")
	f.StringVar(&c.proj, "project", "default", "harness project")
	f.StringVar(&c.name, "pipeline", harness.DefaultName, "harness pipeline name")
	f.StringVar(&c.repoConn, "repo-connector", "", "repository connector")
	f.StringVar(&c.repoName, "repo-name", "", "repository name")
	f.StringVar(&c.kubeConn, "kube-connector", "", "kubernetes connector")
	f.StringVar(&c.kubeName, "kube-namespace", "", "kubernets namespace")
	f.StringVar(&c.dockerConn, "docker-connector", "", "dockerhub connector")
	f.StringVar(&c.defaultImage, "default-image", "alpine", "default image for run step")

	// Infrastructure configuration flags
	f.StringVar(&c.infrastructure, "infrastructure", "cloud", "infrastructure type (cloud, kubernetes, local)")
	f.StringVar(&c.os, "os", "linux", "operating system (linux, mac, windows)")
	f.StringVar(&c.arch, "arch", "amd64", "CPU architecture (amd64, arm64)")
}

func (c *JenkinsJson) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() >= 1 {
		return processArgs(f, c)
	} else {
		log.Println("No file(s) specified")
		return subcommands.ExitFailure
	}
}

func processArgs(f *flag.FlagSet, c *JenkinsJson) subcommands.ExitStatus {
	var fileCount, dirCount int
	totalFiles := 0

	// First, count how many files and directories there are to process
	for _, arg := range f.Args() {
		files, dirs := countFilesAndDirs(arg, c)
		totalFiles += files
		dirCount += dirs
	}

	// Then process files while printing progress (skip progress if output is stdout)
	progress := 0
	for _, arg := range f.Args() {
		filesProcessed, dirsProcessed := recursivelyProcessFiles(arg, c, &progress, totalFiles)
		fileCount += filesProcessed
		dirCount += dirsProcessed
	}

	printSummary(fileCount, dirCount, c)
	return subcommands.ExitSuccess
}

func countFilesAndDirs(filename string, c *JenkinsJson) (int, int) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("Cannot open input file: %s\n%v", filename, err)
	}
	if fileInfo.IsDir() {
		return countFilesInDir(filename, c)
	}
	return 1, 0
}

func countFilesInDir(filename string, c *JenkinsJson) (int, int) {
	dirEntries, err := os.ReadDir(filename)
	if err != nil {
		log.Fatalf("Failed to read directory children: %s\n%v", filename, err)
	}
	fileCount := 0
	dirCount := 1 // Count the directory itself
	for _, entry := range dirEntries {
		absFilepath := filename + "/" + entry.Name()
		files, dirs := countFilesAndDirs(absFilepath, c)
		fileCount += files
		dirCount += dirs
	}
	return fileCount, dirCount
}

func recursivelyProcessFiles(filename string, c *JenkinsJson, progress *int, totalFiles int) (int, int) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("Cannot open input file: %s\n%v", filename, err)
	}

	if fileInfo.IsDir() {
		return recursivelyProcessDir(filename, c, progress, totalFiles)
	} else {
		return processFile(filename, c, progress, totalFiles), 0
	}
}

func processFile(filename string, c *JenkinsJson, progress *int, totalFiles int) int {
	outFile := getOutputFile(c, filename)

	if exitStatus := c.processFile(filename, outFile); exitStatus != subcommands.ExitSuccess {
		return 0
	}

	*progress++
	printProgress(*progress, totalFiles, c)
	return 1
}

func recursivelyProcessDir(filename string, c *JenkinsJson, progress *int, totalFiles int) (int, int) {
	dirEntries, err := os.ReadDir(filename)
	if err != nil {
		log.Fatalf("Failed to read directory children: %s\n%v", filename, err)
	}

	fileCount := 0
	dirCount := 1 // Count the directory itself
	for _, entry := range dirEntries {
		absFilepath := filename + "/" + entry.Name()
		filesProcessed, dirsProcessed := recursivelyProcessFiles(absFilepath, c, progress, totalFiles)
		fileCount += filesProcessed
		dirCount += dirsProcessed
	}

	return fileCount, dirCount
}

func getOutputFile(c *JenkinsJson, filename string) *os.File {
	if isOutputToStdout(c) {
		return os.Stdout
	}
	return createOutputFile(c, filename)
}

func isOutputToStdout(c *JenkinsJson) bool {
	return c.outputDir == ""
}

func createOutputFile(c *JenkinsJson, inputFile string) *os.File {
	fileInfo, err := os.Stat(c.outputDir)
	if os.IsNotExist(err) {
		log.Fatalf("Provided output directory does not exist: %s", fileInfo)
	} else if !fileInfo.IsDir() {
		log.Fatalf("Provided output path is not a directory: %s", fileInfo)
	} else if err != nil {
		log.Fatalln(err)
	}

	base := filepath.Base(inputFile)
	ext := filepath.Ext(inputFile)
	filename := strings.TrimSuffix(base, ext)

	outputFilePath := c.outputDir + "/" + filename + ".yaml"
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalln("Failed to create the output file", err)
	}
	return outputFile
}

func printProgress(current, total int, c *JenkinsJson) {
	if !isOutputToStdout(c) {
		progress := float64(current) / float64(total) * 100
		// Print progress on the same line, overwriting previous output
		// The '\r' character moves the cursor to the beginning of the line
		fmt.Printf("\rProcessing... %.2f%% complete (%d/%d)", progress, current, total)
	}
}

func printSummary(fileCount, dirCount int, c *JenkinsJson) {
	if !isOutputToStdout(c) {
		fmt.Printf("\nSummary: Processed %d files and %d directories\n", fileCount, dirCount)
	}
}

func (c *JenkinsJson) processFile(filePath string, file *os.File) subcommands.ExitStatus {
	path := filePath

	// if the user does not specify the path as
	// a command line arg, assume the default path.
	if path == "" {
		path = "jenkinsjson.json"
	}

	// open the jenkinsjson yaml
	before, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// create converter with options
	options := []jenkinsjson.Option{}
	options = append(options, jenkinsjson.WithUseIntelligence(c.useIntelligence))

	// add infrastructure options if specified
	if c.infrastructure != "" {
		options = append(options, jenkinsjson.WithInfrastructure(c.infrastructure))
	}
	if c.os != "" {
		options = append(options, jenkinsjson.WithOS(c.os))
	}
	if c.arch != "" {
		options = append(options, jenkinsjson.WithArch(c.arch))
	}

	// convert the pipeline yaml from the jenkinsjson format to the harness yaml format
	converter := jenkinsjson.New(options...)
	after, err := converter.ConvertBytes(before)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// downgrade from the v1 harness yaml format to the v0 harness yaml format
	if c.downgrade {
		// downgrade to the v0 yaml
		d := downgrader.New(
			downgrader.WithCodebase(c.repoName, c.repoConn),
			downgrader.WithDockerhub(c.dockerConn),
			downgrader.WithKubernetes(c.kubeName, c.kubeConn),
			downgrader.WithName(c.name),
			downgrader.WithOrganization(c.org),
			downgrader.WithProject(c.proj),
			downgrader.WithDefaultImage(c.defaultImage),
			downgrader.WithIntelligence(c.useIntelligence),
			downgrader.WithRandomId(c.randomId),
		)
		after, err = d.Downgrade(after)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}

	// Write the YAML
	file.WriteString("\n---\n")
	if c.beforeAfter {
		file.Write(before)
		file.WriteString("\n---\n")
	}

	file.Write(after)

	return subcommands.ExitSuccess
}

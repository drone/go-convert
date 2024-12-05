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
	"github.com/drone/go-convert/convert/harness"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	downgrade   bool
	beforeAfter bool
	outputDir   string
}

func (*JenkinsJson) Name() string     { return "jenkinsjson" }
func (*JenkinsJson) Synopsis() string { return "converts a jenkinsjson pipeline" }
func (*JenkinsJson) Usage() string {
	return `jenkinsjson [-downgrade] [jenkinsjson.json]
`
}

func (c *JenkinsJson) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.downgrade, "downgrade", false, "downgrade to the legacy yaml format")
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
	for _, arg := range f.Args() {
		recursivelyProcessFiles(arg, c)
	}
	return subcommands.ExitSuccess
}

func recursivelyProcessFiles(filename string, c *JenkinsJson) subcommands.ExitStatus {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("Cannot open input file: %s\n%v", filename, err)
	}
	if fileInfo.IsDir() {
		recursivelyProcessDir(filename, c)
	} else {
		return processFile(filename, c)
	}

	return subcommands.ExitSuccess
}

func processFile(filename string, c *JenkinsJson) subcommands.ExitStatus {
	outFile := getOutputFile(c, filename)

	if exitStatus := c.processFile(filename, outFile); exitStatus != subcommands.ExitSuccess {
		return exitStatus
	}
	return subcommands.ExitSuccess
}

func recursivelyProcessDir(filename string, c *JenkinsJson) {
	dirEntries, err := os.ReadDir(filename)
	if err != nil {
		log.Fatalf("Failed to read directory children: %s\n%v", filename, err)
	}
	for _, entry := range dirEntries {
		absFilepath := filename + "/" + entry.Name()
		recursivelyProcessFiles(absFilepath, c)
	}
}

func getOutputFile(c *JenkinsJson, filename string) *os.File {
	if c.outputDir == "" {
		return os.Stdout
	}
	return createOutputFile(c, filename)
}

func createOutputFile(c *JenkinsJson, inputFile string) *os.File {
	fileInfo, err := os.Stat(c.outputDir)
	if os.IsNotExist(err) {
		log.Fatalf("Provided output directory does not exist: %s", inputFile)
	} else if !fileInfo.IsDir() {
		log.Fatalf("Provided output path is not a directory: %s", inputFile)
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

	// convert the pipeline yaml from the jenkinsjson
	// format to the harness yaml format.
	converter := jenkinsjson.New(
		jenkinsjson.WithDockerhub(c.dockerConn),
		jenkinsjson.WithKubernetes(c.kubeName, c.kubeConn),
	)
	after, err := converter.ConvertBytes(before)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// downgrade from the v1 harness yaml format
	// to the v0 harness yaml format.
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
		)
		after, err = d.Downgrade(after)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}

	if c.beforeAfter {
		file.WriteString("---\n")
		file.Write(before)
		file.WriteString("\n---\n")
	}

	file.Write(after)

	return subcommands.ExitSuccess
}

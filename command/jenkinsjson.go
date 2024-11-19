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
	"io/ioutil"
	"log"
	"os"

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

	f.StringVar(&c.org, "org", "default", "harness organization")
	f.StringVar(&c.proj, "project", "default", "harness project")
	f.StringVar(&c.name, "pipeline", "default", "harness pipeline name")
	f.StringVar(&c.repoConn, "repo-connector", "", "repository connector")
	f.StringVar(&c.repoName, "repo-name", "", "repository name")
	f.StringVar(&c.kubeConn, "kube-connector", "", "kubernetes connector")
	f.StringVar(&c.kubeName, "kube-namespace", "", "kubernets namespace")
	f.StringVar(&c.dockerConn, "docker-connector", "", "dockerhub connector")
	f.StringVar(&c.defaultImage, "default-image", "alpine", "default image for run step")
}

func (c *JenkinsJson) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)

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
		jenkinsjson.WithDefaultImage(c.defaultImage),
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
		)
		after, err = d.Downgrade(after)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}

	if c.beforeAfter {
		os.Stdout.WriteString("---\n")
		os.Stdout.Write(before)
		os.Stdout.WriteString("\n---\n")
	}

	os.Stdout.Write(after)

	return subcommands.ExitSuccess
}

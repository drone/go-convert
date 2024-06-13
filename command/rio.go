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
	"github.com/drone/go-convert/convert/rio"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/subcommands"
)

type Rio struct {
	name       string
	proj       string
	org        string
	repoName   string
	repoConn   string
	kubeName   string
	kubeConn   string
	dockerConn string
	userGroup  string
	githubConn string

	downgrade   bool
	beforeAfter bool
}

func (*Rio) Name() string     { return "rio" }
func (*Rio) Synopsis() string { return "converts a rio pipeline" }
func (*Rio) Usage() string {
	return `rio [-downgrade] <path to rio.yml>
`
}

func (c *Rio) SetFlags(f *flag.FlagSet) {
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
	f.StringVar(&c.userGroup, "notification-user-group", "", "notification-user-group")
	f.StringVar(&c.githubConn, "github-connector", "", "github-connector")
}

func (c *Rio) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)

	fmt.Println("# v1 harness yaml is not supported for Rio pipelines. Converting to v0...")
	// if the user does not specify the path as
	// a command line arg, assume the default path.
	if path == "" {
		path = ".rio.yml"
	}

	// open the rio yaml
	before, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// convert the pipeline yaml from the rio
	// format to the harness yaml format.
	converter := rio.New(
		rio.WithDockerhub(c.dockerConn),
		rio.WithKubernetes(c.kubeName, c.kubeConn),
		rio.WithName(c.name),
		rio.WithIdentifier(c.name),
		rio.WithOrganization(c.org),
		rio.WithProject(c.proj),
		rio.WithNotifyUserGroup(c.userGroup),
		rio.WithGithubConnector(c.githubConn),
	)
	after, err := converter.ConvertBytes(before)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// downgrade from the v1 harness yaml format
	// to the v0 harness yaml format.
	if c.downgrade {
		fmt.Println("# downgrade for Rio pipeline is not supported")
		//// downgrade to the v0 yaml
		//d := downgrader.New(
		//	downgrader.WithCodebase(c.repoName, c.repoConn),
		//	downgrader.WithDockerhub(c.dockerConn),
		//	downgrader.WithKubernetes(c.kubeName, c.kubeName),
		//	downgrader.WithName(c.name),
		//	downgrader.WithOrganization(c.org),
		//	downgrader.WithProject(c.proj),
		//)
		//after, err = d.Downgrade(after)
		//if err != nil {
		//	log.Println(err)
		//	return subcommands.ExitFailure
		//}
	}

	if c.beforeAfter {
		os.Stdout.WriteString("---\n")
		os.Stdout.Write(before)
		os.Stdout.WriteString("\n---\n")
	}

	os.Stdout.Write(after)

	return subcommands.ExitSuccess
}

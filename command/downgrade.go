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

	"github.com/google/subcommands"
)

type Downgrade struct {
	name       string
	proj       string
	org        string
	repoName   string
	repoConn   string
	kubeName   string
	kubeConn   string
	dockerConn string

	downgrade   bool
	beforeAfter bool
}

func (*Downgrade) Name() string     { return "downgrade" }
func (*Downgrade) Synopsis() string { return "converts a harness pipeline to the v0 format" }
func (*Downgrade) Usage() string {
	return `downgrade <path to harness yaml>
`
}

func (c *Downgrade) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.org, "org", "default", "harness organization")
	f.StringVar(&c.proj, "project", "default", "harness project")
	f.StringVar(&c.name, "pipeline", "default", "harness pipeline name")
	f.StringVar(&c.repoConn, "repo-connector", "", "repository connector")
	f.StringVar(&c.repoName, "repo-name", "", "repository name")
	f.StringVar(&c.kubeConn, "kube-connector", "", "kubernetes connector")
	f.StringVar(&c.kubeName, "kube-namespace", "", "kubernets namespace")
	f.StringVar(&c.dockerConn, "docker-connector", "", "dockerhub connector")
}

func (c *Downgrade) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)

	var before []byte
	var err error

	// if the user provides the yaml path,
	// read the yaml file.
	if path != "" {
		before, err = ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}

	} else {
		// else read the yaml file from stdin
		before, _ = ioutil.ReadAll(os.Stdin)
	}

	// downgrade to the v0 yaml
	d := downgrader.New(
		downgrader.WithCodebase(c.repoName, c.repoConn),
		downgrader.WithDockerhub(c.dockerConn),
		downgrader.WithKubernetes(c.kubeName, c.kubeConn),
		downgrader.WithName(c.name),
		downgrader.WithOrganization(c.org),
		downgrader.WithProject(c.proj),
	)
	after, err := d.Downgrade(before)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	if c.beforeAfter {
		os.Stdout.WriteString("---\n")
		os.Stdout.Write(before)
		os.Stdout.WriteString("\n---\n")
	}

	os.Stdout.Write(after)

	return subcommands.ExitSuccess
}

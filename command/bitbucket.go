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

	"github.com/drone/go-convert/convert/bitbucket"
	"github.com/drone/go-convert/convert/harness/downgrader"

	"github.com/google/subcommands"
)

type Bitbucket struct {
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

func (*Bitbucket) Name() string     { return "bitbucket" }
func (*Bitbucket) Synopsis() string { return "converts a bitbucket pipeline" }
func (*Bitbucket) Usage() string {
	return `bitbucket [-downgrade] <path to bitbucket.yml>
`
}

func (p *Bitbucket) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.downgrade, "downgrade", false, "downgrade to the legacy yaml format")
	f.BoolVar(&p.beforeAfter, "before-after", false, "print the befor and after")

	f.StringVar(&p.org, "org", "default", "harness organization")
	f.StringVar(&p.proj, "project", "default", "harness project")
	f.StringVar(&p.name, "pipeline", "default", "harness pipeline name")
	f.StringVar(&p.repoConn, "repo-connector", "", "repository connector")
	f.StringVar(&p.repoName, "repo-name", "", "repository name")
	f.StringVar(&p.kubeConn, "kube-connector", "", "kubernetes connector")
	f.StringVar(&p.kubeName, "kube-namespace", "", "kubernets namespace")
	f.StringVar(&p.dockerConn, "docker-connector", "", "dockerhub connector")
}

func (p *Bitbucket) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)

	// if the user does not specify the path as
	// a command line arg, assume the default path.
	if path == "" {
		path = "bitbucket-pipelines.yml"
	}

	// open the bitbucket yaml
	before, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// convert the bitbucket yaml from the bitbucket
	// format to the harness format.
	after, err := bitbucket.FromBytes(before)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}

	// downgrade from the v1 harness yaml format
	// to the v0 harness yaml format.
	if p.downgrade {
		// downgrade to the v0 yaml
		d := downgrader.New(
			downgrader.WithCodebase(p.repoName, p.repoConn),
			downgrader.WithDockerhub(p.dockerConn),
			downgrader.WithKubernetes(p.kubeName, p.kubeConn),
			downgrader.WithName(p.name),
			downgrader.WithOrganization(p.org),
			downgrader.WithProject(p.proj),
		)
		after, err = d.Downgrade(after)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}

	if p.beforeAfter {
		os.Stdout.WriteString("---\n")
		os.Stdout.Write(before)
		os.Stdout.WriteString("\n---\n")
	}

	os.Stdout.Write(after)

	return subcommands.ExitSuccess
}

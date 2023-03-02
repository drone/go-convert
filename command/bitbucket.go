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
	"github.com/drone/go-convert/convert/downgrade"

	harness "github.com/drone/spec/dist/go"

	"github.com/google/subcommands"
)

type Bitbucket struct {
	name string
	proj string
	org  string
	repo string
	conn string

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
	f.StringVar(&p.conn, "repo-connector", "", "repository connector")
	f.StringVar(&p.repo, "repo-name", "", "repository name")
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
		// unmarshal to the v1 yaml
		v, err := harness.ParseBytes(after)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
		// downgrade to the v0 yaml
		after, err = downgrade.Downgrade(v, downgrade.Args{
			Name:         p.name,
			Organization: p.org,
			Project:      p.proj,
			Docker:       downgrade.Docker{},     // TODO
			Kubernetes:   downgrade.Kubernetes{}, // TODO
			Codebase: downgrade.Codebase{
				Connector: p.conn,
				Repo:      p.repo,
			},
		})
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}
	}

	if p.beforeAfter {
		os.Stdout.WriteString("---\n")
		os.Stdout.Write(before)
		os.Stdout.WriteString("---\n")
	}

	os.Stdout.Write(after)

	return subcommands.ExitSuccess
}

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

package yaml

// TODO support for undocumented addons (srcclr, ssh_known_hosts)
// https://config.travis-ci.com/ref/job/addons

type Addons struct {
	Apt          *Apt          `yaml:"apt,omitempty"`
	AptPackage   Stringorslice `yaml:"apt_package,omitempty"`
	Artifacts    *Artifacts    `yaml:"artifacts,omitempty"`
	Browserstack *Browserstack `yaml:"browserstack,omitempty"`
	Chrome       string        `yaml:"chrome,omitempty"`
	Codeclimate  *Codeclimate  `yaml:"codeclimate,omitempty"`
	Coverity     *Coverity     `yaml:"coverity_scan,omitempty"`
	Firefox      string        `yaml:"firefox,omitempty"`
	Homebrew     *Homebrew     `yaml:"homebrew,omitempty"`
	Hostname     string        `yaml:"hostname,omitempty"`
	Hosts        Stringorslice `yaml:"hosts,omitempty"`
	Mariadb      string        `yaml:"mariadb,omitempty"`
	Postgres     string        `yaml:"postgresql,omitempty"`
	Postgresql   string        `yaml:"postgres,omitempty"`
	Rethinkdb    string        `yaml:"rethinkdb,omitempty"`
	Sauce        *Sauce        `yaml:"sauce_connect,omitempty"`
	Snaps        *Snaps        `yaml:"snaps,omitempty"`
	Sonarcloud   *Sonarcloud   `yaml:"sonarcloud,omitempty"`
}

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

package travis

import (
	"strings"

	travis "github.com/drone/go-convert/convert/travis/yaml"
	harness "github.com/drone/spec/dist/go"
)

func (d *Converter) convertAddons(ctx *context) []*harness.Step {
	addons := ctx.config.Addons

	// return if no addons defined
	if addons == nil {
		return nil
	}

	// aggregate list of steps
	var dst []*harness.Step

	// TODO support for addons.Chrome
	// TODO support for addons.Firefox
	// TODO support for addons.Hostname
	// TODO support for addons.Hosts

	if v := addons.Apt; v != nil {
		// TODO only run apt step if `os` is linux
		dst = append(dst, d.convertApt(v))
	}
	if v := addons.AptPackage; len(v) != 0 {
		// TODO only run apt step if `os` is linux
		dst = append(dst, d.convertAptPackage(v))
	}
	if v := addons.Artifacts; v != nil {
		// TODO support addons.artifacts
		// https://config.travis-ci.com/ref/job/addons/artifacts
	}
	if v := addons.Browserstack; v != nil {
		// TODO support addons.browsertack
		// https://config.travis-ci.com/ref/job/addons/browserstack
	}
	if v := addons.Codeclimate; v != nil {
		// TODO support addons.codeclimate
		// https://config.travis-ci.com/ref/job/addons/code_climate
	}
	if v := addons.Coverity; v != nil {
		// TODO support addons.coverity
		// https://config.travis-ci.com/ref/job/addons/coverity_scan
	}
	if v := addons.Homebrew; v != nil {
		// TODO only run homebrew step if `os` is macos
		dst = append(dst, d.convertHomebrew(v))
	}
	if v := addons.Sauce; v != nil {
		// TODO support addons.sauce_connect
		// https://config.travis-ci.com/ref/job/addons/sauce_connect
	}
	if v := addons.Snaps; v != nil {
		// TODO support addons.snaps
		// https://config.travis-ci.com/ref/job/addons/snaps
	}
	if v := addons.Sonarcloud; v != nil {
		// TODO support addons.sonarcloud
		// https://config.travis-ci.com/ref/job/addons/sonarcloud
	}

	// return if no steps added
	if len(dst) == 0 {
		return nil
	}
	return dst
}

func (d *Converter) convertApt(apt *travis.Apt) *harness.Step {
	var lines []string
	if apt.Enabled {
		// TODO research behavior of `addons: { apt: true }`
	}
	if apt.Update {
		lines = append(lines, "sudo apt-get -qq update")
	}
	for _, s := range apt.Sources {
		// TODO support addon.apt.sources
		if s == nil {
			continue
		}
	}
	for _, s := range apt.Packages {
		lines = append(lines, "sudo apt-get -y install "+s)
	}
	if s := apt.Dist; s != "" {
		// TODO support addon.apt.dist
		// https: //docs.travis-ci.com/user/installing-dependencies/#adding-apt-packages
	}

	if len(lines) == 0 {
		return nil
	}

	return &harness.Step{
		Name: d.identifiers.Generate("apt"),
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(lines, "\n"),
		},
	}
}

func (d *Converter) convertAptPackage(packages travis.Stringorslice) *harness.Step {
	lines := []string{"sudo apt-get -qq update"}
	for _, s := range packages {
		lines = append(lines, "sudo apt-get -y install "+s)
	}
	return &harness.Step{
		Name: d.identifiers.Generate("apt_packages"),
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(lines, "\n"),
		},
	}
}

func (d *Converter) convertHomebrew(homebrew *travis.Homebrew) *harness.Step {
	var lines []string
	if homebrew.Update {
		lines = append(lines, "brew update")
	}
	for _, v := range homebrew.Taps {
		lines = append(lines, "brew tap "+v)
	}
	for _, v := range homebrew.Casks {
		lines = append(lines, "brew cask install "+v)
	}
	for _, v := range homebrew.Packages {
		lines = append(lines, "brew install "+v)
	}
	if v := homebrew.Brewfile; v != "" {
		lines = append(lines, "brew bundle --file="+v)
	}
	if len(lines) == 0 {
		return nil
	}

	return &harness.Step{
		Name: d.identifiers.Generate("homebrew"),
		Type: "script",
		Spec: &harness.StepExec{
			Run: strings.Join(lines, "\n"),
		},
	}
}

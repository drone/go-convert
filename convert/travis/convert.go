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

// Package travis converts Travis pipelines to Harness pipelines.
package travis

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	travis "github.com/drone/go-convert/convert/travis/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// as we walk the yaml, we store a
// a snapshot of the current node and
// its parents.
type context struct {
	config *travis.Pipeline
}

// Converter converts a Travis pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers
}

// New creates a new Converter that converts a Travis
// pipeline to a Harness v1 pipeline.
func New(options ...Option) *Converter {
	d := new(Converter)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = store.New()

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default kubernetes namespace.
	if d.kubeNamespace == "" {
		d.kubeNamespace = "default"
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	config, err := travis.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(&context{
		config: config,
	})
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.Convert(
		bytes.NewBuffer(b),
	)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.Convert(
		bytes.NewBufferString(s),
	)
}

// ConvertFile downgrades a v1 pipeline.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

// converts converts a Travis pipeline to a Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {

	// create the harness pipeline spec
	pipeline := &harness.Pipeline{}

	// create the harness pipeline resource
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    pipeline,
	}

	// convert the clone
	if v := convertGit(ctx); v != nil {
		pipeline.Options = new(harness.Default)
		pipeline.Options.Clone = v
	}

	// conver pipeilne stages
	pipeline.Stages = append(pipeline.Stages, &harness.Stage{
		Name:     "pipeline",
		Desc:     "converted from travis.yml",
		Type:     "ci",
		Delegate: nil, // No Travis equivalent
		Failure:  nil, // No Travis equivalent
		Strategy: convertStrategy(ctx),
		When:     nil, // TODO convert travis condition (if, branches)
		Spec: &harness.StageCI{
			Cache: convertCache(ctx),
			// TODO support for other env variabes, like TRAVIS_RETHINKDB_VERSION
			Envs:     createMatrixEnvs(ctx),
			Platform: convertPlatform(ctx),
			Runtime:  nil, // TODO convert runtime
			Steps:    d.convertSteps(ctx),
		},
	})

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d *Converter) convertSteps(ctx *context) []*harness.Step {
	var steps []*harness.Step

	// convert addon steps
	steps = append(steps, d.convertAddons(ctx)...)

	// convert services to background steps
	steps = append(steps, d.convertServices(ctx)...)

	// from the job lifecycle documentation
	// https://docs.travis-ci.com/user/job-lifecycle/#the-job-lifecycle
	for _, script := range ctx.config.BeforeInstall {
		steps = append(steps, d.convertStep(ctx, "before_install", script))
	}
	for _, script := range ctx.config.Install {
		steps = append(steps, d.convertStep(ctx, "install", script))
	}
	if len(ctx.config.Install) == 0 {
		// when no install is defined, travis may automatically
		// provide the install based on langauge.
		if script, ok := defaultInstall[strings.ToLower(ctx.config.Language)]; ok {
			steps = append(steps, d.convertStep(ctx, "install", script))
		}
	}
	for _, script := range ctx.config.BeforeScript {
		steps = append(steps, d.convertStep(ctx, "before_script", script))
	}
	for _, script := range ctx.config.Script {
		steps = append(steps, d.convertStep(ctx, "script", script))
	}
	if len(ctx.config.Script) == 0 {
		// when no script is defined, travis may automatically
		// provide the script based on langauge.
		if script, ok := defaultScript[strings.ToLower(ctx.config.Language)]; ok {
			steps = append(steps, d.convertStep(ctx, "script", script))
		}
	}
	for _, script := range ctx.config.BeforeCache {
		steps = append(steps, d.convertStep(ctx, "before_cache", script))
	}
	for _, script := range ctx.config.AfterSuccess {
		steps = append(steps, d.convertStep(ctx, "after_success", script))
	}
	for _, script := range ctx.config.AfterFailure {
		steps = append(steps, d.convertStep(ctx, "after_failure", script))
	}
	for _, script := range ctx.config.BeforeDeploy {
		steps = append(steps, d.convertStep(ctx, "before_deploy", script))
	}
	//
	// TODO support deploy steps
	//
	for _, script := range ctx.config.AfterDeploy {
		steps = append(steps, d.convertStep(ctx, "after_deploy", script))
	}
	for _, script := range ctx.config.AfterScript {
		steps = append(steps, d.convertStep(ctx, "after_script", script))
	}
	return steps
}

func (d *Converter) convertStep(ctx *context, section, command string) *harness.Step {
	return &harness.Step{
		Name: d.identifiers.Generate(section),
		// Desc: "",
		Type: "script",
		// Timeout: 0,
		// When: convertCond(src.When),
		// On: nil,
		Spec: &harness.StepExec{
			Image:     convertImageMaybe(ctx, d.kubeEnabled),
			Connector: d.dockerhubConn,
			// Mount:      convertMounts(src.Volumes),
			// Privileged: src.Privileged,
			// Pull:       convertPull(src.Pull),
			// Shell:      convertShell(),
			// User:       src.User,
			// Group:      src.Group,
			// Network:    "",
			// Entrypoint: convertEntrypoint(src.Entrypoint),
			// Args:       convertArgs(src.Entrypoint, src.Command),
			Run: command,
			// Envs:       convertVariables(src.Environment),
			// Resources:  convertResourceLimits(&src.Resource),
			// Reports:    nil,
		},
	}
}

func convertStrategy(ctx *context) *harness.Strategy {
	// TODO env.matrix
	// TODO jobs
	// TODO dart_tasks

	// https://config.travis-ci.com/matrix_expansion
	spec := &harness.Matrix{}
	spec.Axis = map[string][]string{}

	// helper function to append the axis
	// to the matrix definition.
	appendAxis := func(name string, items []string) {
		// ignore empty matrix
		if len(items) > 0 {
			var temp []string
			for _, item := range items {
				item = strings.ReplaceAll(item, "1.x", "1")
				temp = append(temp, item)
			}
			spec.Axis[name] = temp
		}
	}

	appendAxis("compiler", ctx.config.Compiler)
	appendAxis("crystal", ctx.config.Crystal)
	appendAxis("d", ctx.config.D)
	appendAxis("dart", ctx.config.Dart)
	appendAxis("dotnet", ctx.config.Dotnet)
	appendAxis("mono", ctx.config.DotnetMono)
	appendAxis("solution", ctx.config.DotnetSolution)
	appendAxis("elixir", ctx.config.Elixir)
	appendAxis("elm", ctx.config.Elm)
	appendAxis("otp_release", ctx.config.ErlangOTP)
	appendAxis("go", ctx.config.Go)
	appendAxis("hhvm", ctx.config.HHVM)
	appendAxis("haxe", ctx.config.Haxe)
	appendAxis("ghc", ctx.config.GHC)
	appendAxis("jdk", ctx.config.JDK)
	appendAxis("node_js", ctx.config.Node)
	appendAxis("julia", ctx.config.Julia)
	appendAxis("matlab", ctx.config.Matlab)
	appendAxis("nix", ctx.config.Nix)
	appendAxis("xcode_scheme", ctx.config.XcodeScheme)
	appendAxis("xcode_sdk", ctx.config.XcodeSDK)
	appendAxis("php", ctx.config.PHP)
	appendAxis("perl", ctx.config.Perl)
	appendAxis("perl6", ctx.config.Perl6)
	appendAxis("python", ctx.config.Python)
	appendAxis("r", ctx.config.R)
	appendAxis("rvm", append(ctx.config.RubyRVM, append(ctx.config.Ruby, ctx.config.RubyRBenv...)...)) // ruby, rvm, rbenv
	appendAxis("gemfile", append(ctx.config.RubyGemfile, ctx.config.RubyGemfiles...))                  // gemfile, gemfiles
	appendAxis("rust", ctx.config.Rust)
	appendAxis("scala", ctx.config.Scala)
	appendAxis("smalltalk", ctx.config.Smalltalk)
	appendAxis("smalltalk_config", ctx.config.SmalltalkConfig)
	appendAxis("smalltalk_vm", ctx.config.SmalltalkVM)
	appendAxis("os", ctx.config.OS)
	appendAxis("arch", ctx.config.Arch)
	if len(spec.Axis) == 0 {
		return nil
	}
	return &harness.Strategy{
		Type: "matrix",
		Spec: spec,
	}
}

func createMatrixEnvs(ctx *context) map[string]string {
	// https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
	envs := map[string]string{}

	appendEnvs := func(name, env string, slice []string) {
		switch len(slice) {
		case 0:
		case 1:
			if s := slice[0]; s != "" {
				envs[env] = slice[0]
			}
		default:
			envs[env] = fmt.Sprintf("<+matrix.%s>", name)
		}
	}

	appendEnvs("compiler", "TRAVIS_COMPILER", ctx.config.Compiler)
	appendEnvs("crystal", "TRAVIS_CRYSTAL_VERSION", ctx.config.Crystal)
	appendEnvs("d", "TRAVIS_D_VERSION", ctx.config.D)
	appendEnvs("dart", "TRAVIS_DART_VERSION", ctx.config.Dart)
	appendEnvs("dotnet", "TRAVIS_DOTNET_VERSION", ctx.config.Dotnet)
	appendEnvs("mono", "TRAVIS_MONO_VERSION", ctx.config.DotnetMono)
	appendEnvs("solution", "TRAVIS_SOLUTION_VERSION", ctx.config.DotnetSolution)
	appendEnvs("elixir", "TRAVIS_ELIXIR_VERSION", ctx.config.Elixir)
	appendEnvs("elm", "TRAVIS_ELM_VERSION", ctx.config.Elm)
	appendEnvs("otp_release", "TRAVIS_OTP_RELEASE", ctx.config.ErlangOTP)
	appendEnvs("go", "TRAVIS_GO_VERSION", ctx.config.Go)
	appendEnvs("hhvm", "TRAVIS_HHVM_VERSION", ctx.config.HHVM)
	appendEnvs("haxe", "TRAVIS_HAXE_VERSION", ctx.config.Haxe)
	appendEnvs("gemfile", "TRAVIS_GEMFILE_VERSION", append(ctx.config.RubyGemfile, ctx.config.RubyGemfiles...))
	appendEnvs("ghc", "TRAVIS_GHC_VERSION", ctx.config.GHC)
	appendEnvs("jdk", "TRAVIS_JDK_VERSION", ctx.config.JDK)
	appendEnvs("node_js", "TRAVIS_NODE_VERSION", ctx.config.Node)
	appendEnvs("julia", "TRAVIS_JULIA_VERSION", ctx.config.Julia)
	appendEnvs("matlab", "TRAVIS_MATLAB_VERSION", ctx.config.Matlab)
	appendEnvs("nix", "TRAVIS_NIX_VERSION", ctx.config.Nix)
	appendEnvs("xcode_scheme", "TRAVIS_XCODE_SCHEME", ctx.config.XcodeScheme)
	appendEnvs("xcode_sdk", "TRAVIS_XCODE_SDK", ctx.config.XcodeSDK)
	appendEnvs("php", "TRAVIS_PHP_VERSION", ctx.config.PHP)
	appendEnvs("perl", "TRAVIS_PERL_VERSION", ctx.config.Perl)
	appendEnvs("perl6", "TRAVIS_PERL6_VERSION", ctx.config.Perl6)
	appendEnvs("python", "TRAVIS_PYTHON_VERSION", ctx.config.Python)
	appendEnvs("r", "TRAVIS_R_VERSION", ctx.config.R)
	appendEnvs("rust", "TRAVIS_RUST_VERSION", ctx.config.Rust)
	appendEnvs("rvm", "TRAVIS_RUBY_VERSION", append(ctx.config.Ruby, append(ctx.config.RubyRVM, ctx.config.RubyRBenv...)...))
	appendEnvs("scala", "TRAVIS_SCALA_VERSION", ctx.config.Scala)
	appendEnvs("smalltalk", "TRAVIS_SMALLTALK_VERSION", ctx.config.Smalltalk)
	appendEnvs("smalltalk_config", "TRAVIS_SMALLTALK_CONFIG", ctx.config.SmalltalkConfig)
	appendEnvs("smalltalk_vm", "TRAVIS_SMALLTALK_VM", ctx.config.SmalltalkVM)
	appendEnvs("xcode_project", "TRAVIS_XCODE_PROJECT", []string{ctx.config.XcodeProject})

	// append ruby alias
	if env, ok := envs["TRAVIS_RUBY_VERSION"]; ok {
		envs["TRAVIS_RVM_VERSION"] = env
	}

	if len(envs) == 0 {
		return nil
	}
	return envs
}

func convertPlatform(ctx *context) *harness.Platform {
	var os, arch string

	switch len(ctx.config.OS) {
	case 0:
	case 1:
		os = ctx.config.OS[0]
	default:
		os = "<+matrix.os>"
	}

	switch len(ctx.config.Arch) {
	case 0:
	case 1:
		arch = ctx.config.Arch[0]
	default:
		arch = "<+matrix.arch>"
	}

	// normalize os
	switch os {
	case "mac", "ios", "osx":
		os = "macos"
	}

	// return a nil platform if empty which instructs
	// harness to use the platform defaults.
	if os == "" && arch == "" {
		return nil
	}
	// return &harness.Platform{
	// 	Os: os, Arch: arch,
	// }
	return nil // TODO `os` and `arch` cannot be enums to support matrix
}

func convertGit(ctx *context) *harness.Clone {
	src := ctx.config.Git
	if src == nil {
		return nil
	}
	dst := new(harness.Clone)
	if src.Depth != nil {
		dst.Depth = int64(src.Depth.Value)
	}
	// TODO git support for submodules
	// TODO git support for submodules_depth
	// TODO git support for lfs_skip_smudge
	// TODO git support for sparse_checkout
	// TODO git support for autocrlf
	return dst
}

func convertCache(ctx *context) *harness.Cache {
	src := ctx.config.Cache
	if src == nil {
		return nil
	}
	dst := new(harness.Cache)
	dst.Enabled = true
	dst.Paths = append(dst.Paths, src.Directories...)

	if src.Apt {
		// behavior not documented
		// https://docs.travis-ci.com/user/caching/
	}
	if src.Bundler {
		dst.Paths = append(dst.Paths, "~/.rvm")
		dst.Paths = append(dst.Paths, "vendor/bundle")
	}
	if src.Cargo {
		dst.Paths = append(dst.Paths, "target")
		dst.Paths = append(dst.Paths, "~/.cargo")
	}
	if src.Ccache {
		dst.Paths = append(dst.Paths, "~/.ccache")
	}
	if src.Cocoapods {
		// paths not documented
		// https://docs.travis-ci.com/user/caching/
	}
	if src.Edge {
		// behavior not documented
		// https://docs.travis-ci.com/user/caching/
	}
	if src.Npm {
		dst.Paths = append(dst.Paths, "~/.npm")
		dst.Paths = append(dst.Paths, "node_modules")
	}
	if src.Packages {
		dst.Paths = append(dst.Paths, "~/R/Library")
	}
	if src.Pip {
		dst.Paths = append(dst.Paths, "~/.cache/pip")
	}
	if src.Yarn {
		dst.Paths = append(dst.Paths, "~/.cache/yarn")
	}

	// TODO when caching R packages, set R_LIB_USER=~/R/Library
	// TODO cache support for `branch`
	// TODO cache support for `timeout`

	return dst
}

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

	// create the harness pipeline
	pipeline := &harness.Pipeline{
		Version: 1,
		Options: nil, // TODO
	}

	// conver pipeilne stages
	pipeline.Stages = append(pipeline.Stages, &harness.Stage{
		Name:     "pipeline",
		Desc:     "converted from travis.yml",
		Type:     "ci",
		Delegate: nil, // No Travis equivalent
		On:       nil, // No Travis equivalent
		Strategy: convertStrategy(ctx),
		// When:     convertCond(from.Trigger),
		Spec: &harness.StageCI{
			// Cache:    convertCache(from.Cache)
			// Clone:    convertClone(from.Clone),
			Envs: createMatrixEnvs(ctx),
			// Platform: convertPlatform(from.Platform),
			// Runtime:  convertRuntime(from),
			Steps: d.convertSteps(ctx),
			// Volumes:  convertVolumes(from.Volumes),

			// TODO support for delegate.selectors from from.Node
			// TODO support for stage.variables
		},
	})

	// marshal the harness yaml
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (d *Converter) convertSteps(ctx *context) []*harness.Step {
	var steps []*harness.Step
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
		if script, ok := defaultInstall[ctx.config.Language]; ok {
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
		if script, ok := defaultScript[ctx.config.Language]; ok {
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
	// TODO deploy section
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
	// TODO os
	// TODO arch
	// TODO env.matrix
	// TODO compiler
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
			for _, item := range ctx.config.Go {
				item = strings.ReplaceAll(item, "1.x", "1")
				temp = append(temp, item)
			}
			spec.Axis[name] = temp
		}
	}

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

	if len(ctx.config.Crystal) > 0 {
		envs["TRAVIS_CRYSTAL_VERSION"] = "<+matrix.crystal>"
	}
	if len(ctx.config.D) > 0 {
		envs["TRAVIS_D_VERSION"] = "<+matrix.d>"
	}
	if len(ctx.config.Dart) > 0 {
		envs["TRAVIS_DART_VERSION"] = "<+matrix.dart>"
	}
	if len(ctx.config.Dotnet) > 0 {
		envs["TRAVIS_DOTNET_VERSION"] = "<+matrix.dotnet>"
	}
	if len(ctx.config.DotnetMono) > 0 {
		envs["TRAVIS_MONO_VERSION"] = "<+matrix.mono>"
	}
	if len(ctx.config.DotnetSolution) > 0 {
		envs["TRAVIS_SOLUTION_VERSION"] = "<+matrix.solution>"
	}
	if len(ctx.config.Elixir) > 0 {
		envs["TRAVIS_ELIXIR_VERSION"] = "<+matrix.elixir>"
	}
	if len(ctx.config.Elm) > 0 {
		envs["TRAVIS_ELM_VERSION"] = "<+matrix.elm>"
	}
	if len(ctx.config.ErlangOTP) > 0 {
		envs["TRAVIS_OTP_RELEASE"] = "<+matrix.otp_release>"
	}
	if len(ctx.config.Go) > 0 {
		envs["TRAVIS_GO_VERSION"] = "<+matrix.go>"
	}
	if len(ctx.config.HHVM) > 0 {
		envs["TRAVIS_HHVM_VERSION"] = "<+matrix.hhvm>"
	}
	if len(ctx.config.Haxe) > 0 {
		envs["TRAVIS_HAXE_VERSION"] = "<+matrix.haxe>"
	}
	if len(ctx.config.GHC) > 0 {
		envs["TRAVIS_GHC_VERSION"] = "<+matrix.ghc>"
	}
	if len(ctx.config.JDK) > 0 {
		envs["TRAVIS_JDK_VERSION"] = "<+matrix.jdk>"
	}
	if len(ctx.config.Node) > 0 {
		envs["TRAVIS_NODE_VERSION"] = "<+matrix.node_js>"
	}
	if len(ctx.config.Julia) > 0 {
		envs["TRAVIS_JULIA_VERSION"] = "<+matrix.julia>"
	}
	if len(ctx.config.Matlab) > 0 {
		envs["TRAVIS_MATLAB_VERSION"] = "<+matrix.matlab>"
	}
	if len(ctx.config.Nix) > 0 {
		envs["TRAVIS_NIX_VERSION"] = "<+matrix.nix>"
	}
	if len(ctx.config.XcodeScheme) > 0 {
		envs["TRAVIS_XCODE_SCHEME"] = "<+matrix.xcode_scheme>"
	}
	if len(ctx.config.XcodeSDK) > 0 {
		envs["TRAVIS_XCODE_SDK"] = "<+matrix.xcode_sdk>"
	}
	if s := ctx.config.XcodeProject; s != "" {
		envs["TRAVIS_XCODE_PROJECT"] = s
	}
	if len(ctx.config.PHP) > 0 {
		envs["TRAVIS_PHP_VERSION"] = "<+matrix.php>"
	}
	if len(ctx.config.Perl) > 0 {
		envs["TRAVIS_PERL_VERSION"] = "<+matrix.perl>"
	}
	if len(ctx.config.Perl6) > 0 {
		envs["TRAVIS_PERL6_VERSION"] = "<+matrix.perl6>"
	}
	if len(ctx.config.Python) > 0 {
		envs["TRAVIS_PYTHON_VERSION"] = "<+matrix.python>"
	}
	if len(ctx.config.R) > 0 {
		envs["TRAVIS_R_VERSION"] = "<+matrix.r>"
	}
	if len(ctx.config.Ruby) > 0 || len(ctx.config.RubyRVM) > 0 || len(ctx.config.RubyRBenv) > 0 {
		envs["TRAVIS_RVM_VERSION"] = "<+matrix.rvm>"  // TODO verify rvm environment variable
		envs["TRAVIS_RUBY_VERSION"] = "<+matrix.rvm>" // TODO verify ruby environment variable
	}
	if len(ctx.config.RubyGemfile) > 0 || len(ctx.config.RubyGemfiles) > 0 {
		// TODO verify gemfile environment variable
		envs["TRAVIS_GEMFILE_VERSION"] = "<+matrix.gemfile>"
	}
	if len(ctx.config.Rust) > 0 {
		envs["TRAVIS_RUST_VERSION"] = "<+matrix.rust>"
	}
	if len(ctx.config.Scala) > 0 {
		envs["TRAVIS_SCALA_VERSION"] = "<+matrix.scala>"
	}
	if len(ctx.config.Smalltalk) > 0 {
		// TODO verify smalltalk version environment variable
		envs["TRAVIS_SMALLTALK_VERSION"] = "<+matrix.smalltalk>"
	}
	if len(ctx.config.SmalltalkConfig) > 0 {
		// TODO verify smalltalk config environment variable
		envs["TRAVIS_SMALLTALK_CONFIG"] = "<+matrix.smalltalk_config>"
	}
	if len(ctx.config.SmalltalkVM) > 0 {
		// TODO verify smalltalk vm environment variable
		envs["TRAVIS_SMALLTALK_VM"] = "<+matrix.smalltalk_vm>"
	}
	if len(envs) == 0 {
		return nil
	}
	return envs
}

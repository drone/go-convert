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

type (
	// Pipeline defines a gitlab pipeline.
	// https://config.travis-ci.com/index
	Pipeline struct {
		Language      string         `yaml:"language,omitempty"`
		OS            Stringorslice  `yaml:"os,omitempty"`
		OSXImage      Stringorslice  `yaml:"osx_image,omitempty"`
		Arch          Stringorslice  `yaml:"arch,omitempty"`
		Addons        *Addons        `yaml:"addons,omitempty"`
		Branches      *Branches      `yaml:"branches,omitempty"`
		Cache         *Cache         `yaml:"cache,omitempty"`
		Compiler      Stringorslice  `yaml:"compiler,omitempty"`
		Deploy        interface{}    `yaml:"deploy,omitempty"` // TODO
		Dist          string         `yaml:"dist,omitempty"`
		Env           *Env           `yaml:"env,omitempty"`
		Git           *Git           `yaml:"git,omitempty"`
		If            string         `yaml:"if,omitempty"`
		Import        *Imports       `yaml:"import,omitempty"`
		Jobs          *Jobs          `yaml:"jobs,omitempty"`
		Notifications *Notifications `yaml:"notifications,omitempty"` // TODO
		Services      Stringorslice  `yaml:"services,omitempty"`
		Stages        *Stages        `yaml:"stages,omitempty"`
		Version       string         `yaml:"version,omitempty"`

		BeforeInstall Stringorslice `yaml:"before_install,omitempty"`
		Install       Stringorslice `yaml:"install,omitempty"`
		BeforeScript  Stringorslice `yaml:"before_script,omitempty"`
		Script        Stringorslice `yaml:"script,omitempty"`
		BeforeCache   Stringorslice `yaml:"before_cache,omitempty"`
		AfterSuccess  Stringorslice `yaml:"after_success,omitempty"`
		AfterFailure  Stringorslice `yaml:"after_failure,omitempty"`
		BeforeDeploy  Stringorslice `yaml:"before_deploy,omitempty"`
		AfterDeploy   Stringorslice `yaml:"after_deploy,omitempty"`
		AfterScript   Stringorslice `yaml:"after_script,omitempty"`

		//
		// Language Keywords
		//

		// Android
		// https://config.travis-ci.com/ref/language/android
		Android interface{} `yaml:"android,omitempty"` // matrix expand key (TODO)

		// Clojure
		Lein string `yaml:"lein,omitempty"`

		// Crystal
		Crystal Stringorslice `yaml:"crystal,omitempty"` // matrix expand key

		// D
		D Stringorslice `yaml:"d,omitempty"` // matrix expand key

		// Dart Language
		Dart                 Stringorslice `yaml:"dart,omitempty"`      // matrix expand key
		DartTask             interface{}   `yaml:"dart_task,omitempty"` // matrix expand key (TODO)
		DartWithContentShell bool          `yaml:"with_content_shell,omitempty"`

		// Dotnet Language (C#)
		Dotnet         Stringorslice `yaml:"dotnet,omitempty"`   // matrix expand key
		DotnetMono     Stringorslice `yaml:"mono,omitempty"`     // matrix expand key
		DotnetSolution Stringorslice `yaml:"solution,omitempty"` // matrix expand key

		// Elixir Language
		Elixir Stringorslice `yaml:"elixir,omitempty"` // matrix expand key

		// Elm Language
		Elm       Stringorslice `yaml:"elm,omitempty"` // matrix expand key
		ElmFormat string        `yaml:"elm_format,omitempty"`
		ElmTest   string        `yaml:"elm_test,omitempty"`

		// Erlang Language
		ErlangOTP Stringorslice `yaml:"otp_release,omitempty"` // matrix expand key

		// Go Language
		Go           Stringorslice `yaml:"go,omitempty"` // matrix expand key
		GoBuildArgs  string        `yaml:"gobuild_args,omitempty"`
		GoImportPath string        `yaml:"go_import_path,omitempty"`

		// Hack Language
		HHVM Stringorslice `yaml:"hhvm,omitempty"` // matrix expand key

		// Haxe Language
		Haxe     Stringorslice `yaml:"haxe,omitempty"` // matrix expand key
		HaxeXML  Stringorslice `yaml:"hxml,omitempty"`
		HaxeNeko string        `yaml:"neko,omitempty"`

		// Haskell Language
		GHC Stringorslice `yaml:"ghc,omitempty"` // matrix expand key

		// Java Language
		JDK Stringorslice `yaml:"jdk,omitempty"` // matrix expand key

		// Javascript Language
		Node        Stringorslice `yaml:"node_js,omitempty"` // matrix expand key
		NodeNpmArgs string        `yaml:"npm_args,omitempty"`

		// Julia Language
		Julia Stringorslice `yaml:"julia,omitempty"` // matrix expand key

		// Matlab Language
		Matlab Stringorslice `yaml:"matlab,omitempty"` // matrix expand key

		// Nix Language
		Nix Stringorslice `yaml:"nix,omitempty"` // matrix expand key

		// ObjectiveC
		XcodeScheme      Stringorslice `yaml:"xcode_scheme,omitempty"` // matrix expand key
		XcodeSDK         Stringorslice `yaml:"xcode_sdk,omitempty"`    // matrix expand key
		XcodeDestination string        `yaml:"xcode_destination,omitempty"`
		XcodeProject     string        `yaml:"xcode_project,omitempty"`
		XcodeToolArgs    string        `yaml:"xctool_args,omitempty"`
		XcodePodfile     string        `yaml:"xctool_podfile,omitempty"`

		// Php Language
		PHP             Stringorslice `yaml:"php,omitempty"` // matrix expand key
		PHPComposerArgs string        `yaml:"composer_args,omitempty"`

		// Perl Language
		Perl  Stringorslice `yaml:"perl,omitempty"`  // matrix expand key
		Perl6 Stringorslice `yaml:"perl6,omitempty"` // matrix expand key

		// Python Langauge
		Python                Stringorslice          `yaml:"python,omitempty"` // matrix expand key
		PythonVirtualenv      map[string]interface{} `yaml:"virtualenv,omitempty"`
		PythonVirtualenvAlias map[string]interface{} `yaml:"virtual_env,omitempty"`

		// R Langauge
		// https://config.travis-ci.com/ref/language/r
		R               Stringorslice `yaml:"r,omitempty"` // matrix expand key
		RPackage        Stringorslice `yaml:"r_packages,omitempty"`
		RBinaryPackages Stringorslice `yaml:"r_binary_packages,omitempty"`
		RGithubPackage  Stringorslice `yaml:"r_github_packages,omitempty"`

		// Ruby Language
		Ruby            Stringorslice `yaml:"ruby,omitempty"`  // matrix expand key
		RubyRVM         Stringorslice `yaml:"rvm,omitempty"`   // alias for ruby
		RubyRBenv       Stringorslice `yaml:"rbenv,omitempty"` // alias for ruby
		RubyGemfile     Stringorslice `yaml:"gemfile,omitempty"`
		RubyGemfiles    Stringorslice `yaml:"gemfiles,omitempty"` // alias for gemfile
		RubyBundlerArgs string        `yaml:"bundler_args,omitempty"`

		// Rust Langauge
		Rust Stringorslice `yaml:"rust,omitempty"` // matrix expand key

		// Scala Langauge
		Scala        Stringorslice `yaml:"scala,omitempty"` // matrix expand key
		ScalaSbtArgs string        `yaml:"sbt_args,omitempty"`

		// Smalltalk Language
		Smalltalk       Stringorslice `yaml:"smalltalk,omitempty"`        // matrix expand key
		SmalltalkConfig Stringorslice `yaml:"smalltalk_config,omitempty"` // matrix expand key
		SmalltalkVM     Stringorslice `yaml:"smalltalk_vm,omitempty"`     // matrix expand key
	}
)

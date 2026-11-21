// Copyright 2024 Harness, Inc.
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

package xml

type (
	// Project defines a jenkins project.
	Project struct {
		// TODO: add support for more fields.
		ConcurrentBuild bool `xml:"concurrentBuild,omitempty"`
		Disabled        bool `xml:"disabled,omitempty"`

		Builders *Builders `xml:"builders"`

		// Maven module-set jobs (<maven2-moduleset>) carry their build
		// definition in these top-level elements rather than in a
		// freestyle <builders> block.
		Goals     string `xml:"goals,omitempty"`
		MavenName string `xml:"mavenName,omitempty"`
		JDK       string `xml:"jdk,omitempty"`

		// Parameters declares the job's build parameters.
		Parameters []StringParameter `xml:"properties>hudson.model.ParametersDefinitionProperty>parameterDefinitions>hudson.model.StringParameterDefinition"`

		// SCM is the source control configuration. Only the Git SCM
		// (hudson.plugins.git.GitSCM) is modelled today.
		SCM *SCM `xml:"scm"`
	}

	// SCM models a Jenkins source control configuration. The XML class
	// attribute (for example hudson.plugins.git.GitSCM) is captured so the
	// converter can distinguish git from other SCM types.
	SCM struct {
		Class       string   `xml:"class,attr"`
		RemoteURLs  []string `xml:"userRemoteConfigs>hudson.plugins.git.UserRemoteConfig>url"`
		BranchNames []string `xml:"branches>hudson.plugins.git.BranchSpec>name"`
	}

	// StringParameter is a Jenkins string build parameter
	// (hudson.model.StringParameterDefinition).
	StringParameter struct {
		Name         string `xml:"name"`
		Description  string `xml:"description,omitempty"`
		DefaultValue string `xml:"defaultValue,omitempty"`
	}
)

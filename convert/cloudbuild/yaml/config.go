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

// package yaml provides definitions for the Cloud Build schema.
package yaml

import "time"

type (
	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#resource:-build
	Config struct {
		Steps            []*Step           `yaml:"steps,omitempty"`
		Timeout          time.Duration     `yaml:"timeout,omitempty"`
		Queuettl         time.Duration     `yaml:"queueTtl,omitempty"`
		Logsbucket       string            `yaml:"logsBucket,omitempty"`
		Options          *Options          `yaml:"options,omitempty"`
		Substitutions    map[string]string `yaml:"substitutions,omitempty"`
		Tags             []string          `yaml:"tags,omitempty"`
		Serviceaccount   string            `yaml:"serviceAccount,omitempty"`
		Secrets          []*Secret         `yaml:"secrets,omitempty"`
		Availablesecrets *AvailableSecrets `yaml:"availableSecrets,omitempty"`
		Artifacts        *Artifacts        `yaml:"artifacts,omitempty"`
		Images           []string          `yaml:"images,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#artifacts
	Artifacts struct {
		Objects        *ArtifactObjects `yaml:"objects,omitempty"`
		Mavenartifacts []*MavenArifact  `yaml:"mavenArtifacts,omitempty"`
		Pythonpackages []*PythonPackage `yaml:"pythonPackages,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#artifactobjects
	ArtifactObjects struct {
		Location string   `yaml:"location,omitempty"`
		Paths    []string `yaml:"paths,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#Build.Secrets
	AvailableSecrets struct {
		SecretManager []*SecretManagerSecret `yaml:"secretManager,omitempty"`
		Inline        []*Secret              `yaml:"inline,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#mavenartifact
	MavenArifact struct {
		Repository string `yaml:"repository,omitempty"`
		Path       string `yaml:"path,omitempty"`
		Artifactid string `yaml:"artifactId,omitempty"`
		Groupid    string `yaml:"groupId,omitempty"`
		Version    string `yaml:"version,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#buildoptions
	Options struct {
		Sourceprovenancehash      string    `yaml:"sourceProvenanceHash,omitempty"` // ENUM
		Machinetype               string    `yaml:"machineType,omitempty"`          // ENUM
		Disksizegb                string    `yaml:"diskSizeGb,omitempty"`
		Substitutionoption        string    `yaml:"substitutionOption,omitempty"` // ENUM
		Dynamicsubstitutions      string    `yaml:"dynamicSubstitutions,omitempty"`
		Logstreamingoption        string    `yaml:"logStreamingOption,omitempty"`        // ENUM
		Logging                   string    `yaml:"logging,omitempty"`                   // ENUM
		Defaultlogsbucketbehavior string    `yaml:"defaultLogsBucketBehavior,omitempty"` // ENUM
		Env                       []string  `yaml:"env,omitempty"`
		Secretenv                 []string  `yaml:"secretEnv,omitempty"`
		Volumes                   []*Volume `yaml:"volumes,omitempty"`
		Pool                      *Pool     `yaml:"pool,omitempty"`
		Requestedverifyoption     string    `yaml:"requestedVerifyOption,omitempty"` // ENUM
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#Build.PoolOption
	Pool struct {
		Name string `yaml:"name,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#pythonpackage
	PythonPackage struct {
		Repository string   `yaml:"repository,omitempty"`
		Paths      []string `yaml:"paths,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#secret
	Secret struct {
		KMSKeyName string            `yaml:"kmsKeyName,omitempty"`
		SecretEnv  map[string]string `yaml:"secretEnv,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#secretmanagersecret
	SecretManagerSecret struct {
		VersionName string            `yaml:"versionName,omitempty"`
		Env         map[string]string `yaml:"env,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#buildstep
	Step struct {
		Name           string        `yaml:"name,omitempty"`
		Args           []string      `yaml:"args,omitempty"`
		Env            []string      `yaml:"env,omitempty"`
		Allowfailure   bool          `yaml:"allowFailure,omitempty"`
		Allowexitcodes []string      `yaml:"allowExitCodes,omitempty"`
		Dir            string        `yaml:"dir,omitempty"`
		ID             string        `yaml:"id,omitempty"`
		Waitfor        []string      `yaml:"waitFor,omitempty"`
		Entrypoint     string        `yaml:"entrypoint,omitempty"`
		Secretenv      string        `yaml:"secretEnv,omitempty"`
		Volumes        []*Volume     `yaml:"volumes,omitempty"`
		Timeout        time.Duration `yaml:"timeout,omitempty"`
		Script         string        `yaml:"script,omitempty"`
	}

	// https://cloud.google.com/build/docs/api/reference/rest/v1/projects.builds#volume
	Volume struct {
		Name string `yaml:"name,omitempty"`
		Path string `yaml:"path,omitempty"`
	}
)

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

// Package cloudbuild converts Google Cloud Build pipelines to Harness pipelines.
package cloudbuild

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
	"time"

	cloudbuild "github.com/drone/go-convert/convert/cloudbuild/yaml"
	"github.com/drone/go-convert/internal/store"
	harness "github.com/drone/spec/dist/go"

	"github.com/ghodss/yaml"
)

// Converter converts a Cloud Build pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers

	// // as we walk the yaml, we store a
	// // a snapshot of the current node and
	// // its parents.
	// config *cloudbuild.Pipeline
	// stage  *cloudbuild.Stage
}

// New creates a new Converter that converts a Cloud Build
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
	src, err := cloudbuild.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(src)
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

// converts converts a Cloud Build pipeline to a Harness pipeline.
func (d *Converter) convert(src *cloudbuild.Config) ([]byte, error) {

	// create the harness pipeline
	pipeline := &harness.Pipeline{
		Version: 1,
	}

	spec := &harness.StageCI{
		Cache:    nil, // No Google equivalent
		Platform: nil, // TODO
		Envs:     envMappingJexl,
		Runtime:  d.convertRuntime(src),
		Steps:    d.convertSteps(src),
	}

	// add global environment variables
	if opts := src.Options; opts != nil {
		spec.Envs = convertEnv(opts.Env)

		// opts.Sourceprovenancehash
		// opts.Machinetype
		// opts.Disksizegb
		// opts.Substitutionoption
		// opts.Dynamicsubstitutions
		// opts.Logstreamingoption
		// opts.Logging
		// opts.Defaultlogsbucketbehavior
		// opts.Secretenv
		// opts.Volumes
		// opts.Pool
		// opts.Requestedverifyoption
	}

	// src.Timeout
	// src.Queuettl
	// src.Logsbucket
	// src.Substitutions
	// src.Tags
	// src.Serviceaccount
	// src.Secrets
	// src.Availablesecrets
	// src.Artifacts
	// src.Images

	// conver pipeilne stages
	pipeline.Stages = append(pipeline.Stages, &harness.Stage{
		Name:     "pipeline",
		Desc:     "converted from google cloud build",
		Type:     "ci",
		Delegate: nil, // No Google equivalent
		On:       nil, // No Google equivalent
		When:     nil, // No Google equivalent
		Spec:     spec,
	})

	// marshal the harness yaml
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	// find and replace google cloudbuild variables with
	// the harness equivalents.
	for before, after := range envMapping {
		out = bytes.ReplaceAll(out, []byte(before), []byte(after))
	}

	return out, nil
}

func (d *Converter) convertRuntime(src *cloudbuild.Config) *harness.Runtime {
	if d.kubeEnabled {
		return &harness.Runtime{
			Type: "kubernetes",
			Spec: &harness.RuntimeKube{
				Namespace: d.kubeNamespace,
				Connector: d.kubeConnector,
			},
		}
	}
	return &harness.Runtime{
		Type: "cloud",
		Spec: harness.RuntimeMachine{},
	}
}

func (d *Converter) convertSteps(src *cloudbuild.Config) []*harness.Step {
	var steps []*harness.Step
	for _, step := range src.Steps {
		// skip git clone steps by default
		if strings.HasPrefix(step.Name, "gcr.io/cloud-builders/git") {
			continue
		}
		steps = append(steps, d.convertStep(src, step))
	}
	return steps
}

func (d *Converter) convertStep(src *cloudbuild.Config, srcstep *cloudbuild.Step) *harness.Step {
	return &harness.Step{
		Name: d.identifiers.Generate(
			// extract the last segment of the container name
			// and use as the base name
			path.Base(srcstep.Name),
		),
		Desc:    "",  // No Google equivalent
		When:    nil, // No Google equivalent
		On:      nil, // No Google equivalent
		Type:    "script",
		Timeout: convertTimeout(srcstep.Timeout),
		Spec: &harness.StepExec{
			Image:      srcstep.Name,
			Connector:  d.dockerhubConn,
			Privileged: isPrivileged(srcstep.Name),
			Mount:      nil, // TODO
			Pull:       "",  // No Google equivalent
			Shell:      "",  // No Google equivalent
			User:       "",  // No Google equivalent
			Group:      "",  // No Google equivalent
			Network:    "",  // No Google equivalent
			Entrypoint: srcstep.Entrypoint,
			Args:       srcstep.Args,
			Run:        srcstep.Script,
			Envs:       convertEnv(srcstep.Env),
			Resources:  nil, // No Google equivalent
			Reports:    nil, // No Google equivalent

			// TODO support step.allowFailure
			// TODO support step.allowExitCodes
			// TODO support step.dir
			// TODO support step.waitFor
			// TODO support step.secretEnv
		},
	}
}

// helper function returns true if a container
// should be started in privileged mode.
func isPrivileged(name string) bool {
	// TODO should we mount /var/run/docker.sock for this image?
	// Or does it execute docker-in-docker.
	return strings.HasPrefix(name, "gcr.io/cloud-builders/docker")
}

// helper function returns a timeout string. If there
// is no timeout, a zero value is returned.
func convertTimeout(src time.Duration) string {
	if dst := src.String(); dst == "0s" {
		return ""
	} else {
		return dst
	}
}

// helper function that converts a string slice of
// environment variables in key=value format to a map.
func convertEnv(src []string) map[string]string {
	dst := map[string]string{}
	for _, env := range src {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := parts[0]
		v := parts[1]
		dst[k] = v
	}
	if len(dst) == 0 {
		return nil
	} else {
		return dst
	}
}

// helper function combines one or more maps of environment
// variables into a single map.
func combineEnv(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}

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

	// create the harness pipeline spec
	pipeline := &harness.Pipeline{
		Options: new(harness.Default),
	}

	// create the harness pipeline
	config := &harness.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    pipeline,
	}

	// convert subsitutions to inputs
	if v := src.Substitutions; len(v) != 0 {
		pipeline.Inputs = map[string]*harness.Input{}
		for key, val := range src.Substitutions {
			pipeline.Inputs[key] = &harness.Input{
				Type:    "string",
				Default: val,
			}
		}
	}

	// convert pipeline timeout
	if v := src.Timeout; v != 0 {
		pipeline.Options.Timeout = convertTimeout(v)
	}

	spec := &harness.StageCI{
		Cache: nil, // No Google equivalent
		Envs:  nil,
		Platform: &harness.Platform{
			Os:   harness.OSLinux.String(),
			Arch: harness.ArchAmd64.String(),
		},
		Runtime: d.convertRuntime(src),
		Steps:   d.convertSteps(src),
	}

	// add global environment variables
	uniqueVols := map[string]struct{}{}
	if opts := src.Options; opts != nil {
		spec.Envs = convertEnv(opts.Env)

		// add global volumes
		if vols := opts.Volumes; len(vols) > 0 {
			for _, vol := range vols {
				uniqueVols[vol.Name] = struct{}{}
				spec.Volumes = append(spec.Volumes, &harness.Volume{
					Name: vol.Name,
					Type: "temp",
					Spec: &harness.VolumeTemp{},
				})
			}
		}
	}
	// add step volumes
	for _, step := range src.Steps {
		for _, vol := range step.Volumes {
			// do not add the volume if already exists
			if _, ok := uniqueVols[vol.Name]; ok {
				continue
			}
			uniqueVols[vol.Name] = struct{}{}
			spec.Volumes = append(spec.Volumes, &harness.Volume{
				Name: vol.Name,
				Type: "temp",
				Spec: &harness.VolumeTemp{},
			})
		}
	}

	if d.kubeEnabled {
		spec.Volumes = append(spec.Volumes, &harness.Volume{
			Name: "dockersock",
			Type: "temp",
			Spec: &harness.VolumeTemp{},
		})
	} else {
		spec.Volumes = append(spec.Volumes, &harness.Volume{
			Name: "dockersock",
			Type: "host",
			Spec: &harness.VolumeHost{
				Path: "/var/run/docker.sock",
			},
		})
	}

	// TODO src.Secrets
	// TODO src.Availablesecrets
	// TODO opts.Secretenv

	// append steps to publish artifacts
	if v := src.Artifacts; v != nil {
		// TODO
		// https://cloud.google.com/build/docs/build-config-file-schema#artifacts
		// https://cloud.google.com/build/docs/build-config-file-schema#mavenartifacts
		// https://cloud.google.com/build/docs/build-config-file-schema#pythonpackages
	}

	// append steps to push docker images
	if v := src.Images; len(v) != 0 {
		// TODO
		// https://cloud.google.com/build/docs/build-config-file-schema#images
	}

	// conver pipeilne stages
	pipeline.Stages = append(pipeline.Stages, &harness.Stage{
		Name:     "pipeline",
		Desc:     "converted from google cloud build",
		Type:     "ci",
		Delegate: nil, // No Google equivalent
		Failure:  nil, // No Google equivalent
		When:     nil, // No Google equivalent
		Spec:     spec,
	})

	// replace google cloud build substitution variable
	// with harness jexl expressions
	config, err := replaceAll(
		config,
		combineEnv(
			envMappingJexl,
			mapInputsToExpr(src.Substitutions),
		),
	)
	if err != nil {
		return nil, err
	}

	// map cloud build environment variables to harness
	// environment variables using jexl.
	config.Spec.(*harness.Pipeline).Options.Envs = envMappingJexl

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
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
	spec := new(harness.RuntimeCloud)
	if src.Options != nil {
		spec.Size = convertMachine(src.Options.Machinetype)
	}
	return &harness.Runtime{
		Type: "cloud",
		Spec: spec,
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
			srcstep.ID,
			// fallback to the last sebment of the container
			// name and use as the base name.
			path.Base(srcstep.Name),
		),
		Desc:    "",  // No Google equivalent
		When:    nil, // No Google equivalent
		Failure: createFailurestrategy(srcstep),
		Type:    "script",
		Timeout: convertTimeout(srcstep.Timeout),
		Spec: &harness.StepExec{
			Image:      srcstep.Name,
			Connector:  d.dockerhubConn,
			Privileged: false, // No Google Equivalent
			Pull:       "",    // No Google equivalent
			Shell:      "",    // No Google equivalent
			User:       "",    // No Google equivalent
			Group:      "",    // No Google equivalent
			Network:    "",    // No Google equivalent
			Entrypoint: srcstep.Entrypoint,
			Args:       srcstep.Args,
			Run:        srcstep.Script,
			Envs:       convertEnv(srcstep.Env),
			Resources:  nil, // No Google equivalent
			Reports:    nil, // No Google equivalent
			Mount:      createMounts(src, srcstep),

			// TODO support step.dir
			// TODO support step.secretEnv
		},
	}
}

func createFailurestrategy(src *cloudbuild.Step) *harness.FailureList {
	if src.Allowfailure == false && len(src.Allowexitcodes) == 0 {
		return nil
	}
	return &harness.FailureList{
		Items: []*harness.Failure{
			{
				Errors: []string{"all"},
				Action: &harness.FailureAction{
					Type: "ignore",
					Spec: &harness.Ignore{},
					// TODO exit_codes needs to be re-added to spec
					// ExitCodes: src.Allowexitcodes,
				},
			},
		},
	}
}

func createMounts(src *cloudbuild.Config, srcstep *cloudbuild.Step) []*harness.Mount {
	var mounts = []*harness.Mount{
		{
			Name: "dockersock",
			Path: "/var/run/docker.sock",
		},
	}
	for _, vol := range srcstep.Volumes {
		mounts = append(mounts, &harness.Mount{
			Name: vol.Name,
			Path: vol.Path,
		})
	}
	if src.Options != nil {
		for _, vol := range src.Options.Volumes {
			mounts = append(mounts, &harness.Mount{
				Name: vol.Name,
				Path: vol.Path,
			})
		}
	}
	return mounts
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

// helper function returns a machine size that corresponds
// to the google cloud machine type.
func convertMachine(src string) string {
	switch src {
	case "N1_HIGHCPU_8", "E2_HIGHCPU_8":
		return "standard"
	case "N1_HIGHCPU_32", "E2_HIGHCPU_32":
		return "" // TODO convert 32 core machines
	default:
		return ""
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

// helper function maps input variables to expressions.
func mapInputsToExpr(envs map[string]string) map[string]string {
	out := map[string]string{}
	for k := range envs {
		out[k] = "<+inputs." + k + ">"
	}
	return out
}

func replaceAll(in *harness.Config, envs map[string]string) (*harness.Config, error) {
	// marshal the harness yaml
	b, err := yaml.Marshal(in)
	if err != nil {
		return in, err
	}

	// find and replace google cloudbuild variables with
	// the harness equivalents.
	for before, after := range envs {
		b = bytes.ReplaceAll(b, []byte("${"+before+"}"), []byte(after))
	}

	// unarmarshal the yaml
	out, err := harness.ParseBytes(b)
	if err != nil {
		return in, err
	}
	return out, nil
}

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

package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepBackground converts a v0 Background step to v1 background spec
func ConvertStepBackground(src *v0.Step, ctx *StepConvertContext) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepBackground)
	if !ok {
		return nil
	}

	// Container mapping. See ConvertStepRun for Cloud-containerless rationale.
	var container *v1.Container
	resources := ConvertContainerResources(sp.Resources)
	if ctx.IsCloud() && sp.Image == "" {
		WarnDroppedContainerFieldsOnCloud(src.ID, src.Type, map[string]bool{
			"connectorRef": sp.ConnRef != "",
			"registryRef":  sp.RegistryRef != "",
			"privileged":   sp.Privileged != nil,
			"resources":    resources != nil,
			"entrypoint":   sp.Entrypoint != nil,
			"runAsUser":    sp.RunAsUser != nil,
			"portBindings": sp.PortBindings != nil,
		})
	} else if sp.Image != "" || sp.ConnRef != "" || sp.Privileged != nil || resources != nil || sp.RunAsUser != nil {
		pull := ConvertImagePullPolicy(sp.ImagePullPolicy)

		container = &v1.Container{
			Image:        sp.Image,
			Registry:     sp.RegistryRef,
			Connector:    sp.ConnRef,
			Privileged:   sp.Privileged,
			Pull:         pull,
			Resources:    resources,
			Entrypoint:   sp.Entrypoint,
			User:         sp.RunAsUser,
			PortBindings: sp.PortBindings,
		}
	}

	// Reports mapping (JUnit)
	var report *v1.Reports
	if sp.Reports != nil {
		report = &v1.Reports{}
		report.Type = strings.ToLower(sp.Reports.Type)
		if sp.Reports.Spec != nil {
			report.Paths = sp.Reports.Spec.Paths
		}
	}

	// Shell mapping - lower-case common values
	shell := strings.ToLower(sp.Shell)

	dst := &v1.StepRun{
		Container: container,
		Env:       sp.Env,
		Report:    report,
		Shell:     shell,
	}

	// Add command if present
	if sp.Command != "" {
		dst.Script = v1.Stringorslice{sp.Command}
	}

	return dst
}

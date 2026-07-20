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

// ConvertStepRunSpec converts a v0 Run step to v1 run spec only
func ConvertStepRun(src *v0.Step, ctx *StepConvertContext) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepRun)
	if !ok {
		return nil
	}

	script := sp.Command

	// Container mapping. On Cloud infra a step with no image is "containerless"
	// — it runs directly on the hosted VM — so we must NOT emit a container
	// block even when other container-adjacent fields are set. Any such fields
	// get dropped with a single warn per step.
	var container *v1.Container
	resources := ConvertContainerResources(sp.Resources)
	if ctx.IsCloud() && sp.Image == "" {
		WarnDroppedContainerFieldsOnCloud(src.ID, src.Type, map[string]bool{
			"connectorRef": sp.ConnRef != "",
			"registryRef":  sp.RegistryRef != "",
			"privileged":   sp.Privileged != nil,
			"resources":    resources != nil,
			"runAsUser":    sp.RunAsUser != nil,
		})
	} else if sp.Image != "" || sp.ConnRef != "" || sp.Privileged != nil || resources != nil || sp.RunAsUser != nil {
		pull := ConvertImagePullPolicy(sp.ImagePullPolicy)
		container = &v1.Container{
			Image:      sp.Image,
			Registry:   sp.RegistryRef,
			Connector:  sp.ConnRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			Resources:  resources,
			User:       sp.RunAsUser,
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
	if script != "" {
		// use single string so it marshals as block scalar in YAML
		dst.Script = v1.Stringorslice{script}
	}

	dst.Outputs = ConvertOutputVariables(sp.Outputs)

	if sp.Alias != nil {
		dst.Alias = &v1.OutputAlias{
			Key:   sp.Alias.Key,
			Scope: sp.Alias.Scope,
		}
	}

	return dst
}

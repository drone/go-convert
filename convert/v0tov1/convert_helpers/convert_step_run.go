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
func ConvertStepRun(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepRun)
	if !ok {
		return nil
	}

	script := sp.Command

	// Container mapping
	var container *v1.Container
	if sp.Image != "" || sp.ConnRef != "" || sp.ImagePullPolicy != "" {
		pull := ""
		if strings.EqualFold(sp.ImagePullPolicy, "Always") {
			pull = "always"
		} else if strings.EqualFold(sp.ImagePullPolicy, "Never") {
			pull = "never"
		} else if strings.EqualFold(sp.ImagePullPolicy, "IfNotPresent") {
			pull = "if-not-exists"
		}
		cpu := ""
		memory := ""
		if sp.Resources != nil && sp.Resources.Limits != nil {
			cpu = sp.Resources.Limits.GetCPUString()
			memory = sp.Resources.Limits.GetMemoryString()
		}
		container = &v1.Container{
			Image:      sp.Image,
			Connector:  sp.ConnRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			Cpu:        cpu,
			Memory:     memory,
			User:       sp.RunAsUser,
		}
	}

	// Reports mapping (JUnit)
	var report *v1.ReportList
	if sp.Reports != nil && strings.EqualFold(sp.Reports.Type, "JUnit") && sp.Reports.Spec != nil {
		for _, p := range sp.Reports.Spec.Paths {
			if strings.TrimSpace(p) == "" {
				continue
			}
			r := &v1.Report{Type: "junit", Path: p}
			if report == nil {
				report = &v1.ReportList{}
			}
			*report = append(*report, r)
		}
	}

	// Shell mapping - lower-case common values
	shell := strings.ToLower(sp.Shell)
	if shell == "" {
		shell = "sh"
	}

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

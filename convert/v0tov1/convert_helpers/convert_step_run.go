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
	if sp.Image != "" || sp.ConnRef != "" || sp.Privileged || sp.ImagePullPolicy != "" {
		pull := ""
		if strings.EqualFold(sp.ImagePullPolicy, "Always") {
			pull = "always"
		}
		cpu := ""
		memory := ""
		if sp.Resources != nil && sp.Resources.Limits.CPU != nil {
			cpu = sp.Resources.Limits.CPU.String()
		}
		if sp.Resources != nil && sp.Resources.Limits.Memory != nil {
			memory = sp.Resources.Limits.Memory.String()
		}
		container = &v1.Container{
			Image:      sp.Image,
			Connector:  sp.ConnRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			Cpu:        cpu,
			Memory:     memory,
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
		Env:       map[string]interface{}{},
		Report:    report,
		Shell:     shell,
	}
	if script != "" {
		// use single string so it marshals as block scalar in YAML
		dst.Script = v1.Stringorslice{script}
	}

	outputs := make([]*v1.Output, 0)
	for _, outputVar := range sp.Outputs {
		if outputVar == nil {
			continue
		}
		outputs = append(outputs, &v1.Output{
			Name:  outputVar.Name,
			Type:  outputVar.Type,
			Value: outputVar.Value,
		})
	}
	dst.Outputs = outputs

	// merge envVariables and step-level env into run env
	for k, v := range sp.Env {
		if dst.Env == nil {
			dst.Env = make(map[string]interface{})
		}
		dst.Env[k] = v
	}

	return dst
}



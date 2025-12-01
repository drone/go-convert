package converthelpers

import (

	"strings"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertStepTestIntelligence(src *v0.Step) *v1.StepTest {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepTestIntelligence)
	if !ok {
		return nil
	}

	// Container
	var container *v1.Container
	if sp.Image != "" || sp.ConnRef != "" || sp.Privileged || sp.ImagePullPolicy != "" {
		pull := ""
		if strings.EqualFold(sp.ImagePullPolicy, "Always") {
			pull = "always"
		} else if strings.EqualFold(sp.ImagePullPolicy, "Never") {
			pull = "never"
		} else if strings.EqualFold(sp.ImagePullPolicy, "IfNotPresent") {
			pull = "if-not-present"
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

	// Shell
	shell := strings.ToLower(sp.Shell)
	if shell == "" {
		shell = "sh"
	}

	// outputs
	var outputs []*v1.Output
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

	//glob to match
	var match v1.Stringorslice
	for _, glob := range sp.Globs {
		match = append(match, glob)
	}

	//intelligence
	var intelligence *v1.TestIntelligence
	if !sp.IntelligenceMode {
		intelligence = &v1.TestIntelligence{
			Disabled: true,
		}
	}

	dst := &v1.StepTest{
		Env: sp.Env,
		Script: v1.Stringorslice{sp.Command},
		Shell: shell,
		Intelligence: intelligence,
		Container: container,
		Report: report,
		Outputs: outputs,
		Match: match,
	}
	return dst
}
package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepPlugin converts a v0 Plugin step to v1 run format
func ConvertStepPlugin(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepPlugin)
	if !ok {
		return nil
	}

	// Container mapping
	var container *v1.Container
			if sp.Image != "" || sp.ConnRef != "" || sp.Privileged != nil || (sp.Resources != nil && sp.Resources.Limits != nil )|| sp.RunAsUser != nil {
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
			Entrypoint: sp.Entrypoint,
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

	dst := &v1.StepRun{
		Container: container,
		Env:       sp.Env,
		Report:    report,
		With:      sp.Settings,
	}

	return dst
}

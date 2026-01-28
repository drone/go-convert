package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
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
	if sp.Image != "" || sp.ConnRef != "" || sp.ImagePullPolicy != "" {
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
	outputs := ConvertOutputVariables(sp.Outputs)

	//glob to match
	var match v1.Stringorslice
	for _, glob := range sp.Globs {
		match = append(match, glob)
	}

	//intelligence
	var disabled *flexible.Field[bool]
	if sp.IntelligenceMode != nil {
		if val, ok := sp.IntelligenceMode.AsStruct(); ok {
			// It's a boolean value
			disabled = &flexible.Field[bool]{Value: val}
		} else if expr, ok := sp.IntelligenceMode.AsString(); ok {
			// Convert <+expression> to <+!expression>
			modifiedExpr := "<+!" + expr + ">"
			disabled = &flexible.Field[bool]{}
			disabled.SetExpression(modifiedExpr)
		}
	}
	intelligence := &v1.TestIntelligence{
		Disabled: disabled,
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
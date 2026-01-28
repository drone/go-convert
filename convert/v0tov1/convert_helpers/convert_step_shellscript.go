package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepShellScript converts a v0 ShellScript step to a v1 Run step
func ConvertStepShellScript(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepShellScript)
	if !ok {
		return nil
	}

	script := sp.Source.Spec.Script

	shell := strings.ToLower(strings.TrimSpace(sp.Shell))
	if shell == "" {
		shell = "sh"
	}

	dst := &v1.StepRun{
		Env:       map[string]interface{}{},
		Shell:     shell,
	}
	if script != "" {
		dst.Script = v1.Stringorslice{script}
	}

	for _, ev := range sp.EnvironmentVariables {
		if ev == nil || strings.TrimSpace(ev.Name) == "" {
			continue
		}
		if dst.Env == nil {
			dst.Env = make(map[string]interface{})
		}
		dst.Env[ev.Name] = ev.Value
	}

	dst.Outputs = ConvertOutputVariables(sp.OutputVariables)

	return dst
}

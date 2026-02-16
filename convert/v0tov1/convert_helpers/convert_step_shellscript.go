package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
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

	env_map := map[string]interface{}{}
	var env *flexible.Field[map[string]interface{}]
	for _, ev := range sp.EnvironmentVariables {
		if ev == nil || strings.TrimSpace(ev.Name) == "" {
			continue
		}
		env_map[ev.Name] = ev.Value
	}
    if len(env_map) > 0 {
        env = &flexible.Field[map[string]interface{}]{Value: env_map}
    }

	dst := &v1.StepRun{
		Env:       env,
		Shell:     shell,
	}
	if script != "" {
		dst.Script = v1.Stringorslice{script}
	}

	if sp.Alias != nil {
		dst.Alias = &v1.OutputAlias{
			Key:   sp.Alias.Key,
			Scope: sp.Alias.Scope,
		}
	}
	dst.Outputs = ConvertOutputVariables(sp.OutputVariables)

	return dst
}

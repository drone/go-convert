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

	// Build base script from inline source if present
	scriptLines := []string{}
	if sp.Source != nil && strings.EqualFold(sp.Source.Type, "Inline") {
		if s := strings.TrimSpace(sp.Source.Spec.Script); s != "" {
			scriptLines = append(scriptLines, s)
		}
	}

	// Append output variable exports (echo "name=${VALUE}" >> $HARNESS_OUTPUT)
	if len(sp.OutputVariables) > 0 {
		if len(scriptLines) > 0 {
			// add a spacer comment if there is existing script
			scriptLines = append(scriptLines, "# write output variables to harness")
		}
		for _, ov := range sp.OutputVariables {
			if ov == nil || strings.TrimSpace(ov.Name) == "" || strings.TrimSpace(ov.Value) == "" {
				continue
			}
			line := "echo \"" + ov.Name + "=${" + ov.Value + "}\" >> $HARNESS_OUTPUT"
			scriptLines = append(scriptLines, line)
		}
	}

	var script string
	if len(scriptLines) > 0 {
		script = strings.Join(scriptLines, "\n")
	}

	// Container mapping from optional image fields (not typical for ShellScript, but supported)
	var container *v1.Container
	if sp.Image != "" || sp.ConnRef != "" || sp.ImagePullPolicy != "" {
		pull := ""
		if strings.EqualFold(sp.ImagePullPolicy, "Always") {
			pull = "always"
		}
		container = &v1.Container{
			Image:     sp.Image,
			Connector: sp.ConnRef,
			Pull:      pull,
		}
	}

	// Shell mapping
	shell := strings.ToLower(strings.TrimSpace(sp.Shell))
	if shell == "" {
		shell = "sh"
	}

	dst := &v1.StepRun{
		Container: container,
		Env:       map[string]string{},
		Shell:     shell,
	}
	if script != "" {
		// use single string to render as YAML block scalar
		dst.Script = v1.Stringorslice{script}
	}

	// Merge environmentVariables (array) into env map
	for _, ev := range sp.EnvironmentVariables {
		if ev == nil || strings.TrimSpace(ev.Name) == "" {
			continue
		}
		if dst.Env == nil {
			dst.Env = make(map[string]string)
		}
		dst.Env[ev.Name] = ev.Value
	}

	// Merge step-level env as well
	for k, v := range src.Env {
		if dst.Env == nil {
			dst.Env = make(map[string]string)
		}
		dst.Env[k] = v
	}

	return dst
}

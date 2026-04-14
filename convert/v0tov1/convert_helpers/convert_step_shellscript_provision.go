package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertStepShellScriptProvision converts a v0 ShellScriptProvision step to a v1 Run step
func ConvertStepShellScriptProvision(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepShellScriptProvision)
	if !ok {
		return nil
	}

	var script string
	if sp.Source != nil {
		script = sp.Source.Spec.Script
	}

	// Prepend provisioner setup script
	provisionerSetup := `DIR="./shellScriptProvisioner/<+pipeline.executionId>-<+stage.executionId>-<+step.id>"
PROVISIONER_OUTPUT_PATH="$DIR/output.json"

mkdir -p "$DIR"
touch "$PROVISIONER_OUTPUT_PATH"

`
	if script != "" {
		script = provisionerSetup + script
	} else {
		script = provisionerSetup
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
		Env: env,
	}
	if script != "" {
		dst.Script = v1.Stringorslice{script}
	}

	return dst
}

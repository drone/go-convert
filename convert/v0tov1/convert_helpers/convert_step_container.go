package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepContainer converts a v0 Container step to v1 run step
// Container steps are converted to Run steps
func ConvertStepContainer(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepContainer)
	if !ok {
		return nil
	}

	// TODO: Infrastructure conversion is not yet supported for Container step
	if sp.Infrastructure != nil {
		messagelog.GetMessageLogger().LogWarning(
			"UNSUPPORTED_CONTAINER_INFRASTRUCTURE",
			"infrastructure conversion is not yet supported for Container step; infrastructure configuration will be skipped",
			messagelog.WithStep(src.ID, src.Type),
		)
	}

	script := sp.Command

	// Container mapping
	var container *v1.Container
	if sp.Image != "" || sp.ConnRef != "" || sp.Privileged != nil || sp.RunAsUser != nil {
		pull := ""
		if strings.EqualFold(sp.ImagePullPolicy, "Always") {
			pull = "always"
		} else if strings.EqualFold(sp.ImagePullPolicy, "Never") {
			pull = "never"
		} else if strings.EqualFold(sp.ImagePullPolicy, "IfNotPresent") {
			pull = "if-not-exists"
		}

		container = &v1.Container{
			Image:      sp.Image,
			Connector:  sp.ConnRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			User:       sp.RunAsUser,
			Entrypoint: sp.Entrypoint,
		}
	}

	// Shell mapping - lower-case common values
	shell := strings.ToLower(sp.Shell)

	dst := &v1.StepRun{
		Container: container,
		Env:       sp.Env,
		Shell:     shell,
	}

	if script != "" {
		dst.Script = v1.Stringorslice{script}
	}

	// Convert output variables
	dst.Outputs = ConvertOutputVariables(sp.Outputs)

	return dst
}

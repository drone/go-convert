package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepPlugin converts a v0 Plugin step to v1 run format
func ConvertStepPlugin(src *v0.Step, ctx *StepConvertContext) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepPlugin)
	if !ok {
		return nil
	}

	// Container mapping. See ConvertStepRun for Cloud-containerless rationale.
	var container *v1.Container
	resources := ConvertContainerResources(sp.Resources)
	if ctx.IsCloud() && sp.Image == "" {
		WarnDroppedContainerFieldsOnCloud(src.ID, src.Type, map[string]bool{
			"connectorRef": sp.ConnRef != "",
			"registryRef":  sp.RegistryRef != "",
			"privileged":   sp.Privileged != nil,
			"resources":    resources != nil,
			"entrypoint":   sp.Entrypoint != nil,
			"runAsUser":    sp.RunAsUser != nil,
		})
	} else if sp.Image != "" || sp.ConnRef != "" || sp.Privileged != nil || resources != nil || sp.RunAsUser != nil {
		pull := ConvertImagePullPolicy(sp.ImagePullPolicy)
		container = &v1.Container{
			Image:      sp.Image,
			Registry:   sp.RegistryRef,
			Connector:  sp.ConnRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			Resources:  resources,
			Entrypoint: sp.Entrypoint,
			User:       sp.RunAsUser,
		}
	}

	// Reports mapping (JUnit)
	var report *v1.Reports
	if sp.Reports != nil {
		report = &v1.Reports{}
		report.Type = strings.ToLower(sp.Reports.Type)
		if sp.Reports.Spec != nil {
			report.Paths = sp.Reports.Spec.Paths
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

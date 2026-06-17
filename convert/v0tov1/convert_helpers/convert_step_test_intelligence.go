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
	resources := ConvertContainerResources(sp.Resources)
	if sp.Image != "" || sp.ConnRef != "" || sp.Privileged != nil || resources != nil || sp.RunAsUser != nil {
		pull := ConvertImagePullPolicy(sp.ImagePullPolicy)
		container = &v1.Container{
			Image:      sp.Image,
			Registry:   sp.RegistryRef,
			Connector:  sp.ConnRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			Resources:  resources,
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

	// Shell
	shell := strings.ToLower(sp.Shell)

	// outputs
	outputs := ConvertOutputVariables(sp.Outputs)

	//intelligence
	var intelligence *v1.TestIntelligence
	if sp.IntelligenceMode != nil {
		intelligence = &v1.TestIntelligence{
			Disabled: flexible.NegateBool(sp.IntelligenceMode),
		}
	}

	dst := &v1.StepTest{
		Env:          sp.Env,
		Script:       v1.Stringorslice{sp.Command},
		Shell:        shell,
		Intelligence: intelligence,
		Container:    container,
		Report:       report,
		Outputs:      outputs,
		Match:        sp.Globs,
	}
	return dst
}

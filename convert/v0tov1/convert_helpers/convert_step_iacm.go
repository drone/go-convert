package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertStepIACMTerraformPlugin(src *v0.Step) *v1.StepTemplate {
    if src == nil || src.Spec == nil {
        return nil
    }

	spec, ok := src.Spec.(*v0.StepIACMTerraformPlugin)
	if !ok {
		return nil
	}

	with := map[string]interface{}{}

	with["command"] = spec.Command
	if spec.Target != nil {
		with["target"] = spec.Target
	}
	if spec.Replace != nil {
		with["replace"] = spec.Replace
	}
	if spec.ImportVars != nil {
		with["import"] = spec.ImportVars
	}
    
    return &v1.StepTemplate{
		Uses: v1.StepTypeIACMTerraformPlugin,
        With: with,
    }
}

func ConvertStepIACMOpenTofuPlugin(src *v0.Step) *v1.StepTemplate {
    if src == nil || src.Spec == nil {
        return nil
    }

	spec, ok := src.Spec.(*v0.StepIACMOpenTofuPlugin)
	if !ok {
		return nil
	}

	with := map[string]interface{}{}

	with["command"] = spec.Command
	if spec.Target != nil {
		with["target"] = spec.Target
	}
	if spec.Replace != nil {
		with["replace"] = spec.Replace
	}
	if spec.ImportVars != nil {
		with["import"] = spec.ImportVars
	}
    
    return &v1.StepTemplate{
		Uses: v1.StepTypeIACMOpenTofuPlugin,
        With: with,
    }
}
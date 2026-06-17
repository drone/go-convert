package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertStepFilesUpload(src *v0.Step) *v1.StepUpload {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepFilesUpload)
	if !ok {
		return nil
	}

	inputs := ConvertVariables(sp.InputVariables)
	if inputs == nil {
		return &v1.StepUpload{}
	}

	return &v1.StepUpload{
		Inputs: inputs,
	}
}

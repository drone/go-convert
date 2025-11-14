package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepWait converts a v0 Wait step to a v1 action step
func ConvertStepWait(src *v0.Step) *v1.StepWait {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepWait)
	if !ok || sp == nil {
		return nil
	}

	return &v1.StepWait{
		Duration: sp.Duration,
	}
}
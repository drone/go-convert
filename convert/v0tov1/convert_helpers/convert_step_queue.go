package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertStepQueue(src *v0.Step) *v1.StepQueue{
	if src == nil {
		return nil
	}
	spec, ok := src.Spec.(*v0.StepQueue)
	if !ok {
		return nil
	}
	return &v1.StepQueue{
		Key:   spec.Key,
		Scope: spec.Scope,
	}
	
} 
	
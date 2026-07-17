package pipelineconverter

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func (c *PipelineConverter) ConvertTemplate(src *v0.Template) *v1.Template {
	if src == nil {
		return nil
	}

	dst := &v1.Template{
		ID: src.ID,
		Name: src.Name,
		Description: src.Description,
		Tags: src.Tags,
		Version: src.VersionLabel,
	}
	// Based on the type of template, convert it
	switch src.Type {
	case "Stage":
		if spec, ok := src.Spec.(*v0.Stage); ok {
			dst.Stage = c.convertStage(spec, "")
		}
	case "Step":
		if spec, ok := src.Spec.(*v0.Step); ok {
			dst.Step = c.ConvertSingleStep(spec, false, "", "", "")
		}
	case "StepGroup":
		if spec, ok := src.Spec.(*v0.StepGroup); ok {
			dst.Env = spec.Env
			dst.Group = &v1.StepGroup{
				Steps:  c.ConvertSteps(spec.Steps, false, "", "", ""),
				Inputs: c.convertVariables(spec.Variables),
			}
			dst.OnFailure = convert_helpers.ConvertFailureStrategies(spec.FailureStrategies)
			dst.Strategy = convert_helpers.ConvertStrategy(spec.Strategy)
			dst.Timeout = spec.Timeout
			dst.Delegate = convert_helpers.ConvertDelegate(spec.DelegateSelectors, nil)
			dst.If = convert_helpers.ConvertStepWhen(spec.When, spec.Skip)
		}
	case "Pipeline":
		if spec, ok := src.Spec.(*v0.Pipeline); ok {
			dst.Pipeline = c.ConvertPipeline(spec)
		}
	}

	return dst
}

package pipelineconverter

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// convertStepGroupTemplate converts a v0 step group with template reference to v1 step template format.
// The v1 format uses: uses: templateRef@versionLabel with overlay.step.group containing the converted step group.
func (c *PipelineConverter) convertStepGroupTemplate(src *v0.StepGroup, isRollback bool, groupPath string) *v1.StepTemplate {
	if src == nil || src.Template == nil {
		return nil
	}

	template := src.Template

	// Build uses field: templateRef@versionLabel
	uses := template.TemplateRef
	if template.VersionLabel != "" {
		uses = uses + "@" + template.VersionLabel
	}

	result := &v1.StepTemplate{
		Uses: uses,
	}

	// If templateInputs exists, convert it to overlay.step.group
	if template.TemplateInputs != nil {
		inputs := template.TemplateInputs

		// Build the overlay step with group
		overlayStep := &v1.Step{}

		// Convert the step group contents
		overlayStep.Group = &v1.StepGroup{}

		if inputs.Steps != nil {
			overlayStep.Group.Steps = c.ConvertSteps(inputs.Steps, isRollback, groupPath)
		}

		if inputs.Variables != nil {
			overlayStep.Group.Inputs = c.convertVariables(inputs.Variables)
		}

		// Set the with.overlay.step structure
		result.With = &v1.StepTemplateWith{
			Overlay: &v1.StepTemplateOverlay{
				Step: overlayStep,
			},
		}
	}

	return result
}

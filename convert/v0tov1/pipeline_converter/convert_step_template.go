package pipelineconverter

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// convertStepTemplate converts a v0 step with template reference to v1 step template format.
// The v1 format uses: uses: templateRef@versionLabel with overlay.step containing the converted step.
func (c *PipelineConverter) convertStepTemplate(src *v0.Step, isRollback bool, basePath string) *v1.StepTemplate {
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

	// If templateInputs exists, convert it to overlay.step
	if template.TemplateInputs != nil {
		// Convert the templateInputs (which is a *Step) directly
		overlayStep := c.ConvertSingleStep(template.TemplateInputs, isRollback, basePath)
		if overlayStep != nil {
			result.With = &v1.StepTemplateWith{
				Overlay: &v1.StepTemplateOverlay{
					Step: overlayStep,
				},
			}
		}
	}

	return result
}

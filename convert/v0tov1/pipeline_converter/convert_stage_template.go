package pipelineconverter

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// convertStageTemplate converts a v0 stage with template reference to v1 stage template format.
// The v1 format uses: uses: templateRef@versionLabel with overlay.stage containing the converted stage.
func (c *PipelineConverter) convertStageTemplate(src *v0.Stage, basePath string) *v1.StageTemplate {
	if src == nil || src.Template == nil {
		return nil
	}

	template := src.Template

	// Build uses field: templateRef@versionLabel
	uses := template.TemplateRef
	if template.VersionLabel != "" {
		uses = uses + "@" + template.VersionLabel
	}

	result := &v1.StageTemplate{
		Uses: uses,
	}

	// If templateInputs exists, convert it to overlay.stage
	if template.TemplateInputs != nil {
		// Convert the templateInputs (which is a *Stage) directly
		overlayStage := c.convertStage(template.TemplateInputs, basePath)
		if overlayStage != nil {
			result.With = &v1.StageTemplateWith{
				Overlay: &v1.StageTemplateOverlay{
					Stage: overlayStage,
				},
			}
		}
	}

	return result
}

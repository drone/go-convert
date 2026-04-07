package pipelineconverter

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// convertPipelineTemplate converts a v0 pipeline with template reference to v1 pipeline template format.
// The v1 format uses: uses: templateRef@versionLabel with overlay.pipeline containing the converted pipeline.
func (c *PipelineConverter) convertPipelineTemplate(src *v0.Pipeline) *v1.PipelineTemplate {
	if src == nil || src.Template == nil {
		return nil
	}

	template := src.Template

	// Build uses field: templateRef@versionLabel
	uses := template.TemplateRef
	if template.VersionLabel != "" {
		uses = uses + "@" + template.VersionLabel
	}

	result := &v1.PipelineTemplate{
		Uses: uses,
	}

	// If templateInputs exists, convert it to overlay.pipeline
	if template.TemplateInputs != nil {
		// Convert the templateInputs (which is a *Pipeline) directly
		overlayPipeline := c.ConvertPipeline(template.TemplateInputs)
		if overlayPipeline != nil {
			result.With = &v1.PipelineTemplateWith{
				Overlay: &v1.PipelineTemplateOverlay{
					Pipeline: overlayPipeline,
				},
			}
		}
	}

	return result
}

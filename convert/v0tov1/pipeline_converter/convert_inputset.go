package pipelineconverter

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertInputSet converts a v0 InputSet to v1 InputSet format.
// The v0 inputSet.pipeline is converted to v1 inputs.overlay.
func (c *PipelineConverter) ConvertInputSet(src *v0.InputSet) *v1.InputSet {
	if src == nil {
		return nil
	}

	dst := &v1.InputSet{}

	// Convert the pipeline to overlay
	if src.Pipeline != nil {
		dst.Overlay = c.ConvertPipeline(src.Pipeline)
	}

	return dst
}

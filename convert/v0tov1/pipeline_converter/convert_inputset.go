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
		dst.Variables = liftPipelineVariables(dst.Overlay)
	}

	return dst
}

// liftPipelineVariables pulls pipeline-level inputs out of the overlay
// pipeline into a scalar map (name -> value) and clears them from the
// overlay. This emits pipeline-level variables as scalar key/value pairs at
// the top-level "inputs" map (sibling to "overlay") instead of nested
// Input objects under overlay.inputs. Value takes precedence over Default;
// an empty string is used as a fallback to avoid emitting null.
func liftPipelineVariables(overlay *v1.Pipeline) map[string]interface{} {
	if overlay == nil || len(overlay.Inputs) == 0 {
		return nil
	}
	vars := make(map[string]interface{}, len(overlay.Inputs))
	for name, in := range overlay.Inputs {
		switch {
		case in == nil:
			vars[name] = ""
		case in.Value != nil:
			vars[name] = in.Value
		case in.Default != nil:
			vars[name] = in.Default
		default:
			vars[name] = ""
		}
	}
	overlay.Inputs = nil
	return vars
}

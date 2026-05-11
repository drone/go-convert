package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepGCSUpload converts a v0 GCSUpload step to v1 template format
func ConvertStepGCSUpload(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepGCSUpload)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	if spec.Bucket != "" {
		with["bucket"] = spec.Bucket
	}

	if spec.SourcePath != "" {
		with["source"] = spec.SourcePath
	}

	if spec.Target != "" {
		with["target"] = spec.Target
	}

	return &v1.StepTemplate{
		Uses: "uploadArtifactsToGCS",
		With: with,
	}
}

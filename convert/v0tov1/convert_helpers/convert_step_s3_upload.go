package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepS3Upload converts a v0 S3Upload step to v1 template format
func ConvertStepS3Upload(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepS3Upload)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	if spec.Region != "" {
		with["region"] = spec.Region
	}

	if spec.Bucket != "" {
		with["bucket"] = spec.Bucket
	}

	if spec.SourcePath != "" {
		with["source"] = spec.SourcePath
	}

	if spec.Endpoint != "" {
		with["endpoint"] = spec.Endpoint
	}

	if spec.Target != "" {
		with["target"] = spec.Target
	}

	return &v1.StepTemplate{
		Uses: "uploadArtifactsToS3",
		With: with,
	}
}

package converthelpers

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepGCSUpload converts a v0 GCSUpload step to v1 template format
func ConvertStepGCSUpload(src *v0.Step) *v1.StepTemplate {
	if src == nil {
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

	if spec.SourcePath != "" {
		with["source"] = spec.SourcePath
	}

	// Combine bucket and target into single target field: <bucket>/<target>
	if spec.Bucket != "" && spec.Target != "" {
		target := fmt.Sprintf("%s/%s", spec.Bucket, spec.Target)
		with["target"] = target
	} else if spec.Bucket != "" {
		with["target"] = spec.Bucket
	}

	return &v1.StepTemplate{
		Uses: "uploadArtifactsToGCS",
		With: with,
	}
}

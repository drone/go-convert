package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepSaveCacheS3 converts a v0 SaveCacheS3 step to v1 template format
func ConvertStepSaveCacheS3(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepSaveCacheS3)
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

	if spec.Key != "" {
		with["cachekey"] = spec.Key
	}

	if len(spec.SourcePaths) > 0 {
		with["mount"] = spec.SourcePaths
	}

	if spec.Endpoint != "" {
		with["endpoint"] = spec.Endpoint
	}

	if spec.ArchiveFormat != "" {
		// Convert to lowercase as per v1 spec
		with["archiveformat"] = strings.ToLower(spec.ArchiveFormat)
	}

	with["override"] = spec.Override

	with["pathstyle"] = spec.PathStyle

	return &v1.StepTemplate{
		Uses: "saveCacheToS3",
		With: with,
	}
}

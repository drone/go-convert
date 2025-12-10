package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepSaveCacheGCS converts a v0 SaveCacheGCS step to v1 template format
func ConvertStepSaveCacheGCS(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepSaveCacheGCS)
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

	if spec.Key != "" {
		with["cachekey"] = spec.Key
	}

	if len(spec.SourcePaths) > 0 {
		with["mount"] = spec.SourcePaths
	}

	if spec.ArchiveFormat != "" {
		// Convert to lowercase as per v1 spec
		with["archiveformat"] = strings.ToLower(spec.ArchiveFormat)
	}

	if spec.Override {
		with["override"] = spec.Override
	}

	return &v1.StepTemplate{
		Uses: "saveCacheToGCS",
		With: with,
	}
}

package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

type S3UploadStepWith struct {
	CacheKey                string `json:"cachekey,omitempty"`
	Mount                   string `json:"mount,omitempty"`
	Rebuild                 bool   `json:"rebuild,omitempty"`
	Endpoint                string `json:"endpoint,omitempty"`
	Region                  string `json:"region,omitempty"`
	ArchiveFormat           string `json:"archiveformat,omitempty"`
	PathStyle               string `json:"pathstyle,omitempty"`
	Override                string `json:"override,omitempty"`
	ExitCode                bool   `json:"exitcode,omitempty"`
	Backend                 string `json:"backend,omitempty"`
	BackendOperationTimeout string `json:"backend_operation_timeout,omitempty"`
	Bucket                  string `json:"bucket,omitempty"`
	Connector               string `json:"connector,omitempty"`
	Env         map[string]string `json:"env_vars,omitempty"`
}

func ConvertStepS3Upload(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepS3Upload)
	if !ok {
		return nil
	}

	
	// Build With struct from spec fields
	with := S3UploadStepWith{
		Endpoint:  spec.Endpoint,
		Region:    spec.Region,
		Bucket:    spec.Bucket,
		Connector: spec.ConnectorRef,
	}

	// Build environment variables map from spec and step level
	env := make(map[string]string)

	// Merge spec-level environment variables
	if spec.Env != nil {
		for key, value := range spec.Env {
			if strings.TrimSpace(key) == "" {
				continue
			}
			env[key] = value
		}
	}

	// Merge step-level environment variables
	if src.Env != nil {
		for key, value := range src.Env {
			if strings.TrimSpace(key) == "" {
				continue
			}
			env[key] = value
		}
	}

	dst := &v1.StepTemplate{
		Uses: "saveCacheToS3@1.0.0",
		With: with,
		Env:  env,
	}

	return dst
}

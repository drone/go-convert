package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepGitClone converts a v0 GitClone step to v1 template format
func ConvertStepGitClone(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepGitClone)
	if !ok {
		return nil
	}

	with := make(map[string]interface{})

	// Connector
	if sp.ConnRef != "" {
		with["connector"] = sp.ConnRef
	}

	// Build type - extract branch or tag from Build field
	if sp.BuildType != nil && !sp.BuildType.IsNil() {
		if build, ok := sp.BuildType.AsStruct(); ok {
			// Build is a struct with Type and Spec
			if build.Type == "branch" && build.Spec.Branch != "" {
				with["branch"] = build.Spec.Branch
			} else if build.Type == "tag" && build.Spec.Tag != "" {
				with["tag"] = build.Spec.Tag
			}
		} else if expr, ok := sp.BuildType.AsString(); ok {
			// Build is an expression string
			with["build"] = expr
		}
	}

	// Clone directory
	if sp.CloneDirectory != "" {
		with["cloneDirectory"] = sp.CloneDirectory
	}

	// Repository name
	if sp.Repository != "" {
		with["repoName"] = sp.Repository
	}

	// Depth
	if sp.Depth != nil {
		with["depth"] = sp.Depth
	}

	// SSL Verify
	if sp.SSLVerify != "" {
		with["sslVerify"] = sp.SSLVerify
	}

	dst := &v1.StepTemplate{
		Uses: "gitCloneStep",
		With: with,
	}

	return dst
}

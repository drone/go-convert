package converthelpers

import (
	"fmt"
	"strings"

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

	// Repository name
	if sp.Repository != "" {
		with["repo_name"] = sp.Repository
	}

	// Build type - extract branch, tag, PR, or commit_sha from Build field
	if sp.BuildType != nil && !sp.BuildType.IsNil() {
		if build, ok := sp.BuildType.AsStruct(); ok {
			switch build.Type {
			case "branch":
				with["build_target"] = "Git Branch"
			case "tag":
				with["build_target"] = "Tag"
			case "PR":
				with["build_target"] = "Pull Request"
			case "commitSha":
				with["build_target"] = "Commit"
			default:
				with["build_target"] = build.Type
			}
			if build.Spec.Branch != "" {
				with["branch"] = build.Spec.Branch
			} else if build.Spec.Tag != "" {
				with["tag"] = build.Spec.Tag
			} else if build.Spec.Number != nil {
				if pr, ok := build.Spec.Number.AsString(); ok {
					with["pr"] = pr
				} else if pr, ok := build.Spec.Number.AsStruct(); ok {
					with["pr"] = fmt.Sprintf("%d", pr)
				}
			} else if build.Spec.CommitSha != "" {
				with["commit_sha"] = build.Spec.CommitSha
			}
		}
	}

	// Clone directory
	if sp.CloneDirectory != "" {
		with["clone_directory"] = sp.CloneDirectory
	}

	// Depth
	if sp.Depth != nil && !sp.Depth.IsNil() {
		if expr, ok := sp.Depth.AsString(); ok {
			with["depth"] = expr
		} else if depth, ok := sp.Depth.AsStruct(); ok {
			with["depth"] = fmt.Sprintf("%d", depth)
		}
	}

	// Sparse checkout - convert []string to comma-separated string for v1
	if sp.SparseCheckout != nil && !sp.SparseCheckout.IsNil() {
		if expr, ok := sp.SparseCheckout.AsString(); ok {
			with["sparse_checkout"] = expr
		} else if paths, ok := sp.SparseCheckout.AsStruct(); ok && len(paths) > 0 {
			with["sparse_checkout"] = strings.Join(paths, ",")
		}
	}

	// Pre-fetch command
	if sp.PreFetchCommand != "" {
		with["pre_fetch"] = sp.PreFetchCommand
	}

	// Submodule strategy
	if sp.SubmoduleStrategy != nil && !sp.SubmoduleStrategy.IsNil() {
		if expr, ok := sp.SubmoduleStrategy.AsString(); ok {
			with["submodule_strategy"] = expr
		} else if strategy, ok := sp.SubmoduleStrategy.AsStruct(); ok {
			with["submodule_strategy"] = strategy
		}
	}

	// Output file paths content - convert []string to comma-separated string for v1
	if sp.OutputFilePathsContent != nil && !sp.OutputFilePathsContent.IsNil() {
		if expr, ok := sp.OutputFilePathsContent.AsString(); ok {
			with["file_paths_content"] = expr
		} else if paths, ok := sp.OutputFilePathsContent.AsStruct(); ok && len(paths) > 0 {
			with["file_paths_content"] = strings.Join(paths, ",")
		}
	}

	// LFS enabled
	if sp.Lfs != nil && !sp.Lfs.IsNil() {
		if lfs, ok := sp.Lfs.AsString(); ok {
			with["lfs_enabled"] = lfs
		} else if lfs, ok := sp.Lfs.AsStruct(); ok {
			with["lfs_enabled"] = lfs
		}
	}

	// Debug
	if sp.Debug != nil && !sp.Debug.IsNil() {
		if expr, ok := sp.Debug.AsString(); ok {
			with["debug"] = expr
		} else if debug, ok := sp.Debug.AsStruct(); ok {
			with["debug"] = debug
		}
	}

	// Fetch tags
	if sp.FetchTags != nil && !sp.FetchTags.IsNil() {
		if expr, ok := sp.FetchTags.AsString(); ok {
			with["fetch_tags"] = expr
		} else if fetchTags, ok := sp.FetchTags.AsStruct(); ok {
			with["fetch_tags"] = fetchTags
		}
	}

	// pr_merge_strategy: no v0 StepGitClone field; template default is "Source Branch"
	// copy_file_content: no v0 StepGitClone field; optional template input with no default

	dst := &v1.StepTemplate{
		Uses: "gitCloneStep",
		With: with,
	}

	return dst
}

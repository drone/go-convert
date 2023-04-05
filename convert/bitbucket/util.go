// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bitbucket

import (
	"fmt"
	"time"

	bitbucket "github.com/drone/go-convert/convert/bitbucket/yaml"
	harness "github.com/drone/spec/dist/go"

	"github.com/gotidy/ptr"
)

//
// this file contains helper functions that convert
// bitbucket data structures to harness data structures
// structures. these functions are all stateless
// (they do not rely on snapshots of walk state).
//

// helper function converts a bitbucket size enum to a
// harness resource class.
func convertSize(size bitbucket.Size) string {
	switch size {
	case bitbucket.Size2x: // 8GB
		return "large"
	case bitbucket.Size4x: // 16GB
		return "xlarge"
	case bitbucket.Size8x: // 32GB
		return "xxlarge"
	case bitbucket.Size1x: // 4GB
		return "standard"
	default:
		return ""
	}
}

// helper function converts the bitbucket stage clone
// configuration to a harness stage clone configuration.
func convertClone(stage *bitbucket.Stage) *harness.CloneStage {
	var clones []*bitbucket.Clone

	// loop through the steps and if a step
	// defines clone behavior
	for _, step := range extractSteps(stage) {
		if step.Clone != nil {
			clones = append(clones, step.Clone)
		}
	}

	// if there are not clone configurations at
	// the step-level we can return a nil clone.
	if len(clones) == 0 {
		return nil
	}

	clone := new(harness.CloneStage)
	for _, v := range clones {
		if v.Depth != nil {
			if v.Depth.Value > int(clone.Depth) {
				clone.Depth = int64(v.Depth.Value)
			}
		}
		if v.SkipVerify {
			clone.Insecure = true
		}
		if v.Enabled != nil && !ptr.ToBool(v.Enabled) {
			// TODO
		}
	}

	return clone
}

// helper function converts the bitbucket global clone
// configuration to a global harness clone configuration.
func convertCloneGlobal(clone *bitbucket.Clone) *harness.Clone {
	if clone == nil {
		return nil
	}

	to := new(harness.Clone)
	to.Insecure = clone.SkipVerify

	if clone.Depth != nil {
		to.Depth = int64(clone.Depth.Value)
	}

	// disable cloning globally if the user has
	// explicityly disabled this functionality
	if clone.Enabled != nil && ptr.ToBool(clone.Enabled) == false {
		to.Disabled = true
	}

	return to
}

// helper function converts the bitbucket global cache to
// a harness stage cache, filtered by the list of cache names.
func convertCache(defs *bitbucket.Definitions, caches []string) *harness.Cache {
	if defs == nil || len(defs.Caches) == 0 || len(caches) == 0 {
		return nil
	}

	cache := new(harness.Cache)
	cache.Enabled = true

	var files []string
	var paths []string

	for _, name := range caches {
		src, ok := defs.Caches[name]
		if !ok {
			continue
		}
		paths = append(paths, src.Path)
		if src.Key != nil {
			files = append(files, src.Key.Files...)
		}
	}

	for _, name := range caches {
		switch name {
		case "composer":
			paths = append(paths, "composer")
			paths = append(paths, "~/.composer/cache")
		case "dotnetcore":
			paths = append(paths, "dotnetcore")
			paths = append(paths, "~/.nuget/packages")
		case "gradle":
			paths = append(paths, "gradle")
			paths = append(paths, "~/.gradle/caches")
		case "ivy2":
			paths = append(paths, "ivy2")
			paths = append(paths, "~/.ivy2/cache")
		case "maven":
			paths = append(paths, "maven")
			paths = append(paths, "~/.m2/repository")
		case "node":
			paths = append(paths, "node")
			paths = append(paths, "node_modules")
		case "pip":
			paths = append(paths, "pip")
			paths = append(paths, "~/.cache/pip")
		case "sbt":
			paths = append(paths, "sbt")
			paths = append(paths, "ivy2")
			paths = append(paths, "~/.ivy2/cache")
		}
	}

	cache.Paths = paths
	return cache
}

// helper function converts the bitbucket global defaults
// to the harness global default configuration.
func convertDefault(config *bitbucket.Config) *harness.Default {

	// if the global pipeline configuration sections
	// are empty or nil, return nil
	if config.Clone == nil &&
		config.Image == nil &&
		config.Options == nil {
		return nil
	}

	if config.Image == nil {
		// Username
		// Password
	}
	if config.Options == nil {
		// Docker (bool)
		// MaxTime (int)
		// Size (1x, 2x, 4x, 8x)
		// Credentials ???
	}

	var def *harness.Default

	// if the user has configured global clone defaults,
	// convert this to pipeline-level clone settings.
	if config.Clone != nil {
		// create the default if not already created.
		if def == nil {
			def = new(harness.Default)
		}
		def.Clone = convertCloneGlobal(config.Clone)

		// if the clone is disabled we need to make
		// sure it isn't explicitly enabled for any steps.
		if def.Clone.Disabled {
			for _, step := range extractAllSteps(config.Pipelines.Default) {
				if step.Clone != nil && ptr.ToBool(step.Clone.Enabled) {
					def.Clone.Disabled = false
					break
				}
			}
		}
	}

	return def
}

// helper function converts an integer of minutes to a time
// duration string.
func minuteToDurationString(v int64) string {
	dur := time.Duration(v) * time.Minute
	return fmt.Sprint(dur)
}

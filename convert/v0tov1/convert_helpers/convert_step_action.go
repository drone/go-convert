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

package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepActionSpec converts a v0 step action to v1 action spec only
func ConvertStepAction(src *v0.Step) *v1.StepAction {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Type assert the spec to StepAction
	spec, ok := src.Spec.(*v0.StepAction)
	if !ok {
		return nil
	}

	dst := &v1.StepAction{
		Uses: spec.Uses,
		With: spec.With,
		Env:  spec.Envs,
	}

	// Merge step-level environment variables with action-level environment variables
	if src.Env != nil {
		if dst.Env == nil {
			dst.Env = make(map[string]string)
		}
		for k, v := range src.Env {
			dst.Env[k] = v
		}
	}

	return dst
}


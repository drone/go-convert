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
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertStepActionSpec converts a v0 step action to v1 action spec only
func ConvertStepAction(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Type assert the spec to StepAction
	spec, ok := src.Spec.(*v0.StepAction)
	if !ok {
		return nil
	}

	script := fmt.Sprintf("plugin -kind action -name %v", spec.Uses)
	env_map := map[string]interface{}{}
	var env *flexible.Field[map[string]interface{}]
	// Emit each `with` entry as its own PLUGIN_WITH_<key> env var, matching the
	// V0 and native-V1 serializers (VmActionStepSerializer) and the per-key
	// fallback in drone/plugin's getWith. A single JSON-encoded PLUGIN_WITH
	// blob breaks at runtime: Harness expressions such as
	// <+secrets.getValue(...)> are resolved only after conversion, and
	// multi-line secret values inject raw newlines into the pre-encoded JSON,
	// so the plugin's json.Unmarshal fails with
	// "invalid character '\n' in string literal" (CI-24451).
	for k, v := range spec.With {
		env_map["PLUGIN_WITH_"+k] = v
	}
	for k, v := range spec.Envs {
		env_map[k] = v
	}
	if len(env_map) > 0 {
		env = &flexible.Field[map[string]interface{}]{Value: env_map}
	}
	dst := &v1.StepRun{
		Script: v1.Stringorslice{script},
		Env:    env,
	}

	// dst := &v1.StepAction{
	// 	Uses: spec.Uses,
	// 	With: spec.With,
	// 	Env:  spec.Envs,
	// }

	// // Merge step-level environment variables with action-level environment variables
	// if src.Env != nil {
	// 	if dst.Env == nil {
	// 		dst.Env = make(map[string]string)
	// 	}
	// 	for k, v := range src.Env {
	// 		dst.Env[k] = v
	// 	}
	// }

	return dst
}

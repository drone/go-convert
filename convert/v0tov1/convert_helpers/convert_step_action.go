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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

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
	// Encode the `with` map as a single JSON-valued PLUGIN_WITH env var.
	// The per-key PLUGIN_WITH_<key> form is brittle: values like "=1.20.1"
	// can lose information in YAML/string round-trips, and key names like
	// "go-version" aren't POSIX env names. The plugin runner accepts both
	// forms (see drone/plugin plugin/github/env.go getWith), and the JSON
	// form is unambiguous.
	if len(spec.With) > 0 {
		// Disable HTML escaping so Harness expressions like
		// <+pipeline.variables.checkLatest> aren't mangled into
		// \u003c+...\u003e by the default json.Marshal behavior.
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(spec.With); err != nil {
			return nil
		}
		// Encode appends a trailing newline; trim it.
		env_map["PLUGIN_WITH"] = strings.TrimRight(buf.String(), "\n")
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
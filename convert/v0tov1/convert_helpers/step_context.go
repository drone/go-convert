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
	"sort"
	"strings"

	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// StepConvertContext threads stage-scoped facts into per-step converters.
// Only the fields the step helpers actually need live here; extend cautiously.
type StepConvertContext struct {
	Runtime *v1.Runtime
}

// IsCloud reports whether the enclosing stage's resolved runtime is Cloud
// (hosted VM). Nil-safe: nil ctx or nil runtime → false, so any caller
// without stage context transparently keeps today's behavior.
func (c *StepConvertContext) IsCloud() bool {
	return c != nil && IsCloudRuntime(c.Runtime)
}

// IsCloudRuntime is a nil-safe predicate for v1.Runtime.
func IsCloudRuntime(rt *v1.Runtime) bool {
	return rt != nil && rt.Cloud != nil
}

// WarnDroppedContainerFieldsOnCloud emits a single WARN when a Cloud
// containerless step (image absent) drops container-only fields that would
// otherwise have coerced it into container mode. dropped maps field name →
// whether that field was set on the source v0 spec.
func WarnDroppedContainerFieldsOnCloud(stepID, stepType string, dropped map[string]bool) {
	names := make([]string, 0, len(dropped))
	for k, set := range dropped {
		if set {
			names = append(names, k)
		}
	}
	if len(names) == 0 {
		return
	}
	sort.Strings(names)
	messagelog.GetMessageLogger().LogWarning(
		"CLOUD_CONTAINERLESS_FIELDS_DROPPED",
		fmt.Sprintf("cloud stage: step has no image; dropping container-only fields: %s", strings.Join(names, ",")),
		messagelog.WithStep(stepID, stepType),
	)
}

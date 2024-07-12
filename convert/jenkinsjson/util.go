// Copyright 2023 Harness, Inc.
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

package jenkinsjson

import "regexp"

func SanitizeForId(spanName string, spanId string) string {
	spanName = regexp.MustCompile(`[: ]`).ReplaceAllString(spanName, "_")

	// Replace invalid characters with underscores
	invalidCharRegex := regexp.MustCompile("[^a-zA-Z0-9.-_]+")
	sanitized := invalidCharRegex.ReplaceAllString(spanName, "_")

	// Trim leading and trailing underscores
	sanitized = regexp.MustCompile("^_+|_+$").ReplaceAllString(sanitized, "")

	return sanitized + spanId[:6]
}

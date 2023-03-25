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

package yaml

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// helper function that repairs un unparseable `on` section
// of the yaml. The github parser allows a key witout a
// trailing semi-colon, like this:
//
//     on:
//       push
//
// this function converts to this:
//
//     on:
//       push: {}
//
func repairOn(in io.Reader) bytes.Buffer {
	var out bytes.Buffer

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		// here we look for keywords to replace. this is
		// probably in-efficient but is good enough for now.
		temp := strings.TrimSpace(line)
		for _, keyword := range keywords {
			if temp == keyword {
				line = strings.Replace(line, keyword, keyword+": {}", 1)
				break
			}
		}

		out.WriteString(line)
		out.WriteString("\n")
	}

	return out
}

var keywords = []string{
	"branch_protection_rule",
	"check_run",
	"check_suite",
	"create",
	"delete",
	"deployment",
	"deployment_status",
	"discussion",
	"discussion_comment",
	"fork",
	"gollum",
	"issue_comment",
	"issues",
	"label",
	"member",
	"merge_group",
	"milestone",
	"page_build",
	"project",
	"project_card",
	"project_column",
	"public",
	"pull_request",
	"pull_request_review",
	"pull_request_review_comment",
	"pull_request_target",
	"push",
	"registry_package",
	"repository_dispatch",
	"release",
	"schedule",
	"status",
	"watch",
	"workflow_call",
	"workflow_dispatch",
	"workflow_run,omitempty",
}

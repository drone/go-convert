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

package circle

import (
	"bufio"
	"bytes"
	"strings"
)

// replaceParams finds and replaces circle pipeline
// parameters with harness pipeline parameters.
func replaceParams(in []byte) []byte {
	// optimization to exit early and return the input,
	// unmodified, if there are no parameters present
	// in the yaml.
	if !bytes.Contains(in, []byte("<<")) {
		return in
	}

	var out bytes.Buffer
	scanner := bufio.NewScanner(
		bytes.NewBuffer(in),
	)
	for scanner.Scan() {
		line := scanner.Text()

		// detect paramter bracket start and end position
		a := strings.Index(line, "<<")
		b := strings.LastIndex(line, ">>")

		// skip this line and write to the buffer if no
		// brackets are detected.
		if a == -1 || b == -1 || b < a {
			out.WriteString(line)
			out.WriteString("\n")
			continue
		}

		// extract the string
		s := line[a : b+2]
		s = strings.TrimPrefix(s, "<<")
		s = strings.TrimSuffix(s, ">>")
		s = strings.TrimSpace(s)

		// if the parameter is not found we can skip
		// and leave the parameter as-is.
		p, ok := params[s]
		if !ok {
			// if the parameter is not found check to see
			// if it is a user-defined input parameter.
			if strings.HasPrefix(s, "pipeline.parameters.") ||
				strings.HasPrefix(s, "parameters.") {
				s = strings.ReplaceAll(s, "pipeline.parameters.", "inputs.")
				s = strings.ReplaceAll(s, "parameters.", "inputs.")
				p = s
			} else {
				out.WriteString(line)
				out.WriteString("\n")
				continue
			}
		}

		out.WriteString(line[:a])
		out.WriteString("<+")
		out.WriteString(p)
		out.WriteString(">")
		out.WriteString(line[b+2:])
		out.WriteString("\n")
	}

	return out.Bytes()
}

// helper function extracts a circle parameter.
func extractParam(in string) (out string) {
	// detect paramter bracket start and end position
	a := strings.Index(in, "<<")
	b := strings.LastIndex(in, ">>")

	// skip this line and write to the buffer if no
	// brackets are detected.
	if a == -1 || b == -1 || b < a {
		return
	}

	// extract the string
	out = in[a : b+2]
	out = strings.TrimPrefix(out, "<<")
	out = strings.TrimSuffix(out, ">>")
	out = strings.TrimSpace(out)
	return
}

// map of circle pipeline values to harness pipeine values.
// https://circleci.com/docs/pipeline-variables
var params = map[string]string{
	"pipeline.id":                                            "pipeline.identifier",
	"pipeline.number":                                        "pipeline.sequenceId",
	"pipeline.project.git_url":                               "codebase.repoUrl",
	"pipeline.project.type":                                  "github", // github, bitbucket, etc
	"pipeline.git.tag":                                       "codebase.tag",
	"pipeline.git.branch":                                    "codebase.branch",
	"pipeline.git.revision":                                  "codebase.commitSha",
	"pipeline.git.base_revision":                             "codebase.baseCommitSha",
	"pipeline.in_setup":                                      "",
	"pipeline.trigger_source":                                "",
	"pipeline.schedule.name":                                 "",
	"pipeline.schedule.id":                                   "",
	"pipeline.trigger_parameters.circleci.trigger_id":        "",
	"pipeline.trigger_parameters.circleci.config_source_id":  "",
	"pipeline.trigger_parameters.circleci.trigger_type":      "trigger.type",
	"pipeline.trigger_parameters.circleci.event_time":        "pipeline.startTs",
	"pipeline.trigger_parameters.circleci.event_type":        "codebase.build.type",
	"pipeline.trigger_parameters.circleci.project_id":        "project.identifier",
	"pipeline.trigger_parameters.circleci.actor_id":          "codebase.gitUser",
	"pipeline.trigger_parameters.gitlab.type":                "",
	"pipeline.trigger_parameters.gitlab.project_id":          "",
	"pipeline.trigger_parameters.gitlab.ref":                 "codebase.commitRef",
	"pipeline.trigger_parameters.gitlab.checkout_sha":        "codebase.commitSha",
	"pipeline.trigger_parameters.gitlab.user_id":             "codebase.gitUserId",
	"pipeline.trigger_parameters.gitlab.user_name":           "codebase.gitUser",
	"pipeline.trigger_parameters.gitlab.user_username":       "codebase.gitUser",
	"pipeline.trigger_parameters.gitlab.user_avatar":         "codebase.gitUserAvatar",
	"pipeline.trigger_parameters.gitlab.repo_name":           "",
	"pipeline.trigger_parameters.gitlab.repo_url":            "codebase.repoUrl",
	"pipeline.trigger_parameters.gitlab.web_url":             "codebase.repoUrl",
	"pipeline.trigger_parameters.gitlab.commit_sha":          "codebase.commitSha",
	"pipeline.trigger_parameters.gitlab.commit_title":        "codebase.commitMessage",
	"pipeline.trigger_parameters.gitlab.commit_message":      "codebase.commitMessage",
	"pipeline.trigger_parameters.gitlab.commit_timestamp":    "pipeline.startTs",
	"pipeline.trigger_parameters.gitlab.commit_author_name":  "codebase.gitUser",
	"pipeline.trigger_parameters.gitlab.commit_author_email": "codebase.gitUserEmail",
	"pipeline.trigger_parameters.gitlab.total_commits_count": "",
	"pipeline.trigger_parameters.gitlab.branch":              "codebase.sourceBranch",
	"pipeline.trigger_parameters.gitlab.default_branch":      "codebase.sourceBranch",
	"pipeline.trigger_parameters.gitlab.x_gitlab_event_id":   "codebase.build.type",
}

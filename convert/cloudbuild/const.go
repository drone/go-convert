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

package cloudbuild

var envMapping = map[string]string{
	"PROJECT_ID":     "HARNESS_PROJECT_ID",  // the project ID of the build.
	"PROJECT_NUMBER": "HARNESS_PROJECT_ID",  // (TODO) the project number of the build.
	"LOCATION":       "LOCATION",            // (TODO) the location/region of the build.
	"BUILD_ID":       "DRONE_BUILD_NUMBER",  // the autogenerated ID of the build.
	"REPO_NAME":      "DRONE_REPO_NAME",     // the source repository name specified by RepoSource.
	"BRANCH_NAME":    "DRONE_COMMIT_BRANCH", // the branch name specified by RepoSource.
	"TAG_NAME":       "DRONE_TAG",           // the tag name specified by RepoSource.
	"REVISION_ID":    "DRONE_COMMIT_SHA",    // the commit SHA specified by RepoSource or resolved from the specified branch or tag.
	"COMMIT_SHA":     "DRONE_COMMIT_SHA",    // the commit SHA specified by RepoSource or resolved from the specified branch or tag.
	"SHORT_SHA":      "DRONE_COMMIT_SHA",    // first 7 characters of $REVISION_ID or $COMMIT_SHA.
}

var envMappingJexl = map[string]string{
	// "PROJECT_ID":     "",
	// "PROJECT_NUMBER": "",
	// "LOCATION":       "",
	// "BUILD_ID":       "",
	"REPO_NAME":   "<+trigger.payload.repository.name>",
	"BRANCH_NAME": "<+trigger.branch>",
	// "TAG_NAME":       "",
	"REVISION_ID": "<+trigger.commitSha>",
	"COMMIT_SHA":  "<+trigger.commitSha>",
	"SHORT_SHA":   "<+trigger.commitSha>",
}

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

import "errors"

// Retry defines retry logic.
type Retry struct {
	Max  int           `yaml:"max,omitempty"`
	When Stringorslice `yaml:"when,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Retry) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 int
	var out2 = struct {
		Max  int           `yaml:"max"`
		When Stringorslice `yaml:"when"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Max = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Max = out2.Max
		v.When = out2.When
		return nil
	}

	return errors.New("failed to unmarshal retry")
}

// Enum of retry when types
// https://docs.gitlab.com/ee/ci/yaml/#retrywhen
const (
	RetryAlways                 = "always"
	RetryUnknownFailure         = "unknown_failure"
	RetryScriptFailure          = "script_failure"
	RetryApiFailure             = "api_failure"
	RetrySturckOrTimeoutFailure = "stuck_or_timeout_failure"
	RetryRunnerSystemFailure    = "runner_system_failure"
	RetryRunnerUnsupported      = "runner_unsupported"
	RetryStaleSchedule          = "stale_schedule"
	RetryJobExecutionTimeout    = "job_execution_timeout"
	RetryArchivedFailure        = "archived_failure"
	RetrySchedulerFailure       = "scheduler_failure"
	RetryIntegrityFailure       = "data_integrity_failure"
)

func (v *Retry) MarshalYAML() (interface{}, error) {
	if v.Max > 0 && len(v.When) > 0 {
		return struct {
			Max  int           `yaml:"max,omitempty"`
			When Stringorslice `yaml:"when,omitempty"`
		}{
			Max:  v.Max,
			When: v.When,
		}, nil
	} else if v.Max > 0 {
		return v.Max, nil
	}
	return nil, nil
}

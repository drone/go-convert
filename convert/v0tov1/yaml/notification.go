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

type (
	// Notification defines a v1 notification configuration
	Notification struct {
		ID       string                 `json:"id,omitempty" yaml:"id,omitempty"`
		Name     string                 `json:"name,omitempty" yaml:"name,omitempty"`
		On       []*NotificationOn      `json:"on,omitempty" yaml:"on,omitempty"`
		Uses     string                 `json:"uses,omitempty" yaml:"uses,omitempty"`
		With     map[string]interface{} `json:"with,omitempty" yaml:"with,omitempty"`
		Disabled bool                   `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	}

	// NotificationOn defines when notifications are triggered
	NotificationOn struct {
		Pipeline interface{} `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
		Stage    interface{} `json:"stage,omitempty" yaml:"stage,omitempty"`
		Step     interface{} `json:"step,omitempty" yaml:"step,omitempty"`
	}

	// NotificationStageOn defines stage-specific notification triggers
	NotificationStageOn struct {
		Start   interface{} `json:"start,omitempty" yaml:"start,omitempty"`
		Success interface{} `json:"success,omitempty" yaml:"success,omitempty"`
		Failed  interface{} `json:"failed,omitempty" yaml:"failed,omitempty"`
	}
)

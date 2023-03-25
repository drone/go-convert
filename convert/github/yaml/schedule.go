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

type Schedule struct {
	Items []*ScheduleItem
}

type ScheduleItem struct {
	Cron string `yaml:"cron,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface for WorkflowTriggers.
func (v *Schedule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out1 := new(ScheduleItem)
	out2 := []*ScheduleItem{}

	if err := unmarshal(&out1); err == nil {
		v.Items = append(v.Items, out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Items = append(v.Items, out2...)
		return nil
	}

	return errors.New("failed to unmarshal on.schedule")
}

// MarshalYAML implements the marshal interface.
func (v *Schedule) MarshalYAML() (interface{}, error) {
	return v.Items, nil
}

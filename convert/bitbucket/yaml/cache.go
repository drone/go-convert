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

type (
	// Cache configures a cache path.
	Cache struct {
		Key  *CacheKey
		Path string
	}

	// temporary data structure for unmarshaling and
	// marshaling caches.
	cache struct {
		Key  *CacheKey `json:"key"`
		Path string    `json:"path"`
	}

	CacheKey struct {
		Files []string `json:"files"`
	}
)

// MarshalYAML implements the marshal interface.
func (v *Cache) MarshalYAML() (interface{}, error) {
	if v.Key == nil {
		return v.Path, nil
	}
	return &cache{
		Key:  v.Key,
		Path: v.Path,
	}, nil
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Cache) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *cache
	if err := unmarshal(&out1); err == nil {
		v.Path = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Path = out2.Path
		v.Key = out2.Key
		return nil
	}
	return errors.New("failed to unmarshal cache")
}

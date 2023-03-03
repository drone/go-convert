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

package store

import "strconv"

// Identifiers stores identifiers.
type Identifiers struct {
	store map[string]struct{}
}

// New returns a new identifier store.
func New() *Identifiers {
	return &Identifiers{
		store: map[string]struct{}{},
	}
}

// Register registers a name with the store.
func (s *Identifiers) Register(name string) bool {
	if _, ok := s.store[name]; ok {
		return false
	}
	s.store[name] = struct{}{}
	return true
}

// Generage generates and registeres a unique name with the
// store. If the base name is already registered, a unique
// suffix is appended to the name.
func (s *Identifiers) Generate(name ...string) string {
	var base string
	// choose the first non-empty name.
	for _, s := range name {
		if s != "" {
			base = s
			break
		}
	}

	// register the name as-is
	if s.Register(base) {
		return base
	}
	// append a suffix to the name and register
	// the first unique combination.
	for i := 1; ; i++ {
		next := base + strconv.Itoa(i)
		if s.Register(next) {
			return next
		}
	}
}

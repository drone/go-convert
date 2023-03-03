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

import "testing"

func TestStore(t *testing.T) {
	store := New()
	if ok := store.Register("foo"); !ok {
		t.Errorf("Want name registration")
	}
	if ok := store.Register("foo"); ok {
		t.Errorf("Want name registration failure")
	}
	if name := store.Generate("bar"); name != "bar" {
		t.Errorf("Want name generated name bar, got %s", name)
	}
	if name := store.Generate("bar"); name != "bar1" {
		t.Errorf("Want name generated name bar1, got %s", name)
	}
	if name := store.Generate("bar"); name != "bar2" {
		t.Errorf("Want name generated name bar2, got %s", name)
	}
}

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

package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)


// newEnvGroupCtx returns a minimal StageConversionContext for env-group tests.
func newEnvGroupCtx() *StageConversionContext {
	return NewStageConversionContext()
}

// asGroupMap type-asserts the Group interface{} to map[string]interface{}.
func asGroupMap(t *testing.T, ref *v1.EnvironmentRef) map[string]interface{} {
	t.Helper()
	if ref == nil || ref.Group == nil {
		t.Fatal("expected non-nil EnvironmentRef with a Group field")
	}
	m, ok := ref.Group.(map[string]interface{})
	if !ok {
		t.Fatalf("expected Group to be map[string]interface{}, got %T", ref.Group)
	}
	return m
}

// getGroupID extracts the "id" string from a v1.EnvironmentRef.Group map.
func getGroupID(t *testing.T, ref *v1.EnvironmentRef) string {
	t.Helper()
	m := asGroupMap(t, ref)
	id, ok := m["id"].(string)
	if !ok {
		t.Fatalf("expected Group[\"id\"] to be a string, got %T", m["id"])
	}
	return id
}

// getGroupItems extracts the items slice from a v1.EnvironmentRef.Group map.
func getGroupItems(t *testing.T, ref *v1.EnvironmentRef) []*v1.EnvironmentItem {
	t.Helper()
	m := asGroupMap(t, ref)
	raw, ok := m["items"]
	if !ok {
		t.Fatal("expected Group[\"items\"] to be present")
	}
	items, ok := raw.([]*v1.EnvironmentItem)
	if !ok {
		t.Fatalf("expected Group[\"items\"] to be []*v1.EnvironmentItem, got %T", raw)
	}
	return items
}

// TestConvertEnvironmentGroup_AccountLevel_AddsPrefix verifies that when the
// EnvGroupRef starts with "account.", any environment ref inside the group that
// does NOT already carry the prefix gets it prepended automatically.
func TestConvertEnvironmentGroup_AccountLevel_AddsPrefix(t *testing.T) {
	src := &v0.EnvironmentGroup{
		EnvGroupRef: "account.my_group",
		Environments: &flexible.Field[[]*v0.Environment]{Value: []*v0.Environment{
			{EnvironmentRef: "env_one"},         // missing prefix → must be added
			{EnvironmentRef: "account.env_two"}, // already prefixed → must stay as-is
		}},
	}

	ref := ConvertEnvironmentGroup(src, newEnvGroupCtx())
	if ref == nil {
		t.Fatal("expected non-nil EnvironmentRef")
	}

	items := getGroupItems(t, ref)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	if items[0].Id != "account.env_one" {
		t.Errorf("item[0].Id: got %q, want %q", items[0].Id, "account.env_one")
	}
	if items[1].Id != "account.env_two" {
		t.Errorf("item[1].Id: got %q, want %q", items[1].Id, "account.env_two")
	}
}

// TestConvertEnvironmentGroup_NonAccountLevel_NoPrefix verifies that when the
// EnvGroupRef does NOT start with "account.", environment refs are left unchanged.
func TestConvertEnvironmentGroup_NonAccountLevel_NoPrefix(t *testing.T) {
	src := &v0.EnvironmentGroup{
		EnvGroupRef: "my_org_group",
		Environments: &flexible.Field[[]*v0.Environment]{Value: []*v0.Environment{
			{EnvironmentRef: "env_one"},
			{EnvironmentRef: "env_two"},
		}},
	}

	ref := ConvertEnvironmentGroup(src, newEnvGroupCtx())
	if ref == nil {
		t.Fatal("expected non-nil EnvironmentRef")
	}

	items := getGroupItems(t, ref)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	if items[0].Id != "env_one" {
		t.Errorf("item[0].Id: got %q, want %q", items[0].Id, "env_one")
	}
	if items[1].Id != "env_two" {
		t.Errorf("item[1].Id: got %q, want %q", items[1].Id, "env_two")
	}
}

// TestConvertEnvironmentGroup_AccountLevel_AllAlreadyPrefixed verifies that
// environment refs that already carry "account." are not double-prefixed.
func TestConvertEnvironmentGroup_AccountLevel_AllAlreadyPrefixed(t *testing.T) {
	src := &v0.EnvironmentGroup{
		EnvGroupRef: "account.my_group",
		Environments: &flexible.Field[[]*v0.Environment]{Value: []*v0.Environment{
			{EnvironmentRef: "account.env_a"},
			{EnvironmentRef: "account.env_b"},
		}},
	}

	ref := ConvertEnvironmentGroup(src, newEnvGroupCtx())
	if ref == nil {
		t.Fatal("expected non-nil EnvironmentRef")
	}

	items := getGroupItems(t, ref)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	if items[0].Id != "account.env_a" {
		t.Errorf("item[0].Id: got %q, want %q", items[0].Id, "account.env_a")
	}
	if items[1].Id != "account.env_b" {
		t.Errorf("item[1].Id: got %q, want %q", items[1].Id, "account.env_b")
	}
}

// TestConvertEnvironmentGroup_SimpleGroupRef verifies that a simple group ref
// (no environments array) without an account prefix is left untouched.
func TestConvertEnvironmentGroup_SimpleGroupRef(t *testing.T) {
	src := &v0.EnvironmentGroup{
		EnvGroupRef: "my_simple_group",
	}

	ref := ConvertEnvironmentGroup(src, newEnvGroupCtx())
	if ref == nil {
		t.Fatal("expected non-nil EnvironmentRef")
	}
	if id := getGroupID(t, ref); id != "my_simple_group" {
		t.Errorf("Group id: got %q, want %q", id, "my_simple_group")
	}
	if _, hasItems := asGroupMap(t, ref)["items"]; hasItems {
		t.Error("expected no items in simple group ref")
	}
}

// TestConvertEnvironmentGroup_AccountLevel_SimpleGroupRef verifies that a
// simple account-level group ref (no environments array) is passed through
// unchanged — the "account." prefix logic only applies to the items inside.
func TestConvertEnvironmentGroup_AccountLevel_SimpleGroupRef(t *testing.T) {
	src := &v0.EnvironmentGroup{
		EnvGroupRef: "account.my_account_group",
	}

	ref := ConvertEnvironmentGroup(src, newEnvGroupCtx())
	if ref == nil {
		t.Fatal("expected non-nil EnvironmentRef")
	}
	if id := getGroupID(t, ref); id != "account.my_account_group" {
		t.Errorf("Group id: got %q, want %q", id, "account.my_account_group")
	}
	if _, hasItems := asGroupMap(t, ref)["items"]; hasItems {
		t.Error("expected no items in simple group ref")
	}
}

// TestConvertEnvironmentGroupEnvItems_AccountLevel tests the helper directly.
func TestConvertEnvironmentGroupEnvItems_AccountLevel(t *testing.T) {
	envs := []*v0.Environment{
		{EnvironmentRef: "prod"},
		{EnvironmentRef: "account.staging"},
		{EnvironmentRef: ""},
	}

	items := convertEnvironmentGroupEnvItems(envs, true)

	// empty ref is skipped
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Id != "account.prod" {
		t.Errorf("item[0].Id: got %q, want %q", items[0].Id, "account.prod")
	}
	if items[1].Id != "account.staging" {
		t.Errorf("item[1].Id: got %q, want %q", items[1].Id, "account.staging")
	}
}

// TestConvertEnvironmentGroupEnvItems_NonAccountLevel tests the helper directly
// without account-level flag — refs are left as-is.
func TestConvertEnvironmentGroupEnvItems_NonAccountLevel(t *testing.T) {
	envs := []*v0.Environment{
		{EnvironmentRef: "prod"},
		{EnvironmentRef: "staging"},
	}

	items := convertEnvironmentGroupEnvItems(envs, false)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Id != "prod" {
		t.Errorf("item[0].Id: got %q, want %q", items[0].Id, "prod")
	}
	if items[1].Id != "staging" {
		t.Errorf("item[1].Id: got %q, want %q", items[1].Id, "staging")
	}
}

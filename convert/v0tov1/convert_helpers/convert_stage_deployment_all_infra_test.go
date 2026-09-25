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
	"encoding/json"
	"reflect"
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	"github.com/drone/go-convert/internal/flexible"
)

// v0 expressed "deploy everywhere" with a deployToAll flag. v1 expresses it with the all-infra
// marker, which accepts a literal boolean or an expression, so nothing has to be dropped:
//   - v0 deployToAll: true         -> v1 all-infra: true alongside deploy-to
//   - v0 deployToAll: false        -> v1 deploy-to only
//   - v0 deployToAll unspecified   -> v1 deploy-to only
//   - v0 deployToAll: <+input>     -> v1 all-infra: <+input> alongside deploy-to
//   - v0 deployToAll: <expression> -> v1 all-infra: <expression> alongside deploy-to
//
// The marker is nil rather than false when absent, so no spurious all-infra: false is emitted.
// At group level the same flag maps to all-env, which is emitted alongside items so the
// per-environment config authored in v0 is not lost.

func infraDefs(ids ...string) *flexible.Field[[]*v0.InfrastructureDefinition] {
	defs := make([]*v0.InfrastructureDefinition, 0, len(ids))
	for _, id := range ids {
		defs = append(defs, &v0.InfrastructureDefinition{Identifier: id})
	}
	field := &flexible.Field[[]*v0.InfrastructureDefinition]{}
	field.Set(defs)
	return field
}

func infraDefsExpr(expr string) *flexible.Field[[]*v0.InfrastructureDefinition] {
	field := &flexible.Field[[]*v0.InfrastructureDefinition]{}
	field.SetExpression(expr)
	return field
}

func exprField(expr string) *flexible.Field[bool] {
	field := &flexible.Field[bool]{}
	field.SetExpression(expr)
	return field
}

func TestResolveDeployTo_AllInfraMarker(t *testing.T) {
	tests := []struct {
		name         string
		deployToAll  *flexible.Field[bool]
		infraDefs    *flexible.Field[[]*v0.InfrastructureDefinition]
		wantDeployTo interface{}
		wantAllInfra interface{}
	}{
		{
			name:         "deployToAll true -> marker alongside the infra list",
			deployToAll:  &flexible.Field[bool]{Value: true},
			infraDefs:    infraDefs("infra1", "infra2"),
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: true,
		},
		{
			name:         "deployToAll true with a single infra -> marker alongside the infra",
			deployToAll:  &flexible.Field[bool]{Value: true},
			infraDefs:    infraDefs("infra1"),
			wantDeployTo: "infra1",
			wantAllInfra: true,
		},
		{
			name:         "deployToAll true with no infra list -> marker only",
			deployToAll:  &flexible.Field[bool]{Value: true},
			infraDefs:    nil,
			wantDeployTo: nil,
			wantAllInfra: true,
		},
		{
			name:         "deployToAll false -> deploy-to only, no marker",
			deployToAll:  &flexible.Field[bool]{Value: false},
			infraDefs:    infraDefs("infra1", "infra2"),
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: nil,
		},
		{
			name:         "deployToAll unspecified -> deploy-to only, no marker",
			deployToAll:  nil,
			infraDefs:    infraDefs("infra1"),
			wantDeployTo: "infra1",
			wantAllInfra: nil,
		},
		{
			name:         "deployToAll <+input> -> marker carried across alongside the infra list",
			deployToAll:  exprField("<+input>"),
			infraDefs:    infraDefs("infra1", "infra2"),
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: "<+input>",
		},
		{
			name:         "deployToAll other expression -> marker carried across alongside the infra",
			deployToAll:  exprField("<+pipeline.variables.everywhere>"),
			infraDefs:    infraDefs("infra1"),
			wantDeployTo: "infra1",
			wantAllInfra: "<+pipeline.variables.everywhere>",
		},
		{
			name:         "deployToAll other expression -> marker carried across alongside the infra list",
			deployToAll:  exprField("<+pipeline.variables.everywhere>"),
			infraDefs:    infraDefs("infra1", "infra2"),
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: "<+pipeline.variables.everywhere>",
		},
		{
			name:         "deployToAll true with an expression infra list -> marker alongside the expression",
			deployToAll:  &flexible.Field[bool]{Value: true},
			infraDefs:    infraDefsExpr("<+input>"),
			wantDeployTo: "<+input>",
			wantAllInfra: true,
		},
		{
			name:         "deployToAll <+input> with no infra list -> marker only",
			deployToAll:  exprField("<+input>"),
			infraDefs:    nil,
			wantDeployTo: nil,
			wantAllInfra: "<+input>",
		},
		{
			name:         "deployToAll other expression with no infra list -> marker only",
			deployToAll:  exprField("<+pipeline.variables.everywhere>"),
			infraDefs:    nil,
			wantDeployTo: nil,
			wantAllInfra: "<+pipeline.variables.everywhere>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deployTo, allInfra := resolveDeployTo(tt.deployToAll, tt.infraDefs)
			if allInfra != tt.wantAllInfra {
				t.Fatalf("expected all-infra=%v, got all-infra=%v", tt.wantAllInfra, allInfra)
			}
			if !reflect.DeepEqual(deployTo, tt.wantDeployTo) {
				t.Fatalf("expected deploy-to=%#v, got deploy-to=%#v", tt.wantDeployTo, deployTo)
			}
		})
	}
}

func TestConvertEnvironment_AllInfraMarker(t *testing.T) {
	tests := []struct {
		name         string
		deployToAll  *flexible.Field[bool]
		wantDeployTo interface{}
		wantAllInfra interface{}
	}{
		{
			name:         "deployToAll true -> marker alongside deploy-to",
			deployToAll:  &flexible.Field[bool]{Value: true},
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: true,
		},
		{
			name:         "deployToAll false -> deploy-to only",
			deployToAll:  &flexible.Field[bool]{Value: false},
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: nil,
		},
		{
			name:         "deployToAll <+input> -> marker carried across alongside deploy-to",
			deployToAll:  exprField("<+input>"),
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: "<+input>",
		},
		{
			name:         "deployToAll expression -> marker carried across alongside deploy-to",
			deployToAll:  exprField("<+pipeline.variables.everywhere>"),
			wantDeployTo: []string{"infra1", "infra2"},
			wantAllInfra: "<+pipeline.variables.everywhere>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.Environment{
				EnvironmentRef:            "env1",
				DeployToAll:               tt.deployToAll,
				InfrastructureDefinitions: infraDefs("infra1", "infra2"),
			}
			got := ConvertEnvironment(src, NewStageConversionContext())
			items := got.ItemList()
			if len(items) != 1 {
				t.Fatalf("expected 1 environment item, got %d", len(items))
			}
			if items[0].AllInfra != tt.wantAllInfra {
				t.Fatalf("expected all-infra=%v, got all-infra=%v", tt.wantAllInfra, items[0].AllInfra)
			}
			if !reflect.DeepEqual(items[0].DeployTo, tt.wantDeployTo) {
				t.Fatalf("expected deploy-to=%#v, got deploy-to=%#v", tt.wantDeployTo, items[0].DeployTo)
			}
		})
	}
}

// The pipeline is serialized to JSON before YAML (see yaml.MarshalPipeline), and encoding/json
// omits an interface field only when it is nil - a non-nil interface holding false is written
// out. An absent marker must therefore serialize to no all-infra key at all.
func TestConvertEnvironment_AllInfraMarkerSerialization(t *testing.T) {
	tests := []struct {
		name         string
		deployToAll  *flexible.Field[bool]
		wantPresent  bool
		wantAllInfra interface{}
	}{
		{
			name:        "deployToAll false -> no all-infra key",
			deployToAll: &flexible.Field[bool]{Value: false},
			wantPresent: false,
		},
		{
			name:        "deployToAll unspecified -> no all-infra key",
			deployToAll: nil,
			wantPresent: false,
		},
		{
			name:         "deployToAll true -> all-infra: true",
			deployToAll:  &flexible.Field[bool]{Value: true},
			wantPresent:  true,
			wantAllInfra: true,
		},
		{
			name:         "deployToAll <+input> -> all-infra: <+input>",
			deployToAll:  exprField("<+input>"),
			wantPresent:  true,
			wantAllInfra: "<+input>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.Environment{
				EnvironmentRef:            "env1",
				DeployToAll:               tt.deployToAll,
				InfrastructureDefinitions: infraDefs("infra1", "infra2"),
			}
			items := ConvertEnvironment(src, NewStageConversionContext()).ItemList()
			raw, err := json.Marshal(items[0])
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			var got map[string]interface{}
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			allInfra, present := got["all-infra"]
			if present != tt.wantPresent {
				t.Fatalf("expected all-infra present=%v, got %s", tt.wantPresent, raw)
			}
			if tt.wantPresent && allInfra != tt.wantAllInfra {
				t.Fatalf("expected all-infra=%#v, got %#v", tt.wantAllInfra, allInfra)
			}
			// The deploy-to sibling survives in every case, so the v0 infra list is never lost.
			if _, present := got["deploy-to"]; !present {
				t.Fatalf("expected deploy-to alongside the marker, got %s", raw)
			}
		})
	}
}

func TestConvertEnvironmentGroup_AllEnvMarker(t *testing.T) {
	tests := []struct {
		name        string
		deployToAll *flexible.Field[bool]
		wantAllEnv  interface{}
	}{
		{
			name:        "group deployToAll true -> all-env marker",
			deployToAll: &flexible.Field[bool]{Value: true},
			wantAllEnv:  true,
		},
		{
			name:        "group deployToAll false -> no marker",
			deployToAll: &flexible.Field[bool]{Value: false},
			wantAllEnv:  nil,
		},
		{
			name:        "group deployToAll unspecified -> no marker",
			deployToAll: nil,
			wantAllEnv:  nil,
		},
		{
			name:        "group deployToAll <+input> -> marker carried across",
			deployToAll: exprField("<+input>"),
			wantAllEnv:  "<+input>",
		},
		{
			name:        "group deployToAll expression -> marker carried across",
			deployToAll: exprField("<+pipeline.variables.everywhere>"),
			wantAllEnv:  "<+pipeline.variables.everywhere>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environments := &flexible.Field[[]*v0.Environment]{}
			environments.Set([]*v0.Environment{
				{
					EnvironmentRef:            "env1",
					InfrastructureDefinitions: infraDefs("infra1"),
				},
			})
			src := &v0.EnvironmentGroup{
				EnvGroupRef:  "group1",
				DeployToAll:  tt.deployToAll,
				Environments: environments,
			}
			got := ConvertEnvironmentGroup(src, NewStageConversionContext())
			if got == nil {
				t.Fatalf("expected an EnvironmentRef, got nil")
			}
			groupConfig, ok := got.Group.(map[string]interface{})
			if !ok {
				t.Fatalf("expected a group config map, got %#v", got.Group)
			}
			allEnv, present := groupConfig["all-env"]
			if tt.wantAllEnv != nil {
				if !present || allEnv != tt.wantAllEnv {
					t.Fatalf("expected all-env=%#v, got %#v", tt.wantAllEnv, groupConfig)
				}
			} else if present {
				t.Fatalf("expected no all-env key, got %#v", groupConfig)
			}
			// The marker never replaces items: the per-environment config must survive.
			if _, present := groupConfig["items"]; !present {
				t.Fatalf("expected items alongside the group config, got %#v", groupConfig)
			}
		})
	}
}

// A group ref with no environments still carries the marker, so the server resolves the
// full environment list at execution time.
func TestConvertEnvironmentGroup_AllEnvMarkerWithoutEnvironments(t *testing.T) {
	tests := []struct {
		name        string
		deployToAll *flexible.Field[bool]
		want        map[string]interface{}
	}{
		{
			name:        "group deployToAll true -> marker only",
			deployToAll: &flexible.Field[bool]{Value: true},
			want:        map[string]interface{}{"id": "group1", "all-env": true},
		},
		{
			name:        "group deployToAll <+input> -> marker only",
			deployToAll: exprField("<+input>"),
			want:        map[string]interface{}{"id": "group1", "all-env": "<+input>"},
		},
		{
			name:        "group deployToAll other expression -> marker only",
			deployToAll: exprField("<+pipeline.variables.everywhere>"),
			want:        map[string]interface{}{"id": "group1", "all-env": "<+pipeline.variables.everywhere>"},
		},
		{
			name:        "group deployToAll false -> no marker",
			deployToAll: &flexible.Field[bool]{Value: false},
			want:        map[string]interface{}{"id": "group1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.EnvironmentGroup{
				EnvGroupRef: "group1",
				DeployToAll: tt.deployToAll,
			}
			got := ConvertEnvironmentGroup(src, NewStageConversionContext())
			if got == nil {
				t.Fatalf("expected an EnvironmentRef, got nil")
			}
			if !reflect.DeepEqual(got.Group, tt.want) {
				t.Fatalf("expected group=%#v, got group=%#v", tt.want, got.Group)
			}
		})
	}
}

// Nothing here warns. UNSUPPORTED_EXPRESSION means the conversion was lossy everywhere else
// it is used - services.values and environments.values skip conversion outright - whereas an
// expression is carried across unchanged and is a value v1's all-infra accepts.
func TestResolveDeployTo_AllInfraMarkerNoWarning(t *testing.T) {
	tests := []struct {
		name        string
		deployToAll *flexible.Field[bool]
		infraDefs   *flexible.Field[[]*v0.InfrastructureDefinition]
	}{
		{
			name:        "other expression with an infra list",
			deployToAll: exprField("<+pipeline.variables.everywhere>"),
			infraDefs:   infraDefs("infra1", "infra2"),
		},
		{
			name:        "other expression with no infra list",
			deployToAll: exprField("<+pipeline.variables.everywhere>"),
			infraDefs:   nil,
		},
		{
			name:        "<+input>",
			deployToAll: exprField("<+input>"),
			infraDefs:   infraDefs("infra1", "infra2"),
		},
		{
			name:        "deployToAll true",
			deployToAll: &flexible.Field[bool]{Value: true},
			infraDefs:   infraDefs("infra1", "infra2"),
		},
		{
			name:        "deployToAll false",
			deployToAll: &flexible.Field[bool]{Value: false},
			infraDefs:   infraDefs("infra1", "infra2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messagelog.ResetMessageLogger()
			defer messagelog.ResetMessageLogger()
			logger := messagelog.GetMessageLogger()
			logger.Enable("")
			logger.SetCurrentFile("pipeline.yaml")

			resolveDeployTo(tt.deployToAll, tt.infraDefs)

			if fileLog := logger.GetFileLog("pipeline.yaml"); fileLog != nil {
				t.Fatalf("expected no messages, got %#v", fileLog.Messages)
			}
		})
	}
}

// The group path carries the marker across without warning, for the same reason.
func TestResolveAllEnvMarker_NoWarning(t *testing.T) {
	messagelog.ResetMessageLogger()
	defer messagelog.ResetMessageLogger()
	logger := messagelog.GetMessageLogger()
	logger.Enable("")
	logger.SetCurrentFile("pipeline.yaml")

	resolveAllEnvMarker(exprField("<+pipeline.variables.everywhere>"))

	if fileLog := logger.GetFileLog("pipeline.yaml"); fileLog != nil {
		t.Fatalf("expected no messages, got %#v", fileLog.Messages)
	}
}

func TestConvertDeploymentInfrastructure_AllInfraMarker(t *testing.T) {
	tests := []struct {
		name         string
		infraDef     v0.InfrastructureDefinition
		wantDeployTo interface{}
		wantAllInfra interface{}
	}{
		{
			name:         "no infrastructure definition -> marker only",
			infraDef:     v0.InfrastructureDefinition{},
			wantDeployTo: nil,
			wantAllInfra: true,
		},
		{
			name:         "explicit infrastructure definition -> deploy-to only",
			infraDef:     v0.InfrastructureDefinition{Identifier: "infra1"},
			wantDeployTo: "infra1",
			wantAllInfra: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.DeploymentInfrastructure{
				EnvironmentRef:           "env1",
				InfrastructureDefinition: tt.infraDef,
			}
			got := ConvertDeploymentInfrastructure(src)
			items := got.ItemList()
			if len(items) != 1 {
				t.Fatalf("expected 1 environment item, got %d", len(items))
			}
			if items[0].AllInfra != tt.wantAllInfra {
				t.Fatalf("expected all-infra=%v, got all-infra=%v", tt.wantAllInfra, items[0].AllInfra)
			}
			if !reflect.DeepEqual(items[0].DeployTo, tt.wantDeployTo) {
				t.Fatalf("expected deploy-to=%#v, got deploy-to=%#v", tt.wantDeployTo, items[0].DeployTo)
			}
		})
	}
}

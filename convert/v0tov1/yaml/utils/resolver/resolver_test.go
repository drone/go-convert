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

package resolver

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	schema "github.com/bradrydzewski/spec/yaml"

	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestResolveStep(t *testing.T) {
	out, err := diff("testdata/step1.yaml", "testdata/templates/golang.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	if out != "" {
		t.Log(out)
		t.Errorf("Parsed Yaml did not match expected Yaml")
	}
}

func TestResolveStep_WithInputs(t *testing.T) {
	out, err := diff("testdata/step2.yaml", "testdata/templates/golang.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	if out != "" {
		t.Log(out)
		t.Errorf("Parsed Yaml did not match expected Yaml")
	}
}

func TestResolveStage(t *testing.T) {
	out, err := diff("testdata/stage1.yaml", "testdata/templates/golang.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	if out != "" {
		t.Log(out)
		t.Errorf("Parsed Yaml did not match expected Yaml")
	}
}

func diff(file, template string) (string, error) {
	// decode the yaml file
	parsed, err := schema.ParseFile(file)
	if err != nil {
		return "", err
	}
	// resolve the template
	Resolve(parsed, func(name string) (*schema.Template, error) {
		out, err := schema.ParseFile(template)
		if err != nil {
			return nil, err
		}
		return out.Template, nil
	})
	// re-encode the yaml file
	b1, err := json.Marshal(parsed)
	if err != nil {
		return "", err
	}
	// parse the golden yaml file and convert to json
	b2, err := os.ReadFile(
		strings.ReplaceAll(file, ".yaml", ".yaml.golden"),
	)
	if err != nil {
		return "", err
	}
	b2, err = yaml.YAMLToJSON(b2)
	if err != nil {
		return "", err
	}
	// unmarshal both json files into map structures
	m1 := map[string]interface{}{}
	m2 := map[string]interface{}{}
	if err := json.Unmarshal(b1, &m1); err != nil {
		return "", err
	}
	if err := json.Unmarshal(b2, &m2); err != nil {
		return "", err
	}
	// diff the map structures. if the diff is empty this
	// means they match.
	return cmp.Diff(m1, m2), nil
}

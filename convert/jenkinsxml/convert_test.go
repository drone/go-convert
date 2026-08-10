// Copyright 2024 Harness, Inc.
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

package jenkinsxml

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"

	jenkinsxml "github.com/drone/go-convert/convert/jenkinsxml/xml"
	harness "github.com/drone/spec/dist/go"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

// TestConvertNoBuilders verifies that a job without a freestyle <builders>
// block (for example a maven2-moduleset or scripted flow-definition) converts
// without panicking on a nil Builders dereference.
func TestConvertNoBuilders(t *testing.T) {
	converter := New()
	out, err := converter.ConvertFile("testdata/no-builders.xml")
	if err != nil {
		t.Error(err)
		return
	}

	got := map[string]interface{}{}
	if err := yaml.Unmarshal(out, &got); err != nil {
		t.Error(err)
		return
	}

	// a builders-less job should still yield a pipeline resource.
	if kind, _ := got["kind"].(string); kind != "pipeline" {
		t.Errorf("expected kind: pipeline, got %q", kind)
	}
}

// TestConvertMavenGoals verifies that a maven2-moduleset job converts its
// top-level <goals> into an mvn step instead of an empty stage.
func TestConvertMavenGoals(t *testing.T) {
	out, err := New().ConvertFile("testdata/maven.xml")
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Contains(out, []byte("mvn clean install")) {
		t.Errorf("expected an 'mvn clean install' step in the output, got:\n%s", out)
	}
}

// TestConvertSCM verifies that a Jenkins Git SCM converts to a git clone
// step carrying the remote url and the branch (stripped of the "*/" prefix).
func TestConvertSCM(t *testing.T) {
	out, err := New().ConvertFile("testdata/scm.xml")
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Contains(out, []byte("https://git.example.com/scm/team/order-service.git")) {
		t.Errorf("expected the git url in the output, got:\n%s", out)
	}
	if !bytes.Contains(out, []byte("master")) || bytes.Contains(out, []byte("*/master")) {
		t.Errorf("expected branch \"master\" (prefix stripped) in the output, got:\n%s", out)
	}
}

// TestConvertParameters verifies that Jenkins string build parameters are
// mapped to pipeline inputs with their default values.
func TestConvertParameters(t *testing.T) {
	out, err := New().ConvertFile("testdata/parameters.xml")
	if err != nil {
		t.Error(err)
		return
	}

	got := map[string]interface{}{}
	if err := yaml.Unmarshal(out, &got); err != nil {
		t.Error(err)
		return
	}

	spec, _ := got["spec"].(map[string]interface{})
	inputs, ok := spec["inputs"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected inputs in the converted output, got:\n%s", out)
	}
	release, ok := inputs["RELEASE"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected a RELEASE input, got: %v", inputs)
	}
	if release["default"] != "master" {
		t.Errorf("expected RELEASE default \"master\", got %q", release["default"])
	}
}

// TODO: add more teste in subdirectories, as we have for other providers
func TestConvert(t *testing.T) {
	// convert the XML file from Jenkins to harness
	converter := New()
	tmp1, err := converter.ConvertFile("testdata/hello.xml")
	if err != nil {
		t.Error(err)
		return
	}

	// unmarshal the yaml to a map
	got := map[string]interface{}{}
	if err := yaml.Unmarshal(tmp1, &got); err != nil {
		t.Error(err)
		return
	}

	// parse the golden yaml file
	data, err := ioutil.ReadFile("testdata/hello.xml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	// unmarshal the golden yaml file to a map
	want := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &want); err != nil {
		t.Error(err)
		return
	}

	// compare the converted yaml to the golden file
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected conversion result")
		t.Log(diff)
	}
}

func TestConvertAntTaskToStep(t *testing.T) {
	// task struct to test
	task := &jenkinsxml.Task{
		XMLName: xml.Name{
			Local: "hudson.tasks.Ant",
			Space: "",
		},
		Content: "<targets>one/two/three</targets>",
	}

	got := convertAntTaskToStep(task)

	want := &harness.Step{
		Name: "ant",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "harnesscommunitytest/ant-plugin",
			Inputs: map[string]interface{}{
				"goals": "one/two/three",
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected conversion result")
		t.Log(diff)
	}
}

func TestConvertShellTaskToStep(t *testing.T) {
	// task struct to test
	task := &jenkinsxml.Task{
		XMLName: xml.Name{
			Local: "hudson.tasks.Shell",
			Space: "",
		},
		Content: `
      <command>echo hello</command>
      <configuredLocalRules/>
	`,
	}

	got := convertShellTaskToStep(task)

	want := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "echo hello",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected conversion result")
		t.Log(diff)
	}
}

func TestUnsupportedTaskToStep(t *testing.T) {
	task := "hudson.tasks.Unknown"

	got := unsupportedTaskToStep(task)

	want := &harness.Step{
		Name: "shell",
		Type: "script",
		Spec: &harness.StepExec{
			Run: "echo Unsupported field hudson.tasks.Unknown",
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected conversion result")
		t.Log(diff)
	}
}

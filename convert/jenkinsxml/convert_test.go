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
	"encoding/xml"
	"io/ioutil"
	"testing"

	jenkinsxml "github.com/drone/go-convert/convert/jenkinsxml/xml"
	harness "github.com/drone/spec/dist/go"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

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

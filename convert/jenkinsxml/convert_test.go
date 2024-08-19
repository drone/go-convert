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
	"io/ioutil"
	"testing"

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

	// unmarshal the converted yaml file to a map
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

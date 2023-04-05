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

import (
	"io/ioutil"
	"testing"
)

func TestParseString(t *testing.T) {
	out, _ := ioutil.ReadFile("testdata/sample.1.yaml")
	_, err := ParseString(string(out))
	if err != nil {
		t.Error(err)
	}
}

func TestParseString_Error(t *testing.T) {
	_, err := ParseString("[]")
	if err == nil {
		t.Errorf("Expect error when yaml is invalid")
	}
}

func TestParseFile(t *testing.T) {
	_, err := ParseFile("testdata/sample.1.yaml")
	if err != nil {
		t.Error(err)
	}
}

func TestParseFile_Error(t *testing.T) {
	_, err := ParseFile("testdata/file-does-not-exist.yaml")
	if err == nil {
		t.Errorf("Expect error when file not exists")
	}
}

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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	// unmarshal the test pipeline.
	got, err := ParseFile("testdata/pipeline.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	// unmarshal the expected json struct.
	src, _ := ioutil.ReadFile("testdata/pipeline.golden.json")
	want := new(Config)
	if err := json.Unmarshal(src, want); err != nil {
		t.Error(err)
		return
	}

	// compare the test with the expected output.
	// if they do not match, write the diff to the logs
	// and fail the test.
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf(diff)
	}
}

// dump json data to the test logs.
func debug(t *testing.T, v interface{}) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	enc.Encode(v)
	t.Log(buf.String())
}

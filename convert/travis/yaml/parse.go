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
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse parses the configuration from io.Reader r.
func Parse(r io.Reader) (*Pipeline, error) {
	out := new(Pipeline)
	dec := yaml.NewDecoder(r)
	err := dec.Decode(out)
	return out, err
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) (*Pipeline, error) {
	return Parse(
		bytes.NewBuffer(b),
	)
}

// ParseString parses the configuration from string s.
func ParseString(s string) (*Pipeline, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) (*Pipeline, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

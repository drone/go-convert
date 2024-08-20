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

package xml

import (
	"bytes"
	"io"
	"os"
	"strings"

	"encoding/xml"
)

// Parse parses the configuration from io.Reader r.
func Parse(r io.Reader) (*Project, error) {
	out := new(Project)

	// see https://github.com/golang/go/issues/25755
	// encoding/xml does not support XML 1.1, which jenkins uses
	//
	// TODO: this approach is likely brittle and will need to be revisited
	data, _ := io.ReadAll(r)
	res := strings.Replace(string(data), "<?xml version='1.1", "<?xml version='1.0", 1)

	dec := xml.NewDecoder(strings.NewReader(res))
	err := dec.Decode(out)
	return out, err
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) (*Project, error) {
	return Parse(
		bytes.NewBuffer(b),
	)
}

// ParseString parses the configuration from string s.
func ParseString(s string) (*Project, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) (*Project, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

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
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

// Parse parses the configuration from io.Reader r. Unknown fields (keys that
// do not match any json tag in the Config schema) are silently ignored by the
// JSON decoder; use ParseWithUnknownFields to also surface their JSON paths.
func Parse(r io.Reader) (*Config, error) {
	cfg, _, err := ParseWithUnknownFields(r)
	return cfg, err
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) (*Config, error) {
	cfg, _, err := ParseBytesWithUnknownFields(b)
	return cfg, err
}

// ParseString parses the configuration from string s.
func ParseString(s string) (*Config, error) {
	cfg, _, err := ParseStringWithUnknownFields(s)
	return cfg, err
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) (*Config, error) {
	cfg, _, err := ParseFileWithUnknownFields(p)
	return cfg, err
}

// ParseWithUnknownFields parses the configuration from io.Reader r and, in
// addition to the typed Config, returns a sorted list of JSON paths for keys
// present in the input that have no matching field in the Config schema.
//
// Unknown fields do NOT cause the parse to fail — the returned Config is
// populated exactly as ParseWithUnknownFields's predecessor would produce it.
// Callers can log or surface the unknown-fields list for observability.
func ParseWithUnknownFields(r io.Reader) (*Config, []string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	return ParseBytesWithUnknownFields(b)
}

// ParseBytesWithUnknownFields parses the configuration from bytes b and
// returns unknown-field paths alongside the typed Config. See
// ParseWithUnknownFields for semantics.
func ParseBytesWithUnknownFields(b []byte) (*Config, []string, error) {
	jsonBytes, err := yaml.YAMLToJSON(b)
	if err != nil {
		return nil, nil, err
	}
	out := new(Config)
	if err := json.Unmarshal(jsonBytes, out); err != nil {
		return out, nil, err
	}
	var raw interface{}
	// A second unmarshal into interface{} preserves the original tree shape
	// so we can diff against the schema. This is unavoidable: encoding/json
	// discards unknown keys during the typed decode.
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return out, nil, nil
	}
	return out, collectUnknownFields(out, raw), nil
}

// ParseStringWithUnknownFields parses the configuration from string s.
func ParseStringWithUnknownFields(s string) (*Config, []string, error) {
	return ParseBytesWithUnknownFields([]byte(s))
}

// ParseFileWithUnknownFields parses the configuration from path p.
func ParseFileWithUnknownFields(p string) (*Config, []string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	return ParseBytesWithUnknownFields(b)
}


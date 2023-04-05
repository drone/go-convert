package yaml

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse parses the configuration from io.Reader r.
func Parse(r io.Reader) (*Pipeline, error) {
	buf := repairOn(r)
	out := new(Pipeline)
	dec := yaml.NewDecoder(&buf)
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

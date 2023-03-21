package yaml

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func Parse(r io.Reader) ([]*Pipeline, error) {
	res := []*Pipeline{}
	dec := yaml.NewDecoder(r)
	for {
		out := new(Pipeline)
		err := dec.Decode(out)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, out)
	}
	return res, nil
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) ([]*Pipeline, error) {
	return Parse(
		bytes.NewBuffer(b),
	)
}

// ParseString parses the configuration from string s.
func ParseString(s string) ([]*Pipeline, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) ([]*Pipeline, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

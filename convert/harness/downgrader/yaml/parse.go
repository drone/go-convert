package yaml

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	v1 "github.com/drone/spec/dist/go"
)

// Parse parses the configuration from io.Reader r.
func Parse(r io.Reader) ([]*v1.Config, error) {
	var resources []*v1.Config

	// Read all data from the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	docs := bytes.Split(b, []byte("\n---\n"))

	for _, doc := range docs {
		r := bytes.NewReader(doc)
		resource, err := v1.Parse(r)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) ([]*v1.Config, error) {
	return Parse(
		bytes.NewBuffer(b),
	)
}

// ParseString parses the configuration from string s.
func ParseString(s string) ([]*v1.Config, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) ([]*v1.Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

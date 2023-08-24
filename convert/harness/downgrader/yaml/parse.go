package yaml

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	v1 "github.com/drone/spec/dist/go"
)

// Parse parses the configuration from io.Reader r.
func Parse(r io.Reader) ([]*v1.Pipeline, error) {
	var pipelines []*v1.Pipeline

	// Read all data from the reader
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	docs := bytes.Split(b, []byte("\n---\n"))

	for _, doc := range docs {
		// workaround for this issue
		// https://stackoverflow.com/questions/70849190/golang-how-to-avoid-double-quoted-on-key-on-when-marshaling-struct-to-yaml
		doc = bytes.ReplaceAll(doc, []byte(` "on":`), []byte(` on:`))
		docReader := bytes.NewReader(doc)
		parsedPipeline, err := v1.Parse(docReader)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, parsedPipeline)
	}

	return pipelines, nil
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) ([]*v1.Pipeline, error) {
	return Parse(
		bytes.NewBuffer(b),
	)
}

// ParseString parses the configuration from string s.
func ParseString(s string) ([]*v1.Pipeline, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) ([]*v1.Pipeline, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

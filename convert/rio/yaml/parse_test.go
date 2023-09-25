package yaml

import (
	"testing"
)

func TestShort(t *testing.T) {
	fileName := "testdata/short.yaml"
	_, err := ParseFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLong(t *testing.T) {
	fileName := "testdata/long.yaml"
	_, err := ParseFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
}

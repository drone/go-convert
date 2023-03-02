// Copyright 2023 Harness Inc. All rights reserved.

package slug

import "testing"

func TestSlug(t *testing.T) {
	tests := []struct {
		a, b string
	}{
		{"foo bar", "foobar"},
		{"Foo Bar", "foobar"},
		{"Foo-Bar", "foobar"},
		{"Foo/Bar", "foobar"},
	}
	for _, test := range tests {
		got, want := Create(test.a), test.b
		if got != want {
			t.Errorf("Want slug %q, got %q", want, got)
		}
	}
}

package rand

import (
	"testing"
	"unicode"
)

func TestAlphaLength(t *testing.T) {
	lengths := []int{0, 1, 10, 100}
	for _, n := range lengths {
		s := Alpha(n)
		if len(s) != n {
			t.Errorf("Alpha(%d) returned string of length %d", n, len(s))
		}
	}
}

func TestAlphaNegativeLength(t *testing.T) {
	s := Alpha(-1)
	if s != "" {
		t.Errorf("Calling Alpha() with a negative number should return an empty string")
	}
}

func TestAlphanumericLength(t *testing.T) {
	lengths := []int{0, 1, 10, 100}
	for _, n := range lengths {
		s := Alphanumeric(n)
		if len(s) != n {
			t.Errorf("Alphanumeric(%d) returned string of length %d", n, len(s))
		}
	}
}

func TestAlphaCharset(t *testing.T) {
	s := Alpha(1000)
	for _, ch := range s {
		if !unicode.IsLetter(ch) {
			t.Errorf("Alpha string contains non-letter character: %q", ch)
		}
	}
}

func TestAlphanumericCharset(t *testing.T) {
	s := Alphanumeric(1000)
	for _, ch := range s {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			t.Errorf("Alphanumeric string contains invalid character: %q", ch)
		}
	}
}

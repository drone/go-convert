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

package downgrader

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// this function attempts to convert a freeform name to a
// harness name.
//
// ^[a-zA-Z_][-0-9a-zA-Z_\s]{0,127}$
func convertName(s string) string {
	s = strings.TrimSpace(s)
	s = norm.NFKD.String(s)
	s = strings.Map(safe, s)
	if len(s) > 0 {
		f := string(s[0]) // first letter of the string
		return strings.Map(safeFirstLetter, f) + s[1:]
	}
	if len(s) > 127 {
		return s[:127] // trim if > 127 characters
	}
	return s
}

// helper function maps restricted runes to allowed runes
// for the name.
func safe(r rune) rune {
	switch {
	case unicode.IsSpace(r):
		return ' '
	case unicode.IsNumber(r):
		return r
	case unicode.IsLetter(r):
		return r
	}
	switch r {
	case '_', '-':
		return r
	}
	return -1
}

// helper function maps restricted runes to allowed runes
// for the first letter of the name.
func safeFirstLetter(r rune) rune {
	switch {
	case unicode.IsNumber(r):
		return r
	case unicode.IsLetter(r):
		return r
	}
	switch r {
	case '_':
		return r
	}
	return -1
}

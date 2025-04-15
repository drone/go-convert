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

// Package slug provides utilities for working with slug values.
package slug

import (
	"github.com/drone/go-convert/internal/rand"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var safeRanges = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Number,
}

func safe(r rune) rune {
	switch {
	case unicode.IsOneOf(safeRanges, r):
		return unicode.ToLower(r)
	}
	return -1
}

// Create creates a slug from a string.
func Create(s string) string {
	s = norm.NFKD.String(s)
	s = strings.Map(safe, s)
	return s
}

// CreateWithRandom creates a slug, appending a random identifier.
func CreateWithRandom(s string) string {
	return Create(s) + "_" + rand.Alphanumeric(8)
}

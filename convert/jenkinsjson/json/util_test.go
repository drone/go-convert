// Copyright 2023 Harness, Inc.
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

package json

import (
	"testing"
)

func TestSanitizeForId(t *testing.T) {
	tests := []struct {
		name       string
		spanId     string
		spanName   string
		expectedId string
	}{
		{
			name:       "ReplaceSpacesWithUnderscore",
			spanId:     "123456",
			spanName:   "a b c",
			expectedId: "a_b_c123456",
		},
		{
			name:       "RealSample",
			spanId:     "a2e5df",
			spanName:   "Deploy to DEV-CT",
			expectedId: "Deploy_to_DEV_CTa2e5df",
		},
		{
			name:       "TruncateLongName",
			spanId:     "123456",
			spanName:   "this string_is_longer_than_58_character_gets_cut_off_here_all_of_this_gets_trucated",
			expectedId: "this_string_is_longer_than_58_character_gets_cut_off_here_123456",
		},
		{
			name:       "ReplaceColonWithUnderscore",
			spanId:     "123456",
			spanName:   "a string with a : in it",
			expectedId: "a_string_with_a_in_it123456",
		},
		{
			name:       "RemoveLeadingUnderscores",
			spanId:     "123456",
			spanName:   "? a string",
			expectedId: "a_string123456",
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			id := SanitizeForId(tc.spanName, tc.spanId)

			if id != tc.expectedId {
				t.Errorf("%v failed, got %v expected %v", tc.name, id, tc.expectedId)
			}
		})
	}

}

func TestSanitizeForName(t *testing.T) {
	tests := []struct {
		name         string
		spanName     string
		expectedName string
	}{
		{
			name:         "DoNotReplaceSpaces",
			spanName:     "a b c",
			expectedName: "a b c",
		},
		{
			name:         "TruncateLongName",
			spanName:     "this_long_string_is_longer_than_the_max_allowed_limit_of_128_characters_so_anything_beyond_that_length_gets_cut_off_starting_now_all_of_this_gets_trucated",
			expectedName: "this_long_string_is_longer_than_the_max_allowed_limit_of_128_characters_so_anything_beyond_that_length_gets_cut_off_starting_now",
		},
		{
			name:         "ReplaceColonWithUnderscore",
			spanName:     "a string with a : in it",
			expectedName: "a string with a in it",
		},
		{
			name:         "RemoveLeadingSpacesAndUnderscores",
			spanName:     "? a string",
			expectedName: "a string",
		},
		{
			name:         "RemoveTrailingSpacesAndUnderscores",
			spanName:     "a string + ",
			expectedName: "a string",
		},
		{
			name:         "RemoveRepeatingSpacesAndUnderscores",
			spanName:     "a string   with multiple spaces_____and underscores",
			expectedName: "a string with multiple spaces_and underscores",
		},
		{
			name:         "RemoveRepeatingCharsReplacingWithTheLastUsed",
			spanName:     "group ending with   _underscore, group ending with_____ space, group ending with    -hyphen",
			expectedName: "group ending with_underscore group ending with space group ending with-hyphen",
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			name := SanitizeForName(tc.spanName)

			if name != tc.expectedName {
				t.Errorf("%v failed,\ngot      '%v'\nexpected '%v'", tc.name, name, tc.expectedName)
			}
		})
	}

}

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

package circle

import (
	"testing"
)

func TestReplaceParams(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		// no params
		{
			before: "foo\nbar\nbaz\n",
			after:  "foo\nbar\nbaz\n",
		},
		// open bracket only
		{
			before: "foo\nbar<<baz\n",
			after:  "foo\nbar<<baz\n",
		},
		// close bracket only
		{
			before: "foo\nbar>>baz\n",
			after:  "foo\nbar>>baz\n",
		},
		// open bracket before close bracket
		{
			before: "foo\n>>bar<<baz\n",
			after:  "foo\n>>bar<<baz\n",
		},
		// unknown param
		{
			before: "foo\n<< foo.bar >>\nbaz\n",
			after:  "foo\n<< foo.bar >>\nbaz\n",
		},
		{
			before: "foo\n<< pipeline.id >>\nbaz\n",
			after:  "foo\n<+pipeline.identifier>\nbaz\n",
		},
		// pipeline parameter
		{
			before: "foo\n<< parameters.message >>\nbaz\n",
			after:  "foo\n<+inputs.message>\nbaz\n",
		},
		// job parameter
		{
			before: "foo\n<< pipeline.parameters.message >>\nbaz\n",
			after:  "foo\n<+inputs.message>\nbaz\n",
		},
	}

	for _, test := range tests {
		t.Run(test.before, func(t *testing.T) {
			out := replaceParams([]byte(test.before), params)
			if got, want := string(out), test.after; got != want {
				t.Errorf("Got replace params %q, want %q", got, want)
			}
		})
	}
}

func TestExpandParams(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		// no params
		{
			before: "foo\nbar\nbaz\n",
			after:  "foo\nbar\nbaz\n",
		},
		// open bracket only
		{
			before: "foo\nbar<<baz\n",
			after:  "foo\nbar<<baz\n",
		},
		// close bracket only
		{
			before: "foo\nbar>>baz\n",
			after:  "foo\nbar>>baz\n",
		},
		// open bracket before close bracket
		{
			before: "foo\n>>bar<<baz\n",
			after:  "foo\n>>bar<<baz\n",
		},
		// unknown param
		{
			before: "foo\n<< foo.bar >>\nbaz\n",
			after:  "foo\n<< foo.bar >>\nbaz\n",
		},
		// known parameter
		{
			before: "foo\n<< pipeline.branch >>\nbaz\n",
			after:  "foo\nmain\nbaz\n",
		},
	}

	for _, test := range tests {
		t.Run(test.before, func(t *testing.T) {
			out := expandParams([]byte(test.before), map[string]string{"pipeline.branch": "main"})
			if got, want := string(out), test.after; got != want {
				t.Errorf("Got expanded param %q, want %q", got, want)
			}
		})
	}
}

func TestExractParam(t *testing.T) {
	tests := []struct {
		before string
		after  string
	}{
		// no params
		{
			before: "foo\nbar\nbaz\n",
			after:  "",
		},
		// open bracket only
		{
			before: "foo\nbar<<baz\n",
			after:  "",
		},
		// close bracket only
		{
			before: "foo\nbar>>baz\n",
			after:  "",
		},
		// open bracket before close bracket
		{
			before: "foo\n>>bar<<baz\n",
			after:  "",
		},
		// known parameter
		{
			before: "foo\n<< pipeline.id >>\nbaz\n",
			after:  "pipeline.id",
		},
	}

	for _, test := range tests {
		t.Run(test.before, func(t *testing.T) {
			out := extractParam(test.before)
			if got, want := out, test.after; got != want {
				t.Errorf("Extracted param %q, wanted %q", got, want)
			}
		})
	}
}

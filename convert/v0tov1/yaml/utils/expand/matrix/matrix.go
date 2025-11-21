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

package matrix

import (
	"sort"
	"strings"
)

const (
	limitTags = 10
	limitAxis = 25
)

// Matrix represents the build matrix.
type Matrix map[string][]string

// Axis represents a single permutation of entries from
// the build matrix.
type Axis map[string]string

// String returns a string representation of an Axis as a
// comma-separated list
// of environment variables.
func (a Axis) String() string {
	var envs []string
	for k, v := range a {
		envs = append(envs, k+"="+v)
	}
	// sort the slice to ensure the ordering
	// is deterministic.
	sort.SliceStable(envs, func(i, j int) bool {
		return envs[i] < envs[j]
	})
	return strings.Join(envs, " ")
}

// Calc calculates the matrix.
func Calc(matrix Matrix) []Axis {
	// calculate number of permutations and extract the
	// list of tags (ie go_version, redis_version, etc)
	var perm int
	var tags []string
	for k, v := range matrix {
		perm *= len(v)
		if perm == 0 {
			perm = len(v)
		}
		tags = append(tags, k)
	}

	// structure to hold the transformed result set
	axisList := []Axis{}

	// for each axis calculate the uniqe set of values that
	// should be used.
	for p := 0; p < perm; p++ {
		axis := map[string]string{}
		decr := perm
		for i, tag := range tags {
			elems := matrix[tag]
			decr = decr / len(elems)
			elem := p / decr % len(elems)
			axis[tag] = elems[elem]

			// enforce a maximum number of tags in the
			// build matrix.
			if i > limitTags {
				break
			}
		}

		// append to the list of axis.
		axisList = append(axisList, axis)

		// enforce a maximum number of axis that should
		// be calculated.
		if p > limitAxis {
			break
		}
	}

	// sort the slice to ensure the ordering
	// is deterministic.
	sort.SliceStable(axisList, func(i, j int) bool {
		return axisList[i].String() < axisList[j].String()
	})

	return axisList
}

//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//

package slicehelper

import "filesigner/constraints"

// ******** Public functions ********

// Fill fills a generic slice with a generic value in an efficient way.
func Fill[S ~[]T, T any](a S, value T) {
	aLen := ensureLengthIsCapacity(&a)

	if aLen > 0 {
		// Put the value into the first slice element
		a[0] = value

		// Incrementally duplicate the value into the rest of the slice
		for j := 1; j < aLen; j <<= 1 {
			copy(a[j:], a[:j])
		}
	}
}

// ClearNumber clears an integer type slice.
func ClearNumber[S ~[]T, T constraints.Number](a S) {
	Fill(a, 0)
}

// NewReverse returns a new array with the elements in the reverse order of the argument.
func NewReverse[S ~[]T, T any](a S) S {
	aLen := len(a)
	result := make(S, aLen)

	aLen--
	for i, e := range a {
		result[aLen-i] = e
	}

	return result
}

// Concat returns a new slice concatenating the passed in slices.
// This is a streamlined version of the slices.Concat function from Go V1.22.
func Concat[S ~[]E, E any](slices ...S) S {
	// 1. Calculate total size.
	size := 0
	for _, s := range slices {
		size += len(s)
	}

	// 2. Make new slice with the total size as the capacity and 0 length.
	result := make(S, 0, size)

	// 3. Append all source slices.
	for _, s := range slices {
		result = append(result, s...)
	}

	return result
}

// Copy makes a copy of a slice.
func Copy[S ~[]T, T any](a S) S {
	result := make(S, len(a))
	copy(result, a)
	return result
}

// SetCap sets the capacity of a slice to be at least n.
// If n is negative or too large to allocate the memory, SetCap panics.
func SetCap[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}

	c := cap(s)
	if n -= c; n > 0 {
		s = append(s[:c], make([]E, n)...)[:len(s)]
	}

	return s
}

// Prepend adds an element v at the beginning of a slice s.
func Prepend[S ~[]T, T any](v T, s S) []T {
	return append([]T{v}, s...)
}

// ******** Private functions ********

// ensureLengthIsCapacity ensures that the length of the slice is its capacity.
// We need the address of the slice as the parameter. If the '*' would be missing
// we would get a copy of the slice and not the slice itself.
func ensureLengthIsCapacity[S ~[]T, T any](a *S) int {
	ra := *a
	aLen := len(ra)
	aCap := cap(ra)
	if aLen != aCap {
		*a = ra[:aCap]
		aLen = aCap
	}

	return aLen
}

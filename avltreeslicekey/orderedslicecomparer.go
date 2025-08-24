//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
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
//    2025-08-24: V1.0.0: Created.
//

package avltreeslicekey

import "cmp"

// ******** Private constants ********

const (
	isEqual   = 0
	isLess    = -1
	isGreater = 1
)

// ******** Public functions ********

// Compare is a generic slice comparison function.
// It returns -1 if a < b, 0 if a == b, 1 if a > b.
func Compare[T cmp.Ordered](a, b []T) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	// Compare each element for the smaller of the two lengths.
	for i := 0; i < minLen; i++ {
		result := cmp.Compare(a[i], b[i])
		// Return if a difference is found.
		if result != isEqual {
			return result
		}
	}

	// There was no difference, yet.
	// If the lengths are equal, return 0.
	if len(a) == len(b) {
		return isEqual
	}

	// There was no difference, yet, but the lengths differ.
	// Return result based on the lengths of the slices.
	if len(a) < len(b) {
		return isLess
	} else {
		return isGreater
	}
}

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
// Version: 2.0.0
//
// Change history:
//    2025-08-31: V1.0.0: Created.
//    2026-08-31: V1.1.0: Use the inOrder iterator to build the results.
//    2026-08-31: V2.0.0: Replaced Keys, CountEntries and CountEntry by the
//                        allocation-free All iterator.
//

// Package avltreecounter provides a self-balancing binary counter tree with slice keys.
package avltreecounter

import (
	"cmp"
	"iter"
)

// ******** Public types ********

// AVLTree is a self-balancing binary counter tree with slice keys.
type AVLTree[K cmp.Ordered] struct {
	root  *avlNode[K]
	count int
}

// ******** Public functions ********

// Add inserts a new node into the tree.
func (t *AVLTree[K]) Add(key []K) {
	var madeNewNode bool
	t.root, madeNewNode = t.root.add(key)

	if madeNewNode {
		t.count++
	}
}

// Count returns the number of nodes in the tree.
func (t *AVLTree[K]) Count() int {
	return t.count
}

// Search searches for a node with the given key.
// A result of 0 means that the key was not found.
func (t *AVLTree[K]) Search(key []K) uint64 {
	result := t.root.search(key)

	if result == nil {
		return 0
	}

	return result.Count
}

// Clear clears the tree.
func (t *AVLTree[K]) Clear() {
	t.root = nil
	t.count = 0
}

// All returns an iterator over all keys in the tree and their counts,
// in ascending key order. It allocates nothing.
//
// The yielded key is the tree's own storage and must not be modified,
// as that would break the sort order of the tree.
func (t *AVLTree[K]) All() iter.Seq2[[]K, uint64] {
	return t.root.inOrder()
}

// Dump prints the tree to the console.
func (t *AVLTree[K]) Dump() {
	t.root.print("", false)
}

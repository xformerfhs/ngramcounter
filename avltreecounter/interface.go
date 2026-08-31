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
// Version: 3.0.0
//
// Change history:
//    2025-08-31: V1.0.0: Created.
//    2026-08-31: V1.1.0: Use the inOrder iterator to build the results.
//    2026-08-31: V2.0.0: Created "All" iterator.
//    2026-08-31: V3.0.0: Changed "Search" to "Get", "Count" to "NodeCount" and
//                        added "TotalCount".
//

// Package avltreecounter provides a self-balancing binary counter tree with slice keys.
package avltreecounter

import (
	"cmp"
	"iter"
)

// ******** Public types ********

// AVLTree is a self-balancing binary search tree that stores
// a count for each unique key.
type AVLTree[K cmp.Ordered] struct {
	root       *avlNode[K]
	nodeCount  uint64
	totalCount uint64
}

// ******** Public functions ********

// Add adds key to the tree. If key already exists, its count is incremented.
func (t *AVLTree[K]) Add(key []K) {
	var madeNewNode bool
	t.root, madeNewNode = t.root.add(key)
	t.totalCount++

	if madeNewNode {
		t.nodeCount++
	}
}

// NodeCount returns the number of nodes in the tree.
func (t *AVLTree[K]) NodeCount() uint64 {
	return t.nodeCount
}

// TotalCount returns the number of keys added to the tree.
func (t *AVLTree[K]) TotalCount() uint64 {
	return t.totalCount
}

// Get searches for a node with the given key.
func (t *AVLTree[K]) Get(key []K) (uint64, bool) {
	result := t.root.get(key)

	if result == nil {
		return 0, false
	}

	return result.Count, true
}

// Clear clears the tree.
func (t *AVLTree[K]) Clear() {
	t.root = nil
	t.nodeCount = 0
	t.totalCount = 0
}

// All returns an iterator over all keys in the tree and their counts,
// in ascending key order.
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

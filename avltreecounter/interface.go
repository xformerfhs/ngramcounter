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
//    2025-08-31: V1.0.0: Created.
//

// Package avltreecounter provides a self-balancing binary counter tree with slice keys.
package avltreecounter

import (
	"cmp"
)

// ******** Public types ********

// AVLTree is a self-balancing binary counter tree with slice keys.
type AVLTree[K cmp.Ordered] struct {
	root          *avlNode[K]
	count         int
	lastFoundNode *avlNode[K]
}

// CountEntry is the structure for a count entry.
type CountEntry[K cmp.Ordered] struct {
	Key   []K
	Count uint64
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
func (t *AVLTree[K]) Search(key []K) (uint64, bool) {
	result := t.root.search(key)
	t.lastFoundNode = result

	if result == nil {
		return 0, false
	}

	return result.Count, true
}

// Keys returns all keys in the tree in sorted order.
func (t *AVLTree[K]) Keys() [][]K {
	allNodes := make([]*avlNode[K], 0)
	allNodes = t.root.collectNodes(allNodes)
	result := make([][]K, len(allNodes))

	for i, node := range allNodes {
		result[i] = node.Key
	}

	return result
}

// CountEntries returns all count entries in the tree
// in sorted order by key.
func (t *AVLTree[K]) CountEntries() []CountEntry[K] {
	allNodes := make([]*avlNode[K], 0, t.count)
	allNodes = t.root.collectNodes(allNodes)
	result := make([]CountEntry[K], len(allNodes))

	for i, node := range allNodes {
		result[i] = CountEntry[K]{Key: node.Key, Count: node.Count}
	}

	return result
}

// Dump prints the tree to the console.
func (t *AVLTree[K]) Dump() {
	t.root.print("", false)
}

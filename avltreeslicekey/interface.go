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

// Package avltreeslicekey provides a self-balancing binary search tree with slice keys.
package avltreeslicekey

import (
	"cmp"
)

// ******** Public types ********

// AVLTree is a self-balancing binary search tree with slice keys.
type AVLTree[K cmp.Ordered, V any] struct {
	root          *avlNode[K, V]
	count         int
	lastFoundNode *avlNode[K, V]
}

type KeyValuePair[K cmp.Ordered, V any] struct {
	Key   []K
	Value V
}

// ******** Public functions ********

// Insert inserts a new node into the tree.
func (t *AVLTree[K, V]) Insert(key []K, value V) {
	newNode := newAVLNode(key, value)

	if t.root == nil {
		t.root = newNode
	} else {
		t.root = t.root.insert(newNode)
	}

	t.count++
}

// Count returns the number of nodes in the tree.
func (t *AVLTree[K, V]) Count() int {
	return t.count
}

// Search searches for a node with the given key.
func (t *AVLTree[K, V]) Search(key []K) (V, bool) {
	result := t.root.search(key)
	t.lastFoundNode = result

	if result == nil {
		var zeroValue V
		return zeroValue, false
	}

	return result.Value, true
}

// Set sets the node with the given key to the given value.
func (t *AVLTree[K, V]) Set(key []K, newValue V) bool {
	foundNode := t.root.search(key)

	if foundNode == nil {
		return false
	}

	foundNode.Value = newValue

	return true
}

// SetLastFound sets the last found node to the given value.
func (t *AVLTree[K, V]) SetLastFound(newValue V) bool {
	if t.lastFoundNode != nil {
		t.lastFoundNode.Value = newValue
		return true
	}

	return false
}

// Keys returns all keys in the tree in sorted order.
func (t *AVLTree[K, V]) Keys() [][]K {
	allNodes := make([]*avlNode[K, V], 0)
	allNodes = t.root.collectNodes(allNodes)
	result := make([][]K, len(allNodes))

	for i, node := range allNodes {
		result[i] = node.Key
	}

	return result
}

// KeyValuePairs returns all key-value pairs in the tree
// in sorted order by key.
func (t *AVLTree[K, V]) KeyValuePairs() []KeyValuePair[K, V] {
	allNodes := make([]*avlNode[K, V], 0, t.count)
	allNodes = t.root.collectNodes(allNodes)
	result := make([]KeyValuePair[K, V], len(allNodes))

	for i, node := range allNodes {
		result[i] = KeyValuePair[K, V]{Key: node.Key, Value: node.Value}
	}

	return result
}

// Dump prints the tree to the console.
func (t *AVLTree[K, V]) Dump() {
	t.root.print("", false)
}

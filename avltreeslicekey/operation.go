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

import (
	"cmp"
	"slices"
)

// ******** Private types ********

// avlNode is a node in the AVL tree.
type avlNode[K cmp.Ordered, V any] struct {
	Key    []K
	Value  V
	left   *avlNode[K, V]
	right  *avlNode[K, V]
	height int
}

// ******** Private functions ********

// newAVLNode creates a new AVL node.
func newAVLNode[K cmp.Ordered, V any](key []K, value V) *avlNode[K, V] {
	return &avlNode[K, V]{
		Key:    slices.Clone(key),
		Value:  value,
		left:   nil,
		right:  nil,
		height: 0,
	}
}

// insert inserts the node into the tree.
func (n *avlNode[K, V]) insert(newNode *avlNode[K, V]) *avlNode[K, V] {
	if n == nil {
		return newNode
	}

	comparison := Compare(newNode.Key, n.Key)
	if comparison < 0 {
		n.left = n.left.insert(newNode)
	} else if comparison > 0 {
		n.right = n.right.insert(newNode)
	} else {
		n.Value = newNode.Value
	}

	n.updateHeight()

	return n.rebalance()
}

// search searches for the node with the given key.
func (n *avlNode[K, V]) search(key []K) *avlNode[K, V] {
	for current := n; current != nil; {
		comparison := Compare(key, current.Key)

		if comparison < 0 {
			current = current.left
		} else if comparison > 0 {
			current = current.right
		} else {
			return current
		}
	}

	// Not found.
	return nil
}

// collectNodes collects all the nodes in the tree in sorted key order.
func (n *avlNode[K, V]) collectNodes(nodeCollector []*avlNode[K, V]) []*avlNode[K, V] {
	// If the node is nil, return the node collector.
	if n == nil {
		return nodeCollector
	}

	// Collect the left node, i.e., keys that have lesser values.
	if n.left != nil {
		nodeCollector = n.left.collectNodes(nodeCollector)
	}

	// Add the current node to the node collector.
	nodeCollector = append(nodeCollector, n)

	// Collect the right node, i.e., keys that have greater values.
	if n.right != nil {
		nodeCollector = n.right.collectNodes(nodeCollector)
	}

	return nodeCollector
}

// ******** Helper functions ********

// balanceFactor calculates the balance factor of the node.
func (n *avlNode[K, V]) balanceFactor() int {
	left, right := -1, -1

	if n.left != nil {
		left = n.left.height
	}

	if n.right != nil {
		right = n.right.height
	}

	return left - right
}

// rebalance rebalances the tree.
func (n *avlNode[K, V]) rebalance() *avlNode[K, V] {
	newRoot := n
	bf := n.balanceFactor()

	if bf > 1 {
		if n.left != nil &&
			n.left.balanceFactor() < 0 {
			n.left = n.left.leftRotation()
		}

		newRoot = n.rightRotation()
	} else if bf < -1 {
		if n.right != nil &&
			n.right.balanceFactor() > 0 {
			n.right = n.right.rightRotation()
		}

		newRoot = n.leftRotation()
	}

	return newRoot
}

// rightRotation rotates the node to the right.
func (n *avlNode[K, V]) rightRotation() *avlNode[K, V] {
	nl := n.left
	n.left = nl.right
	nl.right = n

	n.updateHeight()
	nl.updateHeight()

	return nl
}

// leftRotation rotates the node to the left.
func (n *avlNode[K, V]) leftRotation() *avlNode[K, V] {
	nr := n.right
	n.right = nr.left
	nr.left = n

	n.updateHeight()
	nr.updateHeight()

	return nr
}

// updateHeight updates the height of the node.
func (n *avlNode[K, V]) updateHeight() {
	left, right := -1, -1

	if n.left != nil {
		left = n.left.height
	}

	if n.right != nil {
		right = n.right.height
	}

	if left > right {
		n.height = left + 1
	} else {
		n.height = right + 1
	}
}

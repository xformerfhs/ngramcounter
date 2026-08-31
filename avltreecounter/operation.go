//
// SPDX-FileCopyrightText: Copyright 2025-2026 Frank Schwab
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
// Version: 1.2.0
//
// Change history:
//    2025-08-31: V1.0.0: Created.
//    2026-08-30: V1.0.1: Better variable naming.
//    2026-08-31: V1.1.0: Replaced collectNodes by the inOrder iterator.
//    2026-08-31: V1.2.0: Simplified height calculations.
//

package avltreecounter

import (
	"cmp"
	"iter"
	"slices"
)

// ******** Private types ********

// avlNode is a node in the AVL count tree.
type avlNode[K cmp.Ordered] struct {
	Key    []K
	Count  uint64
	left   *avlNode[K]
	right  *avlNode[K]
	height int
}

// ******** Private functions ********

// newAVLNode creates a new AVL node.
func newAVLNode[K cmp.Ordered](key []K) *avlNode[K] {
	return &avlNode[K]{
		Key:    slices.Clone(key),
		Count:  1,
		left:   nil,
		right:  nil,
		height: 0,
	}
}

// add adds the key to the tree.
func (n *avlNode[K]) add(key []K) (*avlNode[K], bool) {
	if n == nil {
		return newAVLNode(key), true
	}

	var madeNewNode bool
	comparison := slices.Compare(key, n.Key)
	if comparison < 0 {
		n.left, madeNewNode = n.left.add(key)
	} else if comparison > 0 {
		n.right, madeNewNode = n.right.add(key)
	} else {
		n.Count++
	}

	n.updateHeight()

	return n.rebalance(), madeNewNode
}

// get searches for the node with the given key.
func (n *avlNode[K]) get(key []K) *avlNode[K] {
	for current := n; current != nil; {
		comparison := slices.Compare(key, current.Key)

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

// inOrder returns an iterator over all the keys and counts in the tree
// in ascending key order.
func (n *avlNode[K]) inOrder() iter.Seq2[[]K, uint64] {
	return func(yield func([]K, uint64) bool) {
		n.yieldInOrder(yield)
	}
}

// ******** Helper functions ********

// yieldInOrder yields the keys and counts of this subtree in ascending key order.
// It returns false as soon as the consumer stops the iteration.
//
// The recursion depth is O(log(n)) because the tree is AVL-balanced.
func (n *avlNode[K]) yieldInOrder(yield func([]K, uint64) bool) bool {
	if n == nil {
		return true
	}

	// Yield the lesser keys, then this key, then the greater keys.
	return n.left.yieldInOrder(yield) &&
		yield(n.Key, n.Count) &&
		n.right.yieldInOrder(yield)
}

// balanceFactor calculates the balance factor of the node.
func (n *avlNode[K]) balanceFactor() int {
	return height(n.left) - height(n.right)
}

// updateHeight updates the height of the node.
func (n *avlNode[K]) updateHeight() {
	n.height = max(height(n.left), height(n.right)) + 1
}

// height returns the height of the node.
func height[K cmp.Ordered](n *avlNode[K]) int {
	if n == nil {
		return -1
	}

	return n.height
}

// rebalance rebalances the tree.
func (n *avlNode[K]) rebalance() *avlNode[K] {
	bf := n.balanceFactor()

	if bf > 1 {
		if n.left.balanceFactor() < 0 {
			n.left = n.left.leftRotation()
		}

		return n.rightRotation()
	} else if bf < -1 {
		if n.right.balanceFactor() > 0 {
			n.right = n.right.rightRotation()
		}

		return n.leftRotation()
	}

	return n
}

// rightRotation rotates the node to the right.
func (n *avlNode[K]) rightRotation() *avlNode[K] {
	leftNode := n.left
	n.left = leftNode.right
	leftNode.right = n

	n.updateHeight()
	leftNode.updateHeight()

	return leftNode
}

// leftRotation rotates the node to the left.
func (n *avlNode[K]) leftRotation() *avlNode[K] {
	rightNode := n.right
	n.right = rightNode.left
	rightNode.left = n

	n.updateHeight()
	rightNode.updateHeight()

	return rightNode
}

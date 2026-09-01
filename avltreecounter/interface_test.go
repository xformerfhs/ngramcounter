//
// SPDX-FileCopyrightText: Copyright 2026 Frank Schwab
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
//    2026-09-01: V1.0.0: Created.
//

package avltreecounter

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// ******** Private constants ********

// keyCount is the number of keys used in the tests that fill a large tree.
const keyCount = 2_000

// ******** Test functions ********

// -------- Empty tree --------

// TestEmptyTree tests that all functions work on a tree that has never been used.
func TestEmptyTree(t *testing.T) {
	tree := &AVLTree[int]{}

	if tree.NodeCount() != 0 {
		t.Errorf(`NodeCount of empty tree is %d, want 0`, tree.NodeCount())
	}

	if tree.TotalCount() != 0 {
		t.Errorf(`TotalCount of empty tree is %d, want 0`, tree.TotalCount())
	}

	if count := tree.Get([]int{1}); count != 0 {
		t.Errorf(`Get on empty tree returned %d, want 0`, count)
	}

	keys, _ := allEntries(tree)
	if len(keys) != 0 {
		t.Errorf(`All on empty tree yielded %d keys, want 0`, len(keys))
	}
}

// -------- Add and Get --------

// TestAddOneKey tests adding a single key.
func TestAddOneKey(t *testing.T) {
	tree := &AVLTree[int]{}
	tree.Add([]int{7, 8, 9})

	if tree.NodeCount() != 1 {
		t.Errorf(`NodeCount is %d, want 1`, tree.NodeCount())
	}

	if tree.TotalCount() != 1 {
		t.Errorf(`TotalCount is %d, want 1`, tree.TotalCount())
	}

	if count := tree.Get([]int{7, 8, 9}); count != 1 {
		t.Errorf(`Get of added key returned %d, want 1`, count)
	}

	checkTree(t, tree)
}

// TestAddDuplicateKey tests that adding an existing key only increments its count.
func TestAddDuplicateKey(t *testing.T) {
	const addCount = 5

	tree := &AVLTree[int]{}
	key := []int{1, 2}

	for range addCount {
		tree.Add(key)
	}

	if tree.NodeCount() != 1 {
		t.Errorf(`NodeCount is %d, want 1`, tree.NodeCount())
	}

	if tree.TotalCount() != addCount {
		t.Errorf(`TotalCount is %d, want %d`, tree.TotalCount(), addCount)
	}

	if count := tree.Get(key); count != addCount {
		t.Errorf(`Get returned %d, want %d`, count, addCount)
	}

	checkTree(t, tree)
}

// TestAddDistinctKeys tests adding keys with different counts.
func TestAddDistinctKeys(t *testing.T) {
	wantCounts := map[string]uint64{
		`[3]`:       1,
		`[1]`:       2,
		`[2]`:       3,
		`[1 1]`:     4,
		`[2 7 1 8]`: 5,
	}

	tree := &AVLTree[int]{}
	wantTotalCount := uint64(0)

	for _, key := range sortedKeys(wantCounts) {
		for range wantCounts[key] {
			tree.Add(parseKey(t, key))
			wantTotalCount++
		}
	}

	if tree.NodeCount() != uint64(len(wantCounts)) {
		t.Errorf(`NodeCount is %d, want %d`, tree.NodeCount(), len(wantCounts))
	}

	if tree.TotalCount() != wantTotalCount {
		t.Errorf(`TotalCount is %d, want %d`, tree.TotalCount(), wantTotalCount)
	}

	for key, wantCount := range wantCounts {
		if count := tree.Get(parseKey(t, key)); count != wantCount {
			t.Errorf(`Get of key %s returned %d, want %d`, key, count, wantCount)
		}
	}

	checkTree(t, tree)
}

// TestGetMissingKeys tests that Get returns 0 for keys that are not in the tree.
func TestGetMissingKeys(t *testing.T) {
	tree := &AVLTree[int]{}
	for _, key := range [][]int{{2}, {4}, {6}} {
		tree.Add(key)
	}

	missingKeys := [][]int{{1}, {3}, {5}, {7}, {2, 0}, {}}
	for _, key := range missingKeys {
		if count := tree.Get(key); count != 0 {
			t.Errorf(`Get of missing key %v returned %d, want 0`, key, count)
		}
	}
}

// TestPrefixKeys tests that a key and its prefixes are stored as separate keys.
func TestPrefixKeys(t *testing.T) {
	tree := &AVLTree[int]{}

	keys := [][]int{{1, 2, 3}, {1, 2}, {1}}
	for _, key := range keys {
		tree.Add(key)
	}

	if tree.NodeCount() != uint64(len(keys)) {
		t.Errorf(`NodeCount is %d, want %d`, tree.NodeCount(), len(keys))
	}

	for _, key := range keys {
		if count := tree.Get(key); count != 1 {
			t.Errorf(`Get of key %v returned %d, want 1`, key, count)
		}
	}

	// A prefix must sort before the longer key.
	gotKeys, _ := allEntries(tree)
	wantKeys := [][]int{{1}, {1, 2}, {1, 2, 3}}
	if !slices.EqualFunc(gotKeys, wantKeys, slices.Equal) {
		t.Errorf(`All yielded %v, want %v`, gotKeys, wantKeys)
	}
}

// TestEmptyKey tests that an empty key is a valid key and that a nil key
// and an empty key are the same key.
func TestEmptyKey(t *testing.T) {
	tree := &AVLTree[int]{}

	tree.Add(nil)
	tree.Add([]int{})

	if tree.NodeCount() != 1 {
		t.Errorf(`NodeCount is %d, want 1`, tree.NodeCount())
	}

	if count := tree.Get(nil); count != 2 {
		t.Errorf(`Get of nil key returned %d, want 2`, count)
	}

	if count := tree.Get([]int{}); count != 2 {
		t.Errorf(`Get of empty key returned %d, want 2`, count)
	}

	// The empty key must sort before all other keys.
	tree.Add([]int{0})

	gotKeys, _ := allEntries(tree)
	if len(gotKeys[0]) != 0 {
		t.Errorf(`First key yielded by All is %v, want the empty key`, gotKeys[0])
	}
}

// TestAddClonesKey tests that Add stores a copy of the key so that later
// modifications of the caller's slice do not change the tree.
func TestAddClonesKey(t *testing.T) {
	tree := &AVLTree[int]{}

	key := []int{1, 2, 3}
	tree.Add(key)
	key[0] = 9

	if count := tree.Get([]int{1, 2, 3}); count != 1 {
		t.Errorf(`Get of the original key returned %d, want 1`, count)
	}

	if count := tree.Get([]int{9, 2, 3}); count != 0 {
		t.Errorf(`Get of the modified key returned %d, want 0`, count)
	}
}

// TestStringElementKeys tests that the tree works with an element type
// other than an integer.
func TestStringElementKeys(t *testing.T) {
	tree := &AVLTree[string]{}

	for _, key := range [][]string{{`b`, `c`}, {`a`}, {`b`}, {`a`}} {
		tree.Add(key)
	}

	if tree.NodeCount() != 3 {
		t.Errorf(`NodeCount is %d, want 3`, tree.NodeCount())
	}

	if count := tree.Get([]string{`a`}); count != 2 {
		t.Errorf(`Get of key [a] returned %d, want 2`, count)
	}

	gotKeys, gotCounts := allEntries(tree)
	wantKeys := [][]string{{`a`}, {`b`}, {`b`, `c`}}
	wantCounts := []uint64{2, 1, 1}

	if !slices.EqualFunc(gotKeys, wantKeys, slices.Equal) {
		t.Errorf(`All yielded keys %v, want %v`, gotKeys, wantKeys)
	}

	if !slices.Equal(gotCounts, wantCounts) {
		t.Errorf(`All yielded counts %v, want %v`, gotCounts, wantCounts)
	}

	checkTree(t, tree)
}

// -------- Balancing --------

// TestTreeStaysBalanced tests that the tree stays balanced and correct,
// whatever the insertion order is.
func TestTreeStaysBalanced(t *testing.T) {
	insertionOrders := map[string]func() []int{
		`ascending`:  func() []int { return ascendingNumbers(keyCount) },
		`descending`: func() []int { return descendingNumbers(keyCount) },
		`shuffled`:   func() []int { return shuffledNumbers(keyCount) },
	}

	for name, makeNumbers := range insertionOrders {
		t.Run(name, func(t *testing.T) {
			tree := &AVLTree[int]{}
			for _, number := range makeNumbers() {
				tree.Add([]int{number})
			}

			if tree.NodeCount() != keyCount {
				t.Errorf(`NodeCount is %d, want %d`, tree.NodeCount(), keyCount)
			}

			if tree.TotalCount() != keyCount {
				t.Errorf(`TotalCount is %d, want %d`, tree.TotalCount(), keyCount)
			}

			nodesFound := checkTree(t, tree)
			if nodesFound != keyCount {
				t.Errorf(`Tree contains %d nodes, want %d`, nodesFound, keyCount)
			}

			checkHeightBound(t, tree)

			// All keys must be present exactly once and in ascending order.
			gotKeys, gotCounts := allEntries(tree)
			if len(gotKeys) != keyCount {
				t.Fatalf(`All yielded %d keys, want %d`, len(gotKeys), keyCount)
			}

			for i, key := range gotKeys {
				if len(key) != 1 || key[0] != i {
					t.Fatalf(`Key #%d is %v, want [%d]`, i, key, i)
				}

				if gotCounts[i] != 1 {
					t.Fatalf(`Count of key %v is %d, want 1`, key, gotCounts[i])
				}
			}
		})
	}
}

// TestRandomKeysAndCounts tests the tree with randomly generated keys
// against a map that holds the expected counts.
func TestRandomKeysAndCounts(t *testing.T) {
	const addCount = 20_000

	rng := rand.New(rand.NewPCG(0xcafe_babe, 0xdead_beef))

	tree := &AVLTree[int]{}
	wantCounts := make(map[string]uint64)

	for range addCount {
		key := randomKey(rng)
		tree.Add(key)
		wantCounts[fmt.Sprint(key)]++
	}

	if tree.NodeCount() != uint64(len(wantCounts)) {
		t.Errorf(`NodeCount is %d, want %d`, tree.NodeCount(), len(wantCounts))
	}

	if tree.TotalCount() != addCount {
		t.Errorf(`TotalCount is %d, want %d`, tree.TotalCount(), addCount)
	}

	checkTree(t, tree)
	checkHeightBound(t, tree)

	for key, wantCount := range wantCounts {
		if count := tree.Get(parseKey(t, key)); count != wantCount {
			t.Errorf(`Get of key %s returned %d, want %d`, key, count, wantCount)
		}
	}

	// The iterator must yield exactly the same keys and counts.
	gotKeys, gotCounts := allEntries(tree)
	if len(gotKeys) != len(wantCounts) {
		t.Fatalf(`All yielded %d keys, want %d`, len(gotKeys), len(wantCounts))
	}

	sumOfCounts := uint64(0)
	for i, key := range gotKeys {
		wantCount, found := wantCounts[fmt.Sprint(key)]
		if !found {
			t.Errorf(`All yielded unknown key %v`, key)
			continue
		}

		if gotCounts[i] != wantCount {
			t.Errorf(`All yielded count %d for key %v, want %d`, gotCounts[i], key, wantCount)
		}

		sumOfCounts += gotCounts[i]
	}

	if sumOfCounts != tree.TotalCount() {
		t.Errorf(`Sum of the yielded counts is %d, want %d`, sumOfCounts, tree.TotalCount())
	}
}

// -------- Clear --------

// TestClear tests that Clear empties the tree and that the tree is usable afterwards.
func TestClear(t *testing.T) {
	tree := &AVLTree[int]{}
	for _, number := range shuffledNumbers(100) {
		tree.Add([]int{number})
	}

	tree.Clear()

	if tree.NodeCount() != 0 {
		t.Errorf(`NodeCount after Clear is %d, want 0`, tree.NodeCount())
	}

	if tree.TotalCount() != 0 {
		t.Errorf(`TotalCount after Clear is %d, want 0`, tree.TotalCount())
	}

	if tree.root != nil {
		t.Error(`Root after Clear is not nil`)
	}

	if count := tree.Get([]int{1}); count != 0 {
		t.Errorf(`Get after Clear returned %d, want 0`, count)
	}

	keys, _ := allEntries(tree)
	if len(keys) != 0 {
		t.Errorf(`All after Clear yielded %d keys, want 0`, len(keys))
	}

	// The tree must be reusable.
	tree.Add([]int{42})

	if tree.NodeCount() != 1 || tree.TotalCount() != 1 {
		t.Errorf(`NodeCount/TotalCount after reuse are %d/%d, want 1/1`, tree.NodeCount(), tree.TotalCount())
	}

	if count := tree.Get([]int{42}); count != 1 {
		t.Errorf(`Get after reuse returned %d, want 1`, count)
	}
}

// TestClearEmptyTree tests that Clear works on an empty tree.
func TestClearEmptyTree(t *testing.T) {
	tree := &AVLTree[int]{}
	tree.Clear()

	if tree.NodeCount() != 0 || tree.TotalCount() != 0 {
		t.Errorf(`NodeCount/TotalCount are %d/%d, want 0/0`, tree.NodeCount(), tree.TotalCount())
	}
}

// -------- All --------

// TestAllStopsEarly tests that the iteration stops as soon as the consumer stops it.
func TestAllStopsEarly(t *testing.T) {
	const stopAfter = 3

	tree := &AVLTree[int]{}
	for _, number := range shuffledNumbers(100) {
		tree.Add([]int{number})
	}

	yieldedKeys := make([]int, 0, stopAfter)
	for key := range tree.All() {
		yieldedKeys = append(yieldedKeys, key[0])

		if len(yieldedKeys) == stopAfter {
			break
		}
	}

	wantKeys := ascendingNumbers(stopAfter)
	if !slices.Equal(yieldedKeys, wantKeys) {
		t.Errorf(`Iteration yielded %v, want %v`, yieldedKeys, wantKeys)
	}
}

// TestAllIsRepeatable tests that the iterator can be used more than once.
func TestAllIsRepeatable(t *testing.T) {
	tree := &AVLTree[int]{}
	for _, number := range shuffledNumbers(10) {
		tree.Add([]int{number})
	}

	all := tree.All()

	firstKeys, firstCounts := entriesOfSeq(all)
	secondKeys, secondCounts := entriesOfSeq(all)

	if !slices.EqualFunc(firstKeys, secondKeys, slices.Equal) {
		t.Errorf(`Second iteration yielded keys %v, want %v`, secondKeys, firstKeys)
	}

	if !slices.Equal(firstCounts, secondCounts) {
		t.Errorf(`Second iteration yielded counts %v, want %v`, secondCounts, firstCounts)
	}
}

// ******** Benchmark functions ********

// BenchmarkAddNewKeys benchmarks adding keys that are not yet in the tree.
func BenchmarkAddNewKeys(b *testing.B) {
	tree := &AVLTree[int]{}

	for i := range b.N {
		tree.Add([]int{i})
	}
}

// BenchmarkAddExistingKeys benchmarks adding keys that are already in the tree.
func BenchmarkAddExistingKeys(b *testing.B) {
	tree := &AVLTree[int]{}
	for _, number := range shuffledNumbers(keyCount) {
		tree.Add([]int{number})
	}

	b.ResetTimer()

	for i := range b.N {
		tree.Add([]int{i % keyCount})
	}
}

// BenchmarkGet benchmarks searching for keys.
func BenchmarkGet(b *testing.B) {
	tree := &AVLTree[int]{}
	for _, number := range shuffledNumbers(keyCount) {
		tree.Add([]int{number})
	}

	b.ResetTimer()

	for i := range b.N {
		tree.Get([]int{i % keyCount})
	}
}

// ******** Helper functions ********

// -------- Tree checks --------

// checkTree checks the AVL invariants of the tree and returns the number of nodes found.
func checkTree[K cmp.Ordered](t *testing.T, tree *AVLTree[K]) uint64 {
	t.Helper()

	nodesFound, _ := checkSubTree(t, tree.root)

	return nodesFound
}

// checkSubTree checks the AVL invariants of the subtree that starts at n and
// returns the number of nodes found and the height of the subtree.
func checkSubTree[K cmp.Ordered](t *testing.T, n *avlNode[K]) (uint64, int8) {
	t.Helper()

	if n == nil {
		return 0, -1
	}

	if n.Count == 0 {
		t.Errorf(`Node with key %v has count 0`, n.Key)
	}

	if n.left != nil && slices.Compare(n.left.Key, n.Key) >= 0 {
		t.Errorf(`Left child key %v is not less than the parent key %v`, n.left.Key, n.Key)
	}

	if n.right != nil && slices.Compare(n.right.Key, n.Key) <= 0 {
		t.Errorf(`Right child key %v is not greater than the parent key %v`, n.right.Key, n.Key)
	}

	leftNodes, leftHeight := checkSubTree(t, n.left)
	rightNodes, rightHeight := checkSubTree(t, n.right)

	wantHeight := max(leftHeight, rightHeight) + 1
	if n.height != wantHeight {
		t.Errorf(`Node with key %v has stored height %d, want %d`, n.Key, n.height, wantHeight)
	}

	if balanceFactor := leftHeight - rightHeight; balanceFactor < -1 || balanceFactor > 1 {
		t.Errorf(`Node with key %v has balance factor %d, want -1, 0 or 1`, n.Key, balanceFactor)
	}

	return leftNodes + rightNodes + 1, wantHeight
}

// checkHeightBound checks that the height of the tree does not exceed
// the theoretical maximum height of an AVL tree with this number of nodes.
func checkHeightBound(t *testing.T, tree *AVLTree[int]) {
	t.Helper()

	nodeCount := tree.NodeCount()
	if nodeCount == 0 {
		return
	}

	maxHeight := int8(math.Floor(1.4405*math.Log2(float64(nodeCount)+2.0) - 1.3277))
	if actualHeight := height(tree.root); actualHeight > maxHeight {
		t.Errorf(`Height of the tree with %d nodes is %d, want at most %d`, nodeCount, actualHeight, maxHeight)
	}
}

// -------- Iteration --------

// allEntries collects all keys and counts of the tree in iteration order.
func allEntries[K cmp.Ordered](tree *AVLTree[K]) ([][]K, []uint64) {
	return entriesOfSeq(tree.All())
}

// entriesOfSeq collects all keys and counts that the sequence yields.
func entriesOfSeq[K cmp.Ordered](seq func(func([]K, uint64) bool)) ([][]K, []uint64) {
	keys := make([][]K, 0, 16)
	counts := make([]uint64, 0, 16)

	for key, count := range seq {
		keys = append(keys, slices.Clone(key))
		counts = append(counts, count)
	}

	return keys, counts
}

// -------- Key generation --------

// ascendingNumbers returns the numbers from 0 to count-1 in ascending order.
func ascendingNumbers(count int) []int {
	result := make([]int, count)
	for i := range result {
		result[i] = i
	}

	return result
}

// descendingNumbers returns the numbers from 0 to count-1 in descending order.
func descendingNumbers(count int) []int {
	result := ascendingNumbers(count)
	slices.Reverse(result)

	return result
}

// shuffledNumbers returns the numbers from 0 to count-1 in a reproducible random order.
func shuffledNumbers(count int) []int {
	result := ascendingNumbers(count)

	rng := rand.New(rand.NewPCG(0x1234_5678, 0x9abc_def0))
	rng.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// randomKey returns a random key with 1 to 3 elements in the range from 0 to 9.
func randomKey(rng *rand.Rand) []int {
	result := make([]int, rng.IntN(3)+1)
	for i := range result {
		result[i] = rng.IntN(10)
	}

	return result
}

// parseKey converts the fmt representation of an integer key back into a slice.
func parseKey(t *testing.T, key string) []int {
	t.Helper()

	fields := strings.Fields(strings.Trim(key, `[]`))

	result := make([]int, len(fields))
	for i, field := range fields {
		if _, err := fmt.Sscan(field, &result[i]); err != nil {
			t.Fatalf(`Could not parse key %q: %v`, key, err)
		}
	}

	return result
}

// sortedKeys returns the keys of the map in ascending order,
// so that the tests are reproducible.
func sortedKeys(m map[string]uint64) []string {
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}

	slices.Sort(result)

	return result
}

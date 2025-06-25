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
// Version: 2.3.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-08: V2.0.0: Use different modes, with "overlapped" being the default.
//    2025-01-11: V2.0.1: Correct error message for incomplete n-grams in sequential mode.
//    2025-01-11: V2.1.0: Simplify preparing next collector state and make it faster.
//    2025-01-12: V2.2.0: Simplify preparing next collector state.
//    2025-06-25: V2.3.0: Simplify "CountNGrams" function.
//

package counters

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"ngramcounter/filehelper"
	"os"
	"unicode"
)

// ******** Public types *********

// NgramCounter contains the encoding data for an NgramCounter.
type NgramCounter struct {
	decoder               *encoding.Decoder
	onlyLettersAndNumbers bool
}

// ******** Public functions ********

// NewNgramCounter returns a new NGramCounter for the given encoding text.
func NewNgramCounter(enc encoding.Encoding, allChars bool) *NgramCounter {
	return &NgramCounter{decoder: enc.NewDecoder(), onlyLettersAndNumbers: !allChars}
}

// CountNGrams counts the n-grams in the file.
func (nc *NgramCounter) CountNGrams(fileName string, ngramSize uint, useSequential bool) (map[string]uint64, uint64, error) {
	br, file, err := nc.openAndWrapFile(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(file)

	result := make(map[string]uint64)
	collector := make([]rune, ngramSize)
	collectorIndex := uint(0)
	ngramCounter := uint64(0)

	for {
		var r rune
		r, _, err = br.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, 0, err
		}

		if nc.shouldSkipRune(r) {
			continue
		}

		collector[collectorIndex] = r
		collectorIndex++

		if collectorIndex == ngramSize {
			ngramCounter++
			ngram := string(collector)
			result[ngram]++
			collectorIndex = prepareCollector(collector, collectorIndex, ngramSize, useSequential)
		}
	}

	if useSequential && collectorIndex != 0 {
		return nil, 0, fmt.Errorf(`File ends with a %d-gram`, collectorIndex)
	}

	return result, ngramCounter, nil
}

// ******** Private functions ********

// openAndWrapFile opens the input file and wraps it in a buffered transforming reader.
func (nc *NgramCounter) openAndWrapFile(fileName string) (*bufio.Reader, *os.File, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	tr := transform.NewReader(f, nc.decoder)
	br := bufio.NewReader(tr)

	return br, f, nil
}

// shouldSkipRune reports whether the supplied rune should be skipped.
func (nc *NgramCounter) shouldSkipRune(r rune) bool {
	if unicode.IsControl(r) {
		return true
	}

	if nc.onlyLettersAndNumbers && !unicode.In(r, unicode.L, unicode.N) {
		return true
	}

	return false
}

// prepareCollector prepares the collector for the next rune.
func prepareCollector(collector []rune, collectorIndex uint, ngramSize uint, useSequential bool) uint {
	if useSequential {
		// Sequential mode reuses the collector from the start.
		return 0
	} else {
		// Overlapped mode copies all elements of the collector one place to the left.
		if ngramSize >= 8 {
			// If there are 8 or more elements in the collector, use the copy function.
			copy(collector, collector[1:ngramSize])
		} else {
			// If there are less than 8 elements, a loop is faster.
			ngramSize--
			_ = collector[ngramSize] // Check index of upper limit only once.
			for i, j := uint(0), uint(1); j <= ngramSize; j++ {
				collector[i] = collector[j]
				i = j
			}
		}

		return collectorIndex - 1
	}
}

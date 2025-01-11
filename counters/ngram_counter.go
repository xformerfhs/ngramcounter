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
// Version: 2.1.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-08: V2.0.0: Use different modes, with "overlapped" being the default.
//    2025-01-11: V2.0.1: Correct error message for incomplete n-grams in sequential mode.
//    2025-01-11: V2.1.0: Simplify preparing next collector state and make it faster.
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
	decoder *encoding.Decoder
}

// ******** Public functions ********

// NewNgramCounter returns a new NGramCounter for the given encoding text.
func NewNgramCounter(enc encoding.Encoding) *NgramCounter {
	return &NgramCounter{decoder: enc.NewDecoder()}
}

// CountNGrams counts the n-grams in the file.
func (nc *NgramCounter) CountNGrams(fileName string, ngramSize uint, useSequential bool) (map[string]uint64, uint64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(f)

	// 1. Wrap the file in a transform that decodes the character encoding to UTF-8.
	tr := transform.NewReader(f, nc.decoder)

	// 2. The buffered reader uses the transform reader.
	br := bufio.NewReader(tr)

	// 3. Now loop over the lines.
	result := make(map[string]uint64)
	collector := make([]rune, ngramSize)
	collectorIndex := uint(0)
	ngramCounter := uint64(0)
	for {
		// Read rune and bail out, if there is an error or EOF.
		var r rune
		r, _, err = br.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, 0, err
			}
		}

		// Only look at letters and numbers.
		if !unicode.In(r, unicode.L, unicode.N) {
			continue
		}

		// Put rune in collector.
		collector[collectorIndex] = r
		collectorIndex++

		// Count n-gram, if collector is full.
		if collectorIndex == ngramSize {
			ngram := string(collector)
			result[ngram]++
			ngramCounter++

			collectorIndex = prepareCollector(collector, collectorIndex, ngramSize, useSequential)
		}
	}

	if useSequential && collectorIndex != 0 {
		return nil,
			0,
			fmt.Errorf(`File ends with a %d-gram`, collectorIndex)
	}

	return result, ngramCounter, nil
}

// ******** Private functions ********

// prepareCollector prepares the collector for the next rune.
func prepareCollector(collector []rune, collectorIndex uint, ngramSize uint, useSequential bool) uint {
	if useSequential {
		return 0
	} else {
		if ngramSize >= 8 {
			copy(collector, collector[1:ngramSize])
		} else {
			explicitCopy(collector, ngramSize)
		}

		return collectorIndex - 1
	}
}

// explicitCopy does an explicit copy of the collector elements.
// This is faster than the copy function and a copy loop for sizes less than 8.
func explicitCopy(collector []rune, ngramSize uint) {
	ngramSize--

	_ = collector[ngramSize] // Make boundary check only once.

	collector[0] = collector[1]
	ngramSize--
	if ngramSize == 0 {
		return
	}

	collector[1] = collector[2]
	ngramSize--
	if ngramSize == 0 {
		return
	}

	collector[2] = collector[3]
	ngramSize--
	if ngramSize == 0 {
		return
	}

	collector[3] = collector[4]
	ngramSize--
	if ngramSize == 0 {
		return
	}

	collector[4] = collector[5]
	ngramSize--
	if ngramSize == 0 {
		return
	}

	collector[5] = collector[6]
}

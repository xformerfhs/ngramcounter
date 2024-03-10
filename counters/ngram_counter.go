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
// Version: 1.0.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//

package counters

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"ngramcounter/filehelper"
	"ngramcounter/texthelper"
	"os"
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
func (nc *NgramCounter) CountNGrams(fileName string, n uint) (map[string]uint64, uint64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(f)

	// 1. Wrap the file in a transform that decodes the character encoding to UTF-8.
	tr := transform.NewReader(f, nc.decoder)

	// 2. The scanner uses the transform reader.
	scanner := bufio.NewScanner(tr)

	// 3. Now loop over the lines.
	result := make(map[string]uint64)
	collector := make([]rune, n)
	ngramCounter := uint64(0)
	ngramCount := ngramCounter
	for scanner.Scan() {
		line := scanner.Text()
		compressedLine := texthelper.CompressText(line)
		ngramCount, err = scanLineForNGrams(compressedLine, result, collector, n)
		if err != nil {
			return nil, 0, err
		}
		ngramCounter += ngramCount
	}

	if err = scanner.Err(); err != nil {
		return nil, 0, err
	}

	return result, ngramCounter, nil
}

// ******** Private functions ********

// scanLineForNGrams scans a file line and counts the n-grams in it.
func scanLineForNGrams(line string, counter map[string]uint64, collector []rune, n uint) (uint64, error) {
	ngramCounter := uint64(0)
	collectorIndex := uint(0)
	for _, r := range line {
		collector[collectorIndex] = r
		collectorIndex++

		if collectorIndex >= n {
			ngram := string(collector)
			counter[ngram]++
			ngramCounter++

			collectorIndex = 0
		}
	}

	if collectorIndex != 0 {
		return 0, fmt.Errorf(`Line ends with an incomplete %d-gram: '%s'`, n, line)
	}

	return ngramCounter, nil
}

//
// SPDX-FileCopyrightText: Copyright 2024-2025 Frank Schwab
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
//    2025-06-23: V1.0.0: Created.
//

package resultwriter

import (
	"fmt"
	"ngramcounter/filehelper"
	"ngramcounter/maphelper"
	"ngramcounter/platform"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ******** Private constants ********

const fieldSeparator = ","

// ******** Public functions ********

// WriteCountersToTextFile writes the counter values to a CSV file.
func WriteCountersToTextFile(
	fileName string,
	total uint64,
	counter map[string]uint64,
	isNGram bool,
) (string, error) {
	outFileName := outputFileName(fileName)
	f, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return ``, err
	}
	defer filehelper.CloseFile(f)

	err = writeHeader(f, isNGram)
	if err != nil {
		return ``, err
	}

	counts, countToNgrams := sortedKeysAndInvertedCounterMap(counter)

	inverseTotal := 1.0 / float64(total)
	for _, count := range counts {
		for _, ngram := range countToNgrams[count] {
			err = writeLine(f, ngram, count, inverseTotal)
			if err != nil {
				return ``, err
			}
		}
	}

	return outFileName, nil
}

// ******** Private functions ********

// outputFileName builds the output file name from the components of the input file.
func outputFileName(fileName string) string {
	dir, base, ext := filehelper.PathComponents(fileName)
	if len(ext) != 0 {
		base = base + `_` + ext[1:]
	}

	return filepath.Join(dir, base+`.txt`)
}

// writeHeader writes the text file header.
func writeHeader(f *os.File, isNGram bool) error {
	var err error

	if isNGram {
		_, err = f.WriteString(`NGram`)
	} else {
		_, err = f.WriteString(`Byte`)
	}
	if err != nil {
		return err
	}

	_, err = f.WriteString(fieldSeparator)
	if err != nil {
		return err
	}
	_, err = f.WriteString(`Count`)
	if err != nil {
		return err
	}
	_, err = f.WriteString(fieldSeparator)
	if err != nil {
		return err
	}
	_, err = f.WriteString(`Share`)
	if err != nil {
		return err
	}
	_, err = f.WriteString(platform.LineEnd)
	if err != nil {
		return err
	}

	return nil
}

// writeLine writes one line of data.
func writeLine(f *os.File, ngram string, count uint64, inverseTotal float64) error {
	err := writeNgram(f, ngram)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fieldSeparator)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprint(count))
	if err != nil {
		return err
	}

	_, err = f.WriteString(fieldSeparator)
	if err != nil {
		return err
	}

	err = writePercentage(f, count, inverseTotal)
	if err != nil {
		return err
	}

	_, err = f.WriteString(platform.LineEnd)
	if err != nil {
		return err
	}

	return nil
}

// writeNgram writes the value of the [ngram] with the correct formatting for Excel.
func writeNgram(f *os.File, ngram string) error {
	_, err := f.WriteString(`"`)
	if err != nil {
		return err
	}

	// Double all double quotes if necessary.
	ngram = doubleDoubleQuotes(ngram)

	_, err = f.WriteString(ngram)
	if err != nil {
		return err
	}

	_, err = f.WriteString(`"`)
	if err != nil {
		return err
	}

	return nil
}

// doubleDoubleQuotes replaces all double quotes by double double quotes.
func doubleDoubleQuotes(s string) string {
	pos := strings.IndexByte(s, '"')
	if pos != -1 {
		return strings.Replace(s, `"`, `""`, -1)
	}
	return s
}

// writePercentage writes the count as a percentage of the total.
func writePercentage(f *os.File, count uint64, inverseTotal float64) error {
	fractionText := percentageTextFromCount(count, inverseTotal)

	_, err := f.WriteString(fractionText)
	if err != nil {
		return err
	}

	_, err = f.WriteString(`%`)
	if err != nil {
		return err
	}

	return nil
}

// percentageTextFromCount builds the percent text from [count] and [inverseTotal]
// and changes the decimal separator depending on [separator].
func percentageTextFromCount(count uint64, inverseTotal float64) string {
	return fmt.Sprint(float64(count) * inverseTotal * 100)
}

// sortedKeysAndInvertedCounterMap creates a map from counts to a slice of alphabetically sorted
// n-grams that have been found this number of times and returns this map together with a sorted
// slice of the keys in this map. The keys are sorted in descending order.
func sortedKeysAndInvertedCounterMap(counter map[string]uint64) ([]uint64, map[uint64][]string) {
	invertedCounters := sortedInvertedCounterMap(counter)

	counts := maphelper.SortedKeys(invertedCounters)
	slices.Reverse(counts)

	return counts, invertedCounters
}

// sortedInvertedCounterMap creates a map from counts to a slice of alphabetically sorted
// n-grams that have been found this number of times.
func sortedInvertedCounterMap(counter map[string]uint64) map[uint64][]string {
	invertedCounters := invertCounterMap(counter)
	for _, v := range invertedCounters {
		if len(v) > 1 {
			slices.Sort(v)
		}
	}

	return invertedCounters
}

// invertCounterMap creates a map from counts to a slice of n-grams that have
// been found this number of times.
func invertCounterMap(counter map[string]uint64) map[uint64][]string {
	invertedCounters := make(map[uint64][]string)
	for k, v := range counter {
		var ngramList []string

		ngramList = invertedCounters[v]
		if ngramList == nil {
			ngramList = make([]string, 0)
		}

		ngramList = append(ngramList, k)
		invertedCounters[v] = ngramList
	}

	return invertedCounters
}

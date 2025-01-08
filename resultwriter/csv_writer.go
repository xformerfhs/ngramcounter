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
// Version: 2.0.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-08: V2.0.0: Return CSV file name.
//

package resultwriter

import (
	"bytes"
	"fmt"
	"ngramcounter/filehelper"
	"ngramcounter/maphelper"
	"ngramcounter/platform"
	"ngramcounter/stringhelper"
	"os"
	"path/filepath"
	"slices"
)

// ******** Private constants ********

// utf8BOM is the UTF-8 BOM byte sequence.
var utf8BOM = []byte{0xef, 0xbb, 0xbf}

// ******** Public functions ********

// WriteCountersToCSV writes the counter values to a CSV file.
func WriteCountersToCSV(fileName string, total uint64, counter map[string]uint64, separator string, isNGram bool) (string, error) {
	outFileName := csvFileName(fileName)
	f, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return ``, err
	}
	defer filehelper.CloseFile(f)

	err = writeHeader(f, isNGram, separator)
	if err != nil {
		return ``, err
	}

	ngramList := maphelper.Keys(counter)
	slices.Sort(ngramList)

	inverseTotal := 1.0 / float64(total)
	for _, ngram := range ngramList {
		err = writeNgram(f, ngram)
		if err != nil {
			return ``, err
		}

		_, err = f.WriteString(separator)
		if err != nil {
			return ``, err
		}

		count := counter[ngram]

		_, err = f.WriteString(fmt.Sprint(count))
		if err != nil {
			return ``, err
		}

		_, err = f.WriteString(separator)
		if err != nil {
			return ``, err
		}

		err = writePercentage(f, count, inverseTotal, separator)
		if err != nil {
			return ``, err
		}

		_, err = f.WriteString(platform.LineEnd)
		if err != nil {
			return ``, err
		}
	}

	return outFileName, nil
}

// ******** Private functions ********

// csvFileName builds the CSV file name from the components of the input file.
func csvFileName(fileName string) string {
	dir, base, ext := filehelper.PathComponents(fileName)
	if len(ext) != 0 {
		base = base + `_` + ext[1:]
	}

	return filepath.Join(dir, base+`.csv`)
}

// writeHeader writes the CSV header.
func writeHeader(f *os.File, isNGram bool, separator string) error {
	// Excel ought to know that this file is UTF-8 encoded.
	_, err := f.Write(utf8BOM)
	if err != nil {
		return err
	}

	if isNGram {
		_, err = f.WriteString(`NGram`)
	} else {
		_, err = f.WriteString(`Byte`)
	}
	if err != nil {
		return err
	}

	_, err = f.WriteString(separator)
	if err != nil {
		return err
	}
	_, err = f.WriteString(`Count`)
	if err != nil {
		return err
	}
	_, err = f.WriteString(separator)
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

// writeNgram writes the value of the [ngram] with the correct formatting for Excel.
func writeNgram(f *os.File, ngram string) error {
	err := writeTextPrefixIfStartsWithDigit(f, ngram)
	if err != nil {
		return err
	}

	_, err = f.WriteString(`"`)
	if err != nil {
		return err
	}

	// A double quote needs to be doubled for Excel to understand it.
	if ngram == `"` {
		ngram = `""`
	}

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

// writeTextPrefixIfStartsWithDigit writes a text prefix if the [ngram] starts with
// a digit character, as Excel otherwise would interpret this as a number, instead of a text.
func writeTextPrefixIfStartsWithDigit(f *os.File, ngram string) error {
	firstByte := ngram[0]
	if firstByte >= '0' && firstByte <= '9' {
		// The "=" is needed so that Excel understands that this is a text, not a number.
		_, err := f.WriteString(`=`)
		if err != nil {
			return err
		}
	}

	return nil
}

// writePercentage writes the count as a percentage of the total.
func writePercentage(f *os.File, count uint64, inverseTotal float64, separator string) error {
	fractionText := percentageTextFromCount(count, inverseTotal, separator)

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
func percentageTextFromCount(count uint64, inverseTotal float64, separator string) string {
	percentage := float64(count) * inverseTotal * 100
	percentageText := fmt.Sprint(percentage)
	if separator[0] != ',' {
		// Replace '.' with ',', if separator is not ','.
		percentageBytes := stringhelper.UnsafeStringBytes(percentageText)
		pos := bytes.IndexByte(percentageBytes, '.')
		if pos >= 0 {
			percentageBytes[pos] = ','
		}
	}

	return percentageText
}

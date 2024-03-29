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

package resultwriter

import (
	"fmt"
	"ngramcounter/filehelper"
	"ngramcounter/platform"
	"ngramcounter/stringhelper"
	"os"
)

// ******** Private constants ********

// utf8BOM is the UTF-8 BOM byte sequence.
var utf8BOM = []byte{0xef, 0xbb, 0xbf}

// ******** Public functions ********

// WriteCountersToCSV writes the counter values to a CSV file.
func WriteCountersToCSV(fileName string, total uint64, counter map[string]uint64, separator string, isNGram bool) error {
	outFileName := fileName + `.csv`
	f, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer filehelper.CloseFile(f)

	_, err = f.Write(utf8BOM)
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

	inverseTotal := 1.0 / float64(total)
	for g, c := range counter {
		firstByte := g[0]
		if firstByte >= '0' && firstByte <= '9' {
			// The "=" is needed so that Excel understands that this is a text, not a number.
			_, err = f.WriteString(`=`)
			if err != nil {
				return err
			}
		}

		_, err = f.WriteString(`"`)
		if err != nil {
			return err
		}

		// A double quote needs to be doubled for Excel to understand it.
		if g == `"` {
			g = `""`
		}
		_, err = f.WriteString(g)
		if err != nil {
			return err
		}

		_, err = f.WriteString(`"`)
		if err != nil {
			return err
		}
		_, err = f.WriteString(separator)
		if err != nil {
			return err
		}

		_, err = f.WriteString(fmt.Sprint(c))
		if err != nil {
			return err
		}

		_, err = f.WriteString(separator)
		if err != nil {
			return err
		}

		// Convert count to fraction.
		fraction := float64(c) * inverseTotal
		fractionText := fmt.Sprint(fraction)
		if separator[0] != ',' {
			// Replace '.' with ',', if separator is not ','.
			fractionBytes := stringhelper.UnsafeStringBytes(fractionText)
			fractionBytes[1] = ','
		}

		_, err = f.WriteString(fractionText)
		if err != nil {
			return err
		}

		_, err = f.WriteString(platform.LineEnd)
		if err != nil {
			return err
		}
	}

	return nil
}

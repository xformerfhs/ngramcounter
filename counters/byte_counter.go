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
	"errors"
	"io"
	"ngramcounter/filehelper"
	"os"
)

// ******** Private constants ********

// bufferSize is the size of the read buffer.
const bufferSize = 64 * 1024

// ******** Private variables ********

// byteHexTable contains the hexadecimal representation of its indizes.
var byteHexTable []string

// ******** Public functions ********

// CountBytes counts how often a byte appears in a file.
func CountBytes(fileName string) (map[string]uint64, uint64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(f)

	buffer := make([]byte, bufferSize)

	total := uint64(0)
	byteCounter := make(map[string]uint64, 256)

	if byteHexTable == nil {
		byteHexTable = buildByteTable()
	}

	for {
		var n int
		n, err = f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, 0, err
		}

		total += uint64(n)
		for i := 0; i < n; i++ {
			byteCounter[byteHexTable[buffer[i]]]++
		}
	}

	return byteCounter, total, nil
}

// ******** Private functions ********

// buildByteTable builds the index to hexdecimal representation table.
func buildByteTable() []string {
	result := make([]string, 256)

	for i := 0; i < 256; i++ {
		result[i] = string(byteAsHex(byte(i)))
	}

	return result
}

// hexDigits contains the hexadecimal digits
var hexDigits = []byte{
	'0', '1', '2', '3',
	'4', '5', '6', '7',
	'8', '9', 'A', 'B',
	'C', 'D', 'E', 'F',
}

// hexBuffer contains the two hexadecimal digits as bytes.
var hexBuffer [2]byte

// byteAsHex converts a byte into a hexadecimal representation.
func byteAsHex(b byte) []byte {
	hexBuffer[0] = hexDigits[(b>>4)&0x0f]
	hexBuffer[1] = hexDigits[b&0x0f]
	return hexBuffer[:]
}

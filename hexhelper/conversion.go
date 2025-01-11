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
// Version: 1.0.1
//
// Change history:
//    2025-01-08: V1.0.0: Created.
//    2025-01-11: V1.0.1: Correct wrong function name.
//

package hexhelper

var byteHexTable []string

func init() {
	byteHexTable = buildByteHexTable()
}

func ByteToString(b byte) string {
	return byteHexTable[b]
}

// buildByteHexTable builds the index to the hexadecimal representation table.
func buildByteHexTable() []string {
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

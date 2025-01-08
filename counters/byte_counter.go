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

// ******** Public functions ********

// CountBytes counts how often a byte appears in a file.
func CountBytes(fileName string) (map[byte]uint64, uint64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(f)

	buffer := make([]byte, bufferSize)

	total := uint64(0)
	byteCounter := make(map[byte]uint64, 256)

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
			byteCounter[buffer[i]]++
		}
	}

	return byteCounter, total, nil
}

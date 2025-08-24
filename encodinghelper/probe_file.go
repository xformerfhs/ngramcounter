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
// Version: 1.2.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-19: V1.1.0: Correct handling of short files.
//    2025-08-23: V1.2.0: Recognize UTF-32.
//

package encodinghelper

import (
	"errors"
	"ngramcounter/filehelper"
	"os"
	"strings"

	"golang.org/x/text/encoding"
)

// ******** Private constants ********

// utf16BeBom contains the bytes of a UTF-16BE BOM.
var utf16BeBom = []byte{0xfe, 0xff}

// utf16LeBom contains the bytes of a UTF-16LE BOM.
var utf16LeBom = []byte{0xff, 0xfe}

// utf8Bom contains the bytes of a UTF-8 BOM.
var utf8Bom = []byte{0xef, 0xbb, 0xbf}

// ******** Public functions ********

// ProbeFile reads the first bytes of a file to check for BOMs.
// If it finds one, it returns the corresponding encoding.
func ProbeFile(fileName string) (encoding.Encoding, string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, ``, err
	}
	defer filehelper.CloseFile(f)

	// 1. Read the first four bytes.
	var readCount int
	miniBuffer := make([]byte, 4)
	readCount, err = f.Read(miniBuffer)
	if err != nil {
		return nil, ``, err
	}

	// 3. Check read bytes.

	// File has only 1 byte. There is no BOM.
	if readCount < 2 {
		return nil, ``, nil
	}

	var encodingName string
	var found bool
	encodingName, found, err = checkBufferForBom(miniBuffer, readCount)

	if found {
		if strings.HasPrefix(encodingName, `utf32`) {
			return nil, encodingName, errors.New(`UTF-32 is not supported`)
		}

		return textToEncoding[encodingName].encoding, encodingName, nil
	}

	return nil, ``, nil
}

// ******** Private functions ********

// checkBufferForBom checks the first four bytes of a buffer for BOMs.
func checkBufferForBom(buffer []byte, count int) (string, bool, error) {
	// Suppress unnecessary bounds checks.
	_ = buffer[3]
	_ = utf8Bom[2]
	_ = utf16BeBom[1]
	_ = utf16LeBom[1]

	switch buffer[0] {
	case utf8Bom[0]:
		if count >= 3 &&
			buffer[1] == utf8Bom[1] &&
			buffer[2] == utf8Bom[2] {
			return `utf8`, true, nil
		}

	case utf16LeBom[0]:
		if buffer[1] == utf16LeBom[1] {
			if count >= 4 &&
				buffer[2] == 0 &&
				buffer[3] == 0 {
				return `utf32le`, true, nil
			}

			return `utf16le`, true, nil
		}

	case utf16BeBom[0]:
		if buffer[1] == utf16BeBom[1] {
			return `utf16be`, true, nil
		}

	case 0:
		if count >= 4 &&
			buffer[1] == 0 &&
			buffer[2] == utf16BeBom[0] &&
			buffer[3] == utf16BeBom[1] {
			return `utf32be`, true, nil
		}
	}

	return ``, false, nil
}

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

package encodinghelper

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

// ******** Private types ********

// encodingInfo holds the name and the encoding.Encoding of an encoding.
type encodingInfo struct {
	name     string
	encoding encoding.Encoding
}

// ******** Private constants ********

// utf16BeEncoding contains a UTF-16BE encoding.
var utf16BeEncoding = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)

// utf16LeEncoding contains a UTF-16LE encoding.
var utf16LeEncoding = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)

// textToEncoding maps an encoding specification to the corresponding encoding information.
var textToEncoding = map[string]encodingInfo{}

// ******** Public functions *********

// InitializeEncoding initializes encoding variables.
func InitializeEncoding() {
	fillEncodingMap()
}

// ******** Private functions *********

// fillEncdingMap fills the encoding map.
func fillEncodingMap() {
	textToEncoding[`utf8`] = encodingInfo{name: `UTF-8`, encoding: unicode.UTF8BOM}
	textToEncoding[`utf16be`] = encodingInfo{name: `UTF-16BE`, encoding: utf16BeEncoding}
	textToEncoding[`utf16le`] = encodingInfo{name: `UTF-16LE`, encoding: utf16LeEncoding}

	for _, enc := range charmap.All {
		cm, isCm := enc.(*charmap.Charmap)
		if isCm {
			charMapName := cm.String()
			textToEncoding[normalizeEncoding(charMapName)] = encodingInfo{name: charMapName, encoding: enc}
		}
	}
}

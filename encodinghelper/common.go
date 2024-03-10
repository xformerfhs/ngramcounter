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
var textToEncoding = map[string]encodingInfo{
	`cp437`:     {name: charmap.CodePage437.String(), encoding: charmap.CodePage437},
	`cp850`:     {name: charmap.CodePage850.String(), encoding: charmap.CodePage850},
	`cp852`:     {name: charmap.CodePage852.String(), encoding: charmap.CodePage852},
	`iso88591`:  {name: charmap.ISO8859_1.String(), encoding: charmap.ISO8859_1},
	`iso885915`: {name: charmap.ISO8859_15.String(), encoding: charmap.ISO8859_15},
	`utf16be`:   {name: `UTF-16BE`, encoding: utf16BeEncoding},
	`utf16le`:   {name: `UTF-16LE`, encoding: utf16LeEncoding},
	`utf8`:      {name: `UTF-8`, encoding: unicode.UTF8BOM},
	`win1250`:   {name: charmap.Windows1250.String(), encoding: charmap.Windows1250},
	`win1252`:   {name: charmap.Windows1252.String(), encoding: charmap.Windows1252},
}

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
// Version: 2.0.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-09: V1.0.1: Simplified sort call.
//    2025-02-16: V1.1.0: Simplified name normalization.
//    2025-08-24: V2.0.0: Changed function name to "EncodingForName".
//

package encodinghelper

import (
	"fmt"
	"ngramcounter/maphelper"
	"ngramcounter/stringhelper"
	"strings"
	"unicode"

	"golang.org/x/text/encoding"
)

// ******** Public functions ********

// EncodingForName translates a character encoding text into an encoding.Encoding.
func EncodingForName(charEncoding string) (encoding.Encoding, string, error) {
	enc, exists := textToEncoding[normalizeEncoding(charEncoding)]
	if !exists {
		return nil, ``, fmt.Errorf(`Invalid character encoding: '%s'`, charEncoding)
	}

	return enc.encoding, enc.name, nil
}

// EncodingTextList returns the list of encodings as text.
func EncodingTextList() []string {
	keys := maphelper.SortedKeys(textToEncoding)

	result := make([]string, len(keys))
	for i, k := range keys {
		result[i] = fmt.Sprintf(`%-13s: %s`, k, textToEncoding[k].name)
	}

	return result
}

// ******** Private functions ********

// normalizeEncoding normalizes the encoding text, so it is all lower case, contains no separators
// and has a unique way to specify 'win' and 'cp'.
func normalizeEncoding(charEncoding string) string {
	normalizedEncoding := cleanEncodingText(charEncoding)

	// Change "macintosh#" to "mac#". No other transformation is needed in this case.
	if strings.HasPrefix(normalizedEncoding, `macintosh`) {
		normalizedEncoding = `mac` + normalizedEncoding[9:]
		return normalizedEncoding
	}

	// The following cases are mutually exclusive, so they are put into a switch statement.
	switch {
	// Remove "ibm" from "ibmcodepage#".
	case strings.HasPrefix(normalizedEncoding, `ibm`):
		normalizedEncoding = normalizedEncoding[3:]

	// Change "windows#" to "win#".
	case strings.HasPrefix(normalizedEncoding, `windows`):
		normalizedEncoding = `win` + normalizedEncoding[7:]
	}

	// The following statements clean up special cases that may be produced by the previous ones.

	// Change "wincodepage#" to "codepage#".
	if strings.HasPrefix(normalizedEncoding, `wincodepage`) {
		normalizedEncoding = normalizedEncoding[3:]
	}

	// Change "codepage#" to "cp#".
	if strings.HasPrefix(normalizedEncoding, `codepage`) {
		normalizedEncoding = `cp` + normalizedEncoding[8:]
	}

	return normalizedEncoding
}

// cleanBuilder is used to build the clean encoding text.
var cleanBuilder stringhelper.Builder

// cleanEncodingText converts an encoding specification into all lower-case and removes all separators.
func cleanEncodingText(charEncoding string) string {
	cleanBuilder.Ensure(len(charEncoding))
	cleanBuilder.Reset()

	for _, r := range charEncoding {
		if r != '-' && r != '_' && r != '.' && r != ' ' {
			_, _ = cleanBuilder.WriteRune(unicode.ToLower(r))
		}
	}

	resultString := cleanBuilder.String()

	// utf16 has to be mapped to utf16le, as this is the default UTF-16 encoding on Windows.
	if resultString == `utf16` {
		resultString = `utf16le`
	}

	return resultString
}

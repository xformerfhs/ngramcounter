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
	"fmt"
	"golang.org/x/text/encoding"
	"ngramcounter/maphelper"
	"slices"
	"strings"
	"unicode"
)

// ******** Public functions ********

// TranslateEncoding translates a character encoding text into an encoding.Encoding.
func TranslateEncoding(charEncoding string) (encoding.Encoding, string, error) {
	enc, exists := textToEncoding[normalizeEncoding(charEncoding)]
	if !exists {
		return nil, ``, fmt.Errorf(`Invalid character encoding: '%s'`, charEncoding)
	}

	return enc.encoding, enc.name, nil
}

// EncodingTextList returns the list of encodings as text.
func EncodingTextList() []string {
	keys := maphelper.Keys(textToEncoding)
	slices.Sort(keys)

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
	normalizedText := cleanEncodingText(charEncoding)

	if strings.HasPrefix(normalizedText, `ibm`) {
		normalizedText = normalizedText[3:]
	}

	if strings.HasPrefix(normalizedText, `macintosh`) {
		normalizedText = strings.Replace(normalizedText, `macintosh`, `mac`, 1)
	}

	if strings.Contains(normalizedText, `codepage`) {
		normalizedText = strings.Replace(normalizedText, `codepage`, `cp`, 1)
	}

	if strings.HasPrefix(normalizedText, `windows`) {
		normalizedText = strings.Replace(normalizedText, `windows`, `win`, 1)
	}

	if strings.HasPrefix(normalizedText, `wincp`) {
		normalizedText = strings.Replace(normalizedText, `wincp`, `cp`, 1)
	}

	return normalizedText
}

// cleanEncodingText converts an encoding specification into all lower-case and removes all separators.
func cleanEncodingText(charEncoding string) string {
	var result strings.Builder
	result.Grow(len(charEncoding))

	for _, r := range charEncoding {
		if r != '-' && r != '_' && r != '.' && r != ' ' {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	resultString := result.String()

	// utf16 has to be mapped to utf16le, as this is the default UTF-16 encoding on Windows.
	if resultString == `utf16` {
		resultString = `utf16le`
	}

	return resultString
}

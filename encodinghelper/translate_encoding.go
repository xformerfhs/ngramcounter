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
		result[i] = fmt.Sprintf(`%-10s: %s`, k, textToEncoding[k].name)
	}

	return result
}

// ******** Private functions ********

// normalizeEncoding normalizes the encoding text, so it is all lower case, contains no separators
// and has a unique way to specify 'win' and 'cp'.
func normalizeEncoding(charEncoding string) string {
	normalizedText := cleanEncodingText(charEncoding)

	if strings.HasPrefix(normalizedText, `codepage`) {
		return strings.Replace(normalizedText, `codepage`, `cp`, 1)
	}

	if strings.HasPrefix(normalizedText, `windows`) {
		return strings.Replace(normalizedText, `windows`, `win`, 1)
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

	// utf16 has to be mapped to utf16be.
	if resultString == `utf16` {
		resultString = `utf16be`
	}

	return resultString
}

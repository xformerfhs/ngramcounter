package texthelper

import (
	"strings"
	"unicode"
)

var compressedLine strings.Builder

func CompressText(line string) string {
	compressedLine.Reset()
	for _, r := range line {
		if !unicode.Is(unicode.Pattern_White_Space, r) {
			_, _ = compressedLine.WriteRune(r)
		}
	}

	return compressedLine.String()
}

package encodinghelper

import (
	"golang.org/x/text/encoding"
	"ngramcounter/filehelper"
	"os"
)

// ******** Private constants ********

// utf16BeBom contains the bytes of an UTF-16BE BOM.
var utf16BeBom = []byte{0xfe, 0xff}

// utf16LeBom contains the bytes of an UTF-16LE BOM.
var utf16LeBom = []byte{0xff, 0xfe}

// utf8Bom contains the bytes of an UTF-8 BOM.
var utf8Bom = []byte{0xef, 0xbb, 0xbf}

// ******** Public functions ********

// ProbeFile reads the first bytes of a file to check for BOMs.
// If it finds one, it returns the corresponding encoding.
// On return the file's position is 0.
func ProbeFile(fileName string) (encoding.Encoding, string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, ``, err
	}
	defer filehelper.CloseFile(f)

	// 1. Read the first three bytes.
	miniBuffer := make([]byte, 3)
	_, err = f.Read(miniBuffer)
	if err != nil {
		return nil, ``, err
	}

	var ei encodingInfo
	// 3. Check read bytes.
	if miniBuffer[0] == utf16BeBom[0] &&
		miniBuffer[1] == utf16BeBom[1] {
		ei = textToEncoding[`utf16be`]
		return ei.encoding, ei.name, nil
	}

	if miniBuffer[0] == utf16LeBom[0] &&
		miniBuffer[1] == utf16LeBom[1] {
		ei = textToEncoding[`utf16le`]
		return ei.encoding, ei.name, nil
	}

	if miniBuffer[0] == utf8Bom[0] &&
		miniBuffer[1] == utf8Bom[1] &&
		miniBuffer[2] == utf8Bom[2] {
		ei = textToEncoding[`utf8`]
		return ei.encoding, ei.name, nil
	}

	return nil, ``, nil
}

package counters

import (
	"errors"
	"io"
	"ngramcounter/filehelper"
	"os"
)

const bufferSize = 64 * 1024

var byteTable []string

func CountBytes(fileName string) (map[string]uint64, uint64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer filehelper.CloseFile(f)

	buffer := make([]byte, bufferSize)

	total := uint64(0)
	byteCounter := make(map[string]uint64, 256)

	if byteTable == nil {
		byteTable = buildByteTable()
	}

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
			byteCounter[byteTable[buffer[i]]]++
		}
	}

	return byteCounter, total, nil
}

func buildByteTable() []string {
	result := make([]string, 256)

	for i := 0; i < 256; i++ {
		result[i] = string(byteAsHex(byte(i)))
	}

	return result
}

var hexChars = []byte{
	'0', '1', '2', '3',
	'4', '5', '6', '7',
	'8', '9', 'A', 'B',
	'C', 'D', 'E', 'F',
}

var hexBuffer [2]byte

func byteAsHex(b byte) []byte {
	hexBuffer[0] = hexChars[(b>>4)&0x0f]
	hexBuffer[1] = hexChars[b&0x0f]
	return hexBuffer[:]
}

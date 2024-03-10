package texthelper

import "bytes"

type LineSplitter struct {
	afterCR bool
}

func (s *LineSplitter) Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if s.afterCR {
		s.afterCR = false
		if data[0] == '\n' {
			// We had a carriage return before, so this newline needs to be skipped.
			return 1, nil, nil
		}
	}

	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		if data[i] == '\n' {
			// We have a full line terminated by a single newline.
			return i + 1, data[0:i], nil
		}

		// We have a full line terminated by either a single carriage return or carriage return and newline.
		advance = i + 1
		if len(data) == i+1 {
			// We are at the end of the input and do not know yet if the next symbol corresponds to the current carriage return or not.
			s.afterCR = true
		} else if data[i+1] == '\n' {
			advance += 1
		}

		return advance, data[0:i], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

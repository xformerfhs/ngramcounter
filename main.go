//
// SPDX-FileCopyrightText: Copyright 2024-2025 Frank Schwab
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
// Version: 3.2.1
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-08: V1.0.1: Removed unnecessary arguments in call of realMain.
//    2025-01-08: V2.0.0: Scan n-grams in overlapped mode and use an option to revert to sequential mode.
//    2025-01-09: V2.0.1: Corrected CSV write error message.
//    2025-01-09: V2.0.2: Simplified sorting.
//    2025-01-11: V2.0.3: Print mode.
//    2025-01-11: V2.0.4: Corrected error message for incomplete n-grams in sequential mode.
//    2025-01-11: V2.1.0: Simplified counting and make it faster.
//    2025-01-12: V2.2.0: Simplified counting.
//    2025-01-19: V2.3.0: Data is written in descending count order.
//    2025-01-19: V2.4.0: Corrected handling of short files.
//    2025-02-16: V2.4.1: Simplified normalization of encoding names.
//    2025-03-18: V2.4.2: Correct counting message for 1-grams.
//    2025-06-22: V2.5.0: Added "allChars" option.
//    2025-06-23: V3.0.0: Write text files, no "separator" option, anymore.
//    2025-06-23: V3.0.1: Refactored the result writer.
//    2025-06-25: V3.0.2: Simplified n-gram counter.
//    2025-08-23: V3.1.0: Recognize UTF-32.
//    2025-08-24: V3.1.1: Simplified handling of n-gram counter, implement maximum n-gram size.
//    2025-08-24: V3.2.0: Use much less memory when counting n-grams.
//    2025-08-24: V3.2.1: Correct handling of "size" option.
//

package main

import (
	"ngramcounter/logger"
	"os"
	"runtime"
)

// myName contains the name of this executable. It is set in init.
var myName string

// myVersion contains the version number of this executable.
const myVersion = `3.2.1`

// ******** Formal main function ********

// main is the main function and only a stub for a real main function.
func main() {
	logger.PrintInfof(11, `Begin %s V%s (%s)`, myName, myVersion, runtime.Version())
	// Hack, so that we have a way to have args as arguments, set the exit code and run defer functions.
	// This is a severe design deficiency of Go 1.
	rc := realMain()
	logger.PrintInfof(12, `End %s V%s`, myName, myVersion)
	os.Exit(rc)
}

// ******** Private functions ********

// realMain is the real main function which obeys defers and sets a return code.
func realMain() int {
	defineCommandLineFlags()

	if useHelp {
		printUsage()
		return rcOK
	}

	rc := checkCommandLineFlags()
	if rc != rcOK {
		return rc
	}

	var err error

	if ngramSize == 0 {
		logger.PrintInfo(13, `Counting bytes`)
		err = countBytes()
	} else {
		if ngramSize > 1 {
			logger.PrintInfof(14, `Counting %d-grams with %s in %s mode`, ngramSize, charsText(), modeText())
		} else {
			logger.PrintInfof(14, `Counting %d-grams with %s`, ngramSize, charsText())
		}

		err = countNGrams(charEncoding, ngramSize, useSequential, allChars)
	}

	if err != nil {
		logger.PrintError(16, err.Error())
		return rcProcessingError
	}

	return rcOK
}

// modeText returns the string representation of the useSequential flag.
func modeText() string {
	if useSequential {
		return `sequential`
	} else {
		return `overlapping`
	}
}

// charsText returns the string representation of the allChars flag.
func charsText() string {
	if allChars {
		return `all characters`
	} else {
		return `only letters and numbers`
	}
}

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
// Version: 2.1.0
//
// Change history:
//    2024-03-10: V1.0.0: Created.
//    2025-01-08: V1.0.1: Remove unnecessary arguments in call of realMain.
//    2025-01-08: V2.0.0: Scan n-grams in overlapped mode and use an option to revert to sequential mode.
//    2025-01-09: V2.0.1: Correct CSV write error message.
//    2025-01-09: V2.0.2: Simplified sorting.
//    2025-01-11: V2.0.3: Print mode.
//    2025-01-11: V2.0.4: Correct error message for incomplete n-grams in sequential mode.
//    2025-01-11: V2.1.0: Simplify counting and make it faster.
//    2025-01-12: V2.2.0: Simplify counting.
//

package main

import (
	"ngramcounter/logger"
	"os"
	"runtime"
)

// myName contains the name of this executable.
var myName string

// myVersion contains the version number of this executable.
const myVersion = `2.2.0`

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

// ******** Private constants ********

// maxNGram is the maximum allowed n-gram length
const maxNGram = 50

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
		err = countBytes(separator)
	} else {
		if ngramSize <= maxNGram {
			logger.PrintInfof(14, `Counting %d-grams in %s mode`, ngramSize, modeText())
			err = countNGrams(charEncoding, ngramSize, separator, useSequential)
		} else {
			logger.PrintErrorf(15, `n-gram count '%d' is too large (max=%d)`, ngramSize, maxNGram)
			return rcCmdLineError
		}
	}

	if err != nil {
		logger.PrintError(16, err.Error())
		return rcProcessingError
	}

	return rcOK
}

// modeText returns the string representation of the mode.
func modeText() string {
	if useSequential {
		return `sequential`
	} else {
		return `overlapping`
	}
}

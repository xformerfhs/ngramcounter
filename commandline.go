//
// SPDX-FileCopyrightText: Copyright 2025 Frank Schwab
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
//    2025-01-08: V1.0.0: Created.
//

package main

import (
	"flag"
	"ngramcounter/encodinghelper"
	"ngramcounter/logger"
)

// Possible return codes.
const (
	rcOK              = 0
	rcCmdLineError    = 1
	rcProcessingError = 2
)

// ******** Private variables ********

// ngramSize is the size of the n-gram.
var ngramSize uint

// charEncoding is the character encoding of the source file.
var charEncoding string

// separator is the output field separator character.
var separator string

// useSequential specifies that the n-grams should be read useSequential and not overlapped.
var useSequential bool

// useHelp specifies that the help should be printed.
var useHelp bool

// ******** Private functions ********

// defineCommandLineFlags defines the command line flags.
func defineCommandLineFlags() {
	flag.UintVar(&ngramSize, `size`, 0, `Scan files as n-grams with the given length (if this is not set, bytes are counted)`)

	flag.StringVar(&charEncoding, `encoding`, encodinghelper.PlatformDefaultEncoding(), `Character encoding for n-grams`)

	flag.StringVar(&separator, `separator`, `;`, `Output field separator (either ',' or ';')`)

	flag.BoolVar(&useSequential, `sequential`, false, `Read n-grams in sequential mode`)

	flag.BoolVar(&useHelp, `help`, false, `Print usage and exit`)

	flag.Usage = printUsage

	flag.Parse()
}

// checkCommandLineFlags checks the command line flags.
func checkCommandLineFlags() int {
	if flag.NArg() == 0 {
		logger.PrintError(21, `File names missing`)
		printUsage()
		return rcCmdLineError
	}

	if separator != `,` && separator != `;` {
		logger.PrintErrorf(22, `Separator must be either ',' or ';' but is '%s'`, separator)
		return rcCmdLineError
	}

	return rcOK
}

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
	"fmt"
	"ngramcounter/encodinghelper"
	"ngramcounter/logger"
	"os"
)

// printAnalysisInfo prints which file is analyzed.
func printAnalysisInfo(fileName string) {
	logger.PrintInfof(31, `Analyzing file '%s'`, fileName)
}

// printCSVInfo prints the output file name.
func printCSVInfo(fileName string) {
	logger.PrintInfof(32, `CSV file: '%s'`, fileName)
}

// makeCountError build an error from an error and a file name for the count phase.
func makeCountError(fileName string, err error) error {
	return fmt.Errorf(`Error analyzing file '%s': %v`, fileName, err)
}

// makeWriteError build an error from an error and a file name for the write phase.
func makeWriteError(fileName string, err error) error {
	return fmt.Errorf(`Error writing count file '%s': %v`, fileName, err)
}

// printUsage prints the usage text
func printUsage() {
	fmt.Println(`A program to count bytes or n-grams in files

Usage:`)

	_, _ = fmt.Fprintf(os.Stderr, "\n%s [-size <count>] [-encoding <encoding>] [-separator <char>] [-sequential] [files...]\n\nwith the following options:\n\n",
		myName)
	flag.PrintDefaults()

	_, _ = fmt.Fprintf(os.Stderr, `
followed by a list of file names.

The results are written to '<filename>.csv', i.e. the filename with '.csv' appended.

The format is a 'character separated value' file which can be imported e.g. by Excel.
`)

	_, _ = fmt.Fprintln(os.Stderr, "\n'encoding' can be one of the following values of the first column:")
	for _, e := range encodinghelper.EncodingTextList() {
		_, _ = fmt.Fprint(os.Stderr, `  `)
		_, _ = fmt.Fprintln(os.Stderr, e)
	}
	_, _ = fmt.Fprintln(os.Stderr, "\n  'utf16' may be used as a synonym for 'utf16le'")

	_, _ = fmt.Fprintln(os.Stderr)
}

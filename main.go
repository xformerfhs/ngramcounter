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

package main

import (
	"flag"
	"fmt"
	"golang.org/x/text/encoding"
	"ngramcounter/counters"
	"ngramcounter/encodinghelper"
	"ngramcounter/filehelper"
	"ngramcounter/logger"
	"ngramcounter/resultwriter"
	"os"
)

// myName contains the name of this executable.
var myName string

// myVersion contains the version number of this executable.
const myVersion = `1.0.0`

// ******** Formal main function ********

// main is the main function and only a stub for a real main function.
func main() {
	myName = filehelper.GetRealBaseName(os.Args[0])
	logger.PrintInfof(11, `Begin %s v%s`, myName, myVersion)
	// Hack, so that we have a way to have args as arguments, set the exit code and run defer functions.
	// This is a severe design deficiency of Go 1
	rc := mainWithReturnCode(os.Args[1:])
	logger.PrintInfof(12, `End %s v%s`, myName, myVersion)
	os.Exit(rc)
}

const (
	rcOk           = 0
	rcCmdLineError = 1
	rcProcessError = 2
)

const maxNGram = 50

// mainWithReturnCode is the real main function which obeys defers and sets a return code.
func mainWithReturnCode(args []string) int {
	var err error

	var ngram uint
	flag.UintVar(&ngram, `ngram`, 0, `Scan files as n-grams with the given length (if this is not set, bytes are counted)`)

	var charEncoding string
	flag.StringVar(&charEncoding, `encoding`, encodinghelper.PlatformDefaultEncoding(), `Character encoding for n-grams`)

	var separator string
	flag.StringVar(&separator, `separator`, `;`, `Output field separator (either ',' or ';')`)

	var printHelp bool
	flag.BoolVar(&printHelp, `help`, false, `Print usage`)

	flag.Usage = printUsage

	flag.Parse()

	if printHelp {
		printUsage()
		return rcOk
	}

	if flag.NArg() == 0 {
		logger.PrintError(13, `File names missing`)
		printUsage()
		return rcCmdLineError
	}

	if separator != `,` && separator != `;` {
		logger.PrintErrorf(14, `Separator must be either ',' or ';' but is '%s'`, separator)
		return rcCmdLineError
	}

	if ngram == 0 {
		logger.PrintInfo(15, `Counting bytes`)
		err = countBytes(separator)
	} else {
		if ngram <= maxNGram {
			logger.PrintInfof(16, `Counting %d-grams`, ngram)
			err = countNGrams(charEncoding, ngram, separator)
		} else {
			logger.PrintErrorf(17, `n-gram count '%d' is too large (max=%d)`, ngram, maxNGram)
			return rcCmdLineError
		}
	}

	if err != nil {
		logger.PrintError(18, err.Error())
		return rcProcessError
	}

	return rcOk
}

// countBytes counts the bytes in all specified files.
func countBytes(separator string) error {
	var err error
	var count map[string]uint64
	var total uint64

	for _, fileName := range flag.Args() {
		printCountInfo(fileName)
		count, total, err = counters.CountBytes(fileName)
		if err != nil {
			return makeCountError(fileName, err)
		}

		err = resultwriter.WriteCountersToCSV(fileName, total, count, separator, false)
		if err != nil {
			return makeWriteError(fileName, err)
		}
	}

	return nil
}

// countNGrams counts n-grams in all specified files.
func countNGrams(charEncoding string, ngram uint, separator string) error {
	var err error
	var count map[string]uint64
	var total uint64

	// 1. Get requested encoding and corresponding n-gRam counter.
	var requestedEncoding encoding.Encoding
	var requestedEncodingName string
	requestedEncoding, requestedEncodingName, err = encodinghelper.TranslateEncoding(charEncoding)
	requestedNgramCounter := counters.NewNgramCounter(requestedEncoding)

	logger.PrintInfof(19, `File encoding is '%s'`, requestedEncodingName)

	// 2. Loop of files.
	for _, fileName := range flag.Args() {
		printCountInfo(fileName)

		// 3. Check if the current file has a byte-order-mark and change counter if it has one.
		var actNgramCounter *counters.NgramCounter
		actNgramCounter, err = chooseCounter(fileName, requestedEncoding, requestedNgramCounter)

		// 4. Count n-grams.
		count, total, err = actNgramCounter.CountNGrams(fileName, ngram)
		if err != nil {
			return makeCountError(fileName, err)
		}

		// 5. Write result.
		err = resultwriter.WriteCountersToCSV(fileName, total, count, separator, true)
		if err != nil {
			return makeWriteError(fileName, err)
		}
	}

	return nil
}

// chooseCounter checks if the file has a byte-order-mark and returns the corresponding counter.
func chooseCounter(fileName string, requestedEncoding encoding.Encoding, requestedNGramCounter *counters.NgramCounter) (*counters.NgramCounter, error) {
	probedEncoding, probedEncodingName, err := encodinghelper.ProbeFile(fileName)
	if err != nil {
		return nil, err
	}

	if probedEncoding != nil && probedEncoding != requestedEncoding {
		logger.PrintInfof(20, `File '%s' has a %s byte-order-mark and is read with this encoding`, fileName, probedEncodingName)
		return counters.NewNgramCounter(probedEncoding), nil
	}

	return requestedNGramCounter, nil
}

// printCountInfo prints which file is counted.
func printCountInfo(fileName string) int {
	logger.PrintInfof(21, `Count file '%s'`, fileName)
	return rcProcessError
}

// makeCountError build an error from an error and a file name for the count phase.
func makeCountError(fileName string, err error) error {
	return fmt.Errorf(`Error counting file '%s': %v`, fileName, err)
}

// makeWriteError build an error from an error and a file name for the write phase.
func makeWriteError(fileName string, err error) error {
	return fmt.Errorf(`Error writing count file for file '%s': %v`, fileName, err)
}

// printUsage prints the usage text
func printUsage() {
	fmt.Println(`A program to count bytes or n-grams in files

Usage:`)

	_, _ = fmt.Fprintf(os.Stderr, "\n%s [-ngram count] [-encoding encoding] [-separator char] [files...]\n\nwith the following options:\n\n",
		myName)
	flag.PrintDefaults()

	_, _ = fmt.Fprintf(os.Stderr, `
followed by a list of file names.

The results are written to 'filename.csv', i.e. the filename with '.csv' appended.

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

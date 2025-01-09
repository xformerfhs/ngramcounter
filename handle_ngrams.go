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
// Version: 1.0.1
//
// Change history:
//    2025-01-08: V1.0.0: Created.
//    2025-01-09: V1.0.1: Correct CSV file error message.
//

package main

import (
	"flag"
	"golang.org/x/text/encoding"
	"ngramcounter/counters"
	"ngramcounter/encodinghelper"
	"ngramcounter/logger"
	"ngramcounter/resultwriter"
)

// countNGrams counts n-grams in all specified files.
func countNGrams(charEncoding string, ngram uint, separator string, useSequential bool) error {
	var err error
	var count map[string]uint64
	var total uint64

	// 1. Get requested encoding and corresponding n-gram counter.
	var requestedEncoding encoding.Encoding
	var requestedEncodingName string
	requestedEncoding, requestedEncodingName, err = encodinghelper.TranslateEncoding(charEncoding)
	if err != nil {
		return err
	}

	requestedNgramCounter := counters.NewNgramCounter(requestedEncoding)

	logger.PrintInfof(19, `File encoding is '%s'`, requestedEncodingName)

	// 2. Loop through files.
	for _, fileName := range flag.Args() {
		printAnalysisInfo(fileName)

		// 3. Check if the current file has a byte-order mark and change counter if it has one.
		var actNgramCounter *counters.NgramCounter
		actNgramCounter, err = chooseCounter(fileName, requestedEncoding, requestedNgramCounter)
		if err != nil {
			return makeCountError(fileName, err)
		}

		// 4. Count n-grams.
		count, total, err = actNgramCounter.CountNGrams(fileName, ngram, useSequential)
		if err != nil {
			return makeCountError(fileName, err)
		}

		// 5. Write result.
		var csvFileName string
		csvFileName, err = resultwriter.WriteCountersToCSV(fileName, total, count, separator, true)
		if err != nil {
			return makeWriteError(csvFileName, err)
		}

		printCSVInfo(csvFileName)
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
		logger.PrintInfof(20, `File '%s' has a %s byte-order mark and is read with this encoding`, fileName, probedEncodingName)
		return counters.NewNgramCounter(probedEncoding), nil
	}

	return requestedNGramCounter, nil
}

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
// Version: 2.0.0
//
// Change history:
//    2025-01-08: V1.0.0: Created.
//    2025-01-09: V1.0.1: Correct CSV file error message.
//    2025-06-23: V2.0.0: Output text file.
//

package main

import (
	"flag"
	"ngramcounter/counters"
	"ngramcounter/hexhelper"
	"ngramcounter/resultwriter"
)

// countBytes counts the bytes in all specified files.
func countBytes() error {
	var err error
	var count map[string]uint64
	var total uint64

	for _, fileName := range flag.Args() {
		printAnalysisInfo(fileName)

		count, total, err = countBytesInFile(fileName)
		if err != nil {
			return makeCountError(fileName, err)
		}

		var outputFileName string
		outputFileName, err = resultwriter.WriteCountersToTextFile(fileName, total, count, false)
		if err != nil {
			return makeWriteError(outputFileName, err)
		}

		printOutputInfo(outputFileName)
	}

	return nil
}

func countBytesInFile(fileName string) (map[string]uint64, uint64, error) {
	count, total, err := counters.CountBytes(fileName)

	return convertByteMapToString(count), total, err
}

func convertByteMapToString(count map[byte]uint64) map[string]uint64 {
	result := make(map[string]uint64, len(count))
	for k, v := range count {
		result[hexhelper.ByteToString(k)] = v
	}

	return result
}

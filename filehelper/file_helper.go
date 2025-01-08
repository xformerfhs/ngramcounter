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
// Version: 2.0.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2025-01-08: V2.0.0: Added PathComponents.
//

package filehelper

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ******** Public functions ********

// CloseFile closes a file and prints an error message if closing failed.
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		printFileOperationError(`clos`, file.Name(), err)
	}
}

// GetRealBaseName gets the base name of a file without the extension.
func GetRealBaseName(filePath string) string {
	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

// PathComponents returns the directory, base name and extension (with leading '.') of the supplied file path.
func PathComponents(filePath string) (string, string, string) {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	base = strings.TrimSuffix(base, ext)
	return dir, base, ext
}

// ******** Private functions ********

// printFileOperationError prints an error message for a file operation.
func printFileOperationError(opName string, filePath string, err error) {
	log.Printf(`Error %sing file '%s': %v`, opName, filePath, err)
}

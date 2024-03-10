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

package texthelper

import (
	"strings"
	"unicode"
)

var compressedLine strings.Builder

func CompressText(line string) string {
	compressedLine.Reset()
	for _, r := range line {
		if !unicode.Is(unicode.Pattern_White_Space, r) {
			_, _ = compressedLine.WriteRune(r)
		}
	}

	return compressedLine.String()
}

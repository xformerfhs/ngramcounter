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
//    2025-02-16: V1.0.0: Created.
//

package stringhelper

import (
	"unicode/utf8"
)

// ******** Public types ********

// A Builder is used to efficiently build a string using [Builder.Write] methods.
// It minimizes memory copying and is intended to be reused and reset.
// The zero value is ready to use. This Builder may be copied.
type Builder struct {
	// External users should never get direct access to this buffer, since
	// the slice at some point will be converted to a string using unsafe, also
	// data between len(buf) and cap(buf) might be uninitialized.
	buf []byte
}

// ******** Public functions ********

// String returns the accumulated string.
func (b *Builder) String() string {
	// This is the first difference: The string is copied from the byte buffer.
	return string(b.buf)
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *Builder) Len() int {
	return len(b.buf)
}

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *Builder) Cap() int {
	return cap(b.buf)
}

// Reset resets the [Builder] to be empty.
func (b *Builder) Reset() {
	// This is the second difference: The byte buffer is reused, instead of allocating a new one.
	b.buf = b.buf[:0]
}

// Ensure ensures that b has at least a capacity of n bytes.
// If n is negative, Ensure panics.
func (b *Builder) Ensure(n int) {
	if n < 0 {
		panic(`negative count`)
	}

	if cap(b.buf) < n {
		b.newBuffer(n)
	}
}

// -------- Write functions --------

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *Builder) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *Builder) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *Builder) WriteRune(r rune) (int, error) {
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *Builder) WriteString(s string) (int, error) {
	b.buf = append(b.buf, s...)
	return len(s), nil
}

// ******** Private functions ********

// newBuffer allocates a new buffer that can hold at least n bytes.
func (b *Builder) newBuffer(n int) {
	buf := make([]byte, len(b.buf), n+(n>>2)) // Allocate 25% more than needed.
	copy(buf, b.buf)
	b.buf = buf
}

/*
 * trie.go
 * Trie
 *
 * Created by Jim Dovey on 16/07/2010.
 *
 * Copyright (c) 2010 Jim Dovey
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *
 * Neither the name of the project's author nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package trie

import (
	"unicode"
	"utf8"
	"container/vector"
	"strings"
)


// Specialized function for TeX-style hyphenation patterns.  Accepts strings of the form '.hy2p'.
// The value it stores is of type vector.IntVector
func (p *Trie) AddPatternString(s string) {
	v := new(vector.IntVector)

	// precompute the Unicode rune for the character '0'
	rune0, _ := utf8.DecodeRune([]byte{'0'})

	strLen := len(s)

	// Using the range keyword will give us each Unicode rune.
	for pos, rune := range s {
		if unicode.IsDigit(rune) {
			if pos == 0 {
				// This is a prefix number
				v.Push(rune - rune0)
			}

			// this is a number referring to the previous character, and has
			// already been handled
			continue
		}

		if pos < strLen-1 {
			// look ahead to see if it's followed by a number
			next := int(s[pos+1])
			if unicode.IsDigit(next) {
				// next char is the hyphenation value for this char
				v.Push(next - rune0)
			} else {
				// hyphenation for this char is an implied zero
				v.Push(0)
			}
		} else {
			// last character gets an implied zero
			v.Push(0)
		}
	}

	pure := strings.Map(func(rune int) int {
		if unicode.IsDigit(rune) {
			return -1
		}
		return rune
	},
		s)
	leaf := p.addRunes(strings.NewReader(pure))
	if leaf == nil {
		return
	}

	leaf.value = v
}

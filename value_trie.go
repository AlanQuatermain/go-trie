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

/*
	The trie package implements a basic character trie type. Instead of using bytes however, it uses
	integer-sized runes as traversal keys.  In Go, this means that each node refers to exactly one Unicode
	character, so the implementation doesn't depend on the particular semantics of UTF-8 byte streams.

	There is an additional specialization, which stores an integer value along with the Unicode character
	on each node.  This is to implement TeX-style hyphenation pattern storage.
*/
package trie

import (
	"strings"
	"container/vector"
	"utf8"
	"sort"
	"unicode"
)

type IntArray []int

// Apparently IntVector.Iter() doesn't exist on my amd64 Mac Pro, but does on my i386 iMac. Er, what?
func iterate(p *vector.IntVector) <-chan int {
	iter := make(chan int, 100)
	go func() {
		for _, i := range *p {
			iter <- i
		}
		close(iter)
	}()
	return iter
}

// Determines whether an input string:IntVector pair has a prefix value.
func hasPrefixValue(s string, v *vector.IntVector) bool {
	return v.Len() > utf8.RuneCountInString(s)
}

// Determines whether an input pattern string has a prefix value.
func patternHasPrefixValue(s string) bool {
	rune, _ := utf8.DecodeRuneInString(s)
	return unicode.IsDigit(rune)
}

// Creates and returns a new ValueTrie instance.
func NewValueTrie() *ValueTrie {
	t := new(ValueTrie)
	t.value = 0
	t.leaf = false
	t.children = make(map[int]*ValueTrie)
	return t
}

// Internal function: adds items to the trie, reading runes from a strings.Reader
func (p *ValueTrie) addRunes(r *strings.Reader, iter <-chan int, hasPrefix bool) {
	rune, _, err := r.ReadRune()
	if err != nil {
		p.leaf = true
		return
	}

	// always read a value from the iterator
	val := <-iter
	n := p.children[rune]

	if n == nil {
		n = NewValueTrie()
		if hasPrefix {
			n.prefixValue = val
			val = <-iter
		}
		n.value = val
		p.children[rune] = n
	}

	// recurse to store sub-runes below the new node
	n.addRunes(r, iter, false) // recursion *never* has a prefix; that's only for the first rune
}

// Adds a string of Unicode characters/runes and their associated values to the ValueTrie. If the string is already
// present, no additional storage happens. Yay!
func (p *ValueTrie) Add(s string, v *vector.IntVector) {
	if len(s) == 0 {
		return
	}

	// append the runes to the trie
	iter := iterate(v)
	p.addRunes(strings.NewReader(s), iter, hasPrefixValue(s, v))
	close(iter)
}

// Adds a TeX-style hyphenation pattern to the ValueTrie.  Accepts string of the form '.hy2p' for example.
func (p *ValueTrie) AddPatternString(s string) {
	iter := make(chan int, 40)
	// precompute the Unicode rune for the character '0'
	rune0, _ := utf8.DecodeRune([]byte{'0'})

	// spawn a goroutine to spit each character's hyphenation value into the channel
	go func() {
		strLen := len(s)

		// Using the range keyword will give us each Unicode rune.
		for pos, rune := range s {
			if unicode.IsDigit(rune) {
				if pos == 0 {
					// This is a prefix number
					iter <- (rune - rune0)
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
					val := next - rune0
					iter <- val
				} else {
					// hyphenation for this char is an implied zero
					iter <- 0
				}
			}
		}

		// close our end of the channel
		close(iter)
	}()

	pure := strings.Map(func(rune int) int {
		if unicode.IsDigit(rune) {
			return -1
		}
		return rune
	},
		s)
	p.addRunes(strings.NewReader(pure), iter, patternHasPrefixValue(s))
	close(iter)
}


// Internal string removal function.  Returns trie if this node is empty following the removal.
func (p *ValueTrie) removeRunes(r *strings.Reader) bool {
	rune, _, err := r.ReadRune()
	if err != nil {
		p.leaf = false
		return len(p.children) == 0
	}

	child, ok := p.children[rune]
	if ok && child.removeRunes(r) {
		// the child is now empty following the removal, so prune it
		p.children[rune] = nil, false
	}

	return len(p.children) == 0
}

// Remove a string from the trie.  Returns true if the Trie is now empty.
func (p *ValueTrie) Remove(s string) bool {
	if len(s) == 0 {
		return len(p.children) == 0
	}

	// remove the runes, returning the final result
	return p.removeRunes(strings.NewReader(s))
}

// Internal string inclusion function.
func (p *ValueTrie) includes(r *strings.Reader) bool {
	rune, _, err := r.ReadRune()
	if err != nil {
		return p.leaf // no more runes + leaf node == the string was present
	}

	child, ok := p.children[rune]
	if !ok {
		return false // no node for this rune was in the trie
	}

	// recurse down to the next node with the remainder of the string
	return child.includes(r)
}

// Test for the inclusion of a particular string in the Trie.
func (p *ValueTrie) Contains(s string) bool {
	if len(s) == 0 {
		return false // empty strings can't be included (how could we add them?)
	}
	return p.includes(strings.NewReader(s))
}

func (p *ValueTrie) getValues(r *strings.Reader, v *vector.IntVector) bool {
	rune, _, err := r.ReadRune()
	if err != nil {
		return p.leaf
	}

	child, ok := p.children[rune]
	if !ok {
		return false
	}
	
	// append the value for this node
	if child.prefixValue != 0 {
		v.Push(child.prefixValue)
	}
	v.Push(child.value)
	return child.getValues(r, v)
}

// Retrieve the values associated with the given string.
func (p *ValueTrie) ValuesForString(s string) (*vector.IntVector, bool) {
	r := strings.NewReader(s)
	v := new(vector.IntVector)
	ok := p.getValues(r, v)
	return v, ok
}


// Internal output-building function used by Members()
func (p *ValueTrie) buildMembers(prefix string, includeValues, includeZeroes bool) *vector.StringVector {
	strList := new(vector.StringVector)

	if p.leaf {
		strList.Push(prefix)
	}

	// for each child, go grab all suffixes
	for rune, child := range p.children {
		buf := make([]byte, 4)
		numChars := utf8.EncodeRune(rune, buf)

		var substr string = prefix + string(buf[0:numChars])
		if includeValues {
			if child.prefixValue != 0 {
				substr += string('0' + child.prefixValue)
			}
			if child.value != 0 || includeZeroes {
				substr += string('0' + child.value)
			}
		}
		strList.AppendVector(child.buildMembers(substr, includeValues, includeZeroes))
	}

	return strList
}

// Retrieves all member strings, in order.
func (p *ValueTrie) Members() (members *vector.StringVector) {
	members = p.buildMembers(``, false, false)
	sort.Sort(members)
	return
}

// Retrieves all the members with their hyphenation values interspersed with the characters.
// The interspersal is optional in the case of zeroes.
func (p *ValueTrie) PatternMembers(includeZeroes bool) (members *vector.StringVector) {
	members = p.buildMembers(``, true, includeZeroes)
	sort.Sort(members)
	return
}

// Introspection -- counts all the nodes of the entire ValueTrie, NOT including the root node.
func (p *ValueTrie) Size() (sz int) {
	sz = len(p.children)

	for _, child := range p.children {
		sz += child.Size()
	}

	return
}

// Return the longest substring of `s` found in the trie.
func (p *ValueTrie) LongestSubstring(s string) (string, *vector.IntVector) {
	v := new(vector.IntVector)

	for pos, rune := range s {
		child, ok := p.children[rune]
		if !ok {
			// out of data in the trie -- return what we've got so far
			return s[0:pos], v
		}

		// append the value for this node
		if child.prefixValue != 0 {
			v.Push(child.prefixValue)
		}
		v.Push(child.value)
		p = child
	}

	// end of the string -- return input string and whatever values we accumulated
	return s, v
}

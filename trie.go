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
)

// Creates and returns a new Trie instance.
func NewTrie() *Trie {
	t := new(Trie)
	t.leaf = false
	t.children = make(map[int]*Trie)
	return t
}

// Internal function: adds items to the trie, reading runes from a strings.Reader
func (p *Trie) addRunes(r *strings.Reader) {
	rune, _, err := r.ReadRune()
	if err != nil {
		p.leaf = true
		return
	}

	n := p.children[rune]
	if n == nil {
		n = NewTrie()
		p.children[rune] = n
	}

	// recurse to store sub-runes below the new node
	n.addRunes(r)
}

// Adds a string to the trie. If the string is already present, no additional storage happens. Yay!
func (p *Trie) Add(s string) {
	if len(s) == 0 {
		return
	}

	// append the runes to the trie
	p.addRunes(strings.NewReader(s))
}

// Internal string removal function.  Returns trie if this node is empty following the removal.
func (p *Trie) removeRunes(r *strings.Reader) bool {
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
func (p *Trie) Remove(s string) bool {
	if len(s) == 0 {
		return len(p.children) == 0
	}

	// remove the runes, returning the final result
	return p.removeRunes(strings.NewReader(s))
}

// Internal string inclusion function.
func (p *Trie) includes(r *strings.Reader) bool {
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
func (p *Trie) Contains(s string) bool {
	if len(s) == 0 {
		return false // empty strings can't be included (how could we add them?)
	}
	return p.includes(strings.NewReader(s))
}

// Internal output-building function used by Members()
func (p *Trie) buildMembers(prefix string) *vector.StringVector {
	strList := new(vector.StringVector)

	if p.leaf {
		strList.Push(prefix)
	}

	// for each child, go grab all suffixes
	for rune, child := range p.children {
		buf := make([]byte, 4)
		numChars := utf8.EncodeRune(rune, buf)
		strList.AppendVector(child.buildMembers(prefix + string(buf[0:numChars])))
	}

	return strList
}

// Retrieves all member strings, in order.
func (p *Trie) Members() (members *vector.StringVector) {
	members = p.buildMembers(``)
	sort.Sort(members)
	return
}

// Introspection -- counts all the nodes of the entire Trie, NOT including the root node.
func (p *Trie) Size() (sz int) {
	sz = len(p.children)

	for _, child := range p.children {
		sz += child.Size()
	}

	return
}

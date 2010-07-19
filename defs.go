/*
 * defs.go
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

// The basic form of a Trie uses runes rather than characters, therefore it works on integer types.
type Trie struct {
	leaf     bool          // whether the node is a leaf (the end of an input string).
	children map[int]*Trie // a map of sub-tries for each child rune value.
}

// The second form stores a rune:integer pair.  This is used in the implementation of TeX hyphenation
// pattern tries.
type ValueTrie struct {
	value       int                // the value for the letter which indexed this node.
	prefixValue int                // some hyphenation strings *begin* with a numeric value. Le sigh.
	leaf        bool               // whether the node is a leaf (where an input string ended).
	children    map[int]*ValueTrie // a map of sub-tries for each child rune value.
}

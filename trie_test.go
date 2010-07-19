/*
 * trie_test.go
 * Trie
 *
 * Created by Jim Dovey on 17/07/2010.
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
	"testing"
	"container/vector"
)

func checkValues(trie *ValueTrie, sMatch, sCheck string, v *vector.IntVector, t *testing.T) {
	var values *vector.IntVector
	var testStr string = sMatch

	if len(sCheck) == 0 {
		var ok bool
		values, ok = trie.ValuesForString(sMatch)
		if !ok {
			t.Fatalf("No values returned for string '%s'", sMatch)
		}
	} else {
		var sub string
		testStr = sCheck
		sub, values = trie.LongestSubstring(sCheck)
		if sub != sMatch {
			t.Fatalf("Longest substring of '%s' should be '%s'", sCheck, sMatch)
		}
	}

	if values.Len() != v.Len() {
		t.Fatalf("Length mismatch: Values for '%s' should be %v, but got %v", testStr, *v, *values)
	}
	for i := 0; i < values.Len(); i++ {
		if values.At(i) != v.At(i) {
			t.Fatalf("Content mismatch: Values for '%s' should be %v, but got %v", testStr, *v, *values)
		}
	}
}


func TestTrie(t *testing.T) {
	trie := NewTrie()

	trie.Add("hello, world!")
	trie.Add("hello, there!")
	trie.Add("this is a sentence.")

	if !trie.Contains("hello, world!") {
		t.Error("trie should contain 'hello, world!'")
	}
	if !trie.Contains("hello, there!") {
		t.Error("trie should contain 'hello, there!'")
	}
	if !trie.Contains("this is a sentence.") {
		t.Error("trie should contain 'this is a sentence.'")
	}
	if trie.Contains("hello, Wisconsin!") {
		t.Error("trie should NOT contain 'hello, Wisconsin!'")
	}

	expectedSize := len("hello, ") + len("world!") + len("there!") + len("this is a sentence.")
	if trie.Size() != expectedSize {
		t.Errorf("trie should contain %d nodes", expectedSize)
	}

	// insert an existing string-- should be no change
	trie.Add("hello, world!")
	if trie.Size() != expectedSize {
		t.Errorf("trie should still contain only %d nodes after re-adding an existing member string", expectedSize)
	}

	// three strings in total
	if trie.Members().Len() != 3 {
		t.Error("trie should contain exactly three member strings")
	}

	// remove a string-- should reduce the size by the number of unique characters in that string
	trie.Remove("hello, world!")
	if trie.Contains("hello, world!") {
		t.Error("trie should no longer contain the string 'hello, world!'")
	}

	expectedSize -= len("world!")
	if trie.Size() != expectedSize {
		t.Errorf("trie should contain %d nodes after removing 'hello, world!'", expectedSize)
	}
}

func TestValueTrie(t *testing.T) {
	trie := NewValueTrie()

	str := "hyphenation"
	hyp := &vector.IntVector{0, 3, 0, 0, 2, 5, 4, 2, 0, 2, 0}

	hyphStr := "hy3phe2n5a4t2io2n"
	fullStr := "h0y3p0h0e2n5a4t2i0o2n0"

	// test addition using separate string and vector
	trie.Add(str, hyp)
	if !trie.Contains(str) {
		t.Error("value trie should contain the word 'hyphenation'")
	}

	if trie.Size() != len(str) {
		t.Errorf("value trie should have %d nodes (the number of characters in 'hyphenation')", len(str))
	}

	if trie.Members().Len() != 1 {
		t.Error("value trie should have only one member string")
	}

	shortPat := trie.PatternMembers(false)
	if shortPat.At(0) != hyphStr {
		t.Errorf("value trie should contain short pattern string '%s', but found '%s'", hyphStr, shortPat.At(0))
	}

	longPat := trie.PatternMembers(true)
	if longPat.At(0) != fullStr {
		t.Errorf("value trie should contain full pattern string '%s', but found '%s'", fullStr, longPat.At(0))
	}

	trie.Remove(str)
	if trie.Contains(str) {
		t.Errorf("value trie should no longer contain the word '%s'", str)
	}
	if trie.Size() != 0 {
		t.Error("value trie should have a node count of zero")
	}

	// test with an interspersed string of the form TeX's patterns use
	trie.AddPatternString(hyphStr)
	if !trie.Contains(str) {
		t.Errorf("value trie should now contain the word '%s'", str)
	}
	if trie.Size() != len(str) {
		t.Errorf("value trie should consist of %d nodes, instead has %d", len(str), trie.Size())
	}
	if trie.Members().Len() != 1 {
		t.Error("value trie should have only one member string")
	}

	mem := trie.Members()
	if mem.At(0) != str {
		t.Errorf("Expected first member string to be '%s', got '%s'", str, mem.At(0))
	}
	shortPat = trie.PatternMembers(false)
	if shortPat.At(0) != hyphStr {
		t.Errorf("value trie should contain short pattern string '%s', but found '%s'", hyphStr,
			shortPat.At(0))
	}
	longPat = trie.PatternMembers(true)
	if longPat.At(0) != fullStr {
		t.Errorf("value trie should contain full pattern string '%s', but found '%s'", fullStr,
			longPat.At(0))
	}

	checkValues(trie, `hyphenation`, ``, hyp, t)
	checkValues(trie, `hyphenation`, `hyphenationisagreatthing`, hyp, t)

	trie.Remove(`hyphenation`)
	if trie.Size() != 0 {
		t.Fail()
	}

	// test prefix values
	prefixedStr := `5emnix` // this is actually a string from the en_US TeX hyphenation trie
	purePrefixedStr := `emnix`
	values := &vector.IntVector{5, 0, 0, 0, 0, 0}
	trie.Add(purePrefixedStr, values)

	if trie.Size() != len(purePrefixedStr) {
		t.Errorf("Size of trie after adding '%s' should be %d, was %d", purePrefixedStr,
			len(purePrefixedStr), trie.Size())
	}

	checkValues(trie, `emnix`, ``, values, t)
	checkValues(trie, `emnix`, `emnixion`, values, t)

	trie.Remove(`emnix`)
	if trie.Size() != 0 {
		t.Fail()
	}

	trie.AddPatternString(prefixedStr)

	if trie.Size() != len(purePrefixedStr) {
		t.Errorf("Size of trie after adding '%s' should be %d, was %d", prefixedStr, len(purePrefixedStr),
			trie.Size())
	}
	
	checkValues(trie, `emnix`, ``, values, t)
	checkValues(trie, `emnix`, `emnixion`, values, t)
}

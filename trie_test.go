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
	"os"
	"io"
	"fmt"
	"scanner"
	"utf8"
)

func checkValues(trie *Trie, s string, v *vector.IntVector, t *testing.T) {
	value, ok := trie.GetValue(s)
	values := value.(*vector.IntVector)
	if !ok {
		t.Fatalf("No value returned for string '%s'", s)
	}

	if values.Len() != v.Len() {
		t.Fatalf("Length mismatch: Values for '%s' should be %v, but got %v", s, *v, *values)
	}
	for i := 0; i < values.Len(); i++ {
		if values.At(i) != v.At(i) {
			t.Fatalf("Content mismatch: Values for '%s' should be %v, but got %v", s, *v, *values)
		}
	}
}

func TestTrie(t *testing.T) {
	trie := NewTrie()

	trie.AddString("hello, world!")
	trie.AddString("hello, there!")
	trie.AddString("this is a sentence.")

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
	trie.AddString("hello, world!")
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

func TestMultiFind(t *testing.T) {
	trie := NewTrie()

	// these are part of the matches for the word 'hyphenation'
	trie.AddString(`hyph`)
	trie.AddString(`hen`)
	trie.AddString(`hena`)
	trie.AddString(`henat`)

	expected := new(vector.StringVector)
	expected.Push(`hyph`)
	found := trie.AllSubstrings(`hyphenation`)
	if found.Len() != expected.Len() {
		t.Errorf("expected %v but found %v", *expected, *found)
	}

	expected.Cut(0, expected.Len())
	expected.Push(`hen`)
	expected.Push(`hena`)
	expected.Push(`henat`)
	found = trie.AllSubstrings(`henation`)
	if found.Len() != expected.Len() {
		t.Errorf("expected %v but found %v", *expected, *found)
	}
}

///////////////////////////////////////////////////////////////
// Trie tests

func TestTrieValues(t *testing.T) {
	trie := NewTrie()

	str := "hyphenation"
	hyp := &vector.IntVector{0, 3, 0, 0, 2, 5, 4, 2, 0, 2, 0}

	hyphStr := "hy3phe2n5a4t2io2n"

	// test addition using separate string and vector
	trie.AddValue(str, hyp)
	if !trie.Contains(str) {
		t.Error("value trie should contain the word 'hyphenation'")
	}

	if trie.Size() != len(str) {
		t.Errorf("value trie should have %d nodes (the number of characters in 'hyphenation')", len(str))
	}

	if trie.Members().Len() != 1 {
		t.Error("value trie should have only one member string")
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

	checkValues(trie, `hyphenation`, hyp, t)

	trie.Remove(`hyphenation`)
	if trie.Size() != 0 {
		t.Fail()
	}

	// test prefix values
	prefixedStr := `5emnix` // this is actually a string from the en_US TeX hyphenation trie
	purePrefixedStr := `emnix`
	values := &vector.IntVector{5, 0, 0, 0, 0, 0}
	trie.AddValue(purePrefixedStr, values)

	if trie.Size() != len(purePrefixedStr) {
		t.Errorf("Size of trie after adding '%s' should be %d, was %d", purePrefixedStr,
			len(purePrefixedStr), trie.Size())
	}

	checkValues(trie, `emnix`, values, t)

	trie.Remove(`emnix`)
	if trie.Size() != 0 {
		t.Fail()
	}

	trie.AddPatternString(prefixedStr)

	if trie.Size() != len(purePrefixedStr) {
		t.Errorf("Size of trie after adding '%s' should be %d, was %d", prefixedStr, len(purePrefixedStr),
			trie.Size())
	}

	checkValues(trie, `emnix`, values, t)
}

func TestMultiFindValue(t *testing.T) {
	trie := NewTrie()

	// these are part of the matches for the word 'hyphenation'
	trie.AddPatternString(`hy3ph`)
	trie.AddPatternString(`he2n`)
	trie.AddPatternString(`hena4`)
	trie.AddPatternString(`hen5at`)

	v1 := &vector.IntVector{0, 3, 0, 0}
	v2 := &vector.IntVector{0, 2, 0}
	v3 := &vector.IntVector{0, 0, 0, 4}
	v4 := &vector.IntVector{0, 0, 5, 0, 0}

	expectStr := new(vector.StringVector)
	expectVal := new(vector.Vector) // contains elements of type *vector.IntVector

	expectStr.Push(`hyph`)
	expectVal.Push(v1)
	found, values := trie.AllSubstringsAndValues(`hyphenation`)
	if found.Len() != expectStr.Len() {
		t.Errorf("expected %v but found %v", *expectStr, *found)
	}
	if values.Len() != expectVal.Len() {
		t.Errorf("Length mismatch: expected %v but found %v", *expectVal, *values)
	}
	for i := 0; i < found.Len(); i++ {
		if found.At(i) != expectStr.At(i) {
			t.Errorf("Strings content mismatch: expected %v but found %v", *expectStr, *found)
			break
		}
	}
	for i := 0; i < values.Len(); i++ {
		ev := expectVal.At(i).(*vector.IntVector)
		fv := values.At(i).(*vector.IntVector)
		if ev.Len() != fv.Len() {
			t.Errorf("Value length mismatch: expected %v but found %v", *ev, *fv)
			break
		}
		for i := 0; i < ev.Len(); i++ {
			if ev.At(i) != fv.At(i) {
				t.Errorf("Value mismatch: expected %v but found %v", *ev, *fv)
				break
			}
		}
	}

	expectStr.Cut(0, expectStr.Len())
	expectVal.Cut(0, expectVal.Len())

	expectStr.AppendVector(&vector.StringVector{`hen`, `hena`, `henat`})
	expectVal.Push(v2)
	expectVal.Push(v3)
	expectVal.Push(v4)
	found, values = trie.AllSubstringsAndValues(`henation`)
	if found.Len() != expectStr.Len() {
		t.Errorf("expected %v but found %v", *expectStr, *found)
	}
	if values.Len() != expectVal.Len() {
		t.Errorf("Length mismatch: expected %v but found %v", *expectVal, *values)
	}
	for i := 0; i < found.Len(); i++ {
		if found.At(i) != expectStr.At(i) {
			t.Errorf("Strings content mismatch: expected %v but found %v", *expectStr, *found)
			break
		}
	}
	for i := 0; i < values.Len(); i++ {
		ev := expectVal.At(i).(*vector.IntVector)
		fv := values.At(i).(*vector.IntVector)
		if ev.Len() != fv.Len() {
			t.Errorf("Value length mismatch: expected %v but found %v", *ev, *fv)
			break
		}
		for i := 0; i < ev.Len(); i++ {
			if ev.At(i) != fv.At(i) {
				t.Errorf("Value mismatch: expected %v but found %v", *ev, *fv)
				break
			}
		}
	}
}

//////////////////////////////////////////////////////////////////
// Benchmarks
// Run like so:
//   cat patterns-en.go | gotest -benchmarks=".*"
// This is because, for some unknown reason, os.Open() always returns 'resource temporarily unavailable'.

func loadPatterns(reader io.Reader) (*Trie, os.Error) {
	trie := NewTrie()
	var s scanner.Scanner
	s.Init(reader)
	s.Mode = scanner.ScanIdents | scanner.ScanRawStrings | scanner.SkipComments

	var which string

	tok := s.Scan()
	for tok != scanner.EOF {
		switch tok {
		case scanner.Ident:
			// we handle two identifiers: 'patterns' and 'exceptions'
			switch ident := s.TokenText(); ident {
			case `patterns`, `exceptions`:
				which = ident
			default:
				return nil, os.ErrorString(fmt.Sprintf("Unrecognized identifier '%s' at position %v",
					ident, s.Pos()))
			}
		case scanner.String, scanner.RawString:
			// trim the quotes from around the string
			tokstr := s.TokenText()
			str := tokstr[1 : len(tokstr)-1]

			switch which {
			case `patterns`:
				trie.AddPatternString(str)
			}
		}
		tok = s.Scan()
	}

	return trie, nil
}

var benchmarkTrie *Trie = nil

func setupTrie() *Trie {
	/*
		filename := "patterns-en.go"
		f, err := os.Open(filename, 0444, os.O_RDONLY)
		if err != nil {
			fmt.Printf("Failed to open file '%s': %s\n", filename, err)
		}
	*/
	if benchmarkTrie == nil {
		var err os.Error
		benchmarkTrie, err = loadPatterns(os.Stdin)
		if err != nil {
			fmt.Printf("Failed to load patterns from Stdin: %s\n", err)
		}
	}
	return benchmarkTrie
}

func BenchmarkTraversal(b *testing.B) {
	b.StopTimer()
	trie := setupTrie()
	if trie == nil {
		return
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		trie.Members()
	}
}

func BenchmarkHyphenation(b *testing.B) {
	b.StopTimer()
	trie := setupTrie()
	if trie == nil {
		return
	}
	testStr := `.hyphenation.`
	v := make([]int, utf8.RuneCountInString(testStr))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for i := 0; i < len(v); i++ {
			v[i] = 0
		}
		vIndex := 0
		for pos, _ := range testStr {
			t := testStr[pos:]
			strs, values := trie.AllSubstringsAndValues(t)
			for i := 0; i < values.Len(); i++ {
				str := strs.At(i)
				val := values.At(i).(*vector.IntVector)

				diff := val.Len() - len(str)
				vs := v[vIndex-diff:]

				for i := 0; i < val.Len(); i++ {
					if val.At(i) > vs[i] {
						vs[i] = val.At(i)
					}
				}
			}
			vIndex++
		}
	}
}

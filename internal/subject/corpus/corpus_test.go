// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package corpus_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/subject/corpus"

	"github.com/c4-project/c4t/internal/subject"
)

// ExampleCorpus_Add is a runnable example for Add.
func ExampleCorpus_Add() {
	c := make(corpus.Corpus)
	_ = c.Add(*subject.NewOrPanic(litmus.NewOrPanic("bar/baz.litmus")).AddName("foo"))
	fmt.Println(c["foo"].Source.Path)

	// We can't add duplicates to a corpus.
	err := c.Add(*subject.NewOrPanic(litmus.NewOrPanic("bar/baz2.litmus")).AddName("foo"))
	fmt.Println(err)

	// Output:
	// bar/baz.litmus
	// duplicate corpus entry: foo
}

// ExampleCorpus_FilterToNames is a runnable example for Corpus.FilterToNames.
func ExampleCorpus_FilterToNames() {
	c := corpus.Mock()

	for _, n := range c.Names() {
		fmt.Println(n, "is in c")
	}

	c2 := c.FilterToNames("foo", "bar")
	for _, n := range c2.Names() {
		fmt.Println(n, "is in c2")
	}

	c3 := c.FilterToNames()
	for _, n := range c3.Names() {
		fmt.Println(n, "is in c3")
	}

	// Output:
	// bar is in c
	// barbaz is in c
	// baz is in c
	// foo is in c
	// bar is in c2
	// foo is in c2
}

// TestCorpus_Copy tests that Corpus.Copy performs a sufficiently deep copy of the corpus.
func TestCorpus_Copy(t *testing.T) {
	c := corpus.Mock()
	cc := c.Copy()

	for n := range c {
		cs, ok := cc[n]
		if !ok {
			t.Errorf("subject %s disappeared in copy", n)
			continue
		}
		if !reflect.DeepEqual(c[n], cs) {
			t.Errorf("subject %s changed in copy: got=%v, want=%v", n, c[n], cs)
			continue
		}

		c[n] = subject.Subject{}
		if reflect.DeepEqual(c[n], cs) {
			t.Errorf("assignment to %s changed copied version (copy shallow?)", n)
		}
	}

	for n := range cc {
		assert.Containsf(t, c, n, "subject %s appeared out of nowhere in copy", n)
	}
}

// TestCorpus_EraseCompilations tests that Corpus.EraseCompilations works properly.
func TestCorpus_EraseCompilations(t *testing.T) {
	c := corpus.Mock()
	c.EraseCompilations()

	for _, s := range c {
		assert.Empty(t, s.Compilations, "compilation should be empty")
	}
}

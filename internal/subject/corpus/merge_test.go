// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package corpus_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/c4t/internal/helper/testhelp"

	"github.com/MattWindsor91/c4t/internal/subject/corpus"
)

// ExampleMerge is a testable example for Merge.
func ExampleMerge() {
	c1 := corpus.New("foo", "bar", "baz")
	c2 := corpus.New("foo", "bar", "baz")

	// Empty merge
	c3, _ := corpus.Merge(map[string]corpus.Corpus{})
	for n := range c3 {
		fmt.Println("empty>", n)
	}

	// Single merge
	c4, _ := corpus.Merge(map[string]corpus.Corpus{"c1": c1})
	for n := range c4 {
		fmt.Println("single>", n)
	}

	// Multi merge
	c5, _ := corpus.Merge(map[string]corpus.Corpus{"c1": c1, "c2": c2})
	for n := range c5 {
		fmt.Println("multi>", n)
	}

	// Unordered output:
	// single> foo
	// single> bar
	// single> baz
	// multi> c1/foo
	// multi> c1/bar
	// multi> c1/baz
	// multi> c2/foo
	// multi> c2/bar
	// multi> c2/baz
}

// TestMerge_errors tests error cases in Merge.
func TestMerge_errors(t *testing.T) {
	t.Parallel()

	// There is only one real error case so far: failure to Add

	c1 := corpus.New("foo")
	c2 := corpus.New("a/foo")

	_, err := corpus.Merge(map[string]corpus.Corpus{"b/a": c1, "b": c2})
	testhelp.ExpectErrorIs(t, err, corpus.ErrCorpusDup, "merge")
}

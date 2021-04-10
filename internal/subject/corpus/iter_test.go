// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package corpus_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"
	"github.com/c4-project/c4t/internal/subject"
	"github.com/c4-project/c4t/internal/subject/corpus"
)

// ExampleCorpus_Each is a runnable example for Each.
func ExampleCorpus_Each() {
	// A perhaps less efficient version of c.Names():

	c := corpus.New("foo", "bar", "baz", "barbaz")
	names := make([]string, 0, len(c))

	_ = c.Each(func(s subject.Named) error {
		names = append(names, s.Name)
		return nil
	})

	sort.Strings(names)
	for _, n := range names {
		fmt.Println(n)
	}

	// Output:
	// bar
	// barbaz
	// baz
	// foo
}

// TestCorpus_Map tests the Map function on a basic exercise.
func TestCorpus_Map(t *testing.T) {
	c := corpus.New("foo", "bar", "baz", "barbaz")
	err := c.Map(func(s *subject.Named) error {
		s.Source.Path = s.Name + ".litmus"
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error in Map: %v", err)
	}

	// Each subject should've been updated according to the function.
	for n, s := range c {
		got := s.Source
		want := n + ".litmus"

		if got.Path != want {
			t.Errorf("Map set source path incorrectly: got=%s; want=%s", got.Path, want)
		}
	}
}

// TestCorpus_Map_rename makes sure Map fails if there is an attempt to rename a subject.
func TestCorpus_Map_rename(t *testing.T) {
	c := corpus.New("foo", "bar", "baz", "barbaz")
	err := c.Map(func(s *subject.Named) error {
		s.Name += "2"
		return nil
	})
	testhelp.ExpectErrorIs(t, err, corpus.ErrMapRename, "renaming in a Map")
}

// TestCorpus_Map_error makes sure Map fails if there is an error inside an invocation.
func TestCorpus_Map_error(t *testing.T) {
	e := errors.New("test error")

	c := corpus.New("foo", "bar", "baz", "barbaz")
	err := c.Map(func(s *subject.Named) error {
		return e
	})
	testhelp.ExpectErrorIs(t, err, e, "Map of function returning error")
}

func makeHugeCorpus() corpus.Corpus {
	names := make([]string, 640)
	for i := range names {
		names[i] = strconv.Itoa(i)
	}
	return corpus.New(names...)
}

// TestCorpus_Par tests the 'happy path' of Par across various sizes of corpus.
func TestCorpus_Par(t *testing.T) {
	cases := map[string]struct {
		n      int
		corpus corpus.Corpus
	}{
		"empty":        {n: 10, corpus: corpus.Corpus{}},
		"empty-single": {n: 1, corpus: corpus.Corpus{}},
		"small":        {n: 10, corpus: corpus.New("foo", "bar", "baz")},
		"small-single": {n: 1, corpus: corpus.New("foo", "bar", "baz")},
		"large":        {n: 10, corpus: makeHugeCorpus()},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var sm sync.Map
			if err := c.corpus.Par(context.Background(), c.n, func(_ context.Context, s subject.Named) error {
				sm.Store(s.Name, true)
				return nil
			}); err != nil {
				t.Errorf("par failed: %v", err)
			}

			for n := range c.corpus {
				if _, ok := sm.Load(n); !ok {
					t.Errorf("par didn't store %s", n)
				}
			}
		})
	}
}

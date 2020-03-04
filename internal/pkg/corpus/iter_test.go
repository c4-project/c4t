// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
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

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"
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
		s.Litmus = s.Name + ".litmus"
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error in Map: %v", err)
	}

	// Each subject should've been updated according to the function.
	for n, s := range c {
		got := s.Litmus
		want := n + ".litmus"

		if got != want {
			t.Errorf("Map set Litmus incorrectly: got=%s; want=%s", got, want)
		}
	}
}

// TestCorpus_Map_Rename makes sure Map fails if there is an attempt to rename a subject.
func TestCorpus_Map_Rename(t *testing.T) {
	c := corpus.New("foo", "bar", "baz", "barbaz")
	err := c.Map(func(s *subject.Named) error {
		s.Name = s.Name + "2"
		return nil
	})
	testhelp.ExpectErrorIs(t, err, corpus.ErrMapRename, "renaming in a Map")
}

// TestCorpus_Map_Error makes sure Map fails if there is an error inside an invocation.
func TestCorpus_Map_Error(t *testing.T) {
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
	cases := map[string]corpus.Corpus{
		"empty": {},
		"small": corpus.New("foo", "bar", "baz"),
		"large": makeHugeCorpus(),
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var sm sync.Map
			if err := c.Par(context.Background(), 10, func(_ context.Context, s subject.Named) error {
				sm.Store(s.Name, true)
				return nil
			}); err != nil {
				t.Errorf("par failed: %v", err)
			}

			for n := range c {
				if _, ok := sm.Load(n); !ok {
					t.Errorf("par didn't store %s", n)
				}
			}
		})
	}
}

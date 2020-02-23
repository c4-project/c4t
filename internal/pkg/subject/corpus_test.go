package subject_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// ExampleCorpus_Add is a runnable example for Add.
func ExampleCorpus_Add() {
	c := make(subject.Corpus)
	_ = c.Add(subject.Named{Name: "foo", Subject: subject.Subject{Litmus: "bar/baz.litmus"}})
	fmt.Println(c["foo"].Litmus)

	// We can't add duplicates to a corpus.
	err := c.Add(subject.Named{Name: "foo", Subject: subject.Subject{Litmus: "bar/baz2.litmus"}})
	fmt.Println(err)

	// Output:
	// bar/baz.litmus
	// duplicate corpus entry: foo
}

// smallCorpi contains test cases for the overly-small-corpus error handling of SampleCorpus.
var smallCorpi = map[string]struct {
	corpus subject.Corpus
	want   int
}{
	"empty+cap": {subject.Corpus{}, 10},
	"empty":     {subject.Corpus{}, 0},
	"small1":    {subject.NewCorpus("foo"), 2},
	"small2":    {subject.NewCorpus("foo", "bar", "baz"), 10},
}

// exactCorpi contains test cases for the pass-through behaviour of SampleCorpus.
var exactCorpi = map[string]struct {
	corpus subject.Corpus
	want   int
}{
	"nocap1": {subject.NewCorpus("foo"), 0},
	"nocap2": {subject.NewCorpus("foo", "bar"), 0},
	"nocap3": {subject.NewCorpus("foo", "bar", "baz"), 0},
	"nocap4": {subject.NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 0},
	"cap1":   {subject.NewCorpus("foo"), 1},
	"cap2":   {subject.NewCorpus("foo", "bar"), 2},
	"cap3":   {subject.NewCorpus("foo", "bar", "baz"), 3},
	"cap4":   {subject.NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 7},
}

// sampleCorpi contains test cases for the 'actually sample' behaviour of SampleCorpus.
var sampleCorpi = map[string]struct {
	corpus subject.Corpus
	want   int
}{
	"sample1": {subject.NewCorpus("foo", "bar"), 1},
	"sample2": {subject.NewCorpus("foo", "bar", "baz"), 1},
	"sample3": {subject.NewCorpus("foo", "bar", "baz"), 2},
	"sample4": {subject.NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 3},
	"sample5": {subject.NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 5},
}

// TestCorpus_Sample_SmallErrors tests that various 'overly small corpus' situations produce an error.
func TestCorpus_Sample_SmallErrors(t *testing.T) {
	for name, c := range smallCorpi {
		t.Run(name, func(t *testing.T) {
			_, err := c.corpus.Sample(1, c.want)
			if err == nil {
				t.Errorf("no error when sampling small corpus (%v, want %d)", c.corpus, c.want)
			} else if !errors.Is(err, subject.ErrSmallCorpus) {
				t.Errorf("wrong error when sampling small corpus (%v, want %d): %v", c.corpus, c.want, err)
			}
		})
	}
}

// TestCorpus_Sample_PassThrough tests that various cases that shouldn't cause sampling don't.
func TestCorpus_Sample_PassThrough(t *testing.T) {
	for name, c := range exactCorpi {
		t.Run(name, func(t *testing.T) {
			smp, err := c.corpus.Sample(1, c.want)
			if err != nil {
				t.Errorf("error when sampling exact corpus (%v, want %d): %v", c.corpus, c.want, err)
			} else if !reflect.DeepEqual(smp, c.corpus) {
				t.Errorf("bad sample of exact corpus (%v, want %d): %v", c.corpus, c.want, smp)
			}
		})
	}
}

// TestCorpus_Sample_ActuallySample tests that sampling behaves correctly.
func TestCorpus_Sample_ActuallySample(t *testing.T) {
	i := 0
	for name, c := range sampleCorpi {
		t.Run(name, func(t *testing.T) {
			smp, err := c.corpus.Sample(int64(i), c.want)
			if err != nil {
				t.Errorf("error when sampling corpus (%v, want %d): %v", c.corpus, c.want, err)
			} else {
				checkCorpusIsSample(t, c.corpus, smp)
			}
		})
		i++
	}
}

func checkCorpusIsSample(t *testing.T, corpus, smp subject.Corpus) {
	// Each item in the sample should be in the corpus.
	for k, got := range smp {
		want, ok := corpus[k]
		if !ok {
			t.Helper()
			t.Fatalf("sample of %v (%v) contains unexpected key: %q", corpus, smp, k)

		}
		if !reflect.DeepEqual(got, want) {
			t.Helper()
			t.Fatalf("sample of %v (%v) maps %q to %v; want %v", corpus, smp, k, got, want)
		}
	}
}

// TestCorpus_Map tests the Map function on a basic exercise.
func TestCorpus_Map(t *testing.T) {
	c := subject.NewCorpus("foo", "bar", "baz", "barbaz")
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

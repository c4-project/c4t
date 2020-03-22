// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus_test

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"
)

// smallCorpi contains test cases for the overly-small-corpus error handling of SampleCorpus.
var smallCorpi = map[string]struct {
	corpus corpus.Corpus
	want   int
}{
	"empty+cap": {corpus.Corpus{}, 10},
	"empty":     {corpus.Corpus{}, 0},
	"small1":    {corpus.New("foo"), 2},
	"small2":    {corpus.New("foo", "bar", "baz"), 10},
}

// exactCorpi contains test cases for the pass-through behaviour of SampleCorpus.
var exactCorpi = map[string]struct {
	corpus corpus.Corpus
	want   int
}{
	"nocap1": {corpus.New("foo"), 0},
	"nocap2": {corpus.New("foo", "bar"), 0},
	"nocap3": {corpus.New("foo", "bar", "baz"), 0},
	"nocap4": {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 0},
	"cap1":   {corpus.New("foo"), 1},
	"cap2":   {corpus.New("foo", "bar"), 2},
	"cap3":   {corpus.New("foo", "bar", "baz"), 3},
	"cap4":   {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 7},
}

// sampleCorpi contains test cases for the 'actually sample' behaviour of SampleCorpus.
var sampleCorpi = map[string]struct {
	corpus corpus.Corpus
	want   int
}{
	"sample1": {corpus.New("foo", "bar"), 1},
	"sample2": {corpus.New("foo", "bar", "baz"), 1},
	"sample3": {corpus.New("foo", "bar", "baz"), 2},
	"sample4": {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 3},
	"sample5": {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 5},
}

// TestCorpus_Sample_SmallErrors tests that various 'overly small corpus' situations produce an error.
func TestCorpus_Sample_SmallErrors(t *testing.T) {
	for name, c := range smallCorpi {
		t.Run(name, func(t *testing.T) {
			_, err := c.corpus.Sample(rand.New(rand.NewSource(1)), c.want)
			if err == nil {
				t.Errorf("no error when sampling small corpus (%v, want %d)", c.corpus, c.want)
			} else if !errors.Is(err, corpus.ErrSmall) {
				t.Errorf("wrong error when sampling small corpus (%v, want %d): %v", c.corpus, c.want, err)
			}
		})
	}
}

// TestCorpus_Sample_PassThrough tests that various cases that shouldn't cause sampling don't.
func TestCorpus_Sample_PassThrough(t *testing.T) {
	for name, c := range exactCorpi {
		t.Run(name, func(t *testing.T) {
			smp, err := c.corpus.Sample(rand.New(rand.NewSource(1)), c.want)
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
	var i int64
	for name, c := range sampleCorpi {
		t.Run(name, func(t *testing.T) {
			smp, err := c.corpus.Sample(rand.New(rand.NewSource(i)), c.want)
			if err != nil {
				t.Errorf("error when sampling corpus (%v, want %d): %v", c.corpus, c.want, err)
			} else {
				checkCorpusIsSample(t, c.corpus, smp)
			}
		})
		i++
	}
}

func checkCorpusIsSample(t *testing.T, corpus, smp corpus.Corpus) {
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

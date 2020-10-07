// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus"
)

// emptyCorpora contains test cases for the overly-small-corpus error handling of SampleCorpus.
var emptyCorpora = map[string]int{
	"empty+cap":     10,
	"empty":         0,
	"empty+bad-cap": -10,
}

// exactCorpora contains test cases for the pass-through behaviour of SampleCorpus.
var exactCorpora = map[string]struct {
	corpus corpus.Corpus
	want   int
}{
	"nocap1":   {corpus.New("foo"), 0},
	"nocap2":   {corpus.New("foo", "bar"), 0},
	"nocap3":   {corpus.New("foo", "bar", "baz"), 0},
	"nocap4":   {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 0},
	"cap1":     {corpus.New("foo"), 1},
	"cap2":     {corpus.New("foo", "bar"), 2},
	"cap3":     {corpus.New("foo", "bar", "baz"), 3},
	"cap4":     {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 7},
	"overcap1": {corpus.New("foo"), 9},
	"overcap2": {corpus.New("foo", "bar"), 9},
	"overcap3": {corpus.New("foo", "bar", "baz"), 9},
	"overcap4": {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 9},
}

// sampleCorpora contains test cases for the 'actually sample' behaviour of SampleCorpus.
var sampleCorpora = map[string]struct {
	corpus corpus.Corpus
	want   int
}{
	"sample1": {corpus.New("foo", "bar"), 1},
	"sample2": {corpus.New("foo", "bar", "baz"), 1},
	"sample3": {corpus.New("foo", "bar", "baz"), 2},
	"sample4": {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 3},
	"sample5": {corpus.New("you're", "going", "to", "have", "a", "bad", "time"), 5},
}

// TestCorpus_Sample_emptyErrors tests that empty corpus situations produce an error.
func TestCorpus_Sample_emptyErrors(t *testing.T) {
	t.Parallel()
	for name, want := range emptyCorpora {
		want := want
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := corpus.Corpus{}.Sample(rand.New(rand.NewSource(1)), want)
			testhelp.ExpectErrorIs(t, err, corpus.ErrNone, "sampling empty corpus")
		})
	}
}

// TestCorpus_Sample_passThrough tests that various cases that shouldn't cause sampling don't.
func TestCorpus_Sample_passThrough(t *testing.T) {
	t.Parallel()
	for name, c := range exactCorpora {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			smp, err := c.corpus.Sample(rand.New(rand.NewSource(1)), c.want)
			if assert.NoError(t, err, "error when sampling exact corpus") {
				assert.Equal(t, c.corpus, smp, "should pass through exact corpus")
			}
		})
	}
}

// TestCorpus_Sample_actuallySample tests that sampling behaves correctly.
func TestCorpus_Sample_actuallySample(t *testing.T) {
	t.Parallel()
	var i int64
	for name, c := range sampleCorpora {
		c := c
		j := i
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			smp, err := c.corpus.Sample(rand.New(rand.NewSource(j)), c.want)
			if assert.NoError(t, err, "error when sampling corpus") {
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

package model

import (
	"errors"
	"reflect"
	"testing"
)

// smallCorpi contains test cases for the overly-small-corpus error handling of SampleCorpus.
var smallCorpi = []struct {
	corpus Corpus
	want   int
}{
	// Empty corpus
	{Corpus{}, 10},
	{Corpus{}, 0},
	// Small corpus
	{NewCorpus("foo"), 2},
	{NewCorpus("foo", "bar", "baz"), 10},
}

// exactCorpi contains test cases for the pass-through behaviour of SampleCorpus.
var exactCorpi = []struct {
	corpus Corpus
	want   int
}{
	// No sampling requested
	{NewCorpus("foo"), 0},
	{NewCorpus("foo", "bar"), 0},
	{NewCorpus("foo", "bar", "baz"), 0},
	{NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 0},
	// Sample size matches input size
	{NewCorpus("foo"), 1},
	{NewCorpus("foo", "bar"), 2},
	{NewCorpus("foo", "bar", "baz"), 3},
	{NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 7},
}

// sampleCorpi contains test cases for the 'actually sample' behaviour of SampleCorpus.
var sampleCorpi = []struct {
	corpus Corpus
	want   int
}{
	{NewCorpus("foo", "bar"), 1},
	{NewCorpus("foo", "bar", "baz"), 1},
	{NewCorpus("foo", "bar", "baz"), 2},
	{NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 3},
	{NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 5},
}

// TestSampleCorpus_SmallErrors tests that various 'overly small corpus' situations produce an error.
func TestSampleCorpus_SmallErrors(t *testing.T) {
	for _, c := range smallCorpi {
		_, err := c.corpus.Sample(1, c.want)
		if err == nil {
			t.Errorf("no error when sampling small corpus (%v, want %d)", c.corpus, c.want)
		} else if !errors.Is(err, ErrSmallCorpus) {
			t.Errorf("wrong error when sampling small corpus (%v, want %d): %v", c.corpus, c.want, err)
		}
	}
}

// TestSampleCorpus_PassThrough tests that various cases that shouldn't cause sampling don't.
func TestSampleCorpus_PassThrough(t *testing.T) {
	for _, c := range exactCorpi {
		smp, err := c.corpus.Sample(1, c.want)
		if err != nil {
			t.Errorf("error when sampling exact corpus (%v, want %d): %v", c.corpus, c.want, err)
		} else if !reflect.DeepEqual(smp, c.corpus) {
			t.Errorf("bad sample of exact corpus (%v, want %d): %v", c.corpus, c.want, smp)
		}
	}
}

// TestSampleCorpus_ActuallySample tests that sampling behaves correctly.
func TestSampleCorpus_ActuallySample(t *testing.T) {
	for i, c := range sampleCorpi {
		smp, err := c.corpus.Sample(int64(i), c.want)
		if err != nil {
			t.Errorf("error when sampling corpus (%v, want %d): %v", c.corpus, c.want, err)
		} else {
			// The sample should contain the items in ascending order of index in the original corpus.
			// To test this, we do a slightly convoluted lock-step search, where we slowly sweep over the corpus trying
			// to find each sample in turn.

			j := 0
		SampleLoop:
			for _, s := range smp {
				for ; j < len(c.corpus); j++ {
					if c.corpus[j].Litmus == s.Litmus {
						continue SampleLoop
					}
				}
				t.Fatalf("sample of %v (%v) contains bad or ill-positioned item %q", c.corpus, smp, s)
			}
		}
	}
}

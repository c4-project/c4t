package planner

import (
	"errors"
	"reflect"
	"testing"
)

// smallCorpi contains test cases for the overly-small-corpus error handling of SampleCorpus.
var smallCorpi = []struct {
	corpus []string
	want   int
}{
	// Empty corpus
	{[]string{}, 10},
	{[]string{}, 0},
	// Small corpus
	{[]string{"foo"}, 2},
	{[]string{"foo", "bar", "baz"}, 10},
}

// exactCorpi contains test cases for the pass-through behaviour of SampleCorpus.
var exactCorpi = []struct {
	corpus []string
	want   int
}{
	// No sampling requested
	{[]string{"foo"}, 0},
	{[]string{"foo", "bar"}, 0},
	{[]string{"foo", "bar", "baz"}, 0},
	{[]string{"you're", "going", "to", "have", "a", "bad", "time"}, 0},
	// Sample size matches input size
	{[]string{"foo"}, 1},
	{[]string{"foo", "bar"}, 2},
	{[]string{"foo", "bar", "baz"}, 3},
	{[]string{"you're", "going", "to", "have", "a", "bad", "time"}, 7},
}

// sampleCorpi contains test cases for the 'actually sample' behaviour of SampleCorpus.
var sampleCorpi = []struct {
	corpus []string
	want   int
}{
	{[]string{"foo", "bar"}, 1},
	{[]string{"foo", "bar", "baz"}, 1},
	{[]string{"foo", "bar", "baz"}, 2},
	{[]string{"you're", "going", "to", "have", "a", "bad", "time"}, 3},
	{[]string{"you're", "going", "to", "have", "a", "bad", "time"}, 5},
}

// TestSampleCorpus_SmallErrors tests that various 'overly small corpus' situations produce an error.
func TestSampleCorpus_SmallErrors(t *testing.T) {
	for _, c := range smallCorpi {
		_, err := SampleCorpus(1, c.corpus, c.want)
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
		smp, err := SampleCorpus(1, c.corpus, c.want)
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
		smp, err := SampleCorpus(int64(i), c.corpus, c.want)
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
					if c.corpus[j] == s {
						continue SampleLoop
					}
				}
				t.Fatalf("sample of %v (%v) contains bad or ill-positioned item %q", c.corpus, smp, s)
			}
		}
	}
}

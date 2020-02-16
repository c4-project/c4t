package subject

import (
	"errors"
	"reflect"
	"testing"
)

// smallCorpi contains test cases for the overly-small-corpus error handling of SampleCorpus.
var smallCorpi = []struct {
	name   string
	corpus Corpus
	want   int
}{
	{"empty+cap", Corpus{}, 10},
	{"empty", Corpus{}, 0},
	{"small1", NewCorpus("foo"), 2},
	{"small2", NewCorpus("foo", "bar", "baz"), 10},
}

// exactCorpi contains test cases for the pass-through behaviour of SampleCorpus.
var exactCorpi = []struct {
	name   string
	corpus Corpus
	want   int
}{
	{"nocap1", NewCorpus("foo"), 0},
	{"nocap2", NewCorpus("foo", "bar"), 0},
	{"nocap3", NewCorpus("foo", "bar", "baz"), 0},
	{"nocap4", NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 0},
	{"cap1", NewCorpus("foo"), 1},
	{"cap2", NewCorpus("foo", "bar"), 2},
	{"cap3", NewCorpus("foo", "bar", "baz"), 3},
	{"cap4", NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 7},
}

// sampleCorpi contains test cases for the 'actually sample' behaviour of SampleCorpus.
var sampleCorpi = []struct {
	name   string
	corpus Corpus
	want   int
}{
	{"sample1", NewCorpus("foo", "bar"), 1},
	{"sample2", NewCorpus("foo", "bar", "baz"), 1},
	{"sample3", NewCorpus("foo", "bar", "baz"), 2},
	{"sample4", NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 3},
	{"sample5", NewCorpus("you're", "going", "to", "have", "a", "bad", "time"), 5},
}

// TestCorpus_Sample_SmallErrors tests that various 'overly small corpus' situations produce an error.
func TestCorpus_Sample_SmallErrors(t *testing.T) {
	for _, c := range smallCorpi {
		t.Run(c.name, func(t *testing.T) {
			_, err := c.corpus.Sample(1, c.want)
			if err == nil {
				t.Errorf("no error when sampling small corpus (%v, want %d)", c.corpus, c.want)
			} else if !errors.Is(err, ErrSmallCorpus) {
				t.Errorf("wrong error when sampling small corpus (%v, want %d): %v", c.corpus, c.want, err)
			}
		})
	}
}

// TestCorpus_Sample_PassThrough tests that various cases that shouldn't cause sampling don't.
func TestCorpus_Sample_PassThrough(t *testing.T) {
	for _, c := range exactCorpi {
		t.Run(c.name, func(t *testing.T) {
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
	for i, c := range sampleCorpi {
		t.Run(c.name, func(t *testing.T) {
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
					t.Fatalf("sample of %v (%v) contains bad or ill-positioned item: %v", c.corpus, smp, s)
				}
			}
		})
	}
}

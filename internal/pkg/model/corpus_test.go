package model

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

// TestSampleCorpus_SmallErrors tests that various 'overly small corpus' situations produce an error.
func TestSampleCorpus_SmallErrors(t *testing.T) {
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

// TestSampleCorpus_PassThrough tests that various cases that shouldn't cause sampling don't.
func TestSampleCorpus_PassThrough(t *testing.T) {
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

// TestSampleCorpus_ActuallySample tests that sampling behaves correctly.
func TestSampleCorpus_ActuallySample(t *testing.T) {
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
					t.Fatalf("sample of %v (%v) contains bad or ill-positioned item %q", c.corpus, smp, s)
				}
			}
		})
	}
}

// chunkCases contains test cases for Chunks.
var chunkCases = []struct {
	name    string
	corpus  Corpus
	nchunks int
	want    []Corpus
}{
	{
		name:    "nil-negative",
		corpus:  nil,
		nchunks: -1,
		want:    []Corpus{nil},
	},
	{
		name:    "nil-0",
		corpus:  nil,
		nchunks: 0,
		want:    []Corpus{nil},
	},
	{
		name:    "nil-1",
		corpus:  nil,
		nchunks: 1,
		want:    []Corpus{nil},
	},
	{
		name:    "nil-over",
		corpus:  nil,
		nchunks: 2,
		want:    []Corpus{nil, nil},
	},
	{
		name:    "empty-negative",
		corpus:  NewCorpus(),
		nchunks: -1,
		want:    []Corpus{NewCorpus()},
	},
	{
		name:    "empty-0",
		corpus:  NewCorpus(),
		nchunks: 0,
		want:    []Corpus{NewCorpus()},
	},
	{
		name:    "empty-1",
		corpus:  NewCorpus(),
		nchunks: 1,
		want:    []Corpus{NewCorpus()},
	},
	{
		name:    "empty-over",
		corpus:  NewCorpus(),
		nchunks: 2,
		want:    []Corpus{NewCorpus(), nil},
	},
	{
		name:    "even-negative",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: -1,
		want:    []Corpus{NewCorpus("foo", "bar", "baz", "barbaz")},
	},
	{
		name:    "even-0",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: 0,
		want:    []Corpus{NewCorpus("foo", "bar", "baz", "barbaz")},
	},
	{
		name:    "even-1",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: 1,
		want:    []Corpus{NewCorpus("foo", "bar", "baz", "barbaz")},
	},
	{
		name:    "even-div",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: 2,
		want:    []Corpus{NewCorpus("foo", "bar"), NewCorpus("baz", "barbaz")},
	},
	{
		name:    "even-single",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: 4,
		want:    []Corpus{NewCorpus("foo"), NewCorpus("bar"), NewCorpus("baz"), NewCorpus("barbaz")},
	},
	{
		name:    "even-over",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: 5,
		want:    []Corpus{NewCorpus("foo"), NewCorpus("bar"), NewCorpus("baz"), NewCorpus("barbaz"), nil},
	},
	{
		name:    "even-misdiv",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz"),
		nchunks: 3,
		want:    []Corpus{NewCorpus("foo", "bar"), NewCorpus("baz"), NewCorpus("barbaz")},
	},
	{
		name:    "prime-misdiv",
		corpus:  NewCorpus("foo", "bar", "baz", "barbaz", "foobar"),
		nchunks: 3,
		want:    []Corpus{NewCorpus("foo", "bar"), NewCorpus("baz", "barbaz"), NewCorpus("foobar")},
	},
}

func TestCorpus_Chunks(t *testing.T) {
	for _, c := range chunkCases {
		t.Run(c.name, func(t *testing.T) {
			got := c.corpus.Chunks(c.nchunks)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("chunks of %s=%v; want %v", c.name, got, c.want)
			}
		})
	}
}

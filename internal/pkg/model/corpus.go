package model

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

var (
	// ErrSmallCorpus occurs when the viable test corpus is smaller than that requested by the user.
	ErrSmallCorpus = errors.New("test corpus too small")

	// ErrNoCorpus is a variant of ErrSmallCorpus that occurs when the viable test corpus is empty.
	ErrNoCorpus = fmt.Errorf("%w: no corpus given", ErrSmallCorpus)
)

// Corpus is the type of test corpi (groups of test subjects).
type Corpus []Subject

// NewCorpus creates a Corpus from a list of input Litmus files.
func NewCorpus(infiles ...string) Corpus {
	corpus := make(Corpus, len(infiles))
	for i, f := range infiles {
		corpus[i] = Subject{Litmus: f}
	}
	return corpus
}

// SampleCorpus tries to select a sample of size want from corpus-slice corpus.
// If want is non-positive, no sampling occurs.
// If sampling occurs, seed becomes the seed for the random number generator used.
func (c Corpus) Sample(seed int64, want int) (Corpus, error) {
	// TODO(@MattWindsor91): test
	got := len(c)

	if got == 0 {
		return nil, ErrNoCorpus
	}

	if want <= 0 || got == want {
		return c, nil
	}
	if got < want {
		return nil, fmt.Errorf("%w: corpus size=%d, want %d", ErrSmallCorpus, got, want)
	}

	rng := rand.New(rand.NewSource(seed))
	return c.actuallySample(rng, want), nil
}

func (c Corpus) actuallySample(r *rand.Rand, want int) Corpus {
	sample := make(Corpus, want)

	js := r.Perm(len(c))[:want]
	sort.Ints(js)

	for i, j := range js {
		sample[i] = c[j]
	}

	return sample
}

package act_tester_plan

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

func (p *Planner) planCorpus(seed int64) ([]string, error) {
	// TODO(@MattWindsor91): perform corpus pruning
	prunedCorpus := p.Corpus

	return SampleCorpus(seed, prunedCorpus, p.CorpusSize)
}

// SampleCorpus tries to select a sample of size want from corpus-slice corpus.
// If want is non-positive, no sampling occurs.
// If sampling occurs, seed becomes the seed for the random number generator used.
func SampleCorpus(seed int64, corpus []string, want int) ([]string, error) {
	// TODO(@MattWindsor91): test
	got := len(corpus)

	if got == 0 {
		return nil, ErrNoCorpus
	}

	if want <= 0 || got == want {
		return corpus, nil
	}
	if got < want {
		return nil, fmt.Errorf("%w: corpus size=%d, want %d", ErrSmallCorpus, got, want)
	}

	rng := rand.New(rand.NewSource(seed))
	return actuallySampleCorpus(rng, corpus, want), nil
}

func actuallySampleCorpus(r *rand.Rand, corpus []string, want int) []string {
	sample := make([]string, want)

	js := r.Perm(len(corpus))[:want]
	sort.Ints(js)

	for i, j := range js {
		sample[i] = corpus[j]
	}

	return sample
}

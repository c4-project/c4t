// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"fmt"
	"math/rand"
	"sort"
)

// Sample tries to select a sample of size want from corpus-slice corpus.
// If want is non-positive, no sampling occurs.
// If sampling occurs, seed becomes the seed for the random number generator used.
func (c Corpus) Sample(rng *rand.Rand, want int) (Corpus, error) {
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

	return c.actuallySample(rng, want), nil
}

func (c Corpus) actuallySample(r *rand.Rand, want int) Corpus {
	sample := make(Corpus, want)

	names := c.Names()

	for _, j := range c.sampleIndices(r, want) {
		n := names[j]
		sample[n] = c[n]
	}

	return sample
}

// sampleIndices produces a random sorted list of n indices into the corpus's name list.
func (c Corpus) sampleIndices(r *rand.Rand, n int) []int {
	indices := r.Perm(len(c))[:n]
	sort.Ints(indices)
	return indices
}

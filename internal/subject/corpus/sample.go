// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"math/rand"
	"sort"
)

// Sample tries to select a sample of size want from this corpus.
// If want is non-positive, or the corpus is smaller than want, no sampling occurs.
func (c Corpus) Sample(rng *rand.Rand, want int) (Corpus, error) {
	got := len(c)

	if got == 0 {
		return nil, ErrNone
	}

	if want <= 0 || got <= want {
		return c, nil
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

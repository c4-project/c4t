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

// Sample tries to select a sample of size want from corpus-slice corpus.
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

// Chunks divides this corpus into n roughly-equal subcorpi.
// If n is less than or equal to 1, it returns the original corpus only.
// If n is greater than the length of the corpus, the remaining chunks will be nil.
// If n doesn't subdivide the corpus exactly, chunks towards the end will be smaller than chunks at the start.
func (c Corpus) Chunks(n int) []Corpus {
	if n <= 1 {
		return []Corpus{c}
	}

	chunks := make([]Corpus, n)

	rem := len(c)
	// Corner case to make sure empty corpi always have the empty corpus at the start;
	// if we didn't do this, the early-out at the top would be noticeably different.
	if rem == 0 {
		chunks[0] = c
		return chunks
	}
	pos := 0
	for i := n; 0 < i && 0 < rem; i-- {
		chsize := rem / i
		if chsize*i < rem {
			chsize++
		}

		chunks[n-i] = c[pos : pos+chsize]

		rem -= chsize
		pos += chsize
	}

	return chunks
}

package subject

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrCorpusDup occurs when we try to Add a subject into a corpus under a name that is already taken.
	ErrCorpusDup = errors.New("duplicate corpus entry")

	// ErrSmallCorpus occurs when the viable test corpus is smaller than that requested by the user.
	ErrSmallCorpus = errors.New("test corpus too small")

	// ErrNoCorpus is a variant of ErrSmallCorpus that occurs when the viable test corpus is empty.
	ErrNoCorpus = fmt.Errorf("%w: no corpus given", ErrSmallCorpus)
)

// Corpus is the type of test corpi (groups of test subjects).
type Corpus map[string]Subject

// NewCorpus creates a blank Corpus from a list of names.
func NewCorpus(names ...string) Corpus {
	corpus := make(Corpus, len(names))
	for _, n := range names {
		corpus[n] = Subject{}
	}
	return corpus
}

// Add tries to add s to the corpus.
// It fails if the corpus already has a subject with the given name.
func (c Corpus) Add(s Named) error {
	_, exists := c[s.Name]
	if exists {
		return fmt.Errorf("%w: %s", ErrCorpusDup, s.Name)
	}

	c[s.Name] = s.Subject
	return nil
}

// Par runs f for every subject in the plan's corpus.
// It threads through a context that will terminate each machine if an error occurs on some other machine.
// It also takes zero or more 'auxiliary' funcs to launch within the same context.
func (c Corpus) Par(ctx context.Context, f func(context.Context, Named) error, aux ...func(context.Context) error) error {
	eg, ectx := errgroup.WithContext(ctx)
	for n, s := range c {
		sc := Named{Name: n, Subject: s}
		eg.Go(func() error { return f(ectx, sc) })
	}
	for _, a := range aux {
		eg.Go(func() error { return a(ectx) })
	}
	return eg.Wait()
}

// Map sequentially maps f over the subjects in this corpus.
// It passes each invocation of f a pointer to a copy of a subject, but propagates any changes made to that copy back to
// the corpus.
// It does not permit making change to the name.
func (c Corpus) Map(f func(*Named) error) error {
	for n := range c {
		sn := Named{Name: n, Subject: c[n]}
		if err := f(&sn); err != nil {
			return err
		}

		if n != sn.Name {
			return fmt.Errorf("name change from %q to %q forbidden", n, sn.Name)
		}
		c[n] = sn.Subject
	}
	return nil
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

// Names returns a sorted list of this corpus's subject names.
func (c Corpus) Names() []string {
	ns := make([]string, len(c))
	i := 0
	for n := range c {
		ns[i] = n
		i++
	}
	sort.Strings(ns)
	return ns
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

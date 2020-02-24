package corpus

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
	"golang.org/x/sync/errgroup"
)

// Par runs f for every subject in the plan's corpus.
// It threads through a context that will terminate each machine if an error occurs on some other machine.
// It also takes zero or more 'auxiliary' funcs to launch within the same context.
func (c Corpus) Par(ctx context.Context, f func(context.Context, subject.Named) error, aux ...func(context.Context) error) error {
	eg, ectx := errgroup.WithContext(ctx)
	for n, s := range c {
		sc := subject.Named{Name: n, Subject: s}
		eg.Go(func() error { return f(ectx, sc) })
	}
	for _, a := range aux {
		eg.Go(func() error { return a(ectx) })
	}
	return eg.Wait()
}

// Each applies f to each subject in the corpus.
// It fails if any invocation of f fails.
func (c Corpus) Each(f func(subject.Named) error) error {
	for n := range c {
		if err := f(subject.Named{Name: n, Subject: c[n]}); err != nil {
			return err
		}
	}
	return nil
}

// Map sequentially maps f over the subjects in this corpus.
// It passes each invocation of f a pointer to a copy of a subject, but propagates any changes made to that copy back to
// the corpus.
// It does not permit making change to the name.
func (c Corpus) Map(f func(*subject.Named) error) error {
	return c.Each(func(sn subject.Named) error {
		n := sn.Name
		if err := f(&sn); err != nil {
			return err
		}

		if n != sn.Name {
			return fmt.Errorf("%w: from %q to %q", ErrMapRename, n, sn.Name)
		}
		c[n] = sn.Subject
		return nil
	})
}

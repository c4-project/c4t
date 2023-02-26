// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/c4-project/c4t/internal/subject"
)

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
// It does not permit making changes to the name.
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

// Par runs f for every subject in the plan's corpus, with a degree of parallelism.
// It threads through a context that will terminate each machine if an error occurs on some other machine.
// It also takes zero or more 'auxiliary' funcs to launch within the same context.
func (c Corpus) Par(ctx context.Context, nworkers int, f func(context.Context, subject.Named) error, aux ...func(context.Context) error) error {
	eg, ectx := errgroup.WithContext(ctx)

	for _, a := range aux {
		a := a
		eg.Go(func() error { return a(ectx) })
	}
	c.parInner(eg, ectx, nworkers, f)
	return eg.Wait()
}

func (c Corpus) parInner(eg *errgroup.Group, ectx context.Context, nworkers int, f func(context.Context, subject.Named) error) {
	switch {
	case len(c) == 0:
		return
	case len(c) < nworkers:
		c.parDirect(eg, ectx, f)
	case nworkers <= 1:
		c.parSingle(eg, ectx, f)
	default:
		c.parWorkers(eg, ectx, nworkers, f)
	}
}

func (c Corpus) parSingle(eg *errgroup.Group, ectx context.Context, f func(context.Context, subject.Named) error) {
	eg.Go(func() error {
		return c.Each(func(sc subject.Named) error {
			return f(ectx, sc)
		})
	})
}

func (c Corpus) parDirect(eg *errgroup.Group, ectx context.Context, f func(context.Context, subject.Named) error) {
	_ = c.Each(func(s subject.Named) error {
		eg.Go(func() error { return f(ectx, s) })
		return nil
	})
}

func (c Corpus) parWorkers(eg *errgroup.Group, ectx context.Context, nworkers int, f func(context.Context, subject.Named) error) {
	// TODO(@MattWindsor91): is this a decent buffer number?
	wch := make(chan subject.Named, nworkers)

	eg.Go(func() error {
		return c.workerSource(wch, ectx)
	})
	for i := 0; i < nworkers; i++ {
		eg.Go(func() error {
			return c.workerSink(wch, f, ectx)
		})
	}
}

func (c Corpus) workerSink(wch <-chan subject.Named, f func(context.Context, subject.Named) error, ectx context.Context) error {
	for {
		select {
		case sc, ok := <-wch:
			if !ok {
				return nil
			}
			if err := f(ectx, sc); err != nil {
				return err
			}
		case <-ectx.Done():
			return ectx.Err()
		}
	}
}

func (c Corpus) workerSource(wch chan<- subject.Named, ectx context.Context) error {
	err := c.Each(func(sc subject.Named) error {
		select {
		case wch <- sc:
			return nil
		case <-ectx.Done():
			return ectx.Err()
		}
	})
	if err != nil {
		return err
	}
	close(wch)
	return nil
}

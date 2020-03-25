// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package lifter contains the part of the tester framework that lifts litmus tests to compilable C.
package lifter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/job"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"
)

var (
	// ErrMakerNil occurs when a lifter runs without a HarnessMaker set.
	ErrMakerNil = errors.New("harness maker nil")

	// ErrNoBackend occurs when backend information is missing.
	ErrNoBackend = errors.New("no backend provided")
)

// HarnessMaker is an interface capturing the ability to make test harnesses.
type HarnessMaker interface {
	// MakeHarness asks the harness maker to make the test harness described by j.
	// It returns a list outFiles of files created (C files, header files, etc.), and/or an error err.
	// Any error output from child processes should be sent to errw, if it is non-nil.
	MakeHarness(ctx context.Context, j job.Harness, errw io.Writer) (outFiles []string, err error)
}

// Lifter holds the main configuration for the lifter part of the tester framework.
type Lifter struct {
	// conf is the configuration used for this lifter.
	conf Config

	// plan is the plan on which this lifter is operating.
	plan plan.Plan

	// l is the logger to use for this lifter.
	l *log.Logger
}

// New constructs a new Lifter given config c and plan p.
func New(c *Config, p *plan.Plan) (*Lifter, error) {
	if c.Paths == nil {
		return nil, iohelp.ErrPathsetNil
	}
	if c.Maker == nil {
		return nil, ErrMakerNil
	}
	if p == nil {
		return nil, plan.ErrNil
	}
	if p.Backend == nil {
		return nil, ErrNoBackend
	}

	l := Lifter{
		conf: *c,
		plan: *p,
		l:    iohelp.EnsureLog(c.Logger),
	}
	return &l, nil
}

// Run runs a lifting job: taking every test subject in a plan and using a backend to lift each into a test harness.
func (l *Lifter) Run(ctx context.Context) (*plan.Plan, error) {
	l.l.Println("preparing directories")
	if err := l.conf.Paths.Prepare(l.plan.Arches(), l.plan.Corpus.Names()); err != nil {
		return nil, err
	}

	err := l.lift(ctx)
	return &l.plan, err
}

func (l *Lifter) lift(ctx context.Context) error {
	l.l.Println("now lifting")

	b, err := builder.NewBuilder(builder.Config{
		Init: l.plan.Corpus,
		Obs:  l.conf.Observer,
		Manifest: builder.Manifest{
			Name:  "lift",
			NReqs: l.count(),
		},
	})
	if err != nil {
		return fmt.Errorf("when making builder: %w", err)
	}
	mrng := l.plan.Header.Rand()

	var lerr error
	l.plan.Corpus, lerr = l.liftInner(ctx, mrng, b)
	return lerr
}

func (l *Lifter) liftInner(ctx context.Context, mrng *rand.Rand, b *builder.Builder) (corpus.Corpus, error) {
	eg, ectx := errgroup.WithContext(ctx)
	var lc corpus.Corpus
	// It's very likely this will be a single element array.
	for _, a := range l.plan.Arches() {
		j := l.makeJob(a, mrng, b.SendCh)
		eg.Go(func() error {
			return j.Lift(ectx)
		})
	}
	eg.Go(func() error {
		var err error
		lc, err = b.Run(ectx)
		return err
	})
	err := eg.Wait()
	return lc, err
}

func (l *Lifter) makeJob(a id.ID, mrng *rand.Rand, resCh chan<- builder.Request) Job {
	return Job{
		Arch:    a,
		Backend: l.plan.Backend,
		Paths:   l.conf.Paths,
		Maker:   l.conf.Maker,
		Corpus:  l.plan.Corpus,
		Rng:     rand.New(rand.NewSource(mrng.Int63())),
		ResCh:   resCh,
		Stderr:  l.conf.Stderr,
	}
}

// count counts the number of liftings that need doing.
func (l *Lifter) count() int {
	return len(l.plan.Arches()) * len(l.plan.Corpus)
}

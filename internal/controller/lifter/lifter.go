// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package lifter contains the part of the tester framework that lifts litmus tests to compilable C.
package lifter

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

var (
	// ErrConfigNil occurs when we try to construct a lifter without config.
	ErrConfigNil = errors.New("config nil")

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
	if err := checkConfig(c); err != nil {
		return nil, err
	}
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	l := Lifter{
		conf: *c,
		plan: *p,
		l:    iohelp.EnsureLog(c.Logger),
	}
	return &l, nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	return c.Check()
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	if p.Backend == nil {
		return ErrNoBackend
	}
	return p.Check()
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

	cfg := builder.Config{
		Init:      l.plan.Corpus,
		Observers: l.conf.Observers,
		Manifest: builder.Manifest{
			Name:  "lift",
			NReqs: l.count(),
		},
	}

	mrng := l.plan.Header.Rand()

	var err error
	// TODO(@MattWindsor91): extract this 20 into configuration.
	l.plan.Corpus, err = builder.ParBuild(ctx, 20, l.plan.Corpus, cfg, func(ctx context.Context, s subject.Named, rq chan<- builder.Request) error {
		j := l.makeJob(s, mrng, rq)
		return j.Lift(ctx)
	})
	return err
}

func (l *Lifter) makeJob(s subject.Named, mrng *rand.Rand, resCh chan<- builder.Request) Job {
	return Job{
		Arches:  l.plan.Arches(),
		Backend: l.plan.Backend,
		Paths:   l.conf.Paths,
		Maker:   l.conf.Maker,
		Subject: s,
		Rng:     rand.New(rand.NewSource(mrng.Int63())),
		ResCh:   resCh,
		Stderr:  l.conf.Stderr,
	}
}

// count counts the number of liftings that need doing.
func (l *Lifter) count() int {
	return len(l.plan.Arches()) * len(l.plan.Corpus)
}

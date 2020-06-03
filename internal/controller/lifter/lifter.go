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
	// ErrDriverNil occurs when a lifter runs without a SingleLifter set.
	ErrDriverNil = errors.New("driver nil")

	// ErrNoBackend occurs when backend information is missing.
	ErrNoBackend = errors.New("no backend provided")
)

// SingleLifter is an interface capturing the ability to lift single jobs into recipes.
type SingleLifter interface {
	// Lift performs the lifting described by j.
	// It returns a list outFiles of files created (C files, header files, etc.), and/or an error err.
	// Any error output from child processes should be sent to errw, if it is non-nil.
	Lift(ctx context.Context, j job.Lifter, errw io.Writer) (outFiles []string, err error)
}

//go:generate mockery -name SingleLifter

// Lifter holds the main configuration for the lifter part of the tester framework.
type Lifter struct {
	// driver is a single-job lifter.
	driver SingleLifter

	// l is the logger to use for this lifter.
	l *log.Logger

	// obs track the lifter's progress across a corpus.
	obs []builder.Observer

	// paths does path resolution and preparation for the incoming lifter.
	paths Pather

	// errw is the writer to which standard error (eg from the lifting backend) should be sent.
	errw io.Writer
}

// New constructs a new Lifter given driver d, path resolver p, and options os.
func New(d SingleLifter, p Pather, os ...Option) (*Lifter, error) {
	if err := checkConfig(d, p); err != nil {
		return nil, err
	}
	l := Lifter{driver: d, paths: p}
	if err := Options(os...)(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

func checkConfig(d SingleLifter, p Pather) error {
	if d == nil {
		return ErrDriverNil
	}
	if p == nil {
		return iohelp.ErrPathsetNil
	}
	return nil
}

// Run runs a lifting job: taking every test subject in p and using a backend to lift each to a compilable recipe.
func (l *Lifter) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	l.l.Println("preparing directories")
	if err := l.prepareDirs(p); err != nil {
		return nil, err
	}

	return l.lift(ctx, *p)
}

func (l *Lifter) prepareDirs(p *plan.Plan) error {
	l.l.Println("preparing directories")
	return l.paths.Prepare(p.Arches(), p.Corpus.Names())
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

func (l *Lifter) lift(ctx context.Context, p plan.Plan) (*plan.Plan, error) {
	l.l.Println("now lifting")

	cfg := builder.Config{
		Init:      p.Corpus,
		Observers: l.obs,
		Manifest: builder.Manifest{
			Name:  "lift",
			NReqs: p.MaxNumRecipes(),
		},
	}

	mrng := p.Metadata.Rand()

	var err error
	// TODO(@MattWindsor91): extract this 20 into configuration.
	p.Corpus, err = builder.ParBuild(ctx, 20, p.Corpus, cfg, func(ctx context.Context, s subject.Named, rq chan<- builder.Request) error {
		j := l.makeJob(&p, s, mrng, rq)
		return j.Lift(ctx)
	})
	return &p, err
}

func (l *Lifter) makeJob(p *plan.Plan, s subject.Named, mrng *rand.Rand, resCh chan<- builder.Request) Job {
	return Job{
		Arches:  p.Arches(),
		Backend: p.Backend,
		Paths:   l.paths,
		Driver:  l.driver,
		Subject: s,
		Rng:     rand.New(rand.NewSource(mrng.Int63())),
		ResCh:   resCh,
		Stderr:  l.errw,
	}
}

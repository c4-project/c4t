// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package lifter contains the part of the tester framework that lifts litmus tests to compilable C.
package lifter

import (
	"context"
	"errors"
	"io"
	"math/rand"

	"github.com/MattWindsor91/c4t/internal/helper/srvrun"

	"github.com/MattWindsor91/c4t/internal/model/service"

	"github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/subject/corpus"

	"github.com/MattWindsor91/c4t/internal/plan/stage"

	"github.com/MattWindsor91/c4t/internal/model/recipe"

	"github.com/MattWindsor91/c4t/internal/subject"

	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"

	"github.com/MattWindsor91/c4t/internal/plan"
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
	// It returns a recipe describing the files (C files, header files, etc.) created and how to use them, or an error.
	// Any external service running should happen by sr.
	Lift(ctx context.Context, j backend.LiftJob, sr service.Runner) (recipe.Recipe, error)
}

//go:generate mockery --name=SingleLifter

// Lifter holds the main configuration for the lifter part of the tester framework.
type Lifter struct {
	// driver is a single-job lifter.
	driver SingleLifter

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
	return p.RunStage(ctx, stage.Lift, l.lift)
}

func (l *Lifter) prepareDirs(p *plan.Plan) error {
	// TODO(@MattWindsor91): observe this?
	return l.paths.Prepare(p.Arches(), p.Corpus.Names())
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	if err := p.Check(); err != nil {
		return err
	}
	if p.Backend == nil {
		return ErrNoBackend
	}
	return p.Metadata.RequireStage(stage.Plan)
}

func (l *Lifter) lift(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := l.prepareDirs(p); err != nil {
		return nil, err
	}

	var err error
	outp := *p
	outp.Corpus, err = l.liftCorpus(ctx, p)
	return &outp, err
}

func (l *Lifter) liftCorpus(ctx context.Context, p *plan.Plan) (corpus.Corpus, error) {
	cfg := builder.Config{
		Init:      p.Corpus,
		Observers: l.obs,
		Manifest: builder.Manifest{
			Name:  "lift",
			NReqs: p.MaxNumRecipes(),
		},
	}
	mrng := p.Metadata.Rand()
	// TODO(@MattWindsor91): extract this 20 into configuration.
	return builder.ParBuild(ctx, 20, p.Corpus, cfg, func(ctx context.Context, s subject.Named, rq chan<- builder.Request) error {
		j := l.makeJob(p, s, mrng, rq)
		return j.Lift(ctx)
	})
}

func (l *Lifter) makeJob(p *plan.Plan, s subject.Named, mrng *rand.Rand, resCh chan<- builder.Request) Instance {
	return Instance{
		Arches:  p.Arches(),
		Backend: p.Backend,
		Paths:   l.paths,
		Driver:  l.driver,
		Subject: s,
		Rng:     rand.New(rand.NewSource(mrng.Int63())),
		ResCh:   resCh,
		// TODO(@MattWindsor91): push this further up
		Runner: srvrun.NewExecRunner(srvrun.StderrTo(l.errw)),
	}
}

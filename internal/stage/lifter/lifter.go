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

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/helper/srvrun"

	"github.com/c4-project/c4t/internal/subject/corpus"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/subject"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/plan"
)

var (
	// ErrDriverNil occurs when a lifter runs without a SingleLifter set.
	ErrDriverNil = errors.New("driver nil")

	// ErrNoBackend occurs when backend information is missing.
	ErrNoBackend = errors.New("no backend provided")
)

// Lifter holds the main configuration for the lifter part of the tester framework.
type Lifter struct {
	// resolver resolves backend specifications.
	resolver backend.Resolver

	// obs track the lifter's progress across a corpus.
	obs []builder.Observer

	// paths does path resolution and preparation for the incoming lifter.
	paths Pather

	// errw is the writer to which standard error (eg from the lifting backend) should be sent.
	errw io.Writer
}

// New constructs a new Lifter given backend resolver r, path resolver p, and options os.
func New(r backend.Resolver, p Pather, os ...Option) (*Lifter, error) {
	if err := checkConfig(r, p); err != nil {
		return nil, err
	}
	l := Lifter{resolver: r, paths: p}
	if err := Options(os...)(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

func checkConfig(r backend.Resolver, p Pather) error {
	if r == nil {
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
	b, err := l.resolver.Resolve(p.Backend.Spec)
	if err != nil {
		return nil, err
	}

	cfg := builder.Config{
		Init:      p.Corpus,
		Observers: l.obs,
		Manifest: builder.Manifest{
			Name:  "lift",
			NReqs: p.MaxNumRecipes(),
		},
	}
	// TODO(@MattWindsor91): extract this 20 into configuration.
	return builder.ParBuild(ctx, 20, p.Corpus, cfg, func(ctx context.Context, s subject.Named, rq chan<- builder.Request) error {
		j := l.makeJob(p, b, s, rq)
		return j.Lift(ctx)
	})
}

func (l *Lifter) makeJob(p *plan.Plan, b backend.SingleLifter, s subject.Named, resCh chan<- builder.Request) Instance {

	return Instance{
		Arches: p.Arches(),
		// TODO(@MattWindsor91): remove this
		Paths:   l.paths,
		Driver:  b,
		Subject: s,
		ResCh:   resCh,
		// TODO(@MattWindsor91): push this further up
		Runner: srvrun.NewExecRunner(srvrun.StderrTo(l.errw)),
	}
}

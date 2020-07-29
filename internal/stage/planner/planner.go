// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Planner holds all configuration for the test planner.
type Planner struct {
	// source contains all of the various sources for a Planner's information.
	source Source
	// filter is the compiler filter to use to select compilers to test.
	filter string
	// l is the logger used by the planner.
	l *log.Logger
	// observers contains the set of observers used to get feedback on the planning action as it completes.
	observers ObserverSet
	// quantities contains quantity information for this planner.
	quantities QuantitySet
	fs         []string
	mach       machine.Named
	seed       int64
}

// New constructs a new planner with the given config, machine information, files, and options.
func New(src Source, mach machine.Named, fs []string, opts ...Option) (*Planner, error) {
	if err := src.Check(); err != nil {
		return nil, err
	}
	// Early out to prevent us from doing any planning if we received no files.
	if len(fs) == 0 {
		return nil, corpus.ErrNone
	}

	p := &Planner{
		source: src,
		fs:     fs,
		mach:   mach,
		seed:   plan.UseDateSeed,
	}
	if err := Options(opts...)(p); err != nil {
		return nil, err
	}
	p.l = iohelp.EnsureLog(p.l)
	return p, nil
}

// Plan runs the test planner p.
func (p *Planner) Plan(ctx context.Context) (*plan.Plan, error) {
	return (&plan.Plan{Machine: p.mach}).RunStage(ctx, stage.Plan, p.planInner)
}

func (p *Planner) planInner(ctx context.Context, pn *plan.Plan) (*plan.Plan, error) {
	hd := plan.NewMetadata(p.seed)
	pn.Metadata = *hd

	rng := hd.Rand()

	p.l.Println("Planning backend...")
	if err := p.planBackend(ctx, pn); err != nil {
		return nil, err
	}

	p.l.Println("Planning compilers...")
	if err := p.planCompilers(ctx, rng, pn); err != nil {
		return nil, err
	}

	p.l.Println("Planning corpus...")
	if err := p.planCorpus(ctx, rng, pn); err != nil {
		return nil, err
	}

	return pn, nil
}

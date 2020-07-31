// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Planner holds all configuration for the test planner.
type Planner struct {
	// source contains all of the various sources for a Planner's information.
	source Source
	// filter is the compiler filter to use to select compilers to test.
	filter string
	// observers contains the set of observers used to get feedback on the planning action as it completes.
	observers []Observer
	// quantities contains quantity information for this planner.
	quantities QuantitySet
	// fs is the set of input corpus files to use for this planner.
	fs []string
	// mach is the machine to use for this planner.
	mach machine.Named
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
	}
	err := Options(opts...)(p)
	return p, err
}

// Plan runs the test planner p.
func (p *Planner) Plan(ctx context.Context) (*plan.Plan, error) {
	return (&plan.Plan{Machine: p.mach}).RunStage(ctx, stage.Plan, p.planInner)
}

func (p *Planner) planInner(ctx context.Context, pn *plan.Plan) (*plan.Plan, error) {
	p.announce(Message{Kind: KindStart, Quantities: &p.quantities})

	hd := plan.NewMetadata(0)
	pn.Metadata = *hd

	p.announce(Message{Kind: KindPlanningBackend})
	if err := p.planBackend(ctx, pn); err != nil {
		return nil, err
	}

	p.announce(Message{Kind: KindPlanningCompilers, MachineID: p.mach.ID})
	if err := p.planCompilers(ctx, pn); err != nil {
		return nil, err
	}

	p.announce(Message{Kind: KindPlanningCorpus})
	if err := p.planCorpus(ctx, pn); err != nil {
		return nil, err
	}

	return pn, nil
}

func (p *Planner) announce(m Message) {
	OnPlan(m, p.observers...)
}

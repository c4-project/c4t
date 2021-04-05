// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"
	"time"

	"github.com/c4-project/c4t/internal/timing"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/subject/corpus"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/plan"
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
	quantities quantity.PlanSet
}

// New constructs a new planner with the given source and options.
func New(src Source, opts ...Option) (*Planner, error) {
	if err := src.Check(); err != nil {
		return nil, err
	}

	p := &Planner{
		source: src,
	}
	err := Options(opts...)(p)
	return p, err
}

// Plan runs the test planner on the given fileset fs and config ms.
// It returns, on success, a map from machine IDs to their test plans.
func (p *Planner) Plan(ctx context.Context, ms machine.ConfigMap, fs ...string) (plan.Map, error) {
	// Early out to prevent us from doing any planning if we received no files.
	if len(fs) == 0 {
		return nil, corpus.ErrNone
	}

	start := time.Now()
	p.announce(Message{Kind: KindStart, Quantities: &p.quantities})

	p.announce(Message{Kind: KindPlanningCorpus})
	corp, err := p.planCorpus(ctx, fs...)
	if err != nil {
		return nil, err
	}

	return p.planWithCorpus(ms, start, corp)
}

func (p *Planner) planWithCorpus(ms machine.ConfigMap, start time.Time, corp corpus.Corpus) (plan.Map, error) {
	ps := make(plan.Map, len(ms))
	var err error
	for n, m := range ms {
		if ps[n], err = p.makeMachinePlan(start, n, m, corp); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func (p *Planner) makeMachinePlan(start time.Time, mid id.ID, m machine.Config, corp corpus.Corpus) (plan.Plan, error) {
	var (
		pn  plan.Plan
		err error
	)

	pn.Machine = machine.Named{ID: mid, Machine: m.Machine}
	pn.Corpus = corp
	pn.Mutation = m.Mutation

	p.announce(Message{Kind: KindPlanningBackend, MachineID: mid})
	pn.Backend, err = p.planBackend()
	if err != nil {
		return pn, err
	}

	p.announce(Message{Kind: KindPlanningCompilers, MachineID: mid})
	pn.Compilers, err = p.planCompilers(m)
	if err != nil {
		return pn, err
	}

	pn.Metadata = *plan.NewMetadata(0)
	pn.Metadata.ConfirmStage(stage.Plan, timing.SpanSince(start))
	return pn, nil
}

func (p *Planner) announce(m Message) {
	OnPlan(m, p.observers...)
}

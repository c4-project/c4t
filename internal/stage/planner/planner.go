// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"
	"time"

	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/MattWindsor91/c4t/internal/plan/stage"

	"github.com/MattWindsor91/c4t/internal/subject/corpus"

	"github.com/MattWindsor91/c4t/internal/machine"

	"github.com/MattWindsor91/c4t/internal/plan"
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

// Plan runs the test planner p.
func (p *Planner) Plan(ctx context.Context, ms machine.ConfigMap, fs ...string) (map[string]plan.Plan, error) {
	// Early out to prevent us from doing any planning if we received no files.
	if len(fs) == 0 {
		return nil, corpus.ErrNone
	}

	start := time.Now()
	p.announce(Message{Kind: KindStart, Quantities: &p.quantities})

	ps := make(map[string]plan.Plan, len(ms))

	p.announce(Message{Kind: KindPlanningCorpus})
	corp, err := p.planCorpus(ctx, fs...)
	if err != nil {
		return nil, err
	}

	for n, m := range ms {
		nid, err := id.TryFromString(n)
		if err != nil {
			return nil, err
		}
		ps[n], err = p.makeMachinePlan(ctx, start, machine.Named{ID: nid, Machine: m.Machine}, corp)
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func (p *Planner) makeMachinePlan(ctx context.Context, start time.Time, nm machine.Named, corp corpus.Corpus) (plan.Plan, error) {

	var (
		pn  plan.Plan
		err error
	)

	pn.Machine = nm
	pn.Corpus = corp

	p.announce(Message{Kind: KindPlanningBackend, MachineID: nm.ID})
	pn.Backend, err = p.planBackend(ctx, nm.ID)
	if err != nil {
		return pn, err
	}

	p.announce(Message{Kind: KindPlanningCompilers, MachineID: nm.ID})
	pn.Compilers, err = p.planCompilers(ctx, nm.ID)
	if err != nil {
		return pn, err
	}

	pn.Metadata = *plan.NewMetadata(0)
	pn.Metadata.ConfirmStage(stage.Plan, start, time.Since(start))
	return pn, nil
}

func (p *Planner) announce(m Message) {
	OnPlan(m, p.observers...)
}

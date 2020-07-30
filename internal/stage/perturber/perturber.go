// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package perturber contains the logic for the test perturber.
package perturber

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Perturber holds all configuration for the test perturber.
type Perturber struct {
	// source contains all of the various sources for a Planner's information.
	source Source
	// observers contains the set of observers used to get feedback on the planning action as it completes.
	observers []Observer
	// quantities contains quantity information for this planner.
	quantities QuantitySet
	seed       int64
}

// New constructs a new perturber with the given source and options.
func New(src Source, opts ...Option) (*Perturber, error) {
	if err := src.Check(); err != nil {
		return nil, err
	}

	p := &Perturber{
		source: src,
		seed:   plan.UseDateSeed,
	}
	if err := Options(opts...)(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Run runs the test perturber on pn.
// It returns a modified plan on success, which is guaranteed to be different from pn.
func (p *Perturber) Run(ctx context.Context, pn *plan.Plan) (*plan.Plan, error) {
	return pn.RunStage(ctx, stage.Perturb, p.perturbInner)
}

func (p *Perturber) perturbInner(_ context.Context, inplan *plan.Plan) (*plan.Plan, error) {
	OnPerturb(Message{Kind: KindStart, Quantities: &p.quantities})

	// Avoid modifying inplan in-place.
	pn := *inplan

	hd := plan.NewMetadata(p.seed)
	pn.Metadata = *hd
	rng := hd.Rand()

	OnPerturb(Message{Kind: KindRandomiseOpts})
	if err := p.perturbCompilers(rng, &pn); err != nil {
		return nil, err
	}

	OnPerturb(Message{Kind: KindSampleCorpus}, p.observers...)
	if err := p.sampleCorpus(rng, &pn); err != nil {
		return nil, err
	}

	return &pn, nil
}

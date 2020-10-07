// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package perturber contains the logic for the test perturber.
package perturber

import (
	"context"
	"errors"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// ErrCInspectorNil occurs if the perturber constructor is passed a nil compiler inspector.
var ErrCInspectorNil = errors.New("compiler inspector nil")

// Perturber holds all configuration for the test perturber.
type Perturber struct {
	// ci contains the inspector used to get possible optimisation levels for compiler randomisation.
	ci compiler.Inspector
	// observers contains the set of observers used to get feedback on the planning action as it completes.
	observers []Observer
	// useFullIDs tells the perturber whether to promote the compiler IDs to full IDs after perturbance.
	useFullIDs bool
	// quantities contains quantity information for this planner.
	quantities quantity.PerturbSet
	seed       int64
}

// New constructs a new perturber with the given compiler inspector and options.
func New(ci compiler.Inspector, opts ...Option) (*Perturber, error) {
	if ci == nil {
		return nil, ErrCInspectorNil
	}

	p := &Perturber{
		ci:   ci,
		seed: plan.UseDateSeed,
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
	// Avoid modifying inplan in-place.
	pn := *inplan

	if err := p.perturbCopy(&pn); err != nil {
		return nil, err
	}
	return &pn, nil
}

func (p *Perturber) perturbCopy(pn *plan.Plan) error {
	p.announce(Message{Kind: KindStart, Quantities: &p.quantities})

	p.perturbMetadata(pn)
	rng := pn.Metadata.Rand()
	p.announce(Message{Kind: KindSeedChanged, Seed: pn.Metadata.Seed})

	p.announce(Message{Kind: KindRandomisingOpts})
	if err := p.perturbCompilers(rng, pn); err != nil {
		return err
	}

	p.announce(Message{Kind: KindSamplingCorpus})
	if err := p.sampleCorpus(rng, pn); err != nil {
		return err
	}

	return nil
}

func (p *Perturber) perturbMetadata(pn *plan.Plan) {
	hd := plan.NewMetadata(p.seed)
	hd.Stages = pn.Metadata.Stages
	pn.Metadata = *hd
}

func (p *Perturber) announce(m Message) {
	OnPerturb(m, p.observers...)
}

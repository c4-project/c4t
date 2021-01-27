// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package runner contains the part of c4t that runs compiled test binaries and interprets their output.
package runner

import (
	"context"

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/quantity"
	"github.com/c4-project/c4t/internal/stage/mach/observer"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/subject"
)

// Runner contains information necessary to run a plan's compiled test cases.
type Runner struct {
	// observers observe the runner's progress across a corpus.
	observers []observer.Observer

	// resolver resolves backend references in the plan.
	resolver backend.Resolver

	// paths contains the pathset used for this runner's outputs.
	paths *Pathset

	// quantities contains quantity configuration for this runner.
	quantities quantity.BatchSet
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(resolver backend.Resolver, paths *Pathset, opts ...Option) (*Runner, error) {
	if resolver == nil {
		return nil, ErrParserNil
	}
	if paths == nil {
		return nil, iohelp.ErrPathsetNil
	}
	r := &Runner{
		resolver: resolver,
		paths:    paths,
	}
	if err := Options(opts...)(r); err != nil {
		return nil, err
	}
	return r, nil
}

// Stage gets the appropriate stage record for the runner.
func (*Runner) Stage() stage.Stage {
	return stage.Run
}

// Close does nothing.
func (*Runner) Close() error {
	return nil
}

// Run runs the runner on the plan p.
func (r *Runner) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	observer.OnRunStart(r.quantities, r.observers...)

	b, err := r.resolver.Resolve(p.Backend.Spec)
	if err != nil {
		return nil, err
	}

	bcfg := r.builderConfig(p)
	c, err := builder.ParBuild(ctx, r.quantities.NWorkers, p.Corpus, bcfg,
		func(ctx context.Context, named subject.Named, requests chan<- builder.Request) error {
			return r.instance(requests, named, b).Run(ctx)
		})
	if err != nil {
		return nil, err
	}

	np := *p
	np.Corpus = c
	return &np, nil
}

func (r *Runner) builderConfig(p *plan.Plan) builder.Config {
	return builder.Config{
		Init:      p.Corpus,
		Observers: observer.LowerToBuilder(r.observers...),
		Manifest: builder.Manifest{
			Name:  "run",
			NReqs: p.NumExpCompilations(),
		},
	}
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	if err := p.Check(); err != nil {
		return err
	}
	return p.Metadata.RequireStage(stage.Compile)
}

func (r *Runner) instance(requests chan<- builder.Request, named subject.Named, backend backend.Backend) *Instance {

	return &Instance{
		backend:    backend,
		quantities: r.quantities,
		resCh:      requests,
		subject:    &named,
	}
}

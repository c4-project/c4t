// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package runner contains the part of act-tester that runs compiled test binaries and interprets their output.
package runner

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
	"github.com/MattWindsor91/act-tester/internal/plan"
)

// ObsParser is the interface of things that can parse test outcomes.
type ObsParser interface {
	// ParseObs parses the observation in reader r into o according to the backend configuration in b.
	// The backend described by b must have been used to produce the testcase outputting r.
	ParseObs(ctx context.Context, b *service.Backend, r io.Reader, o *obs.Obs) error
}

// Runner contains information necessary to run a plan's compiled test cases.
type Runner struct {
	// observers observe the runner's progress across a corpus.
	observers []builder.Observer

	// parser handles the parsing of observations.
	parser ObsParser

	// paths contains the pathset used for this runner's outputs.
	paths *Pathset

	// quantities contains quantity configuration for this runner.
	quantities quantity.BatchSet
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(parser ObsParser, paths *Pathset, opts ...Option) (*Runner, error) {
	if parser == nil {
		return nil, ErrParserNil
	}
	if paths == nil {
		return nil, iohelp.ErrPathsetNil
	}
	r := &Runner{
		parser: parser,
		paths:  paths,
	}
	if err := Options(opts...)(r); err != nil {
		return nil, err
	}
	return r, nil
}

// Run runs the runner on the plan p.
func (r *Runner) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	return p.RunStage(ctx, stage.Run, r.runInner)
}

func (r *Runner) runInner(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	// TODO(@MattWindsor91): port to observers
	// r.quantities.Log(r.l)

	bcfg := r.builderConfig(p)
	c, err := builder.ParBuild(ctx, r.quantities.NWorkers, p.Corpus, bcfg,
		func(ctx context.Context, named subject.Named, requests chan<- builder.Request) error {
			return r.instance(requests, named, p).Run(ctx)
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
		Observers: r.observers,
		Manifest: builder.Manifest{
			Name:  "run",
			NReqs: p.NumExpCompilations(),
		},
	}
}

func checkPlan(p *plan.Plan) error {
	// TODO(@MattWindsor91): require compile stage
	if p == nil {
		return plan.ErrNil
	}
	if err := p.Check(); err != nil {
		return err
	}
	return p.Metadata.RequireStage(stage.Compile)
}

func (r *Runner) instance(requests chan<- builder.Request, named subject.Named, p *plan.Plan) *Instance {
	return &Instance{
		backend:    p.Backend,
		parser:     r.parser,
		quantities: r.quantities,
		resCh:      requests,
		subject:    &named,
	}
}

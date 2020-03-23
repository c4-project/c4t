// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package runner contains the part of act-tester that runs compiled harness binaries and interprets their output.
package runner

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// Runner contains information necessary to run a plan's compiled test cases.
type Runner struct {
	// l is the logger for this runner.
	l *log.Logger

	// plan is the plan on which this runner is operating.
	plan plan.Plan

	// conf is the configuration used to build this runner.
	conf Config
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(c *Config, p *plan.Plan) (*Runner, error) {
	if c == nil {
		return nil, ErrConfigNil
	}
	if err := c.Check(); err != nil {
		return nil, err
	}

	if p == nil {
		return nil, plan.ErrNil
	}

	r := Runner{
		conf: *c,
		plan: *p,
		l:    iohelp.EnsureLog(c.Logger),
	}

	if err := r.check(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Runner) check() error {
	if len(r.plan.Corpus) == 0 {
		return corpus.ErrNone
	}
	return nil
}

// Run runs the runner.
func (r *Runner) Run(ctx context.Context) (*plan.Plan, error) {
	bcfg := builder.Config{
		Init: r.plan.Corpus,
		Obs:  r.conf.Observer,
		Manifest: builder.Manifest{
			Name:  "run",
			NReqs: r.count(),
		},
	}
	b, berr := builder.NewBuilder(bcfg)
	if berr != nil {
		return nil, berr
	}

	r.l.Printf("running across %d worker(s)", r.conf.NWorkers)
	if 0 < r.conf.Timeout {
		r.l.Printf("timeout at %d minute(s)", r.conf.Timeout)
	}

	err := r.plan.Corpus.Par(ctx, r.conf.NWorkers,
		func(ctx context.Context, named subject.Named) error {
			return r.makeJob(b, named).Run(ctx)
		},
		func(ctx context.Context) error {
			var err error
			r.plan.Corpus, err = b.Run(ctx)
			return err
		},
	)
	return &r.plan, err
}

func (r *Runner) makeJob(b *builder.Builder, named subject.Named) *Job {
	return &Job{
		Backend: r.plan.Backend,
		Conf:    &r.conf,
		ResCh:   b.SendCh,
		Subject: &named,
	}
}

// count returns the number of individual runs this runner will do.
func (r *Runner) count() int {
	return len(r.plan.Corpus) * len(r.plan.Compilers)
}

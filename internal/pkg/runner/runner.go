// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package runner contains the part of act-tester that runs compiled harness binaries and interprets their output.
package runner

import (
	"context"
	"errors"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

var (
	// ErrNoBin occurs when a successful compile result	has no binary path attached.
	ErrNoBin = errors.New("no binary in compile result")

	// ErrConfigNil occurs when we try to construct a Runner using a nil Config.
	ErrConfigNil = errors.New("config nil")
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
		return corpus.ErrNoCorpus
	}
	return nil
}

// Run runs the runner.
func (r *Runner) Run(ctx context.Context) (*plan.Plan, error) {
	bcfg := corpus.BuilderConfig{
		Init:  r.plan.Corpus,
		NReqs: r.count(),
		Obs:   r.conf.Observer,
	}
	b, berr := corpus.NewBuilder(bcfg)
	if berr != nil {
		return nil, berr
	}
	err := r.plan.Corpus.Par(ctx,
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

func (r *Runner) makeJob(b *corpus.Builder, named subject.Named) *Job {
	return &Job{
		Backend: r.plan.Backend,
		Parser:  r.conf.Parser,
		ResCh:   b.SendCh,
		Subject: &named,
	}
}

// count returns the number of individual runs this runner will do.
func (r *Runner) count() int {
	return len(r.plan.Corpus) * len(r.plan.Compilers)
}

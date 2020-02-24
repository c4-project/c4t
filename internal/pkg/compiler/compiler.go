// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains a test-plan batch compiler.
// It relies on the existence of a single-binary compiler.
package compiler

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

var (
	// ErrConfigNil occurs when we try to build a compiler with a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrDriverNil occurs when the compiler tries to use the nil pointer as its single-compile driver.
	ErrDriverNil = errors.New("driver nil")

	// ErrNoCompilers occurs when the machine plan being used for compilation has no compilers.
	ErrNoCompilers = errors.New("no compilers on this machine")
)

// Compiler contains the configuration required to compile the harnesses for a single test run.
type Compiler struct {
	// l is the logger for this batch compiler.
	l *log.Logger

	// plan is the plan on which this batch compiler is operating.
	plan plan.Plan

	// mid is the ID of the machine on which this batch compiler is operating.
	mid model.ID

	// mach is a copy of the specific machine (in plan) on which this batch compiler is operating.
	mach plan.MachinePlan

	// conf is the configuration used to build this compiler.
	conf Config
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(c *Config, p *plan.Plan) (*Compiler, error) {
	if p == nil {
		return nil, plan.ErrNil
	}
	if err := checkConfig(c); err != nil {
		return nil, err
	}

	mid, mp, err := p.Machine(c.MachineID)
	if err != nil {
		return nil, err
	}
	if len(mp.Compilers) <= 0 {
		return nil, fmt.Errorf("%w: machine %s", ErrNoCompilers, c.MachineID.String())
	}

	return &Compiler{mid: mid, plan: *p, conf: *c, l: iohelp.EnsureLog(c.Logger), mach: mp}, nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	if c.Driver == nil {
		return ErrDriverNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	return nil
}

// Run runs the batch compiler with context ctx.
// On success, it returns an amended plan, now associating each subject with its compiler results.
func (c *Compiler) Run(ctx context.Context) (*plan.Plan, error) {
	if err := c.prepareDirs(); err != nil {
		return nil, err
	}

	eg, ectx := errgroup.WithContext(ctx)

	// Note that this builder is modifying the plan's corpus in-place; this means we have to provide the job goroutines
	// with a copy rc.
	b, reqCh, berr := corpus.NewBuilder(c.plan.Corpus, c.count())
	if berr != nil {
		return nil, berr
	}

	for ids, cc := range c.mach.Compilers {
		nc := nameCompiler(ids, cc)
		cr := c.makeJob(nc, reqCh)
		eg.Go(func() error {
			return cr.Compile(ectx)
		})
	}

	var newc corpus.Corpus
	eg.Go(func() error {
		var err error
		newc, err = b.Run(ectx)
		return err
	})

	// Need to wait until there are no goroutines accessing the corpus before we copy it over.
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	c.plan.Corpus = newc
	return &c.plan, nil
}

func (c *Compiler) prepareDirs() error {
	c.l.Println("preparing directories")
	return c.conf.Paths.Prepare(c.mach.CompilerIDs())
}

// makeJob makes a job for the named compiler nc, outputting results to resCh.
// It also takes in a read-only copy, rc, of the corpus; this is because the result handling thread will be modifying
// the corpus proper.
func (c *Compiler) makeJob(nc *model.NamedCompiler, resCh chan<- corpus.BuilderReq) *compileJob {
	return &compileJob{
		MachineID: c.mid,
		Compiler:  nc,
		Corpus:    c.plan.Corpus,
		Pathset:   c.conf.Paths,
		Runner:    c.conf.Driver,
		ResCh:     resCh,
	}
}

// nameCompiler sticks the name ids onto the compiler cc.
func nameCompiler(ids string, cc model.Compiler) *model.NamedCompiler {
	return &model.NamedCompiler{
		ID:       model.IDFromString(ids),
		Compiler: cc,
	}
}

// count gets the number of individual compilations the compiler will perform.
func (c *Compiler) count() int {
	return len(c.mach.Compilers) * len(c.plan.Corpus)
}

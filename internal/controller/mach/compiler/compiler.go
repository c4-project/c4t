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
	"log"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

var (
	// ErrConfigNil occurs when we try to build a compiler with a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrDriverNil occurs when the compiler tries to use the nil pointer as its single-compile driver.
	ErrDriverNil = errors.New("driver nil")
)

// Compiler contains the configuration required to compile the harnesses for a single test run.
type Compiler struct {
	// l is the logger for this batch compiler.
	l *log.Logger

	// plan is the plan on which this batch compiler is operating.
	plan plan.Plan

	// mid is the ID of the machine on which this batch compiler is operating.
	mid id.ID

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
	return &Compiler{plan: *p, conf: *c, l: iohelp.EnsureLog(c.Logger)}, nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	return c.Check()
}

// Run runs the batch compiler with context ctx.
// On success, it returns an amended plan, now associating each subject with its compiler results.
func (c *Compiler) Run(ctx context.Context) (*plan.Plan, error) {
	if err := c.prepareDirs(); err != nil {
		return nil, err
	}

	eg, ectx := errgroup.WithContext(ctx)

	b, err := c.makeBuilder()
	if err != nil {
		return nil, err
	}

	c.conf.Quantities.Log(c.l)

	for ids, cc := range c.plan.Compilers {
		nc, err := cc.AddNameString(ids)
		if err != nil {
			return nil, err
		}
		cr := c.makeJob(nc, b.SendCh)
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

func (c *Compiler) makeBuilder() (*builder.Builder, error) {
	bc := builder.Config{
		Init:      c.plan.Corpus,
		Observers: c.conf.Observers,
		Manifest: builder.Manifest{
			Name:  "compile",
			NReqs: c.count(),
		},
	}
	return builder.New(bc)
}

func (c *Compiler) prepareDirs() error {
	c.l.Println("preparing directories")
	cids, err := c.plan.CompilerIDs()
	if err != nil {
		return err
	}
	return c.conf.Paths.Prepare(cids)
}

// makeJob makes a job for the named compiler nc, outputting results to resCh.
// It also takes in a read-only copy, rc, of the corpus; this is because the result handling thread will be modifying
// the corpus proper.
func (c *Compiler) makeJob(nc *compiler.Named, resCh chan<- builder.Request) *Job {
	return &Job{
		MachineID: c.mid,
		Compiler:  nc,
		Corpus:    c.plan.Corpus,
		Conf:      &c.conf,
		ResCh:     resCh,
	}
}

// count gets the number of individual compilations the compiler will perform.
func (c *Compiler) count() int {
	return len(c.plan.Compilers) * len(c.plan.Corpus)
}

// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains a test-plan batch compiler.
// It relies on the existence of a single-binary compiler.
package compiler

import (
	"context"

	"github.com/c4-project/c4t/internal/stage/mach/interpreter"

	"github.com/c4-project/c4t/internal/stage/mach/observer"
	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/subject"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/plan"
)

// SubjectPather is the interface of types that can produce path sets for compilations.
type SubjectPather interface {
	// Prepare sets up the directories ready to serve through SubjectPaths.
	// It takes the compiler IDs that are to be represented in the pathset.
	Prepare(compilers ...id.ID) error

	// SubjectPaths gets the filepaths for the compilation with name sc.
	SubjectPaths(sc compilation.Name) compilation.CompileFileset

	// TODO(@MattWindsor91): should SubjectPaths return an error if the directories are not prepared?
}

//go:generate mockery --name=SubjectPather

// Compiler contains the configuration required to compile the recipes for a single test run.
type Compiler struct {
	// driver is what the compiler should use to run single compiler jobs.
	driver interpreter.Driver

	// observers observe the compiler's progress across a corpus.
	observers []observer.Observer

	// paths is the pathset for this compiler run.
	paths SubjectPather

	// quantities is this compiler stage's quantity set.
	quantities quantity.BatchSet
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(driver interpreter.Driver, paths SubjectPather, opts ...Option) (*Compiler, error) {
	if driver == nil {
		return nil, interpreter.ErrDriverNil
	}
	if paths == nil {
		return nil, iohelp.ErrPathsetNil
	}

	c := &Compiler{driver: driver, paths: paths}
	if err := Options(opts...)(c); err != nil {
		return nil, err
	}
	return c, nil
}

// Stage gets the appropriate stage record for compilation.
func (*Compiler) Stage() stage.Stage {
	return stage.Compile
}

// Close does nothing.
func (*Compiler) Close() error {
	return nil
}

// Run runs the batch compiler with context ctx and plan p.
// On success, it returns an amended plan, now associating each subject with its compiler results.
func (c *Compiler) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	if err := c.prepareDirs(p); err != nil {
		return nil, err
	}

	observer.OnCompileStart(c.quantities, c.observers...)

	newc, err := builder.ParBuild(
		ctx,
		c.quantities.NWorkers,
		p.Corpus,
		c.builderConfig(p),
		func(ctx context.Context, s subject.Named, requests chan<- builder.Request) error {
			return c.instance(requests, s, p).Compile(ctx)
		})
	if err != nil {
		return nil, err
	}

	np := *p
	np.Corpus = newc
	return &np, nil
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

func (c *Compiler) builderConfig(p *plan.Plan) builder.Config {
	return builder.Config{
		Init:      p.Corpus,
		Observers: observer.LowerToBuilder(c.observers...),
		Manifest: builder.Manifest{
			Name:  "compile",
			NReqs: p.NumExpCompilations(),
		},
	}
}

func (c *Compiler) prepareDirs(p *plan.Plan) error {
	// TODO(@MattWindsor91): port this to observers
	// c.l.Println("preparing directories")
	cids, err := p.CompilerIDs()
	if err != nil {
		return err
	}
	return c.paths.Prepare(cids...)
}

// instance makes an instance for the named compiler nc, outputting results to resCh.
// It also takes in a read-only copy, rc, of the corpus; this is because the result handling thread will be modifying
// the corpus proper.
func (c *Compiler) instance(requests chan<- builder.Request, s subject.Named, p *plan.Plan) *Instance {
	return &Instance{
		machineID:  p.Machine.ID,
		subject:    s,
		compilers:  p.Compilers,
		driver:     c.driver,
		paths:      c.paths,
		resCh:      requests,
		quantities: c.quantities,
	}
}

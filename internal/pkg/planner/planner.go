// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Source contains all of the various sources for a Planner's information.
type Source struct {
	// BProbe is the backend prober.
	BProbe BackendFinder

	// CProbe is the compiler prober.
	CProbe CompilerLister

	// SProbe is the subject prober.
	SProbe SubjectProber
}

// Planner holds all configuration for the test planner.
type Planner struct {
	// Source contains all of the various sources for a Planner's information.
	Source Source

	// Filter is the compiler filter to use to select compilers to test.
	Filter string

	// Logger is the logger used by the planner.
	Logger *log.Logger

	// Observers watch the plan's corpus being built.
	Observers []builder.Observer

	// MachineID is the identifier of the target machine for the plan.
	MachineID id.ID

	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the planner will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int
}

// plan runs the test planner p on input file paths fs.
func (p *Planner) Plan(ctx context.Context, fs []string) (*plan.Plan, error) {
	// Early out to prevent us from doing any planning if we received no files.
	if len(fs) == 0 {
		return nil, corpus.ErrNone
	}

	// TODO(@MattWindsor91): separate Planner from MachConfig to avoid this
	p.Logger = iohelp.EnsureLog(p.Logger)

	hd := plan.NewHeader()
	// TODO(@MattWindsor91): allow manual seed override
	rng := hd.Rand()

	pn := plan.Plan{
		Header:    *hd,
		Machine:   model.Machine{},
		Backend:   nil,
		Compilers: nil,
		Corpus:    nil,
	}

	var err error

	// TODO(@MattWindsor91): probe machine
	p.Logger.Println("Planning backend...")
	if pn.Backend, err = p.planBackend(ctx); err != nil {
		return nil, err
	}

	p.Logger.Println("Planning compilers...")
	if pn.Compilers, err = p.planCompilers(ctx); err != nil {
		return nil, err
	}

	p.Logger.Println("Planning corpus...")
	if pn.Corpus, err = p.planCorpus(ctx, rng, fs); err != nil {
		return nil, err
	}

	return &pn, nil
}

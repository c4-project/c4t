// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/sirupsen/logrus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BackendFinder is the interface of things that can find backends for machines.
type BackendFinder interface {
	// FindBackend asks for a backend with the given style on any one of machines,
	// or a default machine if none have such a backend.
	FindBackend(ctx context.Context, style model.ID, machines ...model.ID) (*model.Backend, error)
}

// Source is the composite interface of types that can provide the requisite information a Planner needs about
// backends, compilers, and subjects.
type Source interface {
	BackendFinder
	CompilerLister
	SubjectProber
}

// Planner holds all configuration for the test planner.
type Planner struct {
	// Source is the planner's information source.
	Source Source

	// Filter is the compiler filter to use to select compilers to test.
	Filter string

	// MachineID is the identifier of the target machine for the plan.
	MachineID model.ID

	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the planner will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int

	// InFiles is a list of paths to files that form the incoming test corpus.
	InFiles []string
}

// plan runs the test planner p.
func (p *Planner) Plan(ctx context.Context) (*plan.Plan, error) {
	// Early out to prevent us from doing any planning if we received no files.
	if len(p.InFiles) == 0 {
		return nil, corpus.ErrNoCorpus
	}

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
	logrus.Infoln("Planning backend...")
	if pn.Backend, err = p.planBackend(ctx); err != nil {
		return nil, err
	}

	logrus.Infoln("Planning compilers...")
	if pn.Compilers, err = p.planCompilers(ctx); err != nil {
		return nil, err
	}

	logrus.Infoln("Planning corpus...")
	if pn.Corpus, err = p.planCorpus(ctx, rng); err != nil {
		return nil, err
	}

	return &pn, nil
}

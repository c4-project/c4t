// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"
	"log"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Planner holds all configuration for the test planner.
type Planner struct {
	conf Config
	fs   []string
	l    *log.Logger
	mid  id.ID
	plan plan.Plan
	rng  *rand.Rand
	seed int64
}

// New constructs a new planner with the given config, machine information, files, and seed override.
// If seed is UseDateSeed, it will be ignored and a date-specific seed generated at runtime.
func New(c Config, mid id.ID, mach plan.Machine, fs []string, seed int64) (*Planner, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	// Early out to prevent us from doing any planning if we received no files.
	if len(fs) == 0 {
		return nil, corpus.ErrNone
	}

	p := Planner{
		conf: c,
		fs:   fs,
		l:    iohelp.EnsureLog(c.Logger),
		mid:  mid,
		plan: plan.Plan{
			Machine: mach,
		},
		seed: seed,
	}

	return &p, nil
}

// Plan runs the test planner p.
func (p *Planner) Plan(ctx context.Context) (*plan.Plan, error) {
	hd := plan.NewHeader(p.seed)
	p.plan.Header = *hd

	if p.rng == nil {
		p.rng = hd.Rand()
	}

	var err error

	p.l.Println("Planning backend...")
	if err = p.planBackend(ctx); err != nil {
		return nil, err
	}

	p.l.Println("Planning compilers...")
	if err = p.planCompilers(ctx); err != nil {
		return nil, err
	}

	p.l.Println("Planning corpus...")
	if err = p.planCorpus(ctx); err != nil {
		return nil, err
	}

	return &p.plan, nil
}

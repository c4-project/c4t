// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package planner contains the logic for the test planner.
package planner

import (
	"context"
	"errors"
	"log"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/model/machine"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Planner holds all configuration for the test planner.
type Planner struct {
	conf Config
	fs   []string
	l    *log.Logger
	plan plan.Plan
	rng  *rand.Rand
	seed int64
}

// ErrConfigNil occurs when we pass a nil config when creating a planner.
var ErrConfigNil = errors.New("config nil")

// New constructs a new planner with the given config, machine information, files, and seed override.
// If seed is UseDateSeed, it will be ignored and a date-specific seed generated at runtime.
func New(c *Config, mach machine.Named, fs []string, seed int64) (*Planner, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}
	// Early out to prevent us from doing any planning if we received no files.
	if len(fs) == 0 {
		return nil, corpus.ErrNone
	}

	p := Planner{
		conf: *c,
		fs:   fs,
		l:    iohelp.EnsureLog(c.Logger),
		plan: plan.Plan{
			Machine: mach,
		},
		seed: seed,
	}

	return &p, nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	return c.Check()
}

// Plan runs the test planner p.
func (p *Planner) Plan(ctx context.Context) (*plan.Plan, error) {
	hd := plan.NewMetadata(p.seed)
	p.plan.Metadata = *hd

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

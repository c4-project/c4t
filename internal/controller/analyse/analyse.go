// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analyse represents the stage of the tester that takes a plan, performs various statistics on it, and outputs
// reports.
package analyse

import (
	"context"
	"errors"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"
	"github.com/MattWindsor91/act-tester/internal/controller/analyse/saver"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Analyse represents the state of the plan analyse stage.
type Analyse struct {
	cfg  *Config
	plan *plan.Plan
	save *saver.Saver
}

// New constructs a new query runner on config c and plan p.
func New(c *Config, p *plan.Plan) (*Analyse, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	s, err := maybeNewSave(c)
	if err != nil {
		return nil, err
	}
	return &Analyse{cfg: c, plan: p, save: s}, nil
}

func maybeNewSave(c *Config) (*saver.Saver, error) {
	if c.SavedPaths == nil {
		return nil, nil
	}
	return saver.New(
		c.SavedPaths,
		func(path string) (saver.Archiver, error) {
			return saver.CreateTGZ(path)
		},
		saver.WithObservers(c.Observers...))
}

func checkConfig(c *Config) error {
	if c == nil {
		return errors.New("config nil")
	}
	return nil
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

// Run runs the query, outputting to the configured output writer.
func (q *Analyse) Run(ctx context.Context) (*plan.Plan, error) {
	a, err := q.analyse(ctx)
	if err != nil {
		return nil, err
	}

	observer.OnAnalysis(*a, q.cfg.Observers...)

	if q.save != nil {
		if err := q.save.Run(*a); err != nil {
			return nil, err
		}
	}

	return q.plan, nil
}

func (q *Analyse) analyse(ctx context.Context) (*analysis.Analysis, error) {
	ar, err := NewAnalyser(q.plan, q.cfg.NWorkers)
	if err != nil {
		return nil, err
	}
	return ar.Analyse(ctx)
}

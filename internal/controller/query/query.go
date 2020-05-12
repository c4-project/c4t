// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package query represents the stage of the tester that takes a plan, performs various statistics on it, and outputs
// reports.
package query

import (
	"context"
	"errors"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Query represents the state of the plan querying tool.
type Query struct {
	cfg  *Config
	plan *plan.Plan
	aw   *AnalysisWriter
}

// New constructs a new query runner on config c and plan p.
func New(c *Config, p *plan.Plan) (*Query, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	aw, err := NewAnalysisWriter(c)
	if err != nil {
		return nil, err
	}

	return &Query{cfg: c, plan: p, aw: aw}, nil
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
func (q *Query) Run(ctx context.Context) (*plan.Plan, error) {
	// TODO(@MattWindsor91): allow customisation of nworkers here
	a, err := analysis.Analyse(ctx, q.plan, 20)
	if err != nil {
		return nil, err
	}
	// TODO(@MattWindsor91): merge this with the analysis/saving step in the test director.
	err = q.aw.Write(a)
	return q.plan, err
}

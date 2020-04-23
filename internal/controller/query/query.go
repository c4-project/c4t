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
	"io"
	"text/template"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Query represents the state of the plan querying tool.
type Query struct {
	cfg      *Config
	plan     *plan.Plan
	w        io.Writer
	planTmpl *template.Template
}

// New constructs a new query runner on config c and plan p.
func New(c *Config, p *plan.Plan) (*Query, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	t, err := getTemplate()
	if err != nil {
		return nil, err
	}
	return &Query{cfg: c, plan: p, planTmpl: t, w: iohelp.EnsureWriter(c.Out)}, nil
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
func (q *Query) Run(_ context.Context) (*plan.Plan, error) {
	err := q.planTmpl.ExecuteTemplate(q.w, "plan", q.plan)
	return q.plan, err
}

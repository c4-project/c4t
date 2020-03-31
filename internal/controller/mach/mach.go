// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mach contains the 'machine-dependent' stage of the tester.
// This stage encapsulates the batch-compile and run stages of the tester, and provides common infrastructure for both.
package mach

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/runner"
)

// Mach encapsulates the state needed for the machine-dependent stage.
type Mach struct {
	// compiler is, if non-nil, the configured compiler substage config.
	compiler *compiler.Config
	// plan is the plan to use when running the machine-dependent stage.
	plan *plan.Plan
	// runner is, if non-nil, the configured runner substage config.
	runner *runner.Config
	// json is, if non-nil, a JSON observer;
	// it exists here so that, if we're using JSON mode, errors get trapped and sent over as JSON.
	json *forward.Observer
}

func New(c *Config, p *plan.Plan) (*Mach, error) {
	if c == nil {
		return nil, errors.New("config nil")
	}
	if err := c.Check(); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, plan.ErrNil
	}
	m := Mach{
		compiler: c.makeCompilerConfig(),
		runner:   c.makeRunnerConfig(),
		plan:     p,
	}
	if c.JsonStatus {
		m.json = &forward.Observer{Encoder: json.NewEncoder(c.Stderr)}
	}

	return &m, nil
}

// trap checks to see if this mach is in JSON mode; if it is, it swallows the error and sends it as a JSON message.
func (m *Mach) trap(err error) error {
	if err == nil {
		return nil
	}
	if m.json != nil {
		m.json.Error(err)
		return nil
	}
	return err
}

func (m *Mach) Run(ctx context.Context) (*plan.Plan, error) {
	p, err := m.runInner(ctx)
	return p, m.trap(err)
}

func (m *Mach) runInner(ctx context.Context) (*plan.Plan, error) {
	cp, err := m.runCompiler(ctx, m.plan)
	if err != nil {
		return nil, fmt.Errorf("while running compiler: %w", err)
	}
	rp, err := m.runRunner(ctx, cp)
	if err != nil {
		return nil, fmt.Errorf("while running runner: %w", err)
	}
	return rp, nil
}

// runCompiler runs the batch compiler on plan p, if available.
// If the compiler is nil, runCompiler returns p unmodified.
func (m *Mach) runCompiler(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if m.compiler == nil {
		return p, nil
	}
	return m.compiler.Run(ctx, p)
}

// runRunner runs the batch runner on plan p.
// If c is nil, runRunner returns p unmodified.
func (m *Mach) runRunner(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if m.runner == nil {
		return p, nil
	}
	return m.runner.Run(ctx, p)
}

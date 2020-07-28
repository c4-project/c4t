// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mach contains the 'machine-dependent' stage of the tester.
// This stage encapsulates the batch-compile and run stages of the tester, and provides common infrastructure for both.
package mach

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/plan"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/runner"
)

// Mach encapsulates the state needed for the machine-dependent stage.
type Mach struct {
	// coptions is the set of options used to configure the compiler.
	coptions []compiler.Option
	// roptions is the set of options used to configure the runner.
	roptions []runner.Option
	// skipCompiler is true if the compiler should be skipped.
	skipCompiler bool
	// skipRunner is true if the runner should be skipped.
	skipRunner bool
	// path is the output directory path for both substages.
	path string

	// compiler is, if non-nil, the configured compiler substage.
	compiler *compiler.Compiler
	// runner is, if non-nil, the configured runner substage.
	runner *runner.Runner
	// json is, if non-nil, a JSON observer;
	// it exists here so that, if we're using JSON mode, errors get trapped and sent over as JSON.
	fwd *forward.Observer
}

func New(cdriver compiler.SingleRunner, rdriver runner.ObsParser, opts ...Option) (*Mach, error) {
	// The respective constructors will check that cdriver and rdriver are ok.

	m := &Mach{}
	// Options can introduce compiler and runner options, so they need to run before the compiler/runner constructors.
	if err := Options(opts...)(m); err != nil {
		return nil, err
	}
	return m, m.makeCompilerAndRunner(cdriver, rdriver)
}

func (m *Mach) makeCompilerAndRunner(cdriver compiler.SingleRunner, rdriver runner.ObsParser) error {
	if err := m.makeCompiler(cdriver); err != nil {
		return err
	}
	return m.makeRunner(rdriver)
}

func (m *Mach) makeCompiler(driver compiler.SingleRunner) error {
	if m.skipCompiler {
		return nil
	}
	var err error
	ps := compiler.NewPathset(m.path)
	m.compiler, err = compiler.New(driver, ps, m.coptions...)
	return err
}

func (m *Mach) makeRunner(driver runner.ObsParser) error {
	if m.skipRunner {
		return nil
	}
	var err error
	ps := runner.NewPathset(m.path)
	m.runner, err = runner.New(driver, ps, m.roptions...)
	return err
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

// trap checks to see if this mach is in JSON mode; if it is, it swallows the error and sends it as a JSON message.
func (m *Mach) trap(err error) error {
	if err == nil {
		return nil
	}
	if m.fwd != nil {
		m.fwd.Error(err)
		return nil
	}
	return err
}

func (m *Mach) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}
	p, err := m.runInner(ctx, p)
	return p, m.trap(err)
}

func (m *Mach) runInner(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	cp, err := m.runCompiler(ctx, p)
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

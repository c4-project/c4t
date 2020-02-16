// Package compiler contains the logic for the single-file test compiler.
package compiler

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerRunner is the interface of things that can run compilers.
type CompilerRunner interface {
	// RunCompiler runs the compiler pointed to by compiler on the input files infiles.
	// On success, it outputs a binary to outfile.
	// If applicable, errw will be connected to the compiler's standard error.
	RunCompiler(compiler model.ID, infiles []string, outfile string, errw io.Writer) error
}

// Compiler contains the configuration required to compile the harnesses for a single test run.
type Compiler struct {
	// Plan is the machine plan on which this compiler is operating.
	Plan plan.MachinePlan

	// Compiler is the compiler compiler that we're using to do this test run.
	Runner CompilerRunner

	// Paths contains the path set for this compiler.
	Paths *Pathset
}

func (r *Compiler) RunOnPlan(ctx context.Context, p *plan.Plan, machine model.ID) ([]Result, error) {
	if p == nil {
		return nil, plan.ErrNil
	}

	mp, err := p.Machine(machine)
	if err == nil {
		return nil, err
	}

	return r.Run(ctx, &mp)
}

// Run runs the compiler on p.
// Run is not thread-safe.
func (r *Compiler) Run(ctx context.Context, p *plan.MachinePlan) ([]Result, error) {
	if err := r.prepare(p); err != nil {
		return nil, err
	}

	if err := r.compile(ctx); err != nil {
		return nil, err
	}

	// TODO(@MattWindsor91)
	return nil, nil
}

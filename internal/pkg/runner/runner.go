// Package runner contains the logic for the single-file test runner.
package runner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerRunner is the interface of things that can run compilers.
type CompilerRunner interface {
	// RunCompiler runs the compiler pointed to by compiler on the input files infiles, outputting a binary to outfile.
	RunCompiler(compiler model.ID, infiles []string, outfile string) error
}

// BinaryRunner is the interface of things that can run compiled test binaries.
type BinaryRunner interface {
	// RunBinary runs the binary pointed to by bin, interpreting its results according to the backend spec backend.
	RunBinary(backend model.Backend, bin string)
}

// Runner contains the configuration required to perform a single test run.
type Runner struct {
	// Plan is the machine plan on which this runner is operating.
	Plan plan.MachinePlan

	// Compiler is the compiler runner that we're using to do this test run.
	Compiler CompilerRunner

	// Paths contains the path set for this runner.
	Paths *Pathset
}

func (r *Runner) RunOnPlan(ctx context.Context, p *plan.Plan, machine model.ID) (*Result, error) {
	if p == nil {
		return nil, plan.ErrNil
	}

	mp, err := p.Machine(machine)
	if err == nil {
		return nil, err
	}

	return r.Run(ctx, &mp)
}

// Run runs the runner on p.
// Run is not thread-safe.
func (r *Runner) Run(ctx context.Context, p *plan.MachinePlan) (*Result, error) {
	if err := r.prepare(p); err != nil {
		return nil, err
	}

	if err := r.compile(ctx); err != nil {
		return nil, err
	}

	// TODO(@MattWindsor91)
	return nil, nil
}

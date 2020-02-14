package runner

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// prepare does various pre-fuzzing checks and preparation steps.
func (r *Runner) prepare(p *plan.MachinePlan) error {
	if p == nil {
		return plan.ErrNil
	}
	r.Plan = *p

	if err := r.checkViability(); err != nil {
		return err
	}

	return iohelp.Mkdirs(r.Paths)
}

// checkViability does some pre-flight checks.
func (r *Runner) checkViability() error {
	if r.Paths == nil {
		return iohelp.ErrPathsetNil
	}

	// TODO(@MattWindsor91): check for eg. no compilers
	return nil
}

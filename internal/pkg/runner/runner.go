// Package runner contains the part of act-tester that runs compiled harness binaries and interprets their output.
package runner

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// Runner contains information necessary to run a plan's compiled test cases.
type Runner struct {
	conf Config
	plan plan.Plan
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(c *Config, p *plan.Plan) (*Runner, error) {
	r := Runner{conf: *c, plan: *p}
	return &r, nil
}

// Run runs the runner.
func (r *Runner) Run() error {
	return nil
}

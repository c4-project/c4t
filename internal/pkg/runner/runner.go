// Package runner contains the part of act-tester that runs compiled harness binaries and interprets their output.
package runner

import (
	"context"
	"errors"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// ErrConfigNil occurs when we try to construct a Runner using a nil Config.
var ErrConfigNil = errors.New("config nil")

// Runner contains information necessary to run a plan's compiled test cases.
type Runner struct {
	conf Config
	l    *log.Logger
	plan plan.Plan
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(c *Config, p *plan.Plan) (*Runner, error) {
	if c == nil {
		return nil, ErrConfigNil
	}
	if p == nil {
		return nil, plan.ErrNil
	}

	r := Runner{conf: *c, plan: *p, l: iohelp.EnsureLog(c.Logger)}

	if err := r.check(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Runner) check() error {
	if len(r.plan.Corpus) == 0 {
		return subject.ErrNoCorpus
	}
	return nil
}

// Run runs the runner.
func (r *Runner) Run(ctx context.Context) error {
	for _, s := range r.plan.Corpus {
		if err := r.runSubject(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) runSubject(ctx context.Context, s subject.Subject) error {
	r.l.Println("running subject:", s.Name)
	return nil
}

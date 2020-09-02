// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package analyser represents the stage of the tester that takes a plan, performs various statistics on it, and outputs
// reports.
package analyser

import (
	"context"
	"errors"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/plan/analysis"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Analyser represents the state of the plan analyser stage.
type Analyser struct {
	// errOnBadStatus makes the plan stage return an error if the analyser saw 'bad' statuses.
	errOnBadStatus bool
	// savePaths is the set of paths to which we save failing corpora.
	savePaths *saver.Pathset
	// aopts is the set of options to pass to the underlying analyser.
	aopts []analysis.Option
	// observers is the list of observers to which analyses are sent.
	observers []Observer
	// saveObservers is the list of observers to which archival operations are sent.
	saveObservers []saver.Observer
}

// New constructs a new analyser stage on plan p, with options opts.
func New(opts ...Option) (*Analyser, error) {
	an := new(Analyser)
	err := Options(opts...)(an)
	return an, err
}

func (a *Analyser) newSaver() (*saver.Saver, error) {
	if a.savePaths == nil {
		return nil, nil
	}
	return saver.New(
		a.savePaths,
		func(path string) (saver.Archiver, error) {
			return saver.CreateTGZ(path)
		},
		saver.ObserveWith(a.saveObservers...))
}

// Run runs the analyser on the plan p, outputting to the configured output writer.
func (a *Analyser) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if err := checkPlan(p); err != nil {
		return nil, err
	}

	an, err := a.analyse(ctx, p)
	if err != nil {
		return nil, err
	}

	OnAnalysis(*an, a.observers...)

	if err := a.maybeSave(an); err != nil {
		return nil, err
	}

	return an.Plan, a.statusErr(an)
}

func checkPlan(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	return p.Check()
}

func (a *Analyser) maybeSave(an *analysis.Analysis) error {
	save, err := a.newSaver()
	// save can be nil if we're not supposed to be saving.
	if err != nil || save == nil {
		return err
	}
	return save.Run(*an)
}

func (a *Analyser) analyse(ctx context.Context, p *plan.Plan) (*analysis.Analysis, error) {
	return analysis.Analyse(ctx, p, a.aopts...)
}

// ErrBadStatus is the error reported when the analyser is asked to error on a bad status, and one arrives.
var ErrBadStatus = errors.New("at least one subject reported a bad status")

func (a *Analyser) statusErr(an *analysis.Analysis) error {
	if a.errOnBadStatus && an.HasBadOutcomes() {
		return ErrBadStatus
	}
	return nil
}

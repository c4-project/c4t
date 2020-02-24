// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package ux

import (
	"context"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// StdinFile is the special file path that the plan loader treats as a request to load from stdin instead.
const StdinFile = "-"

// Load loads a plan pointed to by f.
// If f is empty or StdinFile, Load loads from standard input instead.
func LoadPlan(f string) (*plan.Plan, error) {
	var p plan.Plan

	if f == "" || f == StdinFile {
		_, err := toml.DecodeReader(os.Stdin, &p)
		return nil, err
	}
	_, err := toml.DecodeFile(f, &p)
	return &p, err
}

// PlanRunner is the interface of parts of the tester that transform plans.
type PlanRunner interface {
	// Run runs this type's processing stage on the plan pointed to by p.
	// It also takes a context, which can be used to cancel the process.
	// It returns an updated plan (which may or may not be p edited in-place), or an error.
	Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error)
}

// RunOnPlanFile runs r on the plan pointed to by inf, dumping the resulting plan to stdout.
func RunOnPlanFile(ctx context.Context, r PlanRunner, inf string) error {
	p, perr := LoadPlan(inf)
	if perr != nil {
		return perr
	}
	q, qerr := r.Run(ctx, p)
	if qerr != nil {
		return qerr
	}
	return q.Dump()
}

// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package ux

import (
	"context"
	"io"
	"os"

	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"

	"github.com/c4-project/c4t/internal/plan"
)

// StdinFile is the special file path that the plan loader treats as a request to load from stdin instead.
const StdinFile = "-"

// Load loads a plan pointed to by f.
// If f is empty or StdinFile, Load loads from standard input instead.
func LoadPlan(f string) (*plan.Plan, error) {
	var (
		p   plan.Plan
		err error
	)

	if f == "" || f == StdinFile {
		err = plan.Read(os.Stdin, &p)
	} else {
		err = plan.ReadFile(f, &p)
	}
	return &p, err
}

// RunOnPlanFile runs r on the plan pointed to by inf, dumping the resulting plan to outw.
func RunOnPlanFile(ctx context.Context, r plan.Runner, inf string, outw io.Writer) error {
	p, perr := LoadPlan(inf)
	if perr != nil {
		return perr
	}

	q, qerr := p.RunStage(ctx, r)
	if qerr != nil {
		return qerr
	}

	// There might not be a plan to output; this can happen if an error was handled/trapped earlier.
	if q == nil {
		return nil
	}

	return q.Write(outw, plan.WriteHuman)
}

// RunOnCliPlan runs r on the plan pointed to by the arguments of ctx, dumping the resulting plan to outw.
func RunOnCliPlan(ctx *c.Context, r plan.Runner, outw io.Writer) error {
	pf, err := stdflag.PlanFileFromCli(ctx)
	if err != nil {
		return err
	}
	return RunOnPlanFile(ctx.Context, r, pf, outw)
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package view

import (
	"context"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
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
		_, err = toml.DecodeReader(os.Stdin, &p)
	} else {
		_, err = toml.DecodeFile(f, &p)
	}
	return &p, err
}

// RunOnPlanFile runs r on the plan pointed to by inf, dumping the resulting plan to outw.
func RunOnPlanFile(ctx context.Context, r plan.Runner, inf string, outw io.Writer) error {
	p, perr := LoadPlan(inf)
	if perr != nil {
		return perr
	}
	q, qerr := r.Run(ctx, p)
	if qerr != nil {
		return qerr
	}

	// There might not be a plan to output; this can happen if an error was handled/trapped earlier.
	if q == nil {
		return nil
	}

	return q.Dump(outw)
}

// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"context"

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/service"
)

// Instance holds all state needed to do one run of Litmus.
type Instance struct {
	// Job contains the lifting job being processed.
	Job backend.LiftJob

	// RunInfo contains the appropriately resolved run info for running Litmus.
	RunInfo service.RunInfo

	// Runner is the thing that actually runs Litmus.
	Runner service.Runner
}

// Run runs the litmus wrapper according to the configuration c.
func (l *Instance) Run(ctx context.Context) error {
	// TODO(@MattWindsor91): delitmus support

	if err := l.Job.Check(); err != nil {
		return err
	}

	var f Fixset
	f.PopulateFromLitmus(l.Job.In.Litmus)

	if err := l.runLitmus(ctx, f); err != nil {
		return err
	}

	return l.patch(f)
}

func (l *Instance) runLitmus(ctx context.Context, f Fixset) error {
	ji, err := l.jobRunInfo(f)
	if err != nil {
		return err
	}
	return l.Runner.Run(ctx, ji)
}

// jobRunInfo gets the command and arguments needed to run Litmus on the specific job given.
func (l *Instance) jobRunInfo(f Fixset) (service.RunInfo, error) {
	// This is distinct from the earlier override that overlays user-specified Litmus args on the default ones.
	ri := l.RunInfo
	args, err := l.litmusArgs(f)
	if err != nil {
		return ri, err
	}
	ri.Override(service.RunInfo{Args: args})
	return ri, nil
}

// litmusArgs works out the argument vector for Litmus.
func (l *Instance) litmusArgs(f Fixset) ([]string, error) {
	args := append(f.Args(), "-o", l.Job.Out.Dir)

	cargs, err := litmusCommonArgs(l.Job)
	if err != nil {
		return nil, err
	}

	return append(args, cargs...), nil
}

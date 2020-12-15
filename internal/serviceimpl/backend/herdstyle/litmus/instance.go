// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"context"

	"github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/model/service"
)

// Instance holds all state needed to do one run of Litmus.
type Instance struct {
	// Job contains the lifting job being processed.
	Job backend.LiftJob

	// RunInfo contains the appropriately resolved run info for running Litmus.
	RunInfo service.RunInfo

	// Runner is the thing that actually runs Litmus.
	Runner service.Runner

	// Fixset is the set of enabled fixes.
	// It is part of the config to allow the forcing of fixes that the shim would otherwise deem unnecessary.
	Fixset Fixset
}

// Run runs the litmus wrapper according to the configuration c.
func (l *Instance) Run(ctx context.Context) error {
	// TODO(@MattWindsor91): delitmus support

	if err := l.Job.Check(); err != nil {
		return err
	}

	ji, err := l.jobRunInfo()
	if err != nil {
		return err
	}
	if err := l.Runner.Run(ctx, ji); err != nil {
		return err
	}

	// TODO(@MattWindsor91): this probably needs some work
	if l.Job.In.Litmus.IsC() {
		l.Fixset.PopulateFromStats(l.Job.In.Litmus.Stats)
		return l.patch()
	}
	return nil
}

// jobRunInfo gets the command and arguments needed to run Litmus on the specific job given.
func (l *Instance) jobRunInfo() (service.RunInfo, error) {
	// This is distinct from the earlier override that overlays user-specified Litmus args on the default ones.
	ri := l.RunInfo
	args, err := l.litmusArgs()
	if err != nil {
		return ri, err
	}
	ri.Override(service.RunInfo{Args: args})
	return ri, nil
}

// litmusArgs works out the argument vector for Litmus.
func (l *Instance) litmusArgs() ([]string, error) {
	args := l.Fixset.Args()
	args = append(args, "-o", l.Job.Out.Dir)

	cargs, err := litmusCommonArgs(l.Job)
	if err != nil {
		return nil, err
	}

	return append(args, cargs...), nil
}

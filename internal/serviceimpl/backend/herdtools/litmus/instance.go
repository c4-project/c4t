// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Instance holds all state needed to do one run of Litmus.
type Instance struct {
	// Job contains the lifting job being processed.
	Job job.Lifter

	// RunInfo contains the appropriately resolved run info for running Litmus.
	RunInfo service.RunInfo

	// Runner is the thing that actually runs Litmus.
	Runner service.Runner

	// Fixset is the set of enabled fixes.
	// It is part of the config to allow the forcing of fixes that the shim would otherwise deem unnecessary.
	Fixset Fixset

	// Verbose toggles various 'verbose' dumping actions.
	Verbose bool
}

// Run runs the litmus wrapper according to the configuration c.
func (l *Instance) Run(ctx context.Context) error {
	// TODO(@MattWindsor91): delitmus support

	if err := l.check(); err != nil {
		return err
	}

	l.Fixset.PopulateFromStats(l.Job.In.Stats)

	ji, err := l.jobRunInfo()
	if err != nil {
		return err
	}
	if err := l.Runner.Run(ctx, ji); err != nil {
		return err
	}

	return l.patch()
}

// check checks that the configuration makes sense.
func (l *Instance) check() error {
	return l.Job.Check()
}

// jobRunInfo gets the command and arguments needed to run Litmus on the specific job given.
func (l *Instance) jobRunInfo() (service.RunInfo, error) {
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
	carch, err := lookupArch(l.Job.Arch)
	if err != nil {
		return nil, fmt.Errorf("when looking up -carch: %w", err)
	}

	args := l.Fixset.Args()
	args = append(args, "-o", l.Job.OutDir, "-carch", carch, "-c11", "true", l.Job.In.Path)
	return args, nil
}

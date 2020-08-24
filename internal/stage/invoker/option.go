// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/stage/invoker/runner"

	"github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/remote"
)

// Option is the type of options for the invoker.
type Option func(*Invoker) error

// Options bundles the separate options ops into a single option.
func Options(ops ...Option) Option {
	return func(r *Invoker) error {
		for _, op := range ops {
			if err := op(r); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveMachWith adds each observer given to the invoker's machine observer pool.
func ObserveMachWith(obs ...observer.Observer) Option {
	return func(r *Invoker) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		r.machObservers = append(r.machObservers, obs...)
		return nil
	}
}

// ObserveCopiesWith adds each observer given to the invoker's copy observer pool.
func ObserveCopiesWith(obs ...copier.Observer) Option {
	return func(r *Invoker) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		r.copyObservers = append(r.copyObservers, obs...)
		return nil
	}
}

// UsePlanSSH sets the invoker up to read any SSH configuration from the first plan it receives, and, if needed, open
// a SSH connection to use for that and subsequent invocations.
func UsePlanSSH(gc *remote.Config) Option {
	return func(r *Invoker) error {
		return MakeRunnersWith(
			&runner.FromPlanFactory{
				LocalRoot: r.dirLocal,
				Config:    gc,
			},
		)(r)
	}
}

// UseSSH opens a SSH connection according to gc and mc, and sets the invoker up so that it invokes the machine node
// through that connection.
//
// If mc is nil, UseSSH is a no-op.
func UseSSH(gc *remote.Config, mc *remote.MachineConfig) Option {
	return func(r *Invoker) error {
		if mc == nil {
			return nil
		}

		sr, err := runner.NewRemoteFactory(r.dirLocal, gc, mc)
		if err != nil {
			return err
		}
		return MakeRunnersWith(sr)(r)
	}
}

// MakeRunnersWith sets the invoker to use rf to build runners.
func MakeRunnersWith(rf runner.Factory) Option {
	return func(r *Invoker) error {
		r.rfac = rf
		return nil
	}
}

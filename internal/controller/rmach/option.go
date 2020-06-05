// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/controller/rmach/runner"

	copy2 "github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/remote"
)

var (
	// ErrObserverNil occurs when we try to pass a nil observer as an option.
	ErrObserverNil = errors.New("observer nil")
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

// ObserveWith adds each observer given to the invoker's observer pools.
func ObserveWith(obs ...Observer) Option {
	return func(r *Invoker) error {
		r.observers.Append(NewObserverSet(obs...))
		return nil
	}
}

// ObserveCopiesWith adds each observer given to the invoker's copy observer pool.
func ObserveCopiesWith(obs ...copy2.Observer) Option {
	return func(r *Invoker) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
			r.observers.Copy = append(r.observers.Copy, o)
		}
		return nil
	}
}

// ObserveCorpusWith adds each observer given to the invoker's corpus observer pool.
func ObserveCorpusWith(obs ...builder.Observer) Option {
	return func(r *Invoker) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
			r.observers.Corpus = append(r.observers.Corpus, o)
		}
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

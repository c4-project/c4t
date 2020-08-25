// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/MattWindsor91/act-tester/internal/quantity"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/copier"
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

// OverrideBaseQuantities overrides the base quantity set with qs.
func OverrideBaseQuantities(qs quantity.MachNodeSet) Option {
	return func(r *Invoker) error {
		r.baseQuantities.Override(qs)
		return nil
	}
}

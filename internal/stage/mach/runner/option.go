// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/MattWindsor91/act-tester/internal/quantity"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
)

// Option is the type of functional options for the machine node.
type Option func(*Runner) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(m *Runner) error {
		for _, o := range opts {
			if err := o(m); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith adds each observer in obs to the runner's observer list.
func ObserveWith(obs ...observer.Observer) Option {
	return func(r *Runner) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		r.observers = append(r.observers, obs...)
		return nil
	}
}

// OverrideQuantities overrides this runner's quantities with qs.
func OverrideQuantities(qs quantity.BatchSet) Option {
	return func(r *Runner) error {
		r.quantities.Override(qs)
		return nil
	}
}

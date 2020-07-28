// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"github.com/MattWindsor91/act-tester/internal/stage/mach/quantity"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
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
func ObserveWith(obs ...builder.Observer) Option {
	return func(r *Runner) error {
		var err error
		r.observers, err = builder.AppendObservers(r.observers, obs...)
		return err
	}
}

// OverrideQuantities overrides this runner's quantities with qs.
func OverrideQuantities(qs quantity.SingleSet) Option {
	return func(r *Runner) error {
		r.quantities.Override(qs)
		return nil
	}
}

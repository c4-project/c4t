// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/MattWindsor91/act-tester/internal/quantity"
)

// Option is the type of options to the Planner constructor.
type Option func(*Planner) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(p *Planner) error {
		for _, o := range opts {
			if err := o(p); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith adds each observer in obs to the observer pool.
func ObserveWith(obs ...Observer) Option {
	return func(p *Planner) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		p.observers = append(p.observers, obs...)
		return nil
	}
}

// OverrideQuantities overrides this planner's quantities with qs.
func OverrideQuantities(qs quantity.PlanSet) Option {
	return func(p *Planner) error {
		p.quantities.Override(qs)
		return nil
	}
}

// FilterCompilers sets the glob used for filtering compilers.
func FilterCompilers(filter string) Option {
	return func(p *Planner) error {
		p.filter = filter
		return nil
	}
}

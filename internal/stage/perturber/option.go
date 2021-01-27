// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package perturber

import (
	"github.com/c4-project/c4t/internal/observing"
	"github.com/c4-project/c4t/internal/quantity"
)

// Option is the type of options to the Planner constructor.
type Option func(*Perturber) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(p *Perturber) error {
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
	return func(p *Perturber) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		p.observers = append(p.observers, obs...)
		return nil
	}
}

// OverrideQuantities overrides this planner's quantities with qs.
func OverrideQuantities(qs quantity.PerturbSet) Option {
	return func(p *Perturber) error {
		p.quantities.Override(qs)
		return nil
	}
}

// UseSeed overrides the seed used by the planner.
// If seed is UseDateSeed, a date-specific seed is generated at runtime.
func UseSeed(seed int64) Option {
	return func(p *Perturber) error {
		p.seed = seed
		return nil
	}
}

// UseFullCompilerIDs sets whether the perturber should name compilers with their 'full ID'.
// By default, it preserves their original IDs.
func UseFullCompilerIDs(use bool) Option {
	return func(p *Perturber) error {
		p.useFullIDs = use
		return nil
	}
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"log"
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
		return p.observers.Add(obs...)
	}
}

// OverrideQuantities overrides this planner's quantities with qs.
func OverrideQuantities(qs QuantitySet) Option {
	return func(p *Planner) error {
		p.quantities.Override(qs)
		return nil
	}
}

// LogWith sets the logger for the planner to l.
func LogWith(l *log.Logger) Option {
	// TODO(@MattWindsor91): it goes without saying, but this logger should be replaced with observer calls.
	return func(p *Planner) error {
		// EnsureLog is called after all options are parsed.
		p.l = l
		return nil
	}
}

// UseSeed overrides the seed used by the planner.
// If seed is UseDateSeed, a date-specific seed is generated at runtime.
func UseSeed(seed int64) Option {
	return func(p *Planner) error {
		p.seed = seed
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

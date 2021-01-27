// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/c4-project/c4t/internal/observing"
	"github.com/c4-project/c4t/internal/quantity"
	"github.com/c4-project/c4t/internal/stage/mach/observer"
)

// Option is the type of options to the compiler sub-stage constructor.
type Option func(*Compiler) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(c *Compiler) error {
		for _, o := range opts {
			if err := o(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith adds each observer in obs to the runner's observer list.
func ObserveWith(obs ...observer.Observer) Option {
	return func(c *Compiler) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		c.observers = append(c.observers, obs...)
		return nil
	}
}

// OverrideQuantities overrides this runner's quantities with qs.
func OverrideQuantities(qs quantity.BatchSet) Option {
	return func(c *Compiler) error {
		c.quantities.Override(qs)
		return nil
	}
}

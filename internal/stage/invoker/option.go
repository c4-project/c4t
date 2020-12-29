// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/c4-project/c4t/internal/observing"
	"github.com/c4-project/c4t/internal/quantity"
	"github.com/c4-project/c4t/internal/stage/mach/observer"

	"github.com/c4-project/c4t/internal/copier"
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

// OverrideQuantitiesFromPlanThen tells the invoker to override its base quantity set with the quantities in the
// incoming plan, and then override them again using qs.
//
// This is intended for single-shot uses of the invoker, where there is no pre-calculation of the quantity set;
// in the director form of the invoker, the director will cache the expected quantity set, and there is no need to
// consult the plan.
func OverrideQuantitiesFromPlanThen(qs quantity.MachNodeSet) Option {
	return func(r *Invoker) error {
		r.pqo = &ConfigPlanQuantityOverrider{PostPlanOverrides: qs}
		return nil
	}
}

// AllowReinvoke sets whether the invoker should allow the re-invocation of plans that have already been invoked.
func AllowReinvoke(allow bool) Option {
	return func(r *Invoker) error {
		r.allowReinvoke = allow
		return nil
	}
}

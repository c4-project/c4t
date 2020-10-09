// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import "github.com/MattWindsor91/act-tester/internal/helper/iohelp"

// Option is the type of options to supply to the coverage testbed maker's constructor.
type Option func(*Maker) error

// Options applies each option in opts successively.
func Options(opts ...Option) Option {
	return func(maker *Maker) error {
		for _, o := range opts {
			if err := o(maker); err != nil {
				return err
			}
		}
		return nil
	}
}

// OverrideQuantities overrides the maker's quantity set with qs.
func OverrideQuantities(qs QuantitySet) Option {
	return func(maker *Maker) error {
		maker.qs.Override(qs)
		return nil
	}
}

func AddInputs(paths ...string) Option {
	return func(maker *Maker) error {
		ps, err := iohelp.ExpandMany(paths)
		if err != nil {
			return err
		}
		maker.inputs = append(maker.inputs, ps...)

		return nil
	}
}

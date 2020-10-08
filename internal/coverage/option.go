// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import "errors"

var (
	// ErrConfigNil is produced when we supply a null pointer to OptionsFromConfig.
	ErrConfigNil = errors.New("supplied config is nil")
)

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

// OptionsFromConfig sets various options according to the values in Config.
func OptionsFromConfig(cfg *Config) Option {
	return func(maker *Maker) error {
		if cfg == nil {
			return ErrConfigNil
		}
		return Options(
			OverrideQuantities(cfg.Quantities),
		)(maker)
	}
}

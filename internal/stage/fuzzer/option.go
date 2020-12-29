// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"github.com/c4-project/c4t/internal/model/service/fuzzer"
	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"
)

// Option is the type of options passed to the fuzzer constructor.
type Option func(*Fuzzer) error

// Options bundles the separate options ops into a single option.
func Options(ops ...Option) Option {
	return func(r *Fuzzer) error {
		for _, op := range ops {
			if err := op(r); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith adds each observer given to the invoker's observer pools.
func ObserveWith(obs ...builder.Observer) Option {
	return func(r *Fuzzer) error {
		r.observers = append(r.observers, obs...)
		return nil
	}
}

// OverrideQuantities overrides the fuzzer's quantities with qs.
func OverrideQuantities(qs quantity.FuzzSet) Option {
	return func(f *Fuzzer) error {
		f.quantities.Override(qs)
		return nil
	}
}

// UseConfig populates settings for the fuzzer from the configuration cfg.
func UseConfig(cfg *fuzzer.Configuration) Option {
	// TODO(@MattWindsor91): this should probably install specific settings instead of copying itself.
	return func(f *Fuzzer) error {
		f.config = cfg
		return nil
	}
}

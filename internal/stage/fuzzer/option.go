// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"log"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
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

// LogWith sets the fuzzer's logger to l.
func LogWith(l *log.Logger) Option {
	// TODO(@MattWindsor91): replace logger with observer
	return func(f *Fuzzer) error {
		f.l = l
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

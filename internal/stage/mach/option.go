// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"errors"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/stage/mach/forward"

	"github.com/c4-project/c4t/internal/stage/mach/compiler"
	"github.com/c4-project/c4t/internal/stage/mach/runner"
)

// Option is the type of functional options.
type Option func(*Mach) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(m *Mach) error {
		for _, o := range opts {
			if err := o(m); err != nil {
				return err
			}
		}
		return nil
	}
}

// OverrideQuantities overrides the compiler and runner quantities with qs.
func OverrideQuantities(qs quantity.MachNodeSet) Option {
	return Options(
		WithCompilerOptions(compiler.OverrideQuantities(qs.Compiler)),
		WithRunnerOptions(runner.OverrideQuantities(qs.Runner)),
	)
}

// WithCompilerOptions adds opts to the set of options used to configure the compiler.
func WithCompilerOptions(opts ...compiler.Option) Option {
	return func(m *Mach) error {
		m.coptions = append(m.coptions, opts...)
		return nil
	}
}

// WithRunnerOptions adds opts to the set of options used to configure the runner.
func WithRunnerOptions(opts ...runner.Option) Option {
	return func(m *Mach) error {
		m.roptions = append(m.roptions, opts...)
		return nil
	}
}

// ForwardTo tells the machine node to observe with, and forward errors to, fwd.
func ForwardTo(fwd *forward.Observer) Option {
	return func(m *Mach) error {
		if fwd == nil {
			return errors.New("forward observer nil")
		}
		m.fwd = fwd
		if err := Options(
			WithCompilerOptions(compiler.ObserveWith(fwd)),
			WithRunnerOptions(runner.ObserveWith(fwd)),
		)(m); err != nil {
			return err
		}
		return nil
	}
}

// OutputDir sets the output directory for both compiler and runner to path.
func OutputDir(path string) Option {
	return func(m *Mach) error {
		m.path = path
		return nil
	}
}

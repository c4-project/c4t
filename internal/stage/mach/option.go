// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/runner"
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

// WithUserConfig applies all of the settings specified in uc.
func WithUserConfig(uc UserConfig) Option {
	return Options(
		OutputDir(uc.OutDir),
		OverrideQuantities(uc.Quantities),
		SkipCompiler(uc.SkipCompiler),
		SkipRunner(uc.SkipRunner),
	)
}

// SkipCompiler sets whether to skip the compiler.
func SkipCompiler(skip bool) Option {
	return func(mach *Mach) error {
		mach.skipCompiler = skip
		return nil
	}
}

// SkipRunner sets whether to skip the runner.
func SkipRunner(skip bool) Option {
	return func(mach *Mach) error {
		mach.skipRunner = skip
		return nil
	}
}

// OverrideQuantities overrides the compiler and runner quantities with qs.
func OverrideQuantities(qs QuantitySet) Option {
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

// UserConfig contains the part of the machine-stage configuration that can be set by the user,
// either directly or through invoker.
type UserConfig struct {
	// OutDir is the path to the output directory.
	OutDir string
	// SkipCompiler tells the machine-runner to skip compilation.
	SkipCompiler bool
	// SkipRunner tells the machine-runner to skip running.
	SkipRunner bool
	// Quantities contains various tunable quantities for the machine-dependent stage.
	Quantities QuantitySet
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
)

// Option is the type of functional options for the machine node.
type Option func(*Runner) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(m *Runner) error {
		for _, o := range opts {
			if err := o(m); err != nil {
				return err
			}
		}
		return nil
	}
}

// LogTo sets the runner's logger to l.
func LogTo(l *log.Logger) Option {
	// TODO(@MattWindsor91): as elsewhere, logging should be replaced with observing
	return func(r *Runner) error {
		// Logger ensuring is done after all options are processed
		r.l = l
		return nil
	}
}

// ObserveWith adds each observer in obs to the runner's observer list.
func ObserveWith(obs ...builder.Observer) Option {
	return func(r *Runner) error {
		var err error
		r.observers, err = builder.AppendObservers(r.observers, obs...)
		return err
	}
}

// OverrideQuantities overrides this runner's quantities with qs.
func OverrideQuantities(qs QuantitySet) Option {
	return func(r *Runner) error {
		r.quantities.Override(qs)
		return nil
	}
}

// ObsParser is the interface of things that can parse test outcomes.
type ObsParser interface {
	// ParseObs parses the observation in reader r into o according to the backend configuration in b.
	// The backend described by b must have been used to produce the testcase outputting r.
	ParseObs(ctx context.Context, b *service.Backend, r io.Reader, o *obs.Obs) error
}

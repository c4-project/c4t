// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analyser

import (
	"errors"

	"github.com/c4-project/c4t/internal/plan/analysis"

	"github.com/c4-project/c4t/internal/stage/analyser/saver"
)

// ErrObserverNil occurs if we pass a nil Observer to ObserveWith.
var ErrObserverNil = errors.New("observer nil")

// Option is the type of options to the analyser stage constructor.
type Option func(*Analyser) error

// Options applies each option in opts in turn.
func Options(opts ...Option) Option {
	return func(a *Analyser) error {
		for _, o := range opts {
			if err := o(a); err != nil {
				return err
			}
		}
		return nil
	}
}

// Analysis sets some analysis options.
func Analysis(opts ...analysis.Option) Option {
	return func(a *Analyser) error {
		a.aopts = append(a.aopts, opts...)
		return nil
	}
}

// ObserveWith adds each observer in obs to the observer set for analysing.
func ObserveWith(obs ...Observer) Option {
	return func(a *Analyser) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
		}
		a.observers = append(a.observers, obs...)
		return nil
	}
}

// ObserveSaveWith adds each observer in obs to the observer set for saving.
func ObserveSaveWith(obs ...saver.Observer) Option {
	return func(a *Analyser) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
			a.saveObservers = append(a.saveObservers, o)
		}
		return nil
	}
}

// ErrorOnBadStatus makes the analyser fail if it observed any 'bad' statuses, if set to true.
func ErrorOnBadStatus(on bool) Option {
	return func(a *Analyser) error {
		a.errOnBadStatus = on
		return nil
	}
}

// SaveToPathset makes this analyser stage save to the given pathset.
// This can be nil, in which case saving is disabled.
func SaveToPathset(ps *saver.Pathset) Option {
	return func(a *Analyser) error {
		// ps can be nil
		a.savePaths = ps
		return nil
	}
}

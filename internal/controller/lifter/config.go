// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"errors"
	"io"
	"log"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

var (
	// ErrObserverNil occurs when we try to pass a nil observer as an option.
	ErrObserverNil = errors.New("observer nil")
)

// Option is the type of options to pass to New.
type Option func(*Lifter) error

// Options bundles up each option in os into a single option.
func Options(os ...Option) Option {
	return func(l *Lifter) error {
		for _, o := range os {
			if err := o(l); err != nil {
				return err
			}
		}
		return nil
	}
}

// LogTo makes the lifter log its progress to d.
func LogTo(d *log.Logger) Option {
	// TODO(@MattWindsor91): as in everywhere else, I'd rather deprecate logging in favour of observers.
	return func(l *Lifter) error {
		l.l = iohelp.EnsureLog(d)
		return nil
	}
}

// SendStderrTo makes the lifter send any stderr output from its driver to w.
func SendStderrTo(w io.Writer) Option {
	return func(l *Lifter) error {
		l.errw = iohelp.EnsureWriter(w)
		return nil
	}
}

// ObserveWith adds each observer in obs to the lifter's observer list.
func ObserveWith(obs ...builder.Observer) Option {
	return func(l *Lifter) error {
		for _, ob := range obs {
			if ob == nil {
				return ErrObserverNil
			}
		}
		l.obs = append(l.obs, obs...)
		return nil
	}
}

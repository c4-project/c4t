// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"errors"
	"io"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"
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

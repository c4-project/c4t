// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import "github.com/MattWindsor91/act-tester/internal/observing"

// Option is the type of options to New.
type Option func(*Saver) error

// Options applies each option in ops in turn.
func Options(ops ...Option) Option {
	return func(s *Saver) error {
		for _, o := range ops {
			if err := o(s); err != nil {
				return err
			}
		}
		return nil
	}
}

// ObserveWith appends obs to the observer list for this saver.
func ObserveWith(obs ...Observer) Option {
	return func(s *Saver) error {
		if err := observing.CheckObservers(obs); err != nil {
			return err
		}
		s.observers = append(s.observers, obs...)
		return nil
	}
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver

import "github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"

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
func ObserveWith(obs ...observer.Observer) Option {
	return func(s *Saver) error {
		for _, o := range obs {
			if o == nil {
				return ErrObserverNil
			}
			s.observers = append(s.observers, o)
		}
		return nil
	}
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// ErrObserverNil is the error raised when any of the Observe*With functions receive a nil observer.
var ErrObserverNil = errors.New("observer nil")

// Observer groups all of the disparate observer interfaces that make up an ObserverSet.
// Its main purpose is to let all of those interfaces be implemented by one slice.
type Observer interface {
	CompilerObserver
	builder.Observer
}

// ObserverSet groups the various observers used by a planner.
type ObserverSet struct {
	// Corpus contains the corpus-builder observers to be used when building out a plan.
	Corpus []builder.Observer

	// Compiler contains the compiler observers to be used when configuring the compilers to test.
	Compiler []CompilerObserver
}

// AddCorpus adds corpus observers to the observer set.
func (s *ObserverSet) AddCorpus(obs ...builder.Observer) error {
	for _, o := range obs {
		if o == nil {
			return ErrObserverNil
		}
	}
	s.Corpus = append(s.Corpus, obs...)
	return nil
}

// AddCompiler adds corpus observers to the observer set.
func (s *ObserverSet) AddCompiler(obs ...CompilerObserver) error {
	for _, o := range obs {
		if o == nil {
			return ErrObserverNil
		}
	}
	s.Compiler = append(s.Compiler, obs...)
	return nil
}

// Add adds observers to the observer set.
func (s *ObserverSet) Add(obs ...Observer) error {
	for _, o := range obs {
		if err := s.AddCorpus(o); err != nil {
			return err
		}
		if err := s.AddCompiler(o); err != nil {
			return err
		}
	}
	return nil
}

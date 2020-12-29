// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import (
	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/litmus"
	"github.com/c4-project/c4t/internal/model/recipe"
	"github.com/c4-project/c4t/internal/subject/compilation"
)

// New is a convenience constructor for subjects.
func New(origLitmus *litmus.Litmus, opt ...Option) (*Subject, error) {
	s := Subject{Source: *origLitmus}
	return &s, Options(opt...)(&s)
}

// NewOrPanic is like New, but panics if there is an error.
// Use in tests only.
func NewOrPanic(origLitmus *litmus.Litmus, opt ...Option) *Subject {
	n, err := New(origLitmus, opt...)
	if err != nil {
		panic(err)
	}
	return n
}

// Option is the type of (functional) options to the New constructor.
type Option func(*Subject) error

// Options combines the options os into a single option.
func Options(os ...Option) Option {
	return func(s *Subject) error {
		for _, o := range os {
			if err := o(s); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithCompile is an option that tries to preload a compile result for compiler ID cid onto a subject.
func WithCompile(cid id.ID, c compilation.CompileResult) Option {
	return func(s *Subject) error { return s.AddCompileResult(cid, c) }
}

// WithRecipe is an option that tries to preload a recipe for architecture ID arch onto a subject.
func WithRecipe(arch id.ID, r recipe.Recipe) Option {
	return func(s *Subject) error { return s.AddRecipe(arch, r) }
}

// WithRun is an option that tries to preload a run for compiler ID cid onto a subject.
func WithRun(cid id.ID, r compilation.RunResult) Option {
	return func(s *Subject) error { return s.AddRun(cid, r) }
}

// WithFuzz is an option that sets the incoming subject's fuzzer record to fz.
func WithFuzz(fz *Fuzz) Option {
	return func(s *Subject) error {
		s.Fuzz = fz
		return nil
	}
}

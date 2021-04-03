// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"context"

	"github.com/c4-project/c4t/internal/id"
)

// Option is the type of options to the litmus-test record constructor.
type Option func(*Litmus) error

// Options bundles the options in os into one option.
func Options(os ...Option) Option {
	return func(l *Litmus) error {
		for _, o := range os {
			if err := o(l); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithThreads is an option that sets the test's thread count to threads.
func WithThreads(threads int) Option {
	return func(l *Litmus) error {
		l.Stats.Threads = threads
		return nil
	}
}

// ReadArchFromFile is an option that causes the litmus test to populate its architecture from its given file.
func ReadArchFromFile() Option {
	return (*Litmus).PopulateArchFromFile
}

// WithArch is an option that forces the litmus test's architecture to be id.
func WithArch(id id.ID) Option {
	return func(l *Litmus) error {
		if id.IsEmpty() {
			return ErrEmptyArch
		}
		l.Arch = id
		return nil
	}
}

// PopulateStatsFrom is an option that causes the litmus test to populate its statistics set using s and ctx.
func PopulateStatsFrom(ctx context.Context, s StatDumper) Option {
	// TODO(@MattWindsor91): this capturing of the context is a bit messy.
	return func(l *Litmus) error {
		return l.PopulateStats(ctx, s)
	}
}

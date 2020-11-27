// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

// Option is the type of options to the litmus-test record constructor.
type Option func(*Litmus)

// Options bundles the options in os into one option.
func Options(os ...Option) Option {
	return func(l *Litmus) {
		for _, o := range os {
			o(l)
		}
	}
}

// WithThreads is an option that sets the test's thread count to threads.
func WithThreads(threads int) Option {
	return func(l *Litmus) {
		l.Stats.Threads = threads
	}
}

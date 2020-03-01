// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

// BuilderConfig is a configuration for a Builder.
type BuilderConfig struct {
	// Init is the initial corpus.
	// If nil, the Builder starts with a new corpus with capacity equal to NReqs.
	// Otherwise, it copies this corpus.
	Init Corpus

	// NReqs is the number of expected requests to be made to the Builder.
	// The builder will finish listening for requests when this target is reached.
	NReqs int

	// Obs is the observer to notify as the builder performs various tasks.
	Obs BuilderObserver
}

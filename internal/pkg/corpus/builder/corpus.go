// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import "github.com/MattWindsor91/act-tester/internal/pkg/corpus"

// Config is a configuration for a Builder.
type Config struct {
	// Init is the initial corpus.
	// If nil, the Builder starts with a new corpus with capacity equal to NReqs.
	// Otherwise, it copies this corpus.
	Init corpus.Corpus

	// NReqs is the number of expected requests to be made to the Builder.
	// The builder will finish listening for requests when this target is reached.
	NReqs int

	// Obs is the observer to notify as the builder performs various tasks.
	Obs Observer
}

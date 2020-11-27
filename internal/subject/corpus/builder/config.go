// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import "github.com/MattWindsor91/c4t/internal/subject/corpus"

// Config is a configuration for a Builder.
type Config struct {
	// Init is the initial corpus.
	// If nil, the Builder starts with a new corpus with capacity equal to NReqs.
	// Otherwise, it copies this corpus.
	Init corpus.Corpus

	// Manifest gives us the name of the task and the number of requests in it.
	Manifest

	// Obs is the list of observers to notify as the builder performs various tasks.
	Observers []Observer
}

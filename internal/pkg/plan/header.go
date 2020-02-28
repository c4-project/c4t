// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"math/rand"
	"time"
)

// Header is a grouping of plan metadata.
type Header struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `toml:"created"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `toml:"seed"`
}

// NewHeader produces a new header with a seed and creation time initialised from the current time.
func NewHeader() *Header {
	now := time.Now()
	return &Header{
		Creation: now,
		Seed:     now.UnixNano(),
	}
}

// Rand creates a random number generator using this Header's seed.
func (h *Header) Rand() *rand.Rand {
	return rand.New(rand.NewSource(h.Seed))
}

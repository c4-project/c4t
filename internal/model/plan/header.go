// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"math/rand"
	"time"
)

// UseDateSeed is a value for the header constructor's seed parameter that ensures its RNG will be seeded by run date.
const UseDateSeed int64 = -1

// Header is a grouping of plan metadata.
type Header struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `toml:"created"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `toml:"seed"`
}

// NewHeader produces a new header with a seed and creation time initialised from the current time.
// If seed is set to anything other than UseDateSeed, the seed will be set from the creation time.
func NewHeader(seed int64) *Header {
	now := time.Now()
	if seed == UseDateSeed {
		seed = now.UnixNano()
	}
	return &Header{Creation: now, Seed: seed}
}

// Rand creates a random number generator using this Header's seed.
func (h *Header) Rand() *rand.Rand {
	return rand.New(rand.NewSource(h.Seed))
}

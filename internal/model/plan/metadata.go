// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// UseDateSeed is a value for the header constructor's seed parameter that ensures its RNG will be seeded by run date.
const UseDateSeed int64 = -1

// ErrVersionMismatch occurs when the version of a plan loaded into part of a tester doesn't equal CurrentVer.
var ErrVersionMismatch = errors.New("bad plan version")

// Metadata is a grouping of plan metadata.
type Metadata struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `json:"created,omitempty"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `json:"seed"`

	// Version is a version identifier of the form YYYYMMDD, used to check whether the plan format has changed.
	Version Version `json:"version,omitempty"`
}

// NewMetadata produces a new header with a seed and creation time initialised from the current time.
// If seed is set to anything other than UseDateSeed, the seed will be set from the creation time.
func NewMetadata(seed int64) *Metadata {
	now := time.Now()
	if seed == UseDateSeed {
		seed = now.UnixNano()
	}
	return &Metadata{Creation: now, Seed: seed, Version: CurrentVer}
}

// CheckVersion checks to see if this header's plan version is compatible with this tool's version.
func (h Metadata) CheckVersion() error {
	if !h.Version.IsCurrent() {
		return fmt.Errorf("%w: plan version: %d; tool version: %d", ErrVersionMismatch, h.Version, CurrentVer)
	}
	return nil
}

// Rand creates a random number generator using this Metadata's seed.
func (h *Metadata) Rand() *rand.Rand {
	return rand.New(rand.NewSource(h.Seed))
}

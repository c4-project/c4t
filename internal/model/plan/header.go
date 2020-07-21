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

const (
	// UseDateSeed is a value for the header constructor's seed parameter that ensures its RNG will be seeded by run date.
	UseDateSeed int64 = -1

	// CurrentVer is the current plan version.
	// It changes when the interface between various bits of the tester (generally manifested within the plan version)
	// changes.
	CurrentVer uint32 = 2020_07_21

	// Version history since 2020_05_29:
	//
	// 2020_07_21: Added stage information.
	// 2020_05_29: Initial version for which this comment was maintained.
)

var (
	// ErrVersionMismatch occurs when the version of a plan loaded into part of a tester doesn't equal CurrentVer.
	ErrVersionMismatch = errors.New("bad plan version")
)

// Header is a grouping of plan metadata.
type Header struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `toml:"created,omitzero" json:"created,omitempty"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `toml:"seed" json:"seed"`

	// Version is a version identifier of the form YYYYMMDD, used to check whether the plan format has changed.
	Version uint32 `toml:"version,omitzero" json:"version,omitempty"`
}

// NewHeader produces a new header with a seed and creation time initialised from the current time.
// If seed is set to anything other than UseDateSeed, the seed will be set from the creation time.
func NewHeader(seed int64) *Header {
	now := time.Now()
	if seed == UseDateSeed {
		seed = now.UnixNano()
	}
	return &Header{Creation: now, Seed: seed, Version: CurrentVer}
}

// CheckVersion checks to see if this header's plan version is compatible with this tool's version.
func (h Header) CheckVersion() error {
	if h.Version != CurrentVer {
		return fmt.Errorf("%w: plan version: %d; tool version: %d", ErrVersionMismatch, h.Version, CurrentVer)
	}
	return nil
}

// Rand creates a random number generator using this Metadata's seed.
func (h *Header) Rand() *rand.Rand {
	return rand.New(rand.NewSource(h.Seed))
}

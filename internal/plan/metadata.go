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

	"github.com/MattWindsor91/act-tester/internal/plan/stage"
)

// UseDateSeed is a value for the header constructor's seed parameter that ensures its RNG will be seeded by run date.
const UseDateSeed int64 = -1

var (
	// ErrVersionMismatch occurs when the version of a plan loaded into part of a tester doesn't equal CurrentVer.
	ErrVersionMismatch = errors.New("bad plan version")

	// ErrMissingStage occurs when a plan is missing confirmation of a stage on which something depends.
	ErrMissingStage = errors.New("missing stage")
)

// Metadata is a grouping of plan metadata.
type Metadata struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `json:"created,omitempty"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `json:"seed"`

	// Version is a version identifier of the form YYYYMMDD, used to check whether the plan format has changed.
	Version Version `json:"version,omitempty"`

	// Stages contains a record of every stage that has been completed in the plan's development.
	Stages []stage.Record `json:"stages"`
}

// NewMetadata produces metadata with a seed and creation time initialised from the current time.
// If seed is set to anything other than UseDateSeed, the seed will be set from the creation time.
func NewMetadata(seed int64) *Metadata {
	now := time.Now()
	if seed == UseDateSeed {
		seed = now.UnixNano()
	}
	return &Metadata{Creation: now, Seed: seed, Version: CurrentVer}
}

// CheckVersion checks to see if this metadata's plan version is compatible with this tool's version.
// It returns ErrVersionMismatch if not.
func (m *Metadata) CheckVersion() error {
	if !m.Version.IsCurrent() {
		return fmt.Errorf("%w: plan version: %d; tool version: %d", ErrVersionMismatch, m.Version, CurrentVer)
	}
	return nil
}

// ConfirmStage adds a stage confirmation for s, which started at start and ended at end, to this metadata.
func (m *Metadata) ConfirmStage(s stage.Stage, start, end time.Time) {
	m.Stages = append(m.Stages, stage.NewRecord(s, start, end))
}

func (m *Metadata) requireOneStage(s stage.Stage) error {
	for _, r := range m.Stages {
		if r.Stage == s {
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrMissingStage, s)
}

// RequireStage checks to see if this metadata has had each given stage marked completed at least once.
// It returns ErrMissingStage if not.
func (m *Metadata) RequireStage(stages ...stage.Stage) error {
	for _, s := range stages {
		if err := m.requireOneStage(s); err != nil {
			return err
		}
	}
	return nil
}

// Rand creates a random number generator using this Metadata's seed.
func (m *Metadata) Rand() *rand.Rand {
	return rand.New(rand.NewSource(m.Seed))
}

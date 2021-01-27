// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import (
	"time"

	"github.com/c4-project/c4t/internal/model/litmus"
)

// Fuzz is the set of file paths, and other metadata, associated with a fuzzer output.
type Fuzz struct {
	// Duration is the length of time it took to fuzz this file.
	Duration time.Duration `toml:"duration,omitzero" json:"duration,omitempty"`

	// Litmus holds information about this subject's fuzzed Litmus file.
	Litmus litmus.Litmus `toml:"litmus,omitempty" json:"litmus,omitempty"`

	// Trace is the slashpath to this subject's fuzzer trace file.
	Trace string `toml:"trace,omitempty" json:"trace,omitempty"`
}

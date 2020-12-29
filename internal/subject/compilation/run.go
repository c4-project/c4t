// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"github.com/c4-project/c4t/internal/subject/obs"
)

// RunResult represents information about a single run of a subject.
type RunResult struct {
	Result

	// Obs is this run's processed observation, if any.
	Obs *obs.Obs `toml:"obs,omitempty" json:"obs,omitempty"`
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"github.com/MattWindsor91/act-tester/internal/model/obs"
)

// RunResult represents information about a single run of a subject.
type RunResult struct {
	Result

	// Obs is this run's processed observation, if any.
	Obs *obs.Obs `toml:"obs,omitempty" json:"obs,omitempty"`
}

// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
)

// Run represents information about a single run of a subject.
type Run struct {
	// Time is the time at which the run commenced.
	Time time.Time `toml:"time,omitzero"`

	// Duration is the rough duration of the run.
	Duration time.Duration `toml:"duration,omitzero"`

	// Status is the status of the run.
	Status Status `toml:"status"`

	// Obs is this run's processed observation, if any.
	Obs *obs.Obs `toml:"obs,omitempty"`
}

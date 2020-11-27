// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stage

import (
	"fmt"
	"time"
)

// Record is the type of stage completion records.
// These records log the fact that a stage was completed, the time of completion, and the duration.
type Record struct {
	// Stage is the identifier of the stage completed.
	Stage Stage `json:"stage"`

	// CompletedOn is the time at which the stage was completed.
	CompletedOn time.Time `json:"completed_on,omitempty"`

	// Duration is the time it took to complete the stage.
	Duration time.Duration `json:"duration,omitempty"`
}

// String converts a record to a human-readable string.
func (r Record) String() string {
	return fmt.Sprintf("%s completed on %s (took %s)", r.Stage, r.CompletedOn.Format(time.RFC3339), r.Duration)
}

// NewRecord creates a completion record for a stage s that started on start and lasted for dur.
// The completed-on and duration fields are set relative to the current time.
func NewRecord(s Stage, start time.Time, dur time.Duration) Record {
	return Record{
		Stage:       s,
		CompletedOn: start.Add(dur),
		Duration:    dur,
	}
}

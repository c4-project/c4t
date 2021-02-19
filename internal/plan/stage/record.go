// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stage

import (
	"fmt"

	"github.com/c4-project/c4t/internal/timing"
)

// Record is the type of stage completion records.
// These records log the fact that a stage was completed, the time of completion, and the duration.
type Record struct {
	// Stage is the identifier of the stage completed.
	Stage Stage `json:"stage"`

	// Timespan notes the start and end time of the stage.
	Timespan timing.Span `json:"timespan,omitempty"`
}

// String converts a record to a human-readable string.
func (r Record) String() string {
	return fmt.Sprintf("%s: %s", r.Stage, r.Timespan)
}

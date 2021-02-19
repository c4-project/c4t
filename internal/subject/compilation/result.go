// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"github.com/c4-project/c4t/internal/subject/status"
	"github.com/c4-project/c4t/internal/timing"
)

// Result is the base structure for things that represent the result of an external process.
type Result struct {
	// Timespan embeds a timespan for this Result.
	Timespan timing.Span `json:"time_span,omitempty"`

	// Status is the status of the process.
	Status status.Status `json:"status"`
}

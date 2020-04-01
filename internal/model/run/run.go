// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package run contains the Run type and various related functions.
package run

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Run represents information about a test run: its machine ID, iteration number, and start time.
// This information helps disambiguate between test runs in logs.
type Run struct {
	// MachineID is the ID of the machine whose collation is being reported.
	MachineID id.ID

	// Iter is the iteration number of the run to which this collation is associated.
	Iter uint64

	// Start is the start time of the run to which this collation is associated.
	Start time.Time
}

// String returns a string containing the components of this run in a human-readable manner.
func (r Run) String() string {
	return fmt.Sprintf(
		"[%s #%d (%s)]",
		r.MachineID.String(),
		r.Iter,
		r.Start.Format(time.Stamp),
	)
}

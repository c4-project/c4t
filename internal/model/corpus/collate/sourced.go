// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Sourced contains a corpus collation and various bits of information about whence it came.
type Sourced struct {
	// MachineID is the ID of the machine whose collation is being reported.
	MachineID id.ID

	// Iter is the iteration number of the run to which this collation is associated.
	Iter uint64

	// Start is the start time of the run to which this collation is associated.
	Start time.Time

	// Collation is the collation proper.
	Collation *Collation
}

// String formats a log header for this sourced collation.
func (s *Sourced) String() string {
	cstr := "(nil)"
	if s.Collation != nil {
		cstr = s.Collation.String()
	}
	return fmt.Sprintf(
		"[%s #%d (%s)] %s",
		s.MachineID.String(),
		s.Iter,
		s.Start.Format(time.Stamp),
		cstr,
	)
}

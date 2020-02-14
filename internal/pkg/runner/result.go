package runner

import "time"

// Result is the type of results from a single test run.
type Result struct {
	// Start is the time at which the run commenced.
	Start time.Time

	// Machines is the set of results that were reported in this test run.
	// Each key is a stringified form of a machine ID.
	Machines map[string]MachineResult
}

type MachineResult struct{}

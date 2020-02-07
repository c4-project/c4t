package model

// Plan represents a test plan.
// A plan covers an entire campaign of testing.
type Plan struct {
	// Seed is a pseudo-randomly generated unsigned integer that should be used to drive randomiser input.
	Seed uint `json:"seed"`

	// Machines contains the per-machine plans for this overall test plan.
	Machines []MachinePlan `json:"machines"`
}

// MachinePlan represents a test plan for a single machine.
type MachinePlan struct {
	// A MachinePlan subsumes a machine entry.
	Machine

	// Backend represents the backend targeted by this plan.
	Backend Backend `json:"backend"`

	// Compilers represents the compilers to be targeted by this plan.
	Compilers []Compiler `json:"compilers"`
}

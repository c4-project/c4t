package model

import (
	"math/rand"
	"time"
)

// Plan represents a test plan.
// A plan covers an entire campaign of testing.
type Plan struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `toml:"created"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `toml:"seed"`

	// Machines contains the per-machine plans for this overall test plan.
	Machines []MachinePlan `toml:"machines"`

	// Corpus contains the filenames of each test corpus entry chosen for this plan.
	Corpus []string `toml:"corpus"`
}

// Init initialises the creation-sensitive parts of plan p.
// It randomises the seed using the top-level random number generator;
// and also updates the creation time.
func (p *Plan) Init() {
	p.Creation = time.Now()
	p.Seed = rand.Int63()
}

// MachinePlan represents a test plan for a single machine.
type MachinePlan struct {
	// A MachinePlan subsumes a machine entry.
	Machine

	// Backend represents the backend targeted by this plan.
	Backend Backend `toml:"backend"`

	// Compilers represents the compilers to be targeted by this plan.
	Compilers []Compiler `toml:"compilers"`
}

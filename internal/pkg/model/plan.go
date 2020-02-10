package model

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// ErrPlanLoaded occurs when a PlanLoader tries to load a plan multiple times.
var ErrPlanLoaded = errors.New("plan already loaded")

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
	Corpus []Subject `toml:"corpus"`
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

// PlanLoader holds a Plan pointer and a file, and can load in the former from the latter.
type PlanLoader struct {
	// PlanFile contains, if non-empty, the file path of the plan.
	PlanFile string

	// Plan stores the plan after it has been loaded from PlanFile.
	Plan *Plan
}

// LoadPlan loads the plan pointed to by d.PlanFile into d.Plan, replacing any existing plan.
// It returns an error if there is already a plan loaded.
func (p *PlanLoader) LoadPlan() error {
	if p.Plan != nil {
		return ErrPlanLoaded
	}
	if err := p.actuallyLoadPlan(); err != nil {
		return err
	}
	if p == nil {
		return errors.New("plan nil after loading")
	}
	return nil
}

func (p *PlanLoader) actuallyLoadPlan() error {
	if p.PlanFile == "" || p.PlanFile == "-" {
		_, err := toml.DecodeReader(os.Stdin, &p.Plan)
		return err
	}
	_, err := toml.DecodeFile(p.PlanFile, &p.Plan)
	return err
}

package model

import (
	"context"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	// StdinFile is the special file path that the plan loader treats as a request to load from stdin instead.
	StdinFile = "-"
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
	Corpus []Subject `toml:"corpus"`
}

// Init initialises the creation-sensitive parts of plan p.
// It randomises the seed using the top-level random number generator;
// and also updates the creation time.
func (p *Plan) Init() {
	p.Creation = time.Now()
	p.Seed = rand.Int63()
}

// Dump dumps plan p to stdout.
func (p *Plan) Dump() error {
	// TODO(@MattWindsor91): output to other files
	enc := toml.NewEncoder(os.Stdout)
	enc.Indent = "  "
	return enc.Encode(p)
}

// ParMachines runs f for every machine in the plan, threading through a context that will terminate each machine if
// an error occurs on some other machine.
func (p *Plan) ParMachines(ctx context.Context, f func(context.Context, MachinePlan) error) error {
	eg, ectx := errgroup.WithContext(ctx)
	for _, m := range p.Machines {
		eg.Go(func() error { return f(ectx, m) })
	}
	return eg.Wait()
}

// Load loads a plan pointed to by f into p, replacing any existing plan.
// If f is empty or StdinFile, Load loads from standard input instead.
func (p *Plan) Load(f string) error {
	if f == "" || f == StdinFile {
		_, err := toml.DecodeReader(os.Stdin, &p)
		return err
	}
	_, err := toml.DecodeFile(f, &p)
	return err
}
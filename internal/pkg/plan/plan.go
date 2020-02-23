package plan

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"golang.org/x/sync/errgroup"

	"github.com/BurntSushi/toml"
)

var (
	// ErrNil is an error that can be returned if a tester stage gets a nil plan.
	ErrNil = errors.New("plan nil")

	// ErrNoMachine is an error that can be returned if an attempt to get a machine by its CompilerID fails.
	ErrNoMachine = errors.New("can't get machine")
)

// plan represents a test plan.
// A plan covers an entire campaign of testing.
type Plan struct {
	// Creation marks the time at which the plan was created.
	Creation time.Time `toml:"created"`

	// Seed is a pseudo-randomly generated integer that should be used to drive randomiser input.
	Seed int64 `toml:"seed"`

	// Machines contains the per-machine plans for this overall test plan.
	// Each machine is mapped under a stringified form of its CompilerID.
	Machines map[string]MachinePlan `toml:"machines"`

	// Corpus contains each test corpus entry chosen for this plan.
	Corpus subject.Corpus `toml:"corpus"`
}

// New creates a new plan using the given machines and corpus.
// It randomises the seed using the top-level random number generator;
// and also updates the creation time.
func New(ms map[string]MachinePlan, c subject.Corpus) *Plan {
	p := Plan{Machines: ms, Corpus: c}
	p.Creation = time.Now()
	p.Seed = rand.Int63()
	return &p
}

// Dump dumps plan p to stdout.
func (p *Plan) Dump() error {
	// TODO(@MattWindsor91): output to other files
	enc := toml.NewEncoder(os.Stdout)
	enc.Indent = "  "
	return enc.Encode(p)
}

// ParCorpus runs f for every subject in the plan's corpus.
// It threads through a context that will terminate each machine if an error occurs on some other machine.
// It also takes zero or more 'auxiliary' funcs to launch within the same context.
func (p *Plan) ParCorpus(ctx context.Context, f func(context.Context, subject.Named) error, aux ...func(context.Context) error) error {
	return p.Corpus.Par(ctx, f, aux...)
}

// ParMachines runs f for every machine in the plan.
// It threads through a context that will terminate each machine if an error occurs on some other machine.
// It also takes zero or more 'auxiliary' funcs to launch within the same context.
func (p *Plan) ParMachines(ctx context.Context, f func(context.Context, model.ID, MachinePlan) error, aux ...func(context.Context) error) error {
	eg, ectx := errgroup.WithContext(ctx)
	for i, m := range p.Machines {
		mid := model.IDFromString(i)
		mc := m
		eg.Go(func() error { return f(ectx, mid, mc) })
	}
	for _, a := range aux {
		eg.Go(func() error { return a(ectx) })
	}
	return eg.Wait()
}

// Machine gets the plan of the machine with CompilerID id, if it exists.
// If id is empty and the plan contains only one machine, Machine gets that instead.
// Machine also returns the actual ID of the machine.
func (p *Plan) Machine(id model.ID) (model.ID, MachinePlan, error) {
	if id.IsEmpty() {
		var err error
		if id, err = p.singleMachineId(); err != nil {
			return id, MachinePlan{}, err
		}
	}

	ids := id.String()
	mp, ok := p.Machines[ids]
	if !ok {
		return id, MachinePlan{}, fmt.Errorf("%w: no plan for machine %s", ErrNoMachine, ids)
	}
	return id, mp, nil
}

func (p *Plan) singleMachineId() (model.ID, error) {
	rv := reflect.ValueOf(p.Machines)
	keys := rv.MapKeys()
	if len(keys) != 1 {
		return model.ID{}, fmt.Errorf("%w: machine plan doesn't contain exactly one machine", ErrNoMachine)
	}
	return model.IDFromString(keys[0].String()), nil
}

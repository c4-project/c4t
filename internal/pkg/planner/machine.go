package planner

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ErrNoMachineID arises when a compiler arrives at the machine querying logic with no attached machine ID.
var ErrNoMachineID = errors.New("queried compiler has no machine id")

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers given the filter f.
	ListCompilers(f model.CompilerFilter) ([]*model.Compiler, error)
}

func (p *Planner) planMachines() ([]plan.MachinePlan, error) {
	cmap, err := p.queryCompilers()
	if err != nil {
		return nil, err
	}
	return p.planMachinesFromMap(cmap)
}

// queryCompilers asks ACT for the list of compilers,
// then massages them into a map from stringified machine ID to compiler list.
func (p *Planner) queryCompilers() (map[string][]model.Compiler, error) {
	cs, err := p.Source.ListCompilers(p.Filter)
	if err != nil {
		return nil, err
	}

	cmap := make(map[string][]model.Compiler)
	for _, c := range cs {
		key, kerr := scrubMachineID(c)
		if kerr != nil {
			return nil, kerr
		}
		cmap[key] = append(cmap[key], *c)
	}

	return cmap, nil
}

// scrubMachineID tries to take the machine ID off c and return it as a string.
// It fails if there was no machine ID in the first place.
func scrubMachineID(c *model.Compiler) (key string, err error) {
	if c.MachineID == nil {
		return "", ErrNoMachineID
	}
	key, c.MachineID = c.MachineID.String(), nil
	return key, nil
}

// planMachinesFromMap assembles a list of machine plans by taking a compiler map cmap and performing all other machine
// information scraping necessary.
func (p *Planner) planMachinesFromMap(cmap map[string][]model.Compiler) ([]plan.MachinePlan, error) {
	var err error

	plans := make([]plan.MachinePlan, len(cmap))
	i := 0
	for mstr, cs := range cmap {
		mid := model.IDFromString(mstr)
		if plans[i], err = p.planMachine(mid, cs); err != nil {
			return nil, err
		}
		i++
	}

	return plans, nil
}

// planMachine builds a machine plan given machine ID mid and compiler set compilers.
// It performs various further config lookups on the machine, which can cause errors.
func (p *Planner) planMachine(mid model.ID, compilers []model.Compiler) (plan.MachinePlan, error) {
	style := model.IDFromString("litmus")
	backend, err := p.Source.FindBackend(style, mid)
	if err != nil {
		return plan.MachinePlan{}, err
	}

	// TODO(@MattWindsor91): probe cores
	return plan.MachinePlan{
		Machine:   model.Machine{ID: mid},
		Backend:   *backend,
		Compilers: compilers,
	}, nil
}

package act_tester_plan

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

var ErrNoMachineId = errors.New("queried compiler has no machine id")

func (p *Planner) planMachines() ([]model.MachinePlan, error) {
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
		key, kerr := scrubMachineId(c)
		if kerr != nil {
			return nil, kerr
		}
		cmap[key] = append(cmap[key], *c)
	}

	return cmap, nil
}

// scrubMachineId tries to take the machine ID off c and return it as a string.
// It fails if there was no machine ID in the first place.
func scrubMachineId(c *model.Compiler) (key string, err error) {
	if c.MachineId == nil {
		return "", ErrNoMachineId
	}
	key, c.MachineId = c.MachineId.String(), nil
	return key, nil
}

// planMachinesFromMap assembles a list of machine plans by taking a compiler map cmap and performing all other machine
// information scraping necessary.
func (p *Planner) planMachinesFromMap(cmap map[string][]model.Compiler) ([]model.MachinePlan, error) {
	var err error

	plans := make([]model.MachinePlan, len(cmap))
	i := 0
	for mstr, cs := range cmap {
		mid := model.IdFromString(mstr)
		if plans[i], err = p.planMachine(mid, cs); err != nil {
			return nil, err
		}
		i++
	}

	return plans, nil
}

// planMachine builds a machine plan given machine ID mid and compiler set compilers.
// It performs various further config lookups on the machine, which can cause errors.
func (p *Planner) planMachine(mid model.Id, compilers []model.Compiler) (model.MachinePlan, error) {
	style := model.IdFromString("litmus")
	backend, err := p.Source.FindBackend(style, mid)
	if err != nil {
		return model.MachinePlan{}, err
	}

	// TODO(@MattWindsor91): probe cores
	return model.MachinePlan{
		Machine:   model.Machine{Id: mid},
		Backend:   *backend,
		Compilers: compilers,
	}, nil
}

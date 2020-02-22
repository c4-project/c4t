package planner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers given the filter f.
	ListCompilers(ctx context.Context, f model.CompilerFilter) (map[string]map[string]model.Compiler, error)
}

func (p *Planner) planMachines(ctx context.Context) (map[string]plan.MachinePlan, error) {
	cmap, err := p.Source.ListCompilers(ctx, p.Filter)
	if err != nil {
		return nil, err
	}
	return p.planMachinesFromMap(ctx, cmap)
}

// planMachinesFromMap assembles a list of machine plans by taking a compiler map cmap and performing all other machine
// information scraping necessary.
func (p *Planner) planMachinesFromMap(ctx context.Context, cmap map[string]map[string]model.Compiler) (map[string]plan.MachinePlan, error) {
	var err error

	plans := make(map[string]plan.MachinePlan, len(cmap))
	for mstr, cs := range cmap {
		mid := model.IDFromString(mstr)
		if plans[mstr], err = p.planMachine(ctx, mid, cs); err != nil {
			return nil, err
		}
	}

	return plans, nil
}

// planMachine builds a machine plan given machine CompilerID mid and compiler set compilers.
// It performs various further config lookups on the machine, which can cause errors.
func (p *Planner) planMachine(ctx context.Context, mid model.ID, compilers map[string]model.Compiler) (plan.MachinePlan, error) {
	style := model.IDFromString("litmus")
	backend, err := p.Source.FindBackend(ctx, style, mid)
	if err != nil {
		return plan.MachinePlan{}, err
	}

	// TODO(@MattWindsor91): probe cores
	return plan.MachinePlan{
		Machine:   model.Machine{},
		Backend:   *backend,
		Compilers: compilers,
	}, nil
}

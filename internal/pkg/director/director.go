// Package director contains the top-level ACT test director, which manages a full testing campaign.
package director

import (
	"context"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Director contains the main state and configuration for the test director.
type Director struct {
	// PlanFile contains, if non-empty, the file path of the plan.
	PlanFile string

	// Plan contains the pre-computed plan for this test campaign.
	Plan *model.Plan

	OutDir string
}

// LoadPlan loads the plan pointed to by d.PlanFile into d.Plan, replacing any existing plan.
func (d *Director) LoadPlan() error {
	if d.PlanFile == "" || d.PlanFile == "-" {
		_, err := toml.DecodeReader(os.Stdin, &d.Plan)
		return err
	}
	_, err := toml.DecodeFile(d.PlanFile, &d.Plan)
	return err
}

// openFileOpt opens the file at path, or stdin if path is empty or the special path '-'.
func openFileOpt(path string) (io.ReadCloser, error) {
	if !(path == "" || path == "-") {
		return os.Stdin, nil
	}
	return os.Open(path)
}

// Direct runs the director d over whichever test plan is loaded into it.
// If no plan is loaded, it tries to load one from the given file.
func (d *Director) Direct(ctx context.Context) error {
	var err error

	if err = d.loadPlanIfNeeded(); err != nil {
		return err
	}
	logrus.Infof("Using plan from %v.", d.Plan.Creation)

	var ps *Pathset
	if ps, err = d.prepare(); err != nil {
		return err
	}

	var fc []string
	if fc, err = d.fuzzCorpus(ps); err != nil {
		return err
	}

	return d.directMachines(ctx, ps, fc)
}

func (d *Director) loadPlanIfNeeded() error {
	if d.Plan == nil {
		return d.LoadPlan()
	}

	return nil
}

func (d *Director) prepare() (*Pathset, error) {
	// TODO(@MattWindsor91)

	return nil, nil
}

type Pathset struct {
	FuzzDir string
}

func (d *Director) fuzzCorpus(_ *Pathset) (corpusFiles []string, err error) {
	// TODO(@MattWindsor91)

	return nil, nil
}

func (d *Director) directMachines(ctx context.Context, ps *Pathset, fc []string) error {
	eg, ectx := errgroup.WithContext(ctx)
	for _, m := range d.Plan.Machines {
		eg.Go(func() error { return d.directMachine(ectx, ps, fc, m) })
	}

	return nil
}

func (d *Director) directMachine(_ context.Context, _ *Pathset, _ []string, _ model.MachinePlan) error {
	// TODO(@MattWindsor91)
	return nil
}

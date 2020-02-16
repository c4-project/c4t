package main

import (
	"context"
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/compiler"
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const (
	defaultOutDir    = "run_results"
	usageFlagMachine = "specifies a machine in a multi-machine plan to run"
)

func main() {
	var (
		act   interop.ActRunner
		dir   string
		pfile string
		pmach string
	)
	rn := compiler.Compiler{Runner: &act}

	ux.ActRunnerFlags(&act)
	ux.OutDirFlag(&dir, defaultOutDir)
	ux.PlanFileFlag(&pfile)
	flag.StringVar(&pmach, ux.FlagMachine, "", usageFlagMachine)
	flag.Parse()

	err := run(dir, pfile, pmach, &rn)
	ux.LogTopError(err)
}

func run(dir, pfile, pmach string, rn *compiler.Compiler) error {
	p, perr := ux.LoadPlan(pfile)
	if perr != nil {
		return perr
	}

	rn.Paths = compiler.NewPathset(dir)

	// TODO(@MattWindsor91): output results
	_, rerr := rn.RunOnPlan(context.Background(), p, model.IDFromString(pmach))
	if rerr != nil {
		return rerr
	}

	return nil
}

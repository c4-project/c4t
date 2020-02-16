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

	ux.ActRunnerFlags(&act)
	ux.OutDirFlag(&dir, defaultOutDir)
	ux.PlanFileFlag(&pfile)
	flag.StringVar(&pmach, ux.FlagMachine, "", usageFlagMachine)
	flag.Parse()

	cfg := compiler.Config{
		Driver:    &act,
		MachineID: model.IDFromString(pmach),
		Paths:     compiler.NewPathset(dir),
	}
	err := ux.RunOnPlanFile(context.Background(), &cfg, pfile)
	ux.LogTopError(err)
}

package main

import (
	"context"
	"flag"
	"os"

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
	if err := run(os.Args); err != nil {
		ux.LogTopError(err)
	}
}

func run(args []string) error {
	var (
		act   interop.ActRunner
		dir   string
		pfile string
		pmach string
	)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &act)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pfile)
	fs.StringVar(&pmach, ux.FlagMachine, "", usageFlagMachine)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg := compiler.Config{
		Driver:    &act,
		MachineID: model.IDFromString(pmach),
		Paths:     compiler.NewPathset(dir),
	}
	return ux.RunOnPlanFile(context.Background(), &cfg, pfile)
}

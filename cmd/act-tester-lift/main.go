package main

import (
	"context"
	"flag"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

func main() {
	err := run(os.Args)
	ux.LogTopError(err)
}

func run(args []string) error {
	var (
		act interop.ActRunner
		pf  string
	)
	lift := lifter.Lifter{Maker: &act}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &act)
	ux.OutDirFlag(fs, &lift.OutDir, defaultOutDir)
	ux.PlanFileFlag(fs, &pf)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return ux.RunOnPlanFile(context.Background(), &lift, pf)
}

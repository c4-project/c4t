package main

import (
	"context"
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

// defaultOutDir is the default directory used for the results of the lifter.
const defaultOutDir = "lift_results"

func main() {
	var (
		act interop.ActRunner
		pf  string
	)
	lift := lifter.Lifter{Maker: &act}

	ux.ActRunnerFlags(&act)
	ux.OutDirFlag(&lift.OutDir, defaultOutDir)
	ux.PlanFileFlag(&pf)
	flag.Parse()

	err := lift.LiftPlanFile(context.Background(), pf)
	ux.LogTopError(err)
}

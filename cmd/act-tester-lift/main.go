package main

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/lifter"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

func main() {
	var act interop.ActRunner
	lift := lifter.Lifter{Maker: &act}

	ux.ActRunnerFlags(&act)
	ux.PlanLoaderFlags(&lift.PlanLoader)
	ux.OutDirFlag(&lift.OutDir, "lift_results")
	flag.Parse()

	err := lift.Lift()
	ux.LogTopError(err)
}

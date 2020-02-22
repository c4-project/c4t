package main

import (
	"context"
	"flag"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	"github.com/MattWindsor91/act-tester/internal/pkg/planner"
)

const (
	usageCompPred = "predicate `sexp` used to filter compilers for this test plan"
	usageMachPred = "predicate `sexp` used to filter machines for this test plan"
)

func main() {
	err := run(os.Args)
	ux.LogTopError(err)
}

func run(args []string) error {
	var act interop.ActRunner
	plan := planner.Planner{Source: &act}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&plan.Filter.CompPred, "c", "", usageCompPred)
	fs.StringVar(&plan.Filter.MachPred, ux.FlagMachine, "", usageMachPred)
	ux.ActRunnerFlags(fs, &act)
	ux.CorpusSizeFlag(fs, &plan.CorpusSize)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	plan.InFiles = fs.Args()

	return plan.Plan(context.Background())
}

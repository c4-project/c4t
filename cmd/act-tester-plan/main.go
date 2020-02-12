package main

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	"github.com/MattWindsor91/act-tester/internal/pkg/planner"
)

const (
	usageCompPred = "predicate `sexp` used to filter compilers for this test plan"
	usageMachPred = "predicate `sexp` used to filter machines for this test plan"
)

func main() {
	var act interop.ActRunner
	plan := planner.Planner{Source: &act}

	flag.StringVar(&plan.Filter.CompPred, "c", "", usageCompPred)
	flag.StringVar(&plan.Filter.MachPred, "m", "", usageMachPred)
	ux.ActRunnerFlags(&act)
	ux.CorpusSizeFlag(&plan.CorpusSize)
	flag.Parse()
	plan.Corpus = flag.Args()

	err := plan.Plan()
	ux.LogTopError(err)
}

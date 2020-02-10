package main

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"

	"github.com/MattWindsor91/act-tester/internal/pkg/planner"
)

// The configuration parsed from the command-line arguments.
var cfg = planner.Planner{
	Source: &interop.ActRunner{},
}

const (
	corpusSizeUsage = "`number` of corpus files to select for this test plan;\n" +
		"if non-positive, the planner will use all viable provided corpus files"
	compPredUsage = "predicate `sexp` used to filter compilers for this test plan"
	machPredUsage = "predicate `sexp` used to filter machines for this test plan"
)

func init() {
	flag.StringVar(&cfg.Filter.CompPred, "c", "", compPredUsage)

	flag.StringVar(&cfg.Filter.MachPred, "m", "", machPredUsage)

	flag.IntVar(&cfg.CorpusSize, "n", 0, corpusSizeUsage)
}

func main() {
	flag.Parse()
	cfg.Corpus = flag.Args()

	err := cfg.Plan()
	ux.LogTopError(err)
}

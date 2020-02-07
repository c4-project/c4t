package main

import (
	"flag"
	"fmt"
	"github.com/MattWindsor91/act-tester/internal/app/act-tester-plan"
	"os"
)

// The configuration parsed from the command-line arguments.
var cfg act_tester_plan.Planner

const (
	compPredUsage = "The predicate `sexp` used to filter compilers for this test plan."
	machPredUsage = "The predicate `sexp` used to filter machines for this test plan."
)

func init() {
	flag.StringVar(&cfg.Filter.CompPred, "c", "", compPredUsage + " (shorthand)")
	flag.StringVar(&cfg.Filter.MachPred, "m", "", machPredUsage + " (shorthand)")
	flag.StringVar(&cfg.Filter.CompPred, "compiler-predicate", "", compPredUsage)
	flag.StringVar(&cfg.Filter.MachPred, "machine-predicate", "", machPredUsage)
}

func main() {
	flag.Parse()
	cfg.Corpus = flag.Args()

	if err := cfg.Plan(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

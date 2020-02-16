package main

import (
	"context"
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const (
	// defaultOutDir is the default directory used for the results of the lifter.
	defaultOutDir = "fuzz_results"

	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
)

func main() {
	var (
		act interop.ActRunner
		dir string
		pf  string
	)
	cfg := fuzzer.Config{Driver: &act}

	ux.ActRunnerFlags(&act)
	ux.CorpusSizeFlag(&cfg.CorpusSize)
	ux.OutDirFlag(&dir, defaultOutDir)
	ux.PlanFileFlag(&pf)
	flag.IntVar(&cfg.SubjectCycles, "k", fuzzer.DefaultSubjectCycles, usageSubjectCycles)
	flag.Parse()

	cfg.Paths = fuzzer.NewPathset(dir)
	err := ux.RunOnPlanFile(context.Background(), &cfg, pf)
	ux.LogTopError(err)
}

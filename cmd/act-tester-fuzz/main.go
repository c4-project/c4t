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

	usageFuzzWorkers   = "cap number of `workers` to spawn to perform fuzzing"
	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
)

func main() {
	var (
		act interop.ActRunner
		pf  string
	)
	fuzz := fuzzer.Fuzzer{Driver: &act}

	ux.ActRunnerFlags(&act)
	ux.CorpusSizeFlag(&fuzz.CorpusSize)
	ux.OutDirFlag(&fuzz.OutDir, defaultOutDir)
	ux.PlanFileFlag(&pf)
	flag.IntVar(&fuzz.FuzzWorkers, "j", fuzzer.NoChunkLimit, usageFuzzWorkers)
	flag.IntVar(&fuzz.SubjectCycles, "k", fuzzer.DefaultSubjectCycles, usageSubjectCycles)
	flag.Parse()

	err := ux.RunOnPlanFile(context.Background(), &fuzz, pf)
	ux.LogTopError(err)
}

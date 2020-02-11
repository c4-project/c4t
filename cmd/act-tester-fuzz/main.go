package main

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"

func main() {
	var act interop.ActRunner
	fuzz := fuzzer.Fuzzer{Driver: &act}

	ux.ActRunnerFlags(&act)
	ux.PlanLoaderFlags(&fuzz.PlanLoader)
	ux.OutDirFlag(&fuzz.OutDir, "fuzz_results")
	ux.CorpusSizeFlag(&fuzz.CorpusSize)
	flag.IntVar(&fuzz.SubjectCycles, "k", fuzzer.DefaultSubjectCycles, usageSubjectCycles)
	flag.Parse()

	err := fuzz.Fuzz()
	ux.LogTopError(err)
}

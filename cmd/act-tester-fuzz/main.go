package main

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

var (
	act  interop.ActRunner
	fuzz = fuzzer.Fuzzer{Driver: &act}
)

const (
	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
	usageOutDir        = "`directory` to which fuzzer outputs will be written"
)

func init() {
	ux.ActRunnerFlags(&act)
	ux.PlanLoaderFlags(&fuzz.PlanLoader)

	flag.StringVar(&fuzz.OutDir, "d", "fuzz_results", usageOutDir)
	flag.IntVar(&fuzz.SubjectCycles, "k", fuzzer.DefaultSubjectCycles, usageSubjectCycles)
}

func main() {
	flag.Parse()
	err := fuzz.Fuzz()
	ux.LogTopError(err)
}

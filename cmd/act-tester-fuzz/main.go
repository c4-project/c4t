package main

import (
	"context"
	"flag"
	"os"

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
	err := run(os.Args)
	ux.LogTopError(err)
}

func run(args []string) error {
	var (
		act interop.ActRunner
		dir string
		pf  string
	)
	cfg := fuzzer.Config{Driver: &act}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.ActRunnerFlags(fs, &act)
	ux.CorpusSizeFlag(fs, &cfg.CorpusSize)
	ux.OutDirFlag(fs, &dir, defaultOutDir)
	ux.PlanFileFlag(fs, &pf)
	fs.IntVar(&cfg.SubjectCycles, "k", fuzzer.DefaultSubjectCycles, usageSubjectCycles)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg.Paths = fuzzer.NewPathset(dir)
	return ux.RunOnPlanFile(context.Background(), &cfg, pf)
}

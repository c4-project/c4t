// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/act"
)

const (
	// FlagRunTimeout is a short flag for run timeout.
	FlagRunTimeout  = "T"
	flagActConfFile = "A"
	flagConfFile    = "C"
	// FlagUseJSON is a short flag for enabling JSON output where available.
	FlagUseJSON = "J"
	// FlagOutDir is a short flag for specifying an output directory.
	FlagOutDir = "d"
	// FlagInputFile is a standard flag for arguments that suggest an alternative to stdin for commands that read files.
	FlagInputFile = "i"
	// FlagWorkerCount is a standard flag for arguments that set a worker count.
	FlagWorkerCount   = "j"
	flagSubjectCycles = "k"
	// FlagMachine is a standard flag for machine selection arguments.
	FlagMachine = "m"
	// FlagNum is a standard flag for 'number of' arguments.
	FlagNum = "n"
	// FlagCompilerTimeout is a short flag for compiler timeout.
	FlagCompilerTimeout = "t"
	flagActDuneExec     = "x"

	// FlagCPUProfile is a standard flag for specifying a CPU profile output.
	FlagCPUProfile = "cpuprofile"

	FlagSkipCompiler = "skip-compiler"
	FlagSkipRunner   = "skip-runner"

	// FlagCompilerTimeoutLong is a long flag for compiler timeout.
	FlagCompilerTimeoutLong = "compiler-timeout"
	// FlagRunTimeoutLong is a long flag for run timeout.
	FlagRunTimeoutLong = "run-timeout"
	// FlagUseJSONLong is a long flag for JSON emission.
	FlagUseJSONLong = "emit-json"
	// FlagWorkerCountLong is a long flag for arguments that set a worker count.
	FlagWorkerCountLong = "num-workers"

	usageActConfFile = "read ACT config from this `file`"
	usageCorpusSize  = "`number` of corpus files to select for this test plan;\n" +
		"if non-positive, the planner will use all viable provided corpus files"
	usageDuneExec      = "if true, use 'dune exec' to run OCaml ACT binaries"
	usageOutDir        = "`directory` to which outputs will be written"
	usagePlanFile      = "read from this plan `file` instead of stdin"
	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
	usageCPUProfile    = "`file` into which we should dump pprof information"
)

// CorpusSizeFlag sets up a 'target corpus size' flag on fs.
func CorpusSizeFlag(fs *flag.FlagSet, out *int) {
	fs.IntVar(out, FlagNum, 0, usageCorpusSize)
}

// SubjectCycleFlag sets up a 'number of cycles' flag on fs.
func SubjectCycleFlag(fs *flag.FlagSet, out *int) {
	fs.IntVar(out, flagSubjectCycles, fuzzer.DefaultSubjectCycles, usageSubjectCycles)
}

// OutDirFlag sets up an 'output directory' flag on fs.
func OutDirFlag(fs *flag.FlagSet, out *string, defaultdir string) {
	fs.StringVar(out, FlagOutDir, defaultdir, usageOutDir)
}

// ActRunnerFlags sets up a standard set of arguments on fs feeding into the ActRunner a.
func ActRunnerFlags(fs *flag.FlagSet, a *act.Runner) {
	fs.StringVar(&a.ConfFile, flagActConfFile, "", usageActConfFile)
	fs.BoolVar(&a.DuneExec, flagActDuneExec, false, usageDuneExec)
}

// PlanFileFlag sets up a standard argument on fs for loading a plan file into f.
func PlanFileFlag(fs *flag.FlagSet, f *string) {
	fs.StringVar(f, FlagInputFile, "", usagePlanFile)
}

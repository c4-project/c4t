// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package view

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/act"
)

const (
	// FlagInputFile is a standard flag for arguments that suggest an alternative to stdin for commands that read files.
	FlagInputFile = "i"

	// FlagMachine is a standard flag for machine selection arguments.
	FlagMachine = "m"

	// FlagNum is a standard flag for 'number of' arguments.
	FlagNum = "n"

	flagActConfFile   = "A"
	flagConfFile      = "C"
	flagOutDir        = "d"
	flagSubjectCycles = "k"
	flagActDuneExec   = "x"

	usageConfFile    = "The `file` from which to load the tester configuration."
	usageActConfFile = "read ACT config from this `file`"
	usageCorpusSize  = "`number` of corpus files to select for this test plan;\n" +
		"if non-positive, the planner will use all viable provided corpus files"
	usageDuneExec      = "if true, use 'dune exec' to run OCaml ACT binaries"
	usageOutDir        = "`directory` to which outputs will be written"
	usagePlanFile      = "read from this plan `file` instead of stdin"
	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
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
	fs.StringVar(out, flagOutDir, defaultdir, usageOutDir)
}

// ActRunnerFlags sets up a standard set of arguments on fs feeding into the ActRunner a.
func ActRunnerFlags(fs *flag.FlagSet, a *act.Runner) {
	fs.StringVar(&a.ConfFile, flagActConfFile, "", usageActConfFile)
	fs.BoolVar(&a.DuneExec, flagActDuneExec, false, usageDuneExec)
}

// ConfFileFlag sets up a standard argument on fs for loading a configuration file into f.
func ConfFileFlag(fs *flag.FlagSet) *string {
	return fs.String(flagConfFile, "", usageConfFile)
}

// PlanFileFlag sets up a standard argument on fs for loading a plan file into f.
func PlanFileFlag(fs *flag.FlagSet, f *string) {
	fs.StringVar(f, FlagInputFile, "", usagePlanFile)
}

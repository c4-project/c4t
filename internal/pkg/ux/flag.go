package ux

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
)

const (
	// FlagInputFile is a standard flag for arguments that suggest an alternative to stdin for commands that read files.
	FlagInputFile = "i"

	// FlagNum is a standard flag for 'number of' arguments.
	FlagNum = "n"

	flagActConfFile = "C"
	flagActDuneExec = "x"
	flagOutDir      = "d"

	usageConfFile   = "read ACT config from this `file`"
	usageCorpusSize = "`number` of corpus files to select for this test plan;\n" +
		"if non-positive, the planner will use all viable provided corpus files"
	usageDuneExec = "if true, use 'dune exec' to run OCaml ACT binaries"
	usageOutDir   = "`directory` to which outputs will be written"
	usagePlanFile = "read from this plan `file` instead of stdin"
)

// CorpusSizeFlag sets up a 'target corpus size' flag.
func CorpusSizeFlag(out *int) {
	flag.IntVar(out, FlagNum, 0, usageCorpusSize)
}

// OutDirFlag sets up an 'output directory' flag.
func OutDirFlag(out *string, defaultdir string) {
	flag.StringVar(out, flagOutDir, defaultdir, usageOutDir)
}

// ActRunnerFlags sets up a standard set of arguments feeding into the ActRunner a.
func ActRunnerFlags(a *interop.ActRunner) {
	flag.StringVar(&a.ConfFile, flagActConfFile, "", usageConfFile)
	flag.BoolVar(&a.DuneExec, flagActDuneExec, false, usageDuneExec)
}

// PlanFileFlag sets up a standard argument for loading a plan file into f.
func PlanFileFlag(f *string) {
	flag.StringVar(f, FlagInputFile, "", usagePlanFile)
}

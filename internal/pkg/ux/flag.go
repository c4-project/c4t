package ux

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
)

const (
	// FlagInputFile is a standard flag for arguments that suggest an alternative to stdin for commands that read files.
	FlagInputFile = "i"

	// FlagMachine is a standard flag for machine selection arguments.
	FlagMachine = "m"

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

// CorpusSizeFlag sets up a 'target corpus size' flag on fs.
func CorpusSizeFlag(fs *flag.FlagSet, out *int) {
	fs.IntVar(out, FlagNum, 0, usageCorpusSize)
}

// OutDirFlag sets up an 'output directory' flag on fs.
func OutDirFlag(fs *flag.FlagSet, out *string, defaultdir string) {
	fs.StringVar(out, flagOutDir, defaultdir, usageOutDir)
}

// ActRunnerFlags sets up a standard set of arguments on fs feeding into the ActRunner a.
func ActRunnerFlags(fs *flag.FlagSet, a *interop.ActRunner) {
	fs.StringVar(&a.ConfFile, flagActConfFile, "", usageConfFile)
	fs.BoolVar(&a.DuneExec, flagActDuneExec, false, usageDuneExec)
}

// PlanFileFlag sets up a standard argument on fs for loading a plan file into f.
func PlanFileFlag(fs *flag.FlagSet, f *string) {
	fs.StringVar(f, FlagInputFile, "", usagePlanFile)
}

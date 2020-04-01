// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package view

import (
	"flag"
	"io"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"
	// It's 2020, and tools _still_ can't understand the use of 'v2' unless you do silly hacks like this.
	c "github.com/urfave/cli/v2"

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

	// FlagCompilerTimeoutLong is a long flag for compiler timeout.
	FlagCompilerTimeoutLong = "compiler-timeout"
	// FlagRunTimeoutLong is a long flag for run timeout.
	FlagRunTimeoutLong = "run-timeout"
	// FlagUseJSONLong is a long flag for JSON emission.
	FlagUseJSONLong = "emit-json"
	// FlagWorkerCountLong is a long flag for arguments that set a worker count.
	FlagWorkerCountLong = "num-workers"

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
	fs.StringVar(out, FlagOutDir, defaultdir, usageOutDir)
}

// OutDirCliFlag sets up an 'output directory' cli flag.
func OutDirCliFlag(defaultdir string) c.Flag {
	return &c.PathFlag{
		Name:  FlagOutDir,
		Value: defaultdir,
		Usage: usageOutDir,
	}
}

// OutDirFromCli gets the output directory set up by OutDirCliFlag.
func OutDirFromCli(ctx *c.Context) string {
	return ctx.Path(FlagOutDir)
}

// ActRunnerFlags sets up a standard set of arguments on fs feeding into the ActRunner a.
func ActRunnerFlags(fs *flag.FlagSet, a *act.Runner) {
	fs.StringVar(&a.ConfFile, flagActConfFile, "", usageActConfFile)
	fs.BoolVar(&a.DuneExec, flagActDuneExec, false, usageDuneExec)
}

// ActRunnerCliFlags gets the 'cli' flags needed to set up an ACT runner.
func ActRunnerCliFlags() []c.Flag {
	return []c.Flag{
		&c.PathFlag{
			Name:      flagActConfFile,
			Usage:     usageActConfFile,
			TakesFile: true,
		},
		&c.BoolFlag{
			Name:  flagActDuneExec,
			Usage: usageActConfFile,
		},
	}
}

// ActRunnerFromCli makes an ACT runner using the flags previously set up by ActRunnerCliFlags.
func ActRunnerFromCli(ctx *c.Context, errw io.Writer) *act.Runner {
	return &act.Runner{
		DuneExec: ctx.Bool(flagActDuneExec),
		ConfFile: ctx.Path(flagActConfFile),
		Stderr:   errw,
	}
}

// ConfFileCliFlag creates a cli flag for the config file.
func ConfFileCliFlag() c.Flag {
	return &c.PathFlag{
		Name:      flagConfFile,
		Usage:     usageActConfFile,
		TakesFile: true,
	}
}

// ConfFileFromCli sets up a Config using the file flag set up by ConfFileCliFlag.
func ConfFileFromCli(ctx *c.Context) (*config.Config, error) {
	cfile := ctx.Path(flagConfFile)
	return config.Load(cfile)
}

// ConfFileFlag sets up a standard argument on fs for loading a configuration file into f.
func ConfFileFlag(fs *flag.FlagSet) *string {
	return fs.String(flagConfFile, "", usageConfFile)
}

// PlanFileFlag sets up a standard argument on fs for loading a plan file into f.
func PlanFileFlag(fs *flag.FlagSet, f *string) {
	fs.StringVar(f, FlagInputFile, "", usagePlanFile)
}

// PlanFileCliFlag sets up a standard cli flag for loading a plan file into f.
func PlanFileCliFlag() c.Flag {
	return &c.PathFlag{
		Name:      FlagInputFile,
		TakesFile: true,
		Usage:     usagePlanFile,
	}
}

// PlanFileFromCli retrieves a plan file using the file flag set up by PlanFileCliFlag.
func PlanFileFromCli(ctx *c.Context) string {
	return ctx.Path(FlagInputFile)
}

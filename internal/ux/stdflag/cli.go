// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/config"

	// It's 2020, and tools _still_ can't understand the use of 'v2' unless you do silly hacks like this.
	c "github.com/urfave/cli/v2"
)

// OutDirCliFlag sets up an 'output directory' cli flag.
func OutDirCliFlag(defaultdir string) c.Flag {
	return &c.PathFlag{
		Name:  FlagOutDir,
		Value: defaultdir,
		Usage: usageOutDir,
	}
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
			Usage: usageActDuneExec,
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

// OutDirFromCli gets the output directory set up by OutDirCliFlag.
func OutDirFromCli(ctx *c.Context) string {
	return ctx.Path(FlagOutDir)
}

// ErrBadPlanArguments occurs when we expect a plan file argument, but get something else.
var ErrBadPlanArguments = errors.New("expected plan file argument")

// PlanFileFromCli retrieves a plan file (which may be empty) from the arguments of ctx.
// Its corresponding setup function is SetupPlanAppSettings; there is no 'plan file' flag.
func PlanFileFromCli(ctx *c.Context) (string, error) {
	args := ctx.Args()
	narg := args.Len()
	if 1 < narg {
		return "", fmt.Errorf("%w: got %d arguments, expected at most one", ErrBadPlanArguments, narg)
	}
	return args.First(), nil
}

// CorpusSizeCliFlag sets up a 'target corpus size' flag.
func CorpusSizeCliFlag() c.Flag {
	return &c.IntFlag{Name: FlagNum, Value: 0, Usage: usageCorpusSize}
}

// CorpusSizeFromCli retrieves a plan file using the file flag set up by CorpusSizeCliFlag.
func CorpusSizeFromCli(ctx *c.Context) int {
	return ctx.Int(FlagNum)
}

// SubjectCyclesCliFlag sets up a 'number of cycles' flag.
func SubjectCyclesCliFlag() c.Flag {
	return &c.IntFlag{Name: flagSubjectCycles, Value: fuzzer.DefaultSubjectCycles, Usage: usageSubjectCycles}
}

// SubjectCyclesFromCli retrieves a plan file using the file flag set up by SubjectCyclesCliFlag.
func SubjectCyclesFromCli(ctx *c.Context) int {
	return ctx.Int(flagSubjectCycles)
}

// CPUProfileCliFlag sets up a 'cpu profile dumper' flag.
func CPUProfileCliFlag() c.Flag {
	return &c.PathFlag{Name: FlagCPUProfile, Value: "", Usage: usageCPUProfile}
}

// CPUProfileFromCli retrieves the 'cpu profile dumper' set up by CPUProfileCliFlag.
func CPUProfileFromCli(ctx *c.Context) string {
	return ctx.Path(FlagCPUProfile)
}

// WorkerCountCliFlag sets up a worker count flag.
func WorkerCountCliFlag() c.Flag {
	return &c.IntFlag{
		Name:    FlagWorkerCountLong,
		Aliases: []string{FlagWorkerCount},
		Value:   1,
		Usage:   "number of `workers` to run in parallel",
	}
}

// WorkerCountFromCli retrieves a 'worker count' flag set up by WorkerCountCliFlag.
func WorkerCountFromCli(ctx *c.Context) int {
	return ctx.Int(FlagWorkerCountLong)
}
